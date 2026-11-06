package distribution

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

// Release operators run this against the actual local package before publishing.
// It uses public catalog metadata and synthetic asset receipts, never an install.
func TestMyAIPrepublicationManifestContract(t *testing.T) {
	root, baseline := os.Getenv("MYAI_RELEASE_CANDIDATE"), os.Getenv("WORKSPACE_RELEASE_BASELINE")
	if root == "" || baseline == "" {
		t.Skip("set candidate artifact directory and verified public baseline catalog")
	}
	raw, err := os.ReadFile(filepath.Join(root, "workspace-update.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest AppUpdateManifest
	if err = json.Unmarshal(raw, &manifest); err != nil {
		t.Fatal(err)
	}
	base, err := os.ReadFile(baseline)
	if err != nil {
		t.Fatal(err)
	}
	release := Release{Tag: manifest.Tag, Published: time.Now().UTC()}
	seen := map[string]bool{}
	for _, p := range manifest.Packages {
		if seen[p.Asset] {
			continue
		}
		if filepath.Base(p.Asset) != p.Asset {
			t.Fatal("candidate asset must be a basename")
		}
		body, e := os.ReadFile(filepath.Join(root, p.Asset))
		if e != nil {
			t.Fatal(e)
		}
		release.Assets = append(release.Assets, fixtureAsset(manifest.Repository, manifest.Tag, p.Asset, "1", body))
		seen[p.Asset] = true
	}
	for _, minimal := range []bool{false, true} {
		var current catalog.Catalog
		if err = json.Unmarshal(base, &current); err != nil {
			t.Fatal(err)
		}
		if minimal {
			for id := range current.Packages {
				if strings.HasPrefix(id, "myai-") && id != "myai-LICENSE" && id != "myai-NOTICE" && id != "myai-THIRD-PARTY-LICENSES.txt" {
					delete(current.Packages, id)
				}
			}
		}
		candidate := manifest
		if err = RequireManager(&candidate, current.Version, nil); err != nil {
			t.Fatal(err)
		}
		g := updateGroup{"myai", "My AI", "my-ai", []string{"myai", "aiservice"}}
		choice, e := manifestChoice(g, current, manifest.Repository, release, raw)
		if e != nil {
			t.Fatalf("candidate rejected before publication (minimal=%v): %v", minimal, e)
		}
		choice.Available = true
		plan := UpdatePlan{BaseCatalogSHA256: catalogIdentity(current), Choices: []UpdateChoice{choice}}
		target, _, e := SelectUpdates(current, plan, []string{"myai"})
		if e != nil {
			t.Fatal(e)
		}
		for _, app := range current.Apps {
			next, ok := target.Find(app.ID)
			if !ok || (app.ID != "myai" && app.ID != "aiservice" && !reflect.DeepEqual(app, next)) {
				t.Fatal("unrelated app changed", app.ID)
			}
		}
		for id, p := range current.Packages {
			if id != "myai" && !strings.HasPrefix(id, "myai-") && id != "aiservice" && !strings.HasPrefix(id, "aiservice-") && target.Packages[id] != p {
				t.Fatal("unrelated payload changed", id)
			}
		}
	}
}
