package distribution

import (
	"encoding/json"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"testing"
)

func TestCinemaUpdateAcrossDefaultAndFullDocumentLayouts(t *testing.T) {
	for _, tracked := range []bool{false, true} {
		t.Run(fmt.Sprint(tracked), func(t *testing.T) {
			repo, tag := "oneworkspacestudio/my-cinema", "v0.1.0-preview.4"
			current := catalog.Catalog{Schema: 1, Version: "0.3.0-preview.85", Packages: map[string]catalog.Package{"manager": fixturePackage("Workspace.exe", []byte("manager"))}, Apps: []catalog.App{{ID: "cinema", Name: "My Cinema", Kind: "process", Package: "cinema", Executable: "MyCinema.exe", Port: 9020, Route: "/cinema", Version: "0.1.0-preview.2", Source: "https://github.com/" + repo + "/releases/tag/v0.1.0-preview.2"}}}
			files := map[string]string{"cinema": "MyCinema.exe", "cinema-LICENSE": "licenses/MyCinema-LICENSE", "cinema-NOTICE": "licenses/MyCinema-NOTICE", "cinema-THIRD-PARTY-LICENSES.txt": "licenses/MyCinema-THIRD-PARTY-LICENSES.txt", "cinema-contract": "docs/MyCinema/compatibility.json", "cinema-provenance": "docs/MyCinema/SOURCE-PROVENANCE.json"}
			for _, name := range []string{"SETUP.md", "SUPPORT.md", "ADR-001-architecture.md", "SECURITY-CHECKS.md"} {
				files["cinema-doc-"+name] = "docs/MyCinema/" + name
			}
			app := current.Apps[0]
			app.Version = "0.1.0-preview.4"
			app.Source = "https://github.com/" + repo + "/releases/tag/" + tag
			m := AppUpdateManifest{Repository: repo, Tag: tag, Apps: []catalog.App{app}, Packages: map[string]UpdatePackage{}}
			release := Release{Tag: tag}
			for id, file := range files {
				runtime := id == "cinema" || id == "cinema-LICENSE" || id == "cinema-NOTICE" || id == "cinema-THIRD-PARTY-LICENSES.txt"
				if runtime || tracked {
					current.Packages[id] = fixturePackage(file, []byte("old "+id))
				}
				body := []byte("new " + id)
				asset := fixtureAsset(repo, tag, id+".bin", fmt.Sprint(len(release.Assets)+1), body)
				release.Assets = append(release.Assets, asset)
				m.Packages[id] = UpdatePackage{Package: fixturePackage(file, body), Asset: asset.Name}
			}
			if e := RequireManager(&m, "0.3.0-preview.85", nil); e != nil {
				t.Fatal(e)
			}
			g := updateGroup{"cinema", "My Cinema", "my-cinema", []string{"cinema"}}
			raw, _ := json.Marshal(m)
			choice, e := manifestChoice(g, current, repo, release, raw)
			if e != nil {
				t.Fatal(e)
			}
			choice.Available = true
			plan := UpdatePlan{BaseCatalogSHA256: catalogIdentity(current), Choices: []UpdateChoice{choice}}
			target, sources, e := SelectUpdates(current, plan, []string{"cinema"})
			if e != nil {
				t.Fatal(e)
			}
			want := 4
			if tracked {
				want = 10
			}
			if len(choice.Packages) != want || len(sources) != want {
				t.Fatal("wrong scope", len(choice.Packages), len(sources))
			}
			if target.Apps[0].Version != "0.1.0-preview.4" || target.Packages["manager"] != current.Packages["manager"] || target.Version != current.Version {
				t.Fatal("unrelated manager changed")
			}
			for id := range files {
				_, ok := target.Packages[id]
				if !tracked && id != "cinema" && id != "cinema-LICENSE" && id != "cinema-NOTICE" && id != "cinema-THIRD-PARTY-LICENSES.txt" && ok {
					t.Fatal("untracked document installed", id)
				}
			}
			for _, bad := range []struct{ id, file string }{{"cinema-runtime", "docs/MyCinema/SETUP.md"}, {"cinema-doc-SETUP.md", "MyCinema.exe"}, {"cinema-doc-SETUP.md", "docs/Other/SETUP.md"}, {"cinema-provenance", "docs/MyCinema/../SOURCE-PROVENANCE.json"}, {"other-doc-SETUP.md", "docs/MyCinema/SETUP.md"}} {
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
				delete(m.Packages, "cinema-doc-SETUP.md")
				b, _ := json.Marshal(m)
				if _, e := manifestChoice(g, current, repo, release, b); e == nil {
					t.Fatal("tracked guide disappeared")
				}
			}
		})
	}
}
