package catalog

import "testing"

func TestTicketsPayloadSelectionAndInventory(t *testing.T) {
	c := Catalog{Packages: map[string]Package{"tickets": {File: "MyTickets.exe"}, "tickets-contract": {File: "docs/MyTickets/compatibility.json"}, "tickets-LICENSE": {File: "licenses/MyTickets-LICENSE"}}, Apps: []App{{ID: "tickets", Package: "tickets"}}}
	for _, selected := range []bool{false, true} {
		keep := InstallPackageIDs(c, map[string]bool{"tickets": selected})
		if keep["tickets"] != selected || keep["tickets-contract"] != selected || !keep["tickets-LICENSE"] {
			t.Fatal(selected, keep)
		}
	}
	full := Catalog{}
	for _, a := range ManagedApps() {
		full.Apps = append(full.Apps, App{ID: a.ID, Name: a.Name, Route: a.Route, OIDC: a.OIDC})
	}
	if err := full.ValidateManagedInventory(); err != nil {
		t.Fatal(err)
	}
	var omitted []App
	for _, a := range full.Apps {
		if a.ID != "tickets" {
			omitted = append(omitted, a)
		}
	}
	full.Apps = omitted
	if full.ValidateManagedInventory() == nil {
		t.Fatal("release omitted Tickets without refusal")
	}
}
