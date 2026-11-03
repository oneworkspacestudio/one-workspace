package distribution

import (
	"context"
	"encoding/json"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"testing"
)

func TestNewOptionalVoiceRequiresExplicitSelection(t *testing.T) {
	old := updateCatalogFixture()
	old.Version = "0.3.0-preview.110"
	old.Packages["accounts"] = fixturePackage("MyAccounts.exe", []byte("accounts"))
	old.Apps = append(old.Apps, catalog.App{ID: "accounts", Name: "MyAccounts", Version: "fixture", Kind: "process", Package: "accounts", Executable: "MyAccounts.exe", Port: 25101, Source: "https://github.com/oneworkspacestudio/myaccounts/releases/tag/fixture"})
	incoming := old
	incoming.Version = "0.3.0-preview.111"
	incoming.InstallerIntegration = 7
	incoming.Apps = append([]catalog.App(nil), old.Apps...)
	incoming.Packages = map[string]catalog.Package{}
	for k, p := range old.Packages {
		incoming.Packages[k] = p
	}
	incoming.Packages["manager"] = fixturePackage("Workspace.exe", []byte("manager111"))
	incoming.Packages["voice"] = fixturePackage("MyVoice.exe", []byte("native voice"))
	incoming.Packages["voice-contract"] = fixturePackage("docs/MyVoice/compatibility.json", []byte("contract"))
	incoming.Packages["voice-runtime-owned"] = fixturePackage("VoiceSupport.bin", []byte("owned support fixture"))
	voice := catalog.App{ID: "voice", Name: "My Voice", Kind: "process", Package: "voice", Executable: "MyVoice.exe", Port: 9028, Route: "/voice", StripPrefix: true, OIDC: true, Callback: "/voice/auth/callback", Version: "0.1.0-preview.1", Source: "https://github.com/oneworkspacestudio/my-voice/releases/tag/v0.1.0-preview.1"}
	incoming.Apps = append(incoming.Apps, voice)
	offer := markNewApps(old, incoming, UpdateChoice{ID: "voice", Apps: []catalog.App{voice}})
	if !offer.Optional || len(offer.NewApps) != 1 || offer.NewApps[0] != "voice" {
		t.Fatal("missing explicit optional offer", offer)
	}
	for _, selected := range [][]string{nil, {"workspace"}, {"workspace", "voice"}} {
		target, e := SelectBundleUpdates(old, incoming, selected)
		if e != nil {
			t.Fatal(e)
		}
		_, has := target.Find("voice")
		want := len(selected) == 2
		if has != want {
			t.Fatalf("selection %v installed optional=%t", selected, has)
		}
		for _, id := range []string{"voice", "voice-contract", "voice-runtime-owned"} {
			_, has := target.Packages[id]
			if has != want {
				t.Fatal("incomplete optional namespace", selected, id)
			}
		}
		if has {
			a, _ := target.Find("voice")
			if a.Port != 25138 {
				t.Fatal("new optional port did not rebase", a.Port)
			}
		}
	}
	excluded := old
	excluded.Selection = &catalog.AppSelection{Excluded: []string{"voice"}}
	offer = markNewApps(excluded, incoming, UpdateChoice{ID: "voice", Apps: []catalog.App{voice}})
	if len(offer.NewApps) != 0 {
		t.Fatal("setup exclusion was offered again")
	}
	if _, e := SelectBundleUpdates(excluded, incoming, []string{"voice"}); e == nil {
		t.Fatal("setup opt-out silently bypassed")
	}
}

// Exercise discovery before selection: a repaired current manager has no Voice
// app yet, and the suite version can already equal the installed version.
func TestCheckUpdatesDiscoversOptionalVoiceOnCurrentSuite(t *testing.T) {
	current, tr, _ := setupUpdateTransport(t, true)
	var releases []Release
	if e := json.Unmarshal(tr["https://api.github.com/repos/"+Repository+"/releases?per_page=100"], &releases); e != nil {
		t.Fatal(e)
	}
	var latest catalog.Catalog
	if e := json.Unmarshal(tr[releases[0].Assets[0].APIURL], &latest); e != nil {
		t.Fatal(e)
	}
	current.Version = latest.Version
	current.InstallerIntegration = catalog.InstallerIntegration
	latest.InstallerIntegration = catalog.InstallerIntegration
	accounts := catalog.App{ID: "accounts", Name: "MyAccounts", Kind: "process", Package: "accounts", Executable: "MyAccounts.exe", Version: "fixture", Port: 25101, Route: "/accounts", Source: "https://github.com/oneworkspacestudio/myaccounts/releases/tag/fixture"}
	current.Apps = append(current.Apps, accounts)
	latest.Apps = append(latest.Apps, accounts)
	current.Packages["accounts"] = fixturePackage("MyAccounts.exe", []byte("accounts"))
	latest.Packages["accounts"] = current.Packages["accounts"]
	latest.Apps = append(latest.Apps, catalog.App{ID: "voice", Name: "My Voice", Kind: "process", Package: "voice", Executable: "MyVoice.exe", Version: "0.1.0-preview.2", Port: 9028, Route: "/voice", Depends: []string{"accounts"}, Source: "https://github.com/oneworkspacestudio/my-voice/releases/tag/v0.1.0-preview.2"})
	latest.Packages["voice"] = fixturePackage("MyVoice.exe", []byte("voice"))
	raw, _ := json.Marshal(latest)
	asset := releases[0].Assets[0]
	releases[0].Assets[0] = fixtureAsset(Repository, releases[0].Tag, "catalog.json", "1", raw)
	tr[asset.APIURL] = raw
	tr["https://api.github.com/repos/"+Repository+"/releases?per_page=100"], _ = json.Marshal(releases)
	for _, excluded := range []bool{false, true} {
		if excluded {
			current.Selection = &catalog.AppSelection{Excluded: []string{"voice"}}
		}
		plan, e := CheckUpdates(context.Background(), current)
		if e != nil {
			t.Fatal(e)
		}
		var found *UpdateChoice
		for i := range plan.Choices {
			if plan.Choices[i].ID == "voice" {
				found = &plan.Choices[i]
			}
		}
		if excluded {
			if found != nil {
				t.Fatal("explicit setup exclusion was offered")
			}
			continue
		}
		if found == nil || !found.Optional || !found.Available || len(found.NewApps) != 1 {
			t.Fatalf("missing optional Voice discovery: %+v", found)
		}
		target, _, e := SelectUpdates(current, plan, []string{"voice"})
		if e != nil {
			t.Fatal(e)
		}
		if _, ok := target.Find("voice"); !ok {
			t.Fatal("explicit discovered Voice selection missing")
		}

	}
}
