package distribution

import (
	"encoding/json"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"testing"
)

func TestHostingUpdateAcrossDefaultAndFullDocumentLayouts(t *testing.T) {
	for _, tracked := range []bool{false, true} {
		t.Run(fmt.Sprint(tracked), func(t *testing.T) {
			repo, tag := "oneworkspacestudio/my-hosting", "v0.1.0-preview.5"
			current := catalog.Catalog{Schema: 1, Version: "0.3.0-preview.141", Packages: map[string]catalog.Package{"manager": fixturePackage("Workspace.exe", []byte("manager"))}, Apps: []catalog.App{{ID: "hosting", Name: "My Hosting", Kind: "process", Package: "hosting", Executable: "MyHosting.exe", Port: 9007, Route: "/hosting", Version: "0.1.0-preview.2", Source: "https://github.com/" + repo + "/releases/tag/v0.1.0-preview.2"}}}
			files := map[string]string{"hosting": "MyHosting.exe", "hosting-LICENSE": "licenses/MyHosting-LICENSE", "hosting-NOTICE": "licenses/MyHosting-NOTICE", "hosting-THIRD-PARTY-LICENSES.txt": "licenses/MyHosting-THIRD-PARTY-LICENSES.txt", "hosting-compatibility.json": "docs/MyHosting/compatibility.json", "hosting-dependencies.json": "docs/MyHosting/dependencies.json", "hosting-runtime-pins.json": "docs/MyHosting/runtime-pins.json"}
			for _, name := range []string{"ACCEPTANCE.md", "FRAMEWORKS.md", "OPERATIONS.md", "PUBLISHING.md", "openapi.json"} {
				files["hosting-doc-"+name] = "docs/MyHosting/" + name
			}
			app := current.Apps[0]
			app.Version = "0.1.0-preview.5"
			app.Source = "https://github.com/" + repo + "/releases/tag/" + tag
			m := AppUpdateManifest{Repository: repo, Tag: tag, Apps: []catalog.App{app}, Packages: map[string]UpdatePackage{}}
			release := Release{Tag: tag}
			for id, file := range files {
				runtime := id == "hosting" || id == "hosting-LICENSE" || id == "hosting-NOTICE" || id == "hosting-THIRD-PARTY-LICENSES.txt"
				if runtime || tracked {
					current.Packages[id] = fixturePackage(file, []byte("old "+id))
				}
				body := []byte("new " + id)
				asset := fixtureAsset(repo, tag, id+".bin", fmt.Sprint(len(release.Assets)+1), body)
				release.Assets = append(release.Assets, asset)
				m.Packages[id] = UpdatePackage{Package: fixturePackage(file, body), Asset: asset.Name}
			}
			if e := RequireManager(&m, "0.3.0-preview.141", nil); e != nil {
				t.Fatal(e)
			}
			g := updateGroup{"hosting", "My Hosting", "my-hosting", []string{"hosting"}}
			raw, _ := json.Marshal(m)
			choice, e := manifestChoice(g, current, repo, release, raw)
			if e != nil {
				t.Fatal(e)
			}
			choice.Available = true
			plan := UpdatePlan{BaseCatalogSHA256: catalogIdentity(current), Choices: []UpdateChoice{choice}}
			target, sources, e := SelectUpdates(current, plan, []string{"hosting"})
			if e != nil {
				t.Fatal(e)
			}
			want := 4
			if tracked {
				want = 12
			}
			if len(choice.Packages) != want || len(sources) != want {
				t.Fatal("wrong scope", len(choice.Packages), len(sources))
			}
			if target.Apps[0].Version != "0.1.0-preview.5" || target.Packages["manager"] != current.Packages["manager"] || target.Version != current.Version {
				t.Fatal("unrelated manager changed")
			}
			for id := range files {
				_, ok := target.Packages[id]
				if !tracked && id != "hosting" && id != "hosting-LICENSE" && id != "hosting-NOTICE" && id != "hosting-THIRD-PARTY-LICENSES.txt" && ok {
					t.Fatal("untracked document installed", id)
				}
			}
			for _, bad := range []struct{ id, file string }{{"hosting-runtime", "docs/MyHosting/ACCEPTANCE.md"}, {"hosting-doc-ACCEPTANCE.md", "MyHosting.exe"}, {"hosting-doc-ACCEPTANCE.md", "docs/Other/ACCEPTANCE.md"}, {"hosting-runtime-pins.json", "docs/MyHosting/../runtime-pins.json"}, {"other-doc-ACCEPTANCE.md", "docs/MyHosting/ACCEPTANCE.md"}} {
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
				delete(m.Packages, "hosting-doc-ACCEPTANCE.md")
				b, _ := json.Marshal(m)
				if _, e := manifestChoice(g, current, repo, release, b); e == nil {
					t.Fatal("tracked guide disappeared")
				}
			}
		})
	}
}
