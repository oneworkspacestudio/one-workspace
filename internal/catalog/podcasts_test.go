package catalog

import "testing"

func TestPodcastsSelectionOwnsWholePackage(t *testing.T) {
	if DefaultInstall("podcasts") || CategoryFor("podcasts").ID != "development" {
		t.Fatal("Podcast studio must be optional Creative & Development")
	}
	c := Catalog{Packages: map[string]Package{"podcasts": {File: "MyPodcasts.exe"}, "podcasts-contract": {File: "docs/MyPodcasts/compatibility.json"}, "podcasts-manifest": {File: "docs/MyPodcasts/manifest.json"}, "unrelated": {File: "Other.exe"}}, Apps: []App{{ID: "podcasts", Package: "podcasts"}, {ID: "other", Package: "unrelated"}}}
	for _, on := range []bool{true, false} {
		ids := InstallPackageIDs(c, map[string]bool{"podcasts": on, "other": true})
		if ids["podcasts"] != on || ids["podcasts-contract"] != on || ids["podcasts-manifest"] != on || !ids["unrelated"] {
			t.Fatal("package ownership drift", ids)
		}
	}
}
func TestPodcastsRequiredInCompleteRegistry(t *testing.T) {
	var c Catalog
	for _, a := range ManagedApps() {
		if a.ID != "podcasts" {
			c.Apps = append(c.Apps, App{ID: a.ID, Name: a.Name, Route: a.Route, OIDC: a.OIDC})
		}
	}
	if c.ValidateManagedInventory() == nil {
		t.Fatal("release must reject missing Podcasts")
	}
	if InstallerIntegration < 14 {
		t.Fatal("old controllers cannot install Podcasts")
	}
}
