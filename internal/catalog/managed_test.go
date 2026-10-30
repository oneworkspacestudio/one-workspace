package catalog

import "testing"

func TestManagedReleaseCannotOmitRemote(t *testing.T) {
	c := Catalog{}
	for _, a := range ManagedApps() {
		if a.ID != "remote" {
			c.Apps = append(c.Apps, App{ID: a.ID, Name: a.Name, Route: a.Route, OIDC: a.OIDC})
		}
	}
	if c.ValidateManagedInventory() == nil {
		t.Fatal("missing My Remote accepted")
	}
	c.Apps = append(c.Apps, App{ID: "remote", Name: "My Remote", Route: "/remote", OIDC: true})
	if err := c.ValidateManagedInventory(); err != nil {
		t.Fatal(err)
	}
}
