package distribution

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

type updateTransport map[string][]byte

func (tr updateTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	body, ok := tr[r.URL.String()]
	status := 200
	if !ok {
		status = 404
		body = []byte(`{"message":"fixture not found"}`)
	}
	return &http.Response{StatusCode: status, Header: http.Header{}, Body: io.NopCloser(bytes.NewReader(body)), Request: r}, nil
}
func fixtureAsset(repo, tag, name, id string, body []byte) Asset {
	sum := sha256.Sum256(body)
	return Asset{Name: name, Size: int64(len(body)), Digest: "sha256:" + hex.EncodeToString(sum[:]), State: "uploaded", APIURL: "https://api.github.com/repos/" + repo + "/releases/assets/" + id, URL: "https://github.com/" + repo + "/releases/download/" + tag + "/" + name}
}
func fixturePackage(file string, body []byte) catalog.Package {
	sum := sha256.Sum256(body)
	return catalog.Package{File: file, Size: int64(len(body)), SHA256: hex.EncodeToString(sum[:])}
}
func updateCatalogFixture() catalog.Catalog {
	return catalog.Catalog{Schema: 1, Version: "0.3.0-preview.22", InstallerIntegration: 2, Packages: map[string]catalog.Package{"manager": fixturePackage("Workspace.exe", []byte("old-manager")), "notes": fixturePackage("MyNotes.exe", []byte("old-notes")), "notes-LICENSE": fixturePackage("licenses/Notes-LICENSE", []byte("license"))}, Apps: []catalog.App{{ID: "notes", Name: "Notes", Kind: "process", Package: "notes", Executable: "MyNotes.exe", Port: 25112, Route: "/notes", Version: "0.1.0-preview.4", Source: "https://github.com/oneworkspacestudio/notes/releases/tag/v0.1.0-preview.4"}}}
}
func setupUpdateTransport(t *testing.T, withMarker bool) (catalog.Catalog, updateTransport, AppUpdateManifest) {
	t.Helper()
	old := updateCatalogFixture()
	latest := updateCatalogFixture()
	latest.Version = "0.3.0-preview.23"
	latest.Packages["manager"] = fixturePackage("Workspace.exe", []byte("new-manager"))
	raw, _ := json.Marshal(latest)
	suiteTag := "v" + latest.Version
	cat := fixtureAsset(Repository, suiteTag, "catalog.json", "1", raw)
	manager := fixtureAsset(Repository, suiteTag, "Workspace.exe", "2", []byte("new-manager"))
	bundle := fixtureAsset(Repository, suiteTag, "workspace-bundle.zip", "3", []byte("bundle"))
	suiteRelease := Release{Tag: suiteTag, Published: time.Unix(20, 0), Assets: []Asset{cat, manager, bundle}}
	oldRelease := Release{Tag: "v" + old.Version, Published: time.Unix(10, 0)}
	repo := "oneworkspacestudio/notes"
	tag := "v0.1.0-preview.5"
	app := old.Apps[0]
	app.Version = "0.1.0-preview.5"
	app.Source = "https://github.com/" + repo + "/releases/tag/" + tag
	exe := fixtureAsset(repo, tag, "MyNotes.exe", "10", []byte("new-notes"))
	license := fixtureAsset(repo, tag, "LICENSE", "11", []byte("license"))
	manifest := AppUpdateManifest{Schema: 1, Integration: 2, Repository: repo, Tag: tag, Apps: []catalog.App{app}, Packages: map[string]UpdatePackage{"notes": {Package: fixturePackage("MyNotes.exe", []byte("new-notes")), Asset: "MyNotes.exe"}, "notes-LICENSE": {Package: old.Packages["notes-LICENSE"], Asset: "LICENSE"}}}
	markerRaw, _ := json.Marshal(manifest)
	marker := fixtureAsset(repo, tag, "workspace-update.json", "12", markerRaw)
	assets := []Asset{exe, license}
	if withMarker {
		assets = append(assets, marker)
	}
	appRelease := Release{Tag: tag, Published: time.Unix(30, 0), Assets: assets}
	tr := updateTransport{cat.APIURL: raw, manager.APIURL: []byte("new-manager"), exe.APIURL: []byte("new-notes"), license.APIURL: []byte("license"), marker.APIURL: markerRaw}
	tr["https://api.github.com/repos/"+Repository+"/releases?per_page=100"], _ = json.Marshal([]Release{suiteRelease, oldRelease})
	tr["https://api.github.com/repos/"+repo+"/releases?per_page=100"], _ = json.Marshal([]Release{appRelease, {Tag: "v0.1.0-preview.4", Published: time.Unix(5, 0)}})
	original := http.DefaultTransport
	http.DefaultTransport = tr
	t.Cleanup(func() { http.DefaultTransport = original })
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	return old, tr, manifest
}
func TestDiscoverSelectAndStageIndependentAppAndWorkspace(t *testing.T) {
	current, _, _ := setupUpdateTransport(t, true)
	plan, e := CheckUpdates(context.Background(), current)
	if e != nil {
		t.Fatal(e)
	}
	target, sources, e := SelectUpdates(current, plan, []string{"workspace", "notes"})
	if e != nil {
		t.Fatal(e)
	}
	if target.Version != "0.3.0-preview.23" || target.Apps[0].Version != "0.1.0-preview.5" || target.Apps[0].Port != 25112 {
		t.Fatal(target)
	}
	if current.Version != "0.3.0-preview.22" || current.Packages["notes"].SHA256 == target.Packages["notes"].SHA256 {
		t.Fatal("installed catalog mutated")
	}
	if len(sources) != 2 {
		t.Fatal("unchanged license downloaded", sources)
	}
	dir := t.TempDir()
	if e = StageUpdateFiles(context.Background(), target, sources, dir, func(string) {}); e != nil {
		t.Fatal(e)
	}
	if sum, _, e := digestFixture(filepath.Join(dir, "packages", "MyNotes.exe")); e != nil || sum != target.Packages["notes"].SHA256 {
		t.Fatal(sum, e)
	}
}
func digestFixture(path string) (string, int64, error) {
	f, e := os.Open(path)
	if e != nil {
		return "", 0, e
	}
	defer f.Close()
	h := sha256.New()
	n, e := io.Copy(h, f)
	return hex.EncodeToString(h.Sum(nil)), n, e
}
func TestUncontractedNewReleaseIsVisibleButNotInstalled(t *testing.T) {
	current, _, _ := setupUpdateTransport(t, false)
	plan, e := CheckUpdates(context.Background(), current)
	if e != nil {
		t.Fatal(e)
	}
	var notes UpdateChoice
	for _, ch := range plan.Choices {
		if ch.ID == "notes" {
			notes = ch
		}
	}
	if notes.Available || !strings.Contains(notes.Reason, "v0.1.0-preview.5") {
		t.Fatal(notes)
	}
	if _, _, e = SelectUpdates(current, plan, []string{"notes"}); e == nil {
		t.Fatal("uncontracted executable selected")
	}
}
func TestIndependentContractCannotReplaceOtherAppsOrEscapeArchive(t *testing.T) {
	current, tr, m := setupUpdateTransport(t, true)
	var releases []Release
	json.Unmarshal(tr["https://api.github.com/repos/oneworkspacestudio/notes/releases?per_page=100"], &releases)
	r := releases[0]
	g := updateGroup{}; for _, candidate := range updateGroups { if candidate.id == "notes" { g=candidate } }
	for _, mutate := range []func(*AppUpdateManifest){func(m *AppUpdateManifest) { m.Repository = "someone/notes" }, func(m *AppUpdateManifest) { m.Apps[0].ID = "accounts" }, func(m *AppUpdateManifest) {
		p := m.Packages["notes"]
		p.File = "MyAccounts.exe"
		m.Packages["notes"] = p
	}, func(m *AppUpdateManifest) {
		p := m.Packages["notes"]
		p.Member = "../outside.exe"
		m.Packages["notes"] = p
	}, func(m *AppUpdateManifest) { delete(m.Packages, "notes-LICENSE") }} {
		raw, _ := json.Marshal(m)
		var copy AppUpdateManifest
		json.Unmarshal(raw, &copy)
		mutate(&copy)
		raw, _ = json.Marshal(copy)
		if _, e := manifestChoice(g, current, "oneworkspacestudio/notes", r, raw); e == nil {
			t.Fatal("unsafe contract accepted", string(raw))
		}
	}
}
func TestStagedDownloadChecksumAndArchiveMemberChecks(t *testing.T) {
	current, tr, _ := setupUpdateTransport(t, true)
	plan, e := CheckUpdates(context.Background(), current)
	if e != nil {
		t.Fatal(e)
	}
	target, sources, e := SelectUpdates(current, plan, []string{"notes"})
	if e != nil {
		t.Fatal(e)
	}
	d := sources["notes"]
	tr[d.Asset.APIURL] = []byte("tampered")
	if e = StageUpdateFiles(context.Background(), target, sources, t.TempDir(), func(string) {}); e == nil {
		t.Fatal("tampered executable staged")
	}
	for _, duplicate := range []bool{false, true} {
		var buf bytes.Buffer
		z := zip.NewWriter(&buf)
		name := "nested/MyNotes.exe"
		f, _ := z.Create(name)
		f.Write([]byte("new-notes"))
		if duplicate {
			f, _ = z.Create(name)
			f.Write([]byte("new-notes"))
		}
		z.Close()
		a := fixtureAsset("oneworkspacestudio/notes", "v0.1.0-preview.5", "portable.zip", "99", buf.Bytes())
		tr[a.APIURL] = buf.Bytes()
		d.Asset = a
		d.Member = name
		sources["notes"] = d
		e = StageUpdateFiles(context.Background(), target, sources, t.TempDir(), func(string) {})
		if duplicate && e == nil {
			t.Fatal("duplicate member accepted")
		}
		if !duplicate && e != nil {
			t.Fatal(e)
		}
	}
}
func TestWorkspaceVersionAndSelectionGuards(t *testing.T) {
	for _, tc := range []struct {
		next, old string
		want      bool
	}{{"0.3.0-preview.23", "0.3.0-preview.20-fixes", true}, {"0.3.0-preview.22", "0.3.0-preview.23", false}, {"0.4.0", "0.3.0-preview.23", true}, {"0.3.0-preview.23", "development", false}} {
		if got := newerWorkspace(tc.next, tc.old); got != tc.want {
			t.Fatal(tc, got)
		}
	}
	current, _, _ := setupUpdateTransport(t, true)
	plan, e := CheckUpdates(context.Background(), current)
	if e != nil {
		t.Fatal(e)
	}
	for _, selection := range [][]string{{"notes", "notes"}, {"unknown"}, {}} {
		if _, _, e = SelectUpdates(current, plan, selection); e == nil {
			t.Fatal("invalid selection accepted", selection)
		}
	}
	changed := current
	changed.Version = "0.3.0-preview.24"
	if _, _, e = SelectUpdates(changed, plan, []string{"workspace"}); e == nil {
		t.Fatal("stale plan could downgrade installation")
	}
}
func TestPublisherContractRoundTripAndCustomExtensionOwnership(t *testing.T) {
	current, tr, m := setupUpdateTransport(t, true)
	current.Packages["notes-addon"] = fixturePackage("CustomNotes.exe", []byte("custom extension"))
	current.Apps = append(current.Apps, catalog.App{ID: "notes-addon", Name: "Custom notes helper", Version: "1.0.0", Kind: "process", Package: "notes-addon", Executable: "CustomNotes.exe", Port: 25119, Route: "/custom-notes"})
	generated, e := MakeAppUpdateManifest(current, "notes", "0.1.0-preview.5", "oneworkspacestudio/notes", "v0.1.0-preview.5", m.Packages)
	if e != nil {
		t.Fatal(e)
	}
	raw, _ := json.Marshal(generated)
	var releases []Release
	json.Unmarshal(tr["https://api.github.com/repos/oneworkspacestudio/notes/releases?per_page=100"], &releases)
	var notesGroup updateGroup
	for _, group := range updateGroups {
		if group.id == "notes" {
			notesGroup = group
		}
	}
	ch, e := manifestChoice(notesGroup, current, generated.Repository, releases[0], raw)
	if e != nil {
		t.Fatal(e)
	}
	ch.Available = true
	plan := UpdatePlan{BaseCatalogSHA256: catalogIdentity(current), Choices: []UpdateChoice{ch}}
	target, _, e := SelectUpdates(current, plan, []string{"notes"})
	if e != nil {
		t.Fatal(e)
	}
	if target.Packages["notes-addon"] != current.Packages["notes-addon"] || len(target.Apps) != 2 {
		t.Fatal("custom extension was replaced")
	}
}
