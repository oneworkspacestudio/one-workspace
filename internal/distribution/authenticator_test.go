package distribution

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"testing"
)

func TestAuthenticatorUpgradeOwnsCompletePackageNamespace(t *testing.T) {
	old, incoming := fourDefaultCatalogs()
	incoming.Version = "0.3.0-preview.109"
	incoming.InstallerIntegration = 6
	incoming.Packages["authenticator"] = fixturePackage("one-authenticator.exe", []byte("synthetic-authenticator"))
	incoming.Packages["authenticator-contract"] = fixturePackage("docs/OneAuthenticator/compatibility.json", []byte("synthetic-contract"))
	incoming.Packages["authenticator-LICENSE"] = fixturePackage("licenses/OneAuthenticator-LICENSE", []byte("synthetic-license"))
	incoming.Apps = append(incoming.Apps, catalog.App{ID: "authenticator", Name: "One Authenticator", Version: "0.1.0-preview.2", Kind: "process", Package: "authenticator", Executable: "one-authenticator.exe", Port: 9027, Route: "/authenticator", StripPrefix: true, Health: "/healthz", Source: "https://github.com/oneworkspacestudio/one-authenticator/releases/tag/v0.1.0-preview.2"})
	for _, excluded := range []bool{false, true} {
		current := old
		if excluded {
			current.Selection = &catalog.AppSelection{Excluded: []string{"authenticator"}}
		}
		target, err := SelectBundleUpdates(current, incoming, []string{"workspace"})
		if err != nil {
			t.Fatal(err)
		}
		a, found := target.Find("authenticator")
		if found == excluded {
			t.Fatal("new default/exclusion lost", found, excluded)
		}
		if found && (a.Port != 25137 || a.OIDC || a.Route != "/authenticator") {
			t.Fatal("independent public runtime contract lost", a)
		}
		for _, id := range []string{"authenticator", "authenticator-contract"} {
			if _, ok := target.Packages[id]; ok == excluded {
				t.Fatal("owned package missing or unwanted", id, excluded)
			}
		}
	}
	for _, id := range []string{"authenticator", "authenticator-contract", "authenticator-LICENSE"} {
		if updatePackageOwner(id) != "authenticator" {
			t.Fatal("package namespace ownership missing", id)
		}
	}
}
