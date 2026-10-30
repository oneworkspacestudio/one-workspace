package catalog

import "testing"

func TestAuthenticatorInventoryAndSelection(t *testing.T) {
	c := Catalog{Packages: map[string]Package{"authenticator": {File: "one-authenticator.exe"}, "authenticator-contract": {File: "docs/OneAuthenticator/compatibility.json"}, "authenticator-LICENSE": {File: "licenses/OneAuthenticator-LICENSE"}}, Apps: []App{{ID: "authenticator", Package: "authenticator"}}}
	for _, selected := range []bool{false, true} {
		keep := InstallPackageIDs(c, map[string]bool{"authenticator": selected})
		if keep["authenticator"] != selected || keep["authenticator-contract"] != selected || !keep["authenticator-LICENSE"] {
			t.Fatal("namespace selection lost", selected, keep)
		}
	}
	if !DefaultInstall("authenticator") || CategoryFor("authenticator").ID != "administration" {
		t.Fatal("default administration discovery missing")
	}
	full := Catalog{}
	for _, a := range ManagedApps() {
		if a.ID != "authenticator" {
			full.Apps = append(full.Apps, App{ID: a.ID, Name: a.Name, Route: a.Route, OIDC: a.OIDC})
		}
	}
	if full.ValidateManagedInventory() == nil {
		t.Fatal("incomplete release accepted")
	}
}
