package distribution

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"testing"
)

func TestAIServiceUpgradeOwnsCompleteNamespace(t *testing.T) {
	old, incoming := fourDefaultCatalogs()
	incoming.Version = "0.3.0-preview.118"
	incoming.InstallerIntegration = 12
	owned := map[string]string{"aiservice": "AIService.exe", "aiservice-contract": "docs/AIService/compatibility.json", "aiservice-manifest": "docs/AIService/manifest.json", "aiservice-LICENSE": "licenses/AIService-LICENSE", "aiservice-doc-fixture": "docs/AIService/README.md"}
	for id, file := range owned {
		incoming.Packages[id] = fixturePackage(file, []byte("synthetic "+id))
	}
	incoming.Apps = append(incoming.Apps, catalog.App{ID: "aiservice", Name: "AI Service", Version: "0.1.0-preview.25", Kind: "process", Package: "aiservice", Executable: "AIService.exe", Port: 9042, Route: "/aiservice", StripPrefix: true, OIDC: true, Callback: "/aiservice/auth/callback", Depends: []string{"accounts"}, Source: "https://github.com/oneworkspacestudio/my-ai/releases/tag/v0.1.0-preview.25"})
	for _, excluded := range []bool{false, true} {
		previous := old
		if excluded {
			previous.Selection = &catalog.AppSelection{Excluded: []string{"aiservice"}}
		}
		target, err := SelectBundleUpdates(previous, incoming, []string{"workspace"})
		if err != nil {
			t.Fatal(err)
		}
		app, found := target.Find("aiservice")
		if found == excluded {
			t.Fatal("default selection/exclusion lost")
		}
		if found && (app.Port != 25152 || !app.OIDC || app.Callback != "/aiservice/auth/callback") {
			t.Fatal("lost launch/identity contract", app)
		}
		for id := range owned {
			if _, ok := target.Packages[id]; ok == excluded {
				t.Fatal("lost or unwanted owned payload", id)
			}
			if updatePackageOwner(id) != "myai" {
				t.Fatal("missing update ownership", id)
			}
		}
		if target.Packages["accounts"] != old.Packages["accounts"] {
			t.Fatal("unselected account package changed")
		}
	}
}
