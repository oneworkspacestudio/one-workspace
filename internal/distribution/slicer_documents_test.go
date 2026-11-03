package distribution

import (
	"encoding/json"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"testing"
)

func TestSlicerUpdateAcrossLegacyAndFullDocumentLayouts(t *testing.T) {
	for _, tracked := range []bool{false, true} {
		t.Run(fmt.Sprint(tracked), func(t *testing.T) {
			repo, tag := "oneworkspacestudio/my-slicer", "v0.1.0-preview.8"
			current := catalog.Catalog{Schema: 1, Version: "0.3.0-preview.114", Packages: map[string]catalog.Package{"manager": fixturePackage("Workspace.exe", []byte("manager"))}, Apps: []catalog.App{{ID: "slicer", Name: "My Slicer", Kind: "process", Package: "slicer", Executable: "MySlicer.exe", Port: 9019, Route: "/slicer", Version: "0.1.0-preview.7", Source: "https://github.com/" + repo + "/releases/tag/v0.1.0-preview.7"}}}
			files := map[string]string{"slicer": "MySlicer.exe", "slicer-LICENSE": "licenses/MySlicer-LICENSE", "slicer-NOTICE": "licenses/MySlicer-NOTICE", "slicer-THIRD-PARTY-LICENSES.txt": "licenses/MySlicer-THIRD-PARTY-LICENSES.txt", "slicer-contract": "docs/MySlicer/compatibility.json"}
			for _, name := range []string{"SETUP.md", "ACCEPTANCE.md", "ENGINE.md"} {
				files["slicer-doc-"+name] = "docs/MySlicer/" + name
			}
			app := current.Apps[0]
			app.Version = "0.1.0-preview.8"
			app.Source = "https://github.com/" + repo + "/releases/tag/" + tag
			m := AppUpdateManifest{Repository: repo, Tag: tag, Apps: []catalog.App{app}, Packages: map[string]UpdatePackage{}}
			release := Release{Tag: tag}
			for id, file := range files {
				runtime := id == "slicer" || id == "slicer-LICENSE" || id == "slicer-NOTICE" || id == "slicer-THIRD-PARTY-LICENSES.txt"
				if runtime || tracked {
					current.Packages[id] = fixturePackage(file, []byte("old "+id))
				}
				body := []byte("new " + id)
				asset := fixtureAsset(repo, tag, id+".bin", fmt.Sprint(len(release.Assets)+1), body)
				release.Assets = append(release.Assets, asset)
				m.Packages[id] = UpdatePackage{Package: fixturePackage(file, body), Asset: asset.Name}
			}
			if e := RequireManager(&m, "0.3.0-preview.114", nil); e != nil {
				t.Fatal(e)
			}
			g := updateGroup{"slicer", "My Slicer", "my-slicer", []string{"slicer"}}
			raw, _ := json.Marshal(m)
			choice, e := manifestChoice(g, current, repo, release, raw)
			if e != nil {
				t.Fatal(e)
			}
			choice.Available = true
			plan := UpdatePlan{BaseCatalogSHA256: catalogIdentity(current), Choices: []UpdateChoice{choice}}
			target, sources, e := SelectUpdates(current, plan, []string{"slicer"})
			if e != nil {
				t.Fatal(e)
			}
			want := 4
			if tracked {
				want = 8
			}
			if len(choice.Packages) != want || len(sources) != want {
				t.Fatal("wrong scope", len(choice.Packages), len(sources))
			}
			if target.Apps[0].Version != "0.1.0-preview.8" || target.Packages["manager"] != current.Packages["manager"] || target.Version != current.Version {
				t.Fatal("unrelated manager changed")
			}
			for id := range files {
				_, ok := target.Packages[id]
				if !tracked && id != "slicer" && id != "slicer-LICENSE" && id != "slicer-NOTICE" && id != "slicer-THIRD-PARTY-LICENSES.txt" && ok {
					t.Fatal("untracked document installed", id)
				}
			}
			for _, bad := range []struct{ id, file string }{{"slicer-runtime", "docs/MySlicer/SETUP.md"}, {"slicer-doc-SETUP.md", "MySlicer.exe"}, {"slicer-doc-SETUP.md", "docs/Other/SETUP.md"}, {"slicer-doc-ENGINE.md", "docs/MySlicer/../ENGINE.md"}, {"other-doc-SETUP.md", "docs/MySlicer/SETUP.md"}} {
				copy := m
				copy.Packages = map[string]UpdatePackage{}
				for id, p := range m.Packages {
					copy.Packages[id] = p
				}
				copy.Packages[bad.id] = UpdatePackage{Package: fixturePackage(bad.file, []byte("unexpected")), Asset: "unexpected.bin"}
				b, _ := json.Marshal(copy)
				if _, e := manifestChoice(g, current, repo, release, b); e == nil {
					t.Fatal("foreign/runtime/moved package accepted", bad)
				}
			}
			if tracked {
				delete(m.Packages, "slicer-doc-SETUP.md")
				b, _ := json.Marshal(m)
				if _, e := manifestChoice(g, current, repo, release, b); e == nil {
					t.Fatal("tracked guide disappeared")
				}
			}
		})
	}
}
