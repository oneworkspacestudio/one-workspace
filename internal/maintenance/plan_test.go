package maintenance

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/distribution"
)

func reviewedFixture() (catalog.Catalog, distribution.UpdatePlan) {
	pack := func(file, content string) catalog.Package {
		h := sha256.Sum256([]byte(content))
		return catalog.Package{File: file, Size: int64(len(content)), SHA256: hex.EncodeToString(h[:])}
	}
	app := catalog.App{ID: "notes", Name: "Notes", Kind: "process", Package: "notes", Executable: "MyNotes.exe", Port: 25112, Route: "/notes", Version: "0.1.0-preview.4", Source: "https://github.com/oneworkspacestudio/notes/releases/tag/v0.1.0-preview.4"}
	current := catalog.Catalog{Schema: 1, Version: "0.3.0-preview.70", InstallerIntegration: 3, Apps: []catalog.App{app}, Packages: map[string]catalog.Package{"manager": pack("Workspace.exe", "manager"), "notes": pack("MyNotes.exe", "old")}}
	b, _ := json.Marshal(current)
	h := sha256.Sum256(b)
	app.Version = "0.1.0-preview.5"
	app.Source = "https://github.com/oneworkspacestudio/notes/releases/tag/v" + app.Version
	p := pack("MyNotes.exe", "new")
	a := distribution.Asset{Name: "MyNotes.exe", Size: p.Size, Digest: "sha256:" + p.SHA256, State: "uploaded", APIURL: "https://api.github.com/repos/oneworkspacestudio/notes/releases/assets/123", URL: "https://github.com/oneworkspacestudio/notes/releases/download/v" + app.Version + "/MyNotes.exe"}
	plan := distribution.UpdatePlan{CheckedAt: time.Now().Add(-65 * time.Minute), BaseCatalogSHA256: hex.EncodeToString(h[:]), Choices: []distribution.UpdateChoice{{ID: "notes", Name: "Notes", Current: current.Apps[0].Version, Version: app.Version, Source: app.Source, Available: true, Apps: []catalog.App{app}, Packages: map[string]catalog.Package{"notes": p}, Downloads: map[string]distribution.UpdateDownload{"notes": {Asset: a, Repository: "oneworkspacestudio/notes", Tag: "v" + app.Version}}}}}
	i := distribution.Asset{Name: "WorkspaceSetup-windows-amd64.exe", Size: 8, Digest: "sha256:" + strings.Repeat("a", 64), State: "uploaded", APIURL: "https://api.github.com/repos/" + distribution.Repository + "/releases/assets/321", URL: "https://github.com/" + distribution.Repository + "/releases/download/v" + current.Version + "/WorkspaceSetup-windows-amd64.exe"}
	plan.Installer = &distribution.InstallerRelease{Version: current.Version, Source: "https://github.com/" + distribution.Repository + "/releases/tag/v" + current.Version, Download: distribution.UpdateDownload{Asset: i, Repository: distribution.Repository, Tag: "v" + current.Version}}
	return current, plan
}

func clonePlan(p distribution.UpdatePlan) distribution.UpdatePlan {
	b, _ := json.Marshal(p)
	var c distribution.UpdatePlan
	_ = json.Unmarshal(b, &c)
	return c
}

func TestExpiredPlanRechecksExactlyReviewedUpdates(t *testing.T) {
	for _, reinstall := range []bool{false, true} {
		for _, change := range []string{"none", "unselected", "selected", "installer", "installation", "network", "stale-response"} {
			t.Run(change+map[bool]string{true: "-reinstall", false: "-apps"}[reinstall], func(t *testing.T) {
				current, reviewed := reviewedFixture()
				fresh := clonePlan(reviewed)
				fresh.CheckedAt = time.Now()
				switch change {
				case "unselected":
					fresh.Choices = append(fresh.Choices, distribution.UpdateChoice{ID: "mail", Available: true, Version: "new"})
				case "selected":
					p := fresh.Choices[0].Packages["notes"]
					p.SHA256 = strings.Repeat("b", 64)
					fresh.Choices[0].Packages["notes"] = p
				case "installer":
					fresh.Installer.Download.Asset.Digest = "sha256:" + strings.Repeat("b", 64)
				case "installation":
					current.Version = "0.3.0-preview.71"
				case "stale-response":
					fresh.CheckedAt = reviewed.CheckedAt
				}
				calls := 0
				check := func(context.Context, catalog.Catalog) (distribution.UpdatePlan, error) {
					calls++
					if change == "network" {
						return distribution.UpdatePlan{}, errors.New("offline")
					}
					return fresh, nil
				}
				got, refreshed, err := refreshReviewedPlan(context.Background(), current, reviewed, []string{"notes"}, reinstall, check)
				wantError := change == "installation" || change == "network" || change == "stale-response" || (!reinstall && change == "selected") || (reinstall && change == "installer")
				if calls != 1 || (err != nil) != wantError {
					t.Fatalf("calls=%d refreshed=%v err=%v", calls, refreshed, err)
				}
				if change == "network" || change == "stale-response" {
					if refreshed || !got.CheckedAt.Equal(reviewed.CheckedAt) {
						t.Fatal("failed check replaced the reviewed plan")
					}
				} else if !refreshed || !got.CheckedAt.Equal(fresh.CheckedAt) {
					t.Fatal("refreshed choices not returned for display")
				}
			})
		}
	}
}

func TestFreshPlanDoesNotCheckAgain(t *testing.T) {
	current, plan := reviewedFixture()
	plan.CheckedAt = time.Now()
	_, refreshed, err := refreshReviewedPlan(context.Background(), current, plan, []string{"notes"}, false, func(context.Context, catalog.Catalog) (distribution.UpdatePlan, error) {
		t.Fatal("fresh check repeated")
		return distribution.UpdatePlan{}, nil
	})
	if err != nil || refreshed {
		t.Fatal(refreshed, err)
	}
}

func TestPlanFreshRejectsMissingFutureAndExpiredChecks(t *testing.T) {
	now := time.Now()
	for _, checked := range []time.Time{{}, now.Add(2 * time.Minute), now.Add(-31 * time.Minute)} {
		if (distribution.UpdatePlan{CheckedAt: checked}).Fresh(now) {
			t.Fatal("invalid check accepted", checked)
		}
	}
}
