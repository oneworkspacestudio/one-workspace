package catalog

import "testing"

func TestCommuteSelectionOwnsWholePackage(t *testing.T) {
	if DefaultInstall("commute") || CategoryFor("commute").ID != "general" {
		t.Fatal("London app must be an optional General app")
	}
	c := Catalog{Packages: map[string]Package{"commute": {File: "MyCommute.exe"}, "commute-contract": {File: "docs/MyCommute/compatibility.json"}, "commute-manifest": {File: "docs/MyCommute/manifest.json"}, "unrelated": {File: "Other.exe"}}, Apps: []App{{ID: "commute", Package: "commute"}, {ID: "other", Package: "unrelated"}}}
	for _, on := range []bool{true, false} {
		ids := InstallPackageIDs(c, map[string]bool{"commute": on, "other": true})
		if ids["commute"] != on || ids["commute-contract"] != on || ids["commute-manifest"] != on || !ids["unrelated"] {
			t.Fatal("package ownership drift", ids)
		}
	}
}
func TestCommuteRequiredInCompleteRegistry(t *testing.T) {
	var c Catalog
	for _, a := range ManagedApps() {
		if a.ID != "commute" {
			c.Apps = append(c.Apps, App{ID: a.ID, Name: a.Name, Route: a.Route, OIDC: a.OIDC})
		}
	}
	if c.ValidateManagedInventory() == nil {
		t.Fatal("release must reject missing Commute")
	}
	if InstallerIntegration < 10 {
		t.Fatal("old controllers cannot install Commute")
	}
}
