package distribution

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

// Opt in with MYFINANCE_RELEASE_CANDIDATE (artifact directory) and
// WORKSPACE_RELEASE_BASELINE (verified public catalog.json). This reads release
// artifacts only; staging uses synthetic HTTP responses and t.TempDir, never an
// installed account, live credentials, or a network request.
func TestMyFinancePrepublicationManifestContract(t *testing.T) {
	root, baseline := os.Getenv("MYFINANCE_RELEASE_CANDIDATE"), os.Getenv("WORKSPACE_RELEASE_BASELINE")
	if root == "" || baseline == "" {
		t.Skip("set candidate artifact directory and verified public baseline catalog")
	}
	read := func(path string) []byte {
		t.Helper()
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		return body
	}
	raw := read(filepath.Join(root, "workspace-update.json"))
	var manifest AppUpdateManifest
	if err := json.Unmarshal(raw, &manifest); err != nil {
		t.Fatal(err)
	}
	current, err := catalog.Read(read(baseline))
	if err != nil {
		t.Fatal(err)
	}
	owned := []string{"finance", "finance-runtime", "finance-LICENSE", "finance-NOTICE", "finance-THIRD-PARTY-LICENSES.txt"}
	if manifest.Repository != "oneworkspacestudio/my-finance" || len(manifest.Apps) != 1 || manifest.Apps[0].ID != "finance" || len(manifest.Packages) != len(owned) {
		t.Fatal("candidate must describe My Finance and its five owned payloads")
	}
	for _, id := range owned {
		if _, ok := manifest.Packages[id]; !ok {
			t.Fatalf("missing owned candidate payload %s", id)
		}
		if _, ok := current.Packages[id]; !ok {
			t.Fatalf("public baseline must include owned Finance payload %s", id)
		}
	}
	// An installed port may differ from the publisher's default. Preserve it.
	for i := range current.Apps {
		if current.Apps[i].ID == "finance" {
			current.Apps[i].Port = 49024
		}
	}
	if err := current.Validate(); err != nil {
		t.Fatal(err)
	}
	before := catalogIdentity(current)
	release := Release{Tag: manifest.Tag, Published: time.Now().UTC()}
	archives := map[string]*zip.Reader{}
	for _, id := range owned {
		p := manifest.Packages[id]
		if filepath.Base(p.Asset) != p.Asset || p.Member == "" {
			t.Fatalf("invalid candidate archive location for %s", id)
		}
		z := archives[p.Asset]
		if z == nil {
			body := read(filepath.Join(root, p.Asset))
			z, err = zip.NewReader(bytes.NewReader(body), int64(len(body)))
			if err != nil {
				t.Fatal(err)
			}
			archives[p.Asset] = z
			release.Assets = append(release.Assets, fixtureAsset(manifest.Repository, manifest.Tag, p.Asset, fmt.Sprint(len(release.Assets)+1), body))
		}
		matches := 0
		for _, f := range z.File {
			if f.Name != p.Member {
				continue
			}
			matches++
			in, err := f.Open()
			if err != nil {
				t.Fatal(err)
			}
			h := sha256.New()
			n, err := io.Copy(h, in)
			in.Close()
			if err != nil || n != p.Size || hex.EncodeToString(h.Sum(nil)) != p.SHA256 || f.Mode()&os.ModeSymlink != 0 {
				t.Fatalf("candidate bytes do not satisfy %s contract: %v", id, err)
			}
		}
		if matches != 1 {
			t.Fatalf("candidate needs exactly one %s member; found %d", id, matches)
		}
	}
	g := updateGroup{"finance", "My Finance", "my-finance", []string{"finance"}}
	choice, err := manifestChoice(g, current, manifest.Repository, release, raw)
	if err != nil {
		t.Fatalf("candidate rejected before publication: %v", err)
	}
	choice.Available = true
	plan := UpdatePlan{BaseCatalogSHA256: before, Choices: []UpdateChoice{choice}}
	target, sources, err := SelectUpdates(current, plan, []string{"finance"})
	if err != nil {
		t.Fatal(err)
	}
	if len(choice.Packages) != len(owned) || target.Version != current.Version || target.InstallerIntegration != current.InstallerIntegration || len(target.Packages) != len(current.Packages) || len(target.Apps) != len(current.Apps) {
		t.Fatal("Finance-only selection changed the suite contract or lost owned payloads")
	}
	for id, old := range current.Packages {
		if candidate, ok := manifest.Packages[id]; ok {
			if target.Packages[id] != candidate.Package {
				t.Fatalf("selected Finance payload differs: %s", id)
			}
			if _, downloads := sources[id]; downloads != (old != candidate.Package) {
				t.Fatalf("download selection incorrect for %s", id)
			}
		} else if target.Packages[id] != old {
			t.Fatalf("unrelated package changed: %s", id)
		}
	}
	for _, old := range current.Apps {
		next, ok := target.Find(old.ID)
		want := old
		if old.ID == "finance" {
			want = manifest.Apps[0]
			want.Port = old.Port
		}
		if !ok || !reflect.DeepEqual(next, want) {
			t.Fatalf("unexpected app contract change: %s", old.ID)
		}
	}
	if catalogIdentity(current) != before {
		t.Fatal("planning mutated the baseline")
	}
	t.Run("setup-selection-retains-runtime", func(t *testing.T) {
		selected := map[string]bool{}
		for _, app := range target.Apps {
			selected[app.ID] = true
		}
		included := catalog.InstallPackageIDs(target, selected)
		delete(selected, "finance")
		omitted := catalog.InstallPackageIDs(target, selected)
		for _, id := range owned {
			if !included[id] {
				t.Fatalf("selected Finance lost %s", id)
			}
			// Setup retains the suite's licence notices even for omitted apps.
			wantOmitted := strings.HasPrefix(target.Packages[id].File, "licenses/")
			if omitted[id] != wantOmitted {
				t.Fatalf("omitted Finance payload selection incorrect: %s", id)
			}
		}
		for id, keep := range included {
			if _, finance := manifest.Packages[id]; !finance && omitted[id] != keep {
				t.Fatalf("omitting Finance changed unrelated selection: %s", id)
			}
		}
	})
	t.Run("missing-runtime-rejected", func(t *testing.T) {
		var broken AppUpdateManifest
		if err := json.Unmarshal(raw, &broken); err != nil {
			t.Fatal(err)
		}
		delete(broken.Packages, "finance-runtime")
		data, err := json.Marshal(broken)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := manifestChoice(g, current, broken.Repository, release, data); err == nil || !strings.Contains(err.Error(), "finance-runtime") {
			t.Fatalf("missing Actual runtime was not rejected: %v", err)
		}
	})
	t.Run("same-size-tampered-member-rejected", func(t *testing.T) {
		// A correctly hashed archive cannot override its member's contract.
		// Use a small synthetic licence member to exercise production staging.
		p := manifest.Packages["finance-LICENSE"]
		var body bytes.Buffer
		z := zip.NewWriter(&body)
		out, err := z.Create(p.Member)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = out.Write(bytes.Repeat([]byte("!"), int(p.Size))); err != nil {
			t.Fatal(err)
		}
		if err = z.Close(); err != nil {
			t.Fatal(err)
		}
		asset := fixtureAsset(manifest.Repository, manifest.Tag, p.Asset, "99", body.Bytes())
		original := http.DefaultTransport
		http.DefaultTransport = updateTransport{asset.APIURL: body.Bytes(), asset.URL: body.Bytes()}
		t.Cleanup(func() { http.DefaultTransport = original })
		t.Setenv("GH_TOKEN", "")
		t.Setenv("GITHUB_TOKEN", "")
		download := UpdateDownload{Asset: asset, Repository: manifest.Repository, Tag: manifest.Tag, Member: p.Member}
		err = StageUpdateFiles(context.Background(), target, map[string]UpdateDownload{"finance-LICENSE": download}, t.TempDir(), func(string) {})
		if err == nil || !strings.Contains(err.Error(), "checksum") {
			t.Fatalf("same-size corrupted Finance member was not rejected by its checksum: %v", err)
		}
	})
}
