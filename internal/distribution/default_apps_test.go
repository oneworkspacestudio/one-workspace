package distribution

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"testing"
)

func fourDefaultCatalogs() (catalog.Catalog, catalog.Catalog) {
	old := updateCatalogFixture()
	old.Packages["accounts"] = fixturePackage("MyAccounts.exe", []byte("accounts-original"))
	old.Apps = append(old.Apps, catalog.App{ID: "accounts", Name: "MyAccounts", Version: "fixture", Kind: "process", Package: "accounts", Executable: "MyAccounts.exe", Port: 25101, Source: "https://github.com/oneworkspacestudio/myaccounts/releases/tag/fixture"})
	incoming := old
	incoming.Version = "0.3.0-preview.23"
	incoming.Apps = append([]catalog.App(nil), old.Apps...)
	incoming.Packages = map[string]catalog.Package{}
	for id, p := range old.Packages {
		incoming.Packages[id] = p
	}
	incoming.Packages["manager"] = fixturePackage("Workspace.exe", []byte("new-manager"))
	for i, id := range []string{"tasks", "pdf", "monitor", "search", "slicer", "social", "clips"} {
		name := map[string]string{"tasks": "MyTasks", "pdf": "MyPDF", "monitor": "MyMonitor", "search": "MySearch", "slicer": "MySlicer", "social": "MySocial", "clips": "MyClips"}[id]
		incoming.Packages[id] = fixturePackage(name+".exe", []byte("new-"+id))
		for _, notice := range []string{"LICENSE", "NOTICE", "THIRD-PARTY-LICENSES.txt"} {
			incoming.Packages[id+"-"+notice] = fixturePackage("licenses/"+name+"-"+notice, []byte("license-"+notice))
		}
		incoming.Apps = append(incoming.Apps, catalog.App{ID: id, Name: name, Version: "0.1.0-preview.1", Kind: "process", Package: id, Executable: name + ".exe", Port: 9008 + i, Route: "/" + id, StripPrefix: true, OIDC: true, Callback: "/" + id + "/auth/callback", Depends: []string{"accounts"}, Source: "https://github.com/oneworkspacestudio/my-" + id + "/releases/tag/v0.1.0-preview.1"})
	}
	return old, incoming
}
func TestManagerUpdateAutomaticallyAddsNewDefaultsAndPreservesOptOuts(t *testing.T) {
	old, incoming := fourDefaultCatalogs()
	for _, excluded := range [][]string{nil, {"pdf", "search"}} {
		current := old
		if excluded != nil {
			current.Selection = &catalog.AppSelection{Excluded: excluded}
		}
		target, e := SelectBundleUpdates(current, incoming, []string{"workspace"})
		if e != nil {
			t.Fatal(e)
		}
		for i, id := range []string{"tasks", "pdf", "monitor", "search", "slicer", "social", "clips"} {
			app, installed := target.Find(id)
			want := !(id == "pdf" || id == "search") || excluded == nil
			if installed != want {
				t.Fatalf("%s installed=%t wanted=%t", id, installed, want)
			}
			if installed {
				if app.Port != 25118+i || !app.OIDC || app.Callback != "/"+id+"/auth/callback" {
					t.Fatal("new app launch/OIDC contract", app)
				}
				for _, notice := range []string{"LICENSE", "NOTICE", "THIRD-PARTY-LICENSES.txt"} {
					if _, ok := target.Packages[id+"-"+notice]; !ok {
						t.Fatal("notice missing", id, notice)
					}
				}
			}
		}
		notes, _ := target.Find("notes")
		if notes.Port != old.Apps[0].Port {
			t.Fatal("changed existing custom port")
		}
		if len(old.Apps) != 2 || old.Selection != nil {
			t.Fatal("modified installed input")
		}
	}
}
func TestOnlinePlanDiscoversSameFourNewApplications(t *testing.T) {
	_, tr, _ := setupUpdateTransport(t, true)
	old, incoming := fourDefaultCatalogs()
	var releases []Release
	uri := "https://api.github.com/repos/" + Repository + "/releases?per_page=100"
	json.Unmarshal(tr[uri], &releases)
	raw, _ := json.Marshal(incoming)
	asset := fixtureAsset(Repository, releases[0].Tag, "catalog.json", "1", raw)
	releases[0].Assets[0] = asset
	tr[asset.APIURL] = raw
	tr[uri], _ = json.Marshal(releases)
	tr["https://api.github.com/repos/oneworkspacestudio/myaccounts/releases?per_page=100"] = []byte(`[]`)
	plan, e := CheckUpdates(context.Background(), old)
	if e != nil {
		t.Fatal(e)
	}
	target, downloads, e := SelectUpdates(old, plan, []string{"workspace"})
	if e != nil {
		t.Fatal(e)
	}
	for _, id := range []string{"tasks", "pdf", "monitor", "search", "slicer", "social", "clips"} {
		if _, ok := target.Find(id); !ok {
			t.Fatal("new app missing", id)
		}
		if _, ok := downloads[id]; !ok {
			t.Fatal("new executable not downloaded", id)
		}
	}
	if len(downloads) < 29 {
		t.Fatal(fmt.Sprintf("expected manager plus executable/notices for seven apps; got %d", len(downloads)))
	}
}
