package catalog

import "testing"

func TestInstallRemoteCarriesOwnedRuntimeAndRespectsExclusion(t *testing.T) {
	c := Catalog{Apps: []App{{ID: "remote", Package: "remote"}, {ID: "other", Package: "other"}}, Packages: map[string]Package{
		"remote": {File: "MyRemote.exe"}, "remote-runtime-node": {File: "remote-runtime/node.exe"},
		"remote-runtime-host": {File: "remote-runtime/host/tvnserver.exe"}, "remote-contract": {File: "docs/MyRemote/compatibility.json"},
		"remote-LICENSE": {File: "licenses/MyRemote-LICENSE"}, "other-runtime-node": {File: "other-runtime/node.exe"},
	}}
	for _, selected := range []bool{false, true} {
		keep := InstallPackageIDs(c, map[string]bool{"remote": selected})
		for _, id := range []string{"remote", "remote-runtime-node", "remote-runtime-host", "remote-contract"} {
			if keep[id] != selected {
				t.Errorf("selected=%v package=%s retained=%v", selected, id, keep[id])
			}
		}
		if !keep["remote-LICENSE"] || keep["other-runtime-node"] {
			t.Fatal("licence/unselected runtime policy changed")
		}
	}
}
