package distribution

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

func TestFutureCatalogOffersVerifiedInstallerWithoutInterpretingApps(t *testing.T) {
	current, tr, _ := setupUpdateTransport(t, true)
	var releases []Release
	url := "https://api.github.com/repos/" + Repository + "/releases?per_page=100"
	json.Unmarshal(tr[url], &releases)
	raw, _ := json.Marshal(map[string]any{"schema": 1, "installerIntegration": catalog.InstallerIntegration + 1, "version": "0.3.0-preview.23", "futureApps": "unreadable by this installer"})
	r := &releases[0]
	r.Assets[0] = fixtureAsset(Repository, r.Tag, "catalog.json", "1", raw)
	installer := fixtureAsset(Repository, r.Tag, installerName, "20", []byte("verified-installer"))
	r.Assets = append(r.Assets, installer)
	tr[r.Assets[0].APIURL] = raw
	tr[installer.APIURL] = []byte("verified-installer")
	tr[url], _ = json.Marshal(releases)
	plan, e := CheckUpdates(context.Background(), current)
	if e != nil || plan.Installer == nil || !plan.Installer.Required || len(plan.Choices) != 0 {
		t.Fatal(plan, e)
	}
	if e = ValidateInstallerPlan(current, plan); e != nil {
		t.Fatal(e)
	}
	dir := t.TempDir()
	path, e := DownloadInstaller(context.Background(), *plan.Installer, dir, nil)
	if e != nil {
		t.Fatal(e)
	}
	b, _ := os.ReadFile(path)
	if string(b) != "verified-installer" {
		t.Fatal(string(b))
	}
	tr[installer.APIURL] = []byte("tampered-installer")
	dir = t.TempDir()
	if _, e = DownloadInstaller(context.Background(), *plan.Installer, dir, nil); e == nil {
		t.Fatal("tampered installer accepted")
	}
	if _, e = os.Stat(filepath.Join(dir, installerName)); !os.IsNotExist(e) {
		t.Fatal("bad executable retained", e)
	}
	plan.CheckedAt = time.Now().Add(-31 * time.Minute)
	if ValidateInstallerPlan(current, plan) == nil {
		t.Fatal("expired plan accepted")
	}
	plan.CheckedAt = time.Now()
	plan.BaseCatalogSHA256 = "changed"
	if ValidateInstallerPlan(current, plan) == nil {
		t.Fatal("changed installation accepted")
	}
	plan.BaseCatalogSHA256 = catalogIdentity(current)
	plan.Installer.Download.Asset.URL = "https://attacker.example/installer.exe"
	if ValidateInstallerPlan(current, plan) == nil {
		t.Fatal("untrusted executable URL accepted")
	}
}

func TestCompatibleAndMissingOrInvalidInstallerMetadata(t *testing.T) {
	current, tr, _ := setupUpdateTransport(t, true)
	plan, e := CheckUpdates(context.Background(), current)
	if e != nil || plan.Installer != nil || len(plan.Choices) == 0 {
		t.Fatal(plan, e)
	}
	url := "https://api.github.com/repos/" + Repository + "/releases?per_page=100"
	var releases []Release
	json.Unmarshal(tr[url], &releases)
	for _, bad := range []bool{false, true} {
		raw, _ := json.Marshal(map[string]any{"schema": 1, "installerIntegration": catalog.InstallerIntegration + 1, "version": "0.3.0-preview.23"})
		r := releases[0]
		r.Assets[0] = fixtureAsset(Repository, r.Tag, "catalog.json", "1", raw)
		if bad {
			a := fixtureAsset(Repository, r.Tag, installerName, "22", []byte("x"))
			a.Digest = "sha256:bad"
			r.Assets = append(r.Assets, a)
		}
		tr[r.Assets[0].APIURL] = raw
		tr[url], _ = json.Marshal([]Release{r})
		if _, e = CheckUpdates(context.Background(), current); e == nil {
			t.Fatal("future catalog without trusted installer accepted")
		}
	}
}
