package distribution

import (
	"context"
	"encoding/json"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"testing"
)

func TestPublicSuiteAdvancesOmittedPrivateHistoryWithoutDowngrade(t *testing.T) {
	for _, tc := range []struct {
		name, version string
		available     bool
	}{
		{"newer", "0.1.0-preview.5", true},
		{"older", "0.1.0-preview.3", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			current, tr, _ := setupUpdateTransport(t, true)
			endpoint := "https://api.github.com/repos/" + Repository + "/releases?per_page=100"
			var releases []Release
			if err := json.Unmarshal(tr[endpoint], &releases); err != nil {
				t.Fatal(err)
			}
			asset, _ := releaseAsset(releases[0], "catalog.json")
			var latest catalog.Catalog
			if err := json.Unmarshal(tr[asset.APIURL], &latest); err != nil {
				t.Fatal(err)
			}
			latest.Apps[0].Version = tc.version
			latest.Apps[0].Source = "https://github.com/oneworkspacestudio/notes/releases/tag/v" + tc.version
			latest.Packages["notes"] = fixturePackage("MyNotes.exe", []byte("suite-private-program"))
			raw, _ := json.Marshal(latest)
			releases[0].Assets[0] = fixtureAsset(Repository, releases[0].Tag, "catalog.json", "1", raw)
			tr[asset.APIURL] = raw
			tr[endpoint], _ = json.Marshal(releases)
			tr["https://api.github.com/repos/oneworkspacestudio/notes/releases?per_page=100"] = []byte("[]")
			plan, err := CheckUpdates(context.Background(), current)
			if err != nil {
				t.Fatal(err)
			}
			for _, choice := range plan.Choices {
				if choice.ID == "notes" {
					if choice.Available != tc.available {
						t.Fatalf("available=%v, want %v: %+v", choice.Available, tc.available, choice)
					}
					if tc.available {
						target, downloads, err := SelectUpdates(current, plan, []string{"notes"})
						if err != nil {
							t.Fatal(err)
						}
						if target.Version != current.Version || target.Apps[0].Version != tc.version || len(downloads) != 1 {
							t.Fatalf("wrong scoped update: %+v downloads=%d", target, len(downloads))
						}
					}
					return
				}
			}
			t.Fatal("app update choice missing")
		})
	}
}
