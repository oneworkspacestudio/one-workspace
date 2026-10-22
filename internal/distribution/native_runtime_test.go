package distribution

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"testing"
)

func TestExistingCalendarReceivesNewRuntimePrerequisites(t *testing.T) {
	app := catalog.App{ID: "calendar", Package: "calendar"}
	current := catalog.Catalog{Apps: []catalog.App{app}, Packages: map[string]catalog.Package{"calendar": {File: "Calendar.exe"}}}
	choice := UpdateChoice{ID: "workspace", Packages: map[string]catalog.Package{"msvc-msvcp140": {File: "runtime/msvc/msvcp140.dll"}}, Downloads: map[string]UpdateDownload{}}
	got := scopeUpdateChoice(current, choice)
	if _, ok := got.Packages["msvc-msvcp140"]; !ok {
		t.Fatal("existing Calendar update dropped its new prerequisite")
	}
	if updatePackageOwner("msvc-msvcp140") != "workspace" {
		t.Fatal("runtime is outside the protected manager repair group")
	}
}
