package distribution

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

// Opt-in release acceptance, using public catalog/asset metadata rather than an
// installed user's configuration. No owner data or credentials enter the test.
func TestSlicerPublicationScope(t *testing.T) {
	root := os.Getenv("WORKSPACE_SLICER_PUBLICATION_ROOT")
	if root == "" {
		t.Skip("public release acceptance inputs required")
	}
	read := func(name string) []byte {
		b, e := os.ReadFile(filepath.Join(root, name))
		if e != nil {
			t.Fatal(e)
		}
		return b
	}
	var base catalog.Catalog
	if e := json.Unmarshal(read("catalog113.json"), &base); e != nil {
		t.Fatal(e)
	}
	var release Release
	if e := json.Unmarshal(read("metadata-release.json"), &release); e != nil {
		t.Fatal(e)
	}
	var group updateGroup
	for _, g := range updateGroups {
		if g.id == "slicer" {
			group = g
		}
	}
	choice, e := manifestChoice(group, base, "oneworkspacestudio/my-slicer", release, read("assets8/workspace-update.json"))
	if e != nil {
		t.Fatal(e)
	}
	if choice.Version != "0.1.0-preview.8" || len(choice.Packages) != 8 {
		t.Fatal("incomplete Slicer ownership")
	}
	choice.Available = true
	target, downloads, e := SelectUpdates(base, UpdatePlan{BaseCatalogSHA256: catalogIdentity(base), Choices: []UpdateChoice{choice}}, []string{"slicer"})
	if e != nil {
		t.Fatal(e)
	}
	if target.Version != base.Version || len(downloads) < 1 || len(downloads) > 8 {
		t.Fatal("wrong app-only scope")
	}
	for _, app := range base.Apps {
		if app.ID == "slicer" {
			continue
		}
		now, ok := target.Find(app.ID)
		if !ok || !reflect.DeepEqual(app, now) {
			t.Fatalf("unrelated app changed: %s", app.ID)
		}
	}
	for key, pkg := range base.Packages {
		if _, selected := choice.Packages[key]; !selected && !reflect.DeepEqual(pkg, target.Packages[key]) {
			t.Fatalf("unrelated package changed: %s", key)
		}
	}
	t.Log("All eight Slicer packages selected; other36 app contracts, packages and manager retained")
	if os.Getenv("WORKSPACE_SLICER_ANONYMOUS") != "1" {
		return
	}
	t.Setenv("GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", "")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	plan, e := CheckUpdates(ctx, base)
	if e != nil {
		t.Fatal(e)
	}
	found := false
	for _, c := range plan.Choices {
		if c.ID == "slicer" {
			if !c.Available || c.Version != "0.1.0-preview.8" {
				t.Fatal("public channel did not select Slicer8", c.Reason)
			}
			found = true
		}
	}
	if !found {
		t.Fatal("public channel omitted Slicer")
	}
	selected, files, e := SelectUpdates(base, plan, []string{"slicer"})
	if e != nil {
		t.Fatal(e)
	}
	tempBase := os.Getenv("WORKSPACE_SLICER_TEST_TEMP")
	if tempBase == "" {
		t.Fatal("named isolated NTFS staging required")
	}
	dir, e := os.MkdirTemp(tempBase, "public-download-")
	if e != nil {
		t.Fatal(e)
	}
	defer os.RemoveAll(dir)
	if e = StageUpdateFiles(ctx, selected, files, dir, func(string) {}); e != nil {
		t.Fatal(e)
	}
	proof := map[string]any{"schema": 1, "passed": true, "authentication": "anonymous public channel and release assets", "slicer": "0.1.0-preview.8", "ownedPackages": 8, "selectedDownloads": len(files), "managerPreserved": selected.Version == base.Version, "installed": false}
	b, _ := json.MarshalIndent(proof, "", "  ")
	if e = os.WriteFile(filepath.Join(root, "anonymous-public-check.json"), b, 0600); e != nil {
		t.Fatal(e)
	}
	t.Log("Anonymous channel, asset download and every selected payload digest verified")
}
