package catalog

import "testing"

func TestMyMusicReleaseAndOptionalPayload(t *testing.T) {
	if DefaultInstall("mymusic") || CategoryFor("mymusic").ID != "development" {
		t.Fatal("Vuno must be optional under Creative & Development")
	}
	c := Catalog{InstallerIntegration: MyMusicIntegration}
	for _, app := range ManagedApps() {
		if !app.IndependentDistribution && app.SinceIntegration < MyMusicIntegration {
			c.Apps = append(c.Apps, App{ID: app.ID, Name: app.Name, Route: app.Route, OIDC: app.OIDC, Package: app.ID})
		}
	}
	if c.ValidateManagedInventory() == nil {
		t.Fatal("accepted a current release without Vuno")
	}
	c.InstallerIntegration = IndependentAppIntegration
	if err := c.ValidateManagedInventory(); err != nil {
		t.Fatal("old release failed inventory validation:", err)
	}
	c.InstallerIntegration = MyMusicIntegration
	c.Apps = append(c.Apps, App{ID: "mymusic", Name: "Vuno", Route: "/mymusic", OIDC: true, Package: "mymusic"})
	c.Packages = map[string]Package{
		"accounts":         {File: "MyAccounts.exe"},
		"mymusic":          {File: "MyMusic.exe"},
		"mymusic-contract": {File: "docs/MyMusic/compatibility.json"},
	}
	if err := c.ValidateManagedInventory(); err != nil {
		t.Fatal(err)
	}
	c.Version = "0.3.0-preview.142"
	c.Apps[len(c.Apps)-1].Name = "My Music"
	if err := c.ValidateManagedInventory(); err != nil {
		t.Fatal("published preview.142 name must remain valid:", err)
	}
	c.Version = "0.3.0-preview.143"
	if c.ValidateManagedInventory() == nil {
		t.Fatal("a new release must use the Vuno display name")
	}
	for _, chosen := range []bool{false, true} {
		selected := InstallPackageIDs(c, map[string]bool{"accounts": true, "mymusic": chosen})
		if !selected["accounts"] || selected["mymusic"] != chosen || selected["mymusic-contract"] != chosen {
			t.Fatal("incorrect payload selection:", selected)
		}
	}
}
