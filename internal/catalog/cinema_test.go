package catalog

import "testing"

func TestCinemaRuntimeFollowsSelection(t *testing.T) {
	c := Catalog{Apps: []App{{ID: "cinema", Package: "cinema"}}, Packages: map[string]Package{"cinema": {File: "MyCinema.exe"}, "cinema-runtime": {File: "runtime/mycinema/jellyfin.zip"}, "cinema-contract": {File: "docs/MyCinema/compatibility.json"}}}
	for _, selected := range []bool{false, true} {
		keep := InstallPackageIDs(c, map[string]bool{"cinema": selected})
		for _, id := range []string{"cinema", "cinema-runtime", "cinema-contract"} {
			if keep[id] != selected {
				t.Fatal(id, selected)
			}
		}
	}
}
