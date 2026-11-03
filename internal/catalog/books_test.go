package catalog

import "testing"

func TestBooksDefaultSelectionAndManagedOmission(t *testing.T) {
	c := Catalog{Packages: map[string]Package{"books": {File: "MyBooks.exe"}, "books-contract": {File: "docs/MyBooks/compatibility.json"}, "books-manifest": {File: "docs/MyBooks/manifest.json"}}, Apps: []App{{ID: "accounts"}, {ID: "books", Package: "books", Depends: []string{"accounts"}}}}
	if !DefaultInstall("books") || CategoryFor("books").ID != "general" {
		t.Fatal("Books must be a General app")
	}
	wanted := c.WantedApps(Catalog{Apps: []App{{ID: "accounts"}}})
	if len(wanted) != 2 {
		t.Fatal("existing installations fail to discover Books", wanted)
	}
	for _, on := range []bool{false, true} {
		keep := InstallPackageIDs(c, map[string]bool{"accounts": true, "books": on})
		if keep["books"] != on || keep["books-contract"] != on || keep["books-manifest"] != on {
			t.Fatal("Books selection split its owned contract")
		}
	}
	full := Catalog{}
	for _, a := range ManagedApps() {
		if a.ID != "books" {
			full.Apps = append(full.Apps, App{ID: a.ID, Name: a.Name, Route: a.Route, OIDC: a.OIDC})
		}
	}
	if full.ValidateManagedInventory() == nil {
		t.Fatal("release can omit Books")
	}
}
