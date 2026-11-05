package catalog

import "testing"

func TestNewsDefaultSelectionAndManagedOmission(t *testing.T) {
	c := Catalog{Packages: map[string]Package{"news": {File: "MyNews.exe"}, "news-contract": {File: "docs/MyNews/compatibility.json"}, "news-manifest": {File: "docs/MyNews/manifest.json"}}, Apps: []App{{ID: "accounts"}, {ID: "news", Package: "news", Depends: []string{"accounts"}}}}
	if !DefaultInstall("news") || CategoryFor("news").ID != "general" {
		t.Fatal("News must be a General app")
	}
	wanted := c.WantedApps(Catalog{Apps: []App{{ID: "accounts"}}})
	if len(wanted) != 2 {
		t.Fatal("existing installations fail to discover News", wanted)
	}
	for _, on := range []bool{false, true} {
		keep := InstallPackageIDs(c, map[string]bool{"accounts": true, "news": on})
		if keep["news"] != on || keep["news-contract"] != on || keep["news-manifest"] != on {
			t.Fatal("News selection split its owned contract")
		}
	}
	full := Catalog{}
	for _, a := range ManagedApps() {
		if a.ID != "news" {
			full.Apps = append(full.Apps, App{ID: a.ID, Name: a.Name, Route: a.Route, OIDC: a.OIDC})
		}
	}
	if full.ValidateManagedInventory() == nil {
		t.Fatal("release can omit News")
	}
}
