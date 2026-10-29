package distribution

import (
	"context"
	"encoding/json"
	"testing"
)

func TestSuiteUsesOnlyHashMatchingPackageAssets(t *testing.T) {
	for _, corrupt := range []bool{false, true} {
		t.Run(map[bool]string{false: "verified", true: "mismatch"}[corrupt], func(t *testing.T) {
			current, tr, _ := setupUpdateTransport(t, true)
			endpoint := "https://api.github.com/repos/" + Repository + "/releases?per_page=100"
			var releases []Release
			if err := json.Unmarshal(tr[endpoint], &releases); err != nil {
				t.Fatal(err)
			}
			r := &releases[0]
			assets := r.Assets[:0]
			for _, a := range r.Assets {
				if a.Name != "Workspace.exe" {
					assets = append(assets, a)
				}
			}
			r.Assets = assets
			raw := []byte("new-manager")
			if corrupt {
				raw = []byte("tampered-manager")
			}
			a := fixtureAsset(Repository, r.Tag, "workspace-package-manager.bin", "77", raw)
			r.Assets = append(r.Assets, a)
			tr[endpoint], _ = json.Marshal(releases)
			plan, err := CheckUpdates(context.Background(), current)
			if err != nil {
				t.Fatal(err)
			}
			_, sources, err := SelectUpdates(current, plan, []string{"workspace"})
			if err != nil {
				t.Fatal(err)
			}
			d := sources["manager"]
			if !corrupt && (d.Asset.Name != a.Name || d.Member != "") {
				t.Fatal("small package asset not used", d)
			}
			if corrupt && d.Asset.Name == a.Name {
				t.Fatal("mismatching package asset accepted")
			}
		})
	}
}
