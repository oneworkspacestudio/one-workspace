package catalog

import "testing"

func TestGameStudioManagedSelection(t *testing.T) {
	if DefaultInstall("gamestudio") || CategoryFor("gamestudio").ID != "development" {
		t.Fatal("Game Studio must be optional in Creative & Development")
	}
	c := Catalog{
		Apps: []App{{ID: "accounts", Package: "accounts"}, {ID: "gamestudio", Package: "gamestudio", Depends: []string{"accounts"}}},
		Packages: map[string]Package{
			"accounts":            {File: "MyAccounts.exe"},
			"gamestudio":          {File: "GameStudio.exe"},
			"gamestudio-contract": {File: "docs/GameStudio/compatibility.json"},
		},
	}
	for _, selected := range []bool{true, false} {
		packages := InstallPackageIDs(c, map[string]bool{"accounts": true, "gamestudio": selected})
		if !packages["accounts"] || packages["gamestudio"] != selected || packages["gamestudio-contract"] != selected {
			t.Fatalf("Game Studio selection %t produced %v", selected, packages)
		}
	}
}
