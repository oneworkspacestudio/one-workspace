package distribution

import (
	"encoding/json"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"testing"
)

func TestMyAIUpdateAcrossDefaultAndFullDocumentLayouts(t *testing.T) {
	for _, tracked := range []bool{false, true} {
		t.Run(fmt.Sprint(tracked), func(t *testing.T) {
			repo, tag := "oneworkspacestudio/my-ai", "v0.1.0-preview.4"
			current := catalog.Catalog{Schema: 1, Version: "0.3.0-preview.87", Packages: map[string]catalog.Package{"manager": fixturePackage("Workspace.exe", []byte("manager"))}, Apps: []catalog.App{{ID: "myai", Name: "My AI", Kind: "process", Package: "myai", Executable: "MyAI.exe", Port: 9019, Route: "/myai", Version: "0.1.0-preview.2", Source: "https://github.com/" + repo + "/releases/tag/v0.1.0-preview.2"}}}
			files := map[string]string{"myai": "MyAI.exe", "myai-LICENSE": "licenses/MyAI-LICENSE", "myai-NOTICE": "licenses/MyAI-NOTICE", "myai-THIRD-PARTY-LICENSES.txt": "licenses/MyAI-THIRD-PARTY-LICENSES.txt", "myai-contract": "docs/MyAI/compatibility.json"}
			for _, name := range []string{"SETUP.md", "SUPPORT.md", "ADR-001-architecture.md", "SECURITY-CHECKS.md"} {
				files["myai-doc-"+name] = "docs/MyAI/" + name
			}
			app := current.Apps[0]
			app.Version = "0.1.0-preview.4"
			app.Source = "https://github.com/" + repo + "/releases/tag/" + tag
			m := AppUpdateManifest{Repository: repo, Tag: tag, Apps: []catalog.App{app}, Packages: map[string]UpdatePackage{}}
			release := Release{Tag: tag}
			for id, file := range files {
				runtime := id == "myai" || id == "myai-LICENSE" || id == "myai-NOTICE" || id == "myai-THIRD-PARTY-LICENSES.txt"
				if runtime || tracked {
					current.Packages[id] = fixturePackage(file, []byte("old "+id))
				}
				body := []byte("new " + id)
				asset := fixtureAsset(repo, tag, id+".bin", fmt.Sprint(len(release.Assets)+1), body)
				release.Assets = append(release.Assets, asset)
				m.Packages[id] = UpdatePackage{Package: fixturePackage(file, body), Asset: asset.Name}
			}
			if e := RequireManager(&m, "0.3.0-preview.87", nil); e != nil {
				t.Fatal(e)
			}
			g := updateGroup{"myai", "My AI", "my-ai", []string{"myai"}}
			raw, _ := json.Marshal(m)
			choice, e := manifestChoice(g, current, repo, release, raw)
			if e != nil {
				t.Fatal(e)
			}
			choice.Available = true
			plan := UpdatePlan{BaseCatalogSHA256: catalogIdentity(current), Choices: []UpdateChoice{choice}}
			target, sources, e := SelectUpdates(current, plan, []string{"myai"})
			if e != nil {
				t.Fatal(e)
			}
			want := 4
			if tracked {
				want = 9
			}
			if len(choice.Packages) != want || len(sources) != want {
				t.Fatal("wrong scope", len(choice.Packages), len(sources))
			}
			if target.Apps[0].Version != "0.1.0-preview.4" || target.Packages["manager"] != current.Packages["manager"] || target.Version != current.Version {
				t.Fatal("unrelated manager changed")
			}
			for id := range files {
				_, ok := target.Packages[id]
				if !tracked && id != "myai" && id != "myai-LICENSE" && id != "myai-NOTICE" && id != "myai-THIRD-PARTY-LICENSES.txt" && ok {
					t.Fatal("untracked document installed", id)
				}
			}
			for _, bad := range []struct{ id, file string }{{"myai-runtime", "docs/MyAI/SETUP.md"}, {"myai-doc-SETUP.md", "MyAI.exe"}, {"myai-doc-SETUP.md", "docs/Other/SETUP.md"}, {"myai-doc-SUPPORT.md", "docs/MyAI/../SUPPORT.md"}, {"other-doc-SETUP.md", "docs/MyAI/SETUP.md"}} {
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
				delete(m.Packages, "myai-doc-SETUP.md")
				b, _ := json.Marshal(m)
				if _, e := manifestChoice(g, current, repo, release, b); e == nil {
					t.Fatal("tracked guide disappeared")
				}
			}
		})
	}
}
