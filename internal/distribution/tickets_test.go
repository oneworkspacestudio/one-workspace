package distribution

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"testing"
)

func TestTicketsAddedByManagerUpgradeAndExplicitlyOmittable(t *testing.T) {
	old, incoming := fourDefaultCatalogs()
	incoming.Version = "0.3.0-preview.105"
	incoming.InstallerIntegration = 5
	incoming.Packages["tickets"] = fixturePackage("MyTickets.exe", []byte("isolated-tickets-fixture"))
	incoming.Packages["tickets-contract"] = fixturePackage("docs/MyTickets/compatibility.json", []byte("contract"))
	incoming.Apps = append(incoming.Apps, catalog.App{ID: "tickets", Name: "My Tickets", Version: "0.1.0-preview.1", Kind: "process", Package: "tickets", Executable: "MyTickets.exe", Port: 9026, Route: "/tickets", StripPrefix: true, OIDC: true, Callback: "/tickets/auth/callback", Depends: []string{"accounts"}, Source: "https://github.com/oneworkspacestudio/my-tickets/releases/tag/v0.1.0-preview.1"})
	for _, excluded := range []bool{false, true} {
		previous := old
		if excluded {
			previous.Selection = &catalog.AppSelection{Excluded: []string{"tickets"}}
		}
		target, err := SelectBundleUpdates(previous, incoming, []string{"workspace"})
		if err != nil {
			t.Fatal(err)
		}
		app, found := target.Find("tickets")
		if found == excluded {
			t.Fatal("new default/exclusion policy", found, excluded)
		}
		if found && (!app.OIDC || app.Route != "/tickets" || app.Callback != "/tickets/auth/callback") {
			t.Fatal("lost native identity contract")
		}
		if _, found = target.Packages["tickets-contract"]; found == excluded {
			t.Fatal("missing or unwanted owned contract")
		}
	}
}

func TestTicketsRejectsOlderManager(t *testing.T) {
	m := AppUpdateManifest{Schema: 2, Integration: 3, MinimumManager: "0.3.0-preview.105", RequiredCapabilities: []string{"native-tickets-v1"}}
	if checkManifestRequirements(m, "0.3.0-preview.104") == nil {
		t.Fatal("old manager accepted Tickets")
	}
	if err := checkManifestRequirements(m, "0.3.0-preview.105"); err != nil {
		t.Fatal(err)
	}
}
