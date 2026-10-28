package catalog

import "testing"

func TestFinanceRuntimeFollowsInstallSelection(t *testing.T) {
	c := Catalog{Apps: []App{{ID: "finance", Package: "finance"}}, Packages: map[string]Package{"finance": {File: "MyFinance.exe"}, "finance-runtime": {File: "finance-runtime.zip"}, "finance-LICENSE": {File: "licenses/MyFinance-LICENSE"}}}
	for _, selected := range []bool{false, true} {
		keep := InstallPackageIDs(c, map[string]bool{"finance": selected})
		if keep["finance-runtime"] != selected || keep["finance"] != selected || !keep["finance-LICENSE"] {
			t.Fatal("Finance runtime must follow the app; notices remain available", selected, keep)
		}
	}
}
