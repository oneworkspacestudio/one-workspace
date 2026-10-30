package distribution

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"testing"
)

func TestManagedAppsHaveUpdateGroups(t *testing.T) {
	for _, app := range catalog.ManagedApps() {
		found := false
		for _, g := range updateGroups {
			for _, id := range g.apps {
				if id == app.ID {
					found = true
					if app.UpdateGroup != g.id {
						t.Errorf("%s update registry disagrees: %s vs %s", id, app.UpdateGroup, g.id)
					}
				}
			}
		}
		if !found {
			t.Errorf("%s missing update group", app.ID)
		}
	}
}

func TestWorkspaceUpdateRetainsNewInstallerGuard(t *testing.T) {
	current := updateCatalogFixture()
	incoming := current
	incoming.Version = "0.3.0-preview.104"
	incoming.InstallerIntegration = catalog.InstallerIntegration
	incoming.Packages = map[string]catalog.Package{}
	for id, p := range current.Packages {
		incoming.Packages[id] = p
	}
	incoming.Packages["manager"] = fixturePackage("Workspace.exe", []byte("manager104"))
	target, err := SelectBundleUpdates(current, incoming, []string{"workspace"})
	if err != nil {
		t.Fatal(err)
	}
	if target.InstallerIntegration != incoming.InstallerIntegration {
		t.Fatal("updated installation lost its required installer guard")
	}
}
