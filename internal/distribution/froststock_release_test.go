package distribution

import (
	"context"
	"encoding/json"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

func froststockUpdateFixture(t *testing.T, appVersion string, installed bool) (catalog.Catalog, updateTransport, Release, AppUpdateManifest) {
	t.Helper()
	current, tr, _ := setupUpdateTransport(t, false)
	current.Version = "0.3.0-preview.140"
	current.InstallerIntegration = catalog.IndependentAppIntegration
	for i, id := range []string{"accounts", "drive"} {
		app := catalog.App{ID: id, Name: id, Kind: "process", Package: id, Executable: id + ".exe", Port: 9091 + i, Route: "/" + id, Version: "0.1.0"}
		current.Apps = append(current.Apps, app)
		current.Packages[id] = fixturePackage(app.Executable, []byte(id))
	}
	if !installed {
		current.Selection = &catalog.AppSelection{Excluded: []string{"pet", "froststock"}}
	}
	latest := current
	latest.Version = "0.3.0-preview.141"
	latest.Packages = make(map[string]catalog.Package)
	for id, pkg := range current.Packages {
		latest.Packages[id] = pkg
	}
	latest.Packages["manager"] = fixturePackage("Workspace.exe", []byte("other workspace manager"))
	raw, _ := json.Marshal(latest)
	suiteTag := "v" + latest.Version
	catAsset := fixtureAsset(Repository, suiteTag, "catalog.json", "21", raw)
	bundle := fixtureAsset(Repository, suiteTag, "workspace-bundle.zip", "22", []byte("bundle"))
	tr[catAsset.APIURL] = raw
	tr["https://api.github.com/repos/"+Repository+"/releases?per_page=100"], _ = json.Marshal([]Release{
		{Tag: suiteTag, Published: time.Unix(60, 0), Assets: []Asset{catAsset, bundle}},
	})
	if installed {
		old := froststockLaunch
		old.Port, old.Version, old.Source = 9148, "0.1.0-preview.7", "https://github.com/"+Repository+"/releases/tag/v0.3.0-preview.139"
		current.Apps = append(current.Apps, old)
		for id, path := range froststockPackagePaths {
			current.Packages[id] = fixturePackage(path, []byte("bundle "+id))
		}
	}
	app := froststockLaunch
	app.Version = appVersion
	tag := "v" + appVersion
	app.Source = "https://github.com/" + froststockRepository + "/releases/tag/" + tag
	manifest := AppUpdateManifest{Schema: 2, Integration: 3, MinimumManager: current.Version,
		Apps: []catalog.App{app}, Packages: map[string]UpdatePackage{}, Tag: tag, Repository: froststockRepository}
	r := Release{Tag: tag, Published: time.Unix(65, 0)}
	for id, path := range froststockPackagePaths {
		name := filepath.Base(path)
		payload := []byte("independent " + id)
		d := fixtureAsset(froststockRepository, tag, name, strconv.Itoa(30+len(r.Assets)), payload)
		tr[d.APIURL] = payload
		manifest.Packages[id] = UpdatePackage{Package: fixturePackage(path, payload), Asset: name}
		r.Assets = append(r.Assets, d)
	}
	markerRaw, _ := json.Marshal(manifest)
	marker := fixtureAsset(froststockRepository, tag, "workspace-update.json", "40", markerRaw)
	tr[marker.APIURL] = markerRaw
	r.Assets = append(r.Assets, marker)
	tr["https://api.github.com/repos/"+froststockRepository+"/releases?per_page=100"], _ = json.Marshal([]Release{r})
	return current, tr, r, manifest
}

func TestExternalFroststockCanMigrateOrInstallWithoutWorkspaceUpdate(t *testing.T) {
	for _, installed := range []bool{false, true} {
		label := "new install with earlier opt-out"
		if installed {
			label = "migration from bundled app"
		}
		t.Run(label, func(t *testing.T) {
			current, _, _, _ := froststockUpdateFixture(t, "0.1.0-preview.8", installed)
			if err := current.Validate(); err != nil {
				t.Fatal(err)
			}
			plan, err := CheckUpdates(context.Background(), current)
			if err != nil {
				t.Fatal(err)
			}
			var workspaceAvailable bool
			for _, c := range plan.Choices {
				if c.ID == "workspace" {
					workspaceAvailable = c.Available
				}
				if c.ID == "froststock" && (!c.Available || c.Optional == installed || len(c.NewApps) == 0 && !installed) {
					t.Fatalf("incorrect Froststock offer: %+v", c)
				}
			}
			if !workspaceAvailable {
				t.Fatal("fixture should offer a new Workspace version")
			}
			target, sources, err := SelectUpdates(current, plan, []string{"froststock"})
			if err != nil {
				t.Fatal(err)
			}
			app, ok := target.Find("froststock")
			if !ok || app.Version != "0.1.0-preview.8" || app.Port != 9148 || target.Selection.Excludes("froststock") {
				t.Fatalf("incorrect installed Froststock: %+v", app)
			}
			if target.Version != current.Version || target.Packages["manager"] != current.Packages["manager"] ||
				len(sources) != len(froststockPackagePaths) || (!installed && !target.Selection.Excludes("pet")) {
				t.Fatal("independent selection changed the manager or an unrelated app")
			}
			for id, download := range sources {
				if download.Repository != froststockRepository || froststockPackagePaths[id] == "" {
					t.Fatalf("unexpected app download %s: %+v", id, download)
				}
			}
			if err = StageUpdateFiles(context.Background(), target, sources, t.TempDir(), func(string) {}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestExternalFroststockRejectsOlderVersionsAndUnapprovedLaunch(t *testing.T) {
	current, _, r, manifest := froststockUpdateFixture(t, "0.1.0-preview.8", true)
	group := updateGroup{"froststock", "Froststock", "froststock", []string{"froststock"}}
	raw, _ := json.Marshal(manifest)
	newer := current
	newer.Apps = append([]catalog.App(nil), current.Apps...)
	newer.Apps[len(newer.Apps)-1].Version = "0.1.0-preview.9"
	if _, err := froststockManifestChoice(group, newer, r, raw); err == nil {
		t.Fatal("older Froststock release was accepted")
	}
	olderManager := current
	olderManager.Version = "0.3.0-preview.139"
	if _, err := froststockManifestChoice(group, olderManager, r, raw); err == nil {
		t.Fatal("Froststock release was accepted by an older Workspace version")
	}
	manifest.Apps[0].Args = []string{"--data", "C:/other"}
	raw, _ = json.Marshal(manifest)
	if _, err := froststockManifestChoice(group, current, r, raw); err == nil {
		t.Fatal("changed managed launch arguments were accepted")
	}
}

func TestWorkspaceInstallerRetainsExistingFroststockOutsideBundle(t *testing.T) {
	current, _, _, _ := froststockUpdateFixture(t, "0.1.0-preview.8", true)
	current.Version, current.InstallerIntegration = "0.3.0-preview.139", catalog.ManagerPreflightIntegration
	incoming := current
	incoming.Version, incoming.InstallerIntegration = "0.3.0-preview.140", catalog.IndependentAppIntegration
	incoming.Packages = map[string]catalog.Package{"manager": fixturePackage("Workspace.exe", []byte("compatibility manager"))}
	incoming.Apps = nil
	for _, app := range current.Apps {
		if app.ID != "froststock" {
			incoming.Apps = append(incoming.Apps, app)
		}
	}
	target, err := SelectInstallerUpgrade(current, incoming)
	if err != nil {
		t.Fatal(err)
	}
	old, _ := current.Find("froststock")
	new, ok := target.Find("froststock")
	if !ok || !reflect.DeepEqual(old, new) || !reflect.DeepEqual(groupPackages(updateGroup{id: "froststock"}, current), groupPackages(updateGroup{id: "froststock"}, target)) {
		t.Fatal("Workspace upgrade changed the installed Froststock")
	}
}
