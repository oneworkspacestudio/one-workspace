package catalog

import "testing"

func TestPrivateMailAndSharedConsumersKeepNativeRuntime(t *testing.T) {
	for _, id := range []string{"mail", "calendar", "contacts"} {
		c := Catalog{Apps: []App{{ID: id, Package: id, Env: map[string]string{"WORKSPACE_MAIL_BACKEND": "private"}}}, Packages: map[string]Package{"msvc-msvcp140": {File: "runtime/msvc/msvcp140.dll"}}}
		keep := InstallPackageIDs(c, map[string]bool{id: true})
		if !keep["msvc-msvcp140"] {
			t.Fatal(id, "omitted the engine runtime")
		}
	}
}
