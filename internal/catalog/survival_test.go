package catalog

import "testing"

func TestSurvivalSelectionAndReleaseInventory(t *testing.T) {
	if !DefaultInstall("survival") || CategoryFor("survival").ID != "general" {
		t.Fatal("Survival should be a General default app")
	}
	c := Catalog{Packages: map[string]Package{"survival": {File: "MySurvival.exe"}, "survival-contract": {File: "docs/MySurvival/compatibility.json"}}, Apps: []App{{ID: "accounts"}, {ID: "survival", Package: "survival", Depends: []string{"accounts"}}}}
	for _, on := range []bool{true, false} {
		keep := InstallPackageIDs(c, map[string]bool{"accounts": true, "survival": on})
		if keep["survival"] != on || keep["survival-contract"] != on {
			t.Fatal("split app payload")
		}
	}
	full := Catalog{}
	for _, m := range ManagedApps() {
		full.Apps = append(full.Apps, App{ID: m.ID, Name: m.Name, Route: m.Route, OIDC: m.OIDC})
	}
	if full.ValidateManagedInventory() != nil {
		t.Fatal("complete inventory rejected")
	}
	for i, a := range full.Apps {
		if a.ID == "survival" {
			full.Apps = append(full.Apps[:i], full.Apps[i+1:]...)
			break
		}
	}
	if full.ValidateManagedInventory() == nil {
		t.Fatal("missing Survival accepted")
	}
}
