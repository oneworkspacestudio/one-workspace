package catalog

import "testing"

func TestAIServiceInventorySelectionAndDefaults(t *testing.T) {
	if !DefaultInstall("aiservice") || CategoryFor("aiservice").ID != "administration" {
		t.Fatal("service setup policy")
	}
	c := Catalog{Packages: map[string]Package{"aiservice": {File: "AIService.exe"}, "aiservice-contract": {File: "docs/AIService/compatibility.json"}}, Apps: []App{{ID: "accounts"}, {ID: "aiservice", Package: "aiservice", Depends: []string{"accounts"}}}}
	for _, on := range []bool{true, false} {
		p := InstallPackageIDs(c, map[string]bool{"accounts": true, "aiservice": on})
		if p["aiservice"] != on || p["aiservice-contract"] != on {
			t.Fatal("split service package")
		}
	}
	full := Catalog{}
	for _, a := range ManagedApps() {
		if a.ID != "aiservice" {
			full.Apps = append(full.Apps, App{ID: a.ID, Name: a.Name, Route: a.Route, OIDC: a.OIDC})
		}
	}
	if full.ValidateManagedInventory() == nil {
		t.Fatal("release omitted AI Service")
	}
}
