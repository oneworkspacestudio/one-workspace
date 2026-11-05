package distribution

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"testing"
)

func TestSurvivalUpgradeOwnsCompleteNamespace(t *testing.T) {
	old, incoming := fourDefaultCatalogs()
	incoming.Version = "0.3.0-preview.117"
	incoming.InstallerIntegration = 11
	owned := map[string]string{"survival": "MySurvival.exe", "survival-contract": "docs/MySurvival/compatibility.json", "survival-manifest": "docs/MySurvival/manifest.json", "survival-LICENSE": "licenses/MySurvival-LICENSE", "survival-doc-fixture": "docs/MySurvival/README.md"}
	for id, file := range owned {
		incoming.Packages[id] = fixturePackage(file, []byte("synthetic "+id))
	}
	incoming.Apps = append(incoming.Apps, catalog.App{ID: "survival", Name: "My Survival", Version: "0.1.0-preview.1", Kind: "process", Package: "survival", Executable: "MySurvival.exe", Port: 9041, Route: "/survival", StripPrefix: true, OIDC: true, Callback: "/survival/auth/callback", Depends: []string{"accounts"}, Source: "https://github.com/oneworkspacestudio/my-survival/releases/tag/v0.1.0-preview.1"})
	for _, excluded := range []bool{false, true} {
		previous := old
		if excluded {
			previous.Selection = &catalog.AppSelection{Excluded: []string{"survival"}}
		}
		target, err := SelectBundleUpdates(previous, incoming, []string{"workspace"})
		if err != nil {
			t.Fatal(err)
		}
		app, found := target.Find("survival")
		if found == excluded {
			t.Fatal("default selection/exclusion lost")
		}
		if found && (app.Port != 25151 || !app.OIDC || app.Callback != "/survival/auth/callback") {
			t.Fatal("lost launch/identity contract", app)
		}
		for id := range owned {
			if _, ok := target.Packages[id]; ok == excluded {
				t.Fatal("lost or unwanted owned payload", id)
			}
			if updatePackageOwner(id) != "survival" {
				t.Fatal("missing update ownership", id)
			}
		}
		if target.Packages["accounts"] != old.Packages["accounts"] {
			t.Fatal("unselected account package changed")
		}
	}
}
