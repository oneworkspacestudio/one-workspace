package suite

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

func TestReinstallReviewIncompleteRetentionAndFiles(t *testing.T) {
	root := filepath.Clean(t.TempDir())
	r := Request{Folder: root}
	if err := r.ResolveStorage(); err != nil {
		t.Fatal(err)
	}
	for _, dir := range []string{r.Data, r.Programs, r.Files} {
		os.MkdirAll(dir, 0700)
	}
	os.WriteFile(filepath.Join(r.Data, "unfinished.db"), []byte("previous account data"), 0600)
	os.WriteFile(filepath.Join(r.Programs, "failed.exe"), []byte("partial program"), 0600)
	os.WriteFile(filepath.Join(r.Files, "Keep.txt"), []byte("ordinary file"), 0600)
	review, err := InspectSetupStorage(r)
	if err != nil || review.Kind != "restart" {
		t.Fatal(review, err)
	}
	if _, err = r.PreflightStorage(); err == nil {
		t.Fatal("fresh install accepted occupied folder")
	}
	os.WriteFile(filepath.Join(r.Data, "new.db"), []byte("changed"), 0600)
	if _, err = retainIncompleteSetup(r, review); err == nil {
		t.Fatal("stale review accepted")
	}
	review, err = InspectSetupStorage(r)
	if err != nil {
		t.Fatal(err)
	}
	backup, err := retainIncompleteSetup(r, review)
	if err != nil {
		t.Fatal(err)
	}
	for name, want := range map[string]string{filepath.Join(backup, "Private", "unfinished.db"): "previous account data", filepath.Join(backup, "Private", "new.db"): "changed", filepath.Join(backup, "Apps", "failed.exe"): "partial program", filepath.Join(r.Files, "Keep.txt"): "ordinary file"} {
		b, e := os.ReadFile(name)
		if e != nil || string(b) != want {
			t.Fatal(name, string(b), e)
		}
	}
	if _, err = InspectSetupStorage(r); err != nil {
		t.Fatal(err)
	}
	if _, err = r.PreflightStorage(); err != nil {
		t.Fatal("retained Files prevented fresh install", err)
	}
}

func TestIncompleteReinstallRefusesActiveSetupAndJunction(t *testing.T) {
	r := Request{Folder: filepath.Clean(t.TempDir())}
	r.ResolveStorage()
	os.MkdirAll(r.Data, 0700)
	os.WriteFile(filepath.Join(r.Data, "unfinished"), []byte("keep"), 0600)
	review, err := InspectSetupStorage(r)
	if err != nil {
		t.Fatal(err)
	}
	lock, err := platform.Lock(filepath.Join(r.Data, "install.lock"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = retainIncompleteSetup(r, review); err == nil {
		t.Fatal("active install accepted")
	}
	lock.Close()
	outside := t.TempDir()
	if err = execJunction(filepath.Join(r.Data, "linked"), outside); err != nil {
		t.Fatal(err)
	}
	if _, err = InspectSetupStorage(r); err == nil {
		t.Fatal("junction accepted")
	}
	// Remove only the task-created empty junction; never traverse its target.
	if err = os.Remove(filepath.Join(r.Data, "linked")); err != nil {
		t.Fatal(err)
	}
}

func TestForcedReinstallRepairsSameVersionAndMissingManager(t *testing.T) {
	c, s := updateFixture(t)
	stage := t.TempDir()
	for _, p := range c.Catalog.Packages {
		if err := platform.CopyFile(filepath.Join(c.Programs, p.File), filepath.Join(stage, p.File)); err != nil {
			t.Fatal(err)
		}
	}
	os.Remove(filepath.Join(c.Data, "Workspace.exe"))
	os.WriteFile(filepath.Join(c.Programs, "MyNotes.exe"), []byte("damaged program"), 0600)
	r := Request{Programs: c.Programs, Data: c.Data, Files: c.Files}
	review, err := InspectSetupStorage(r)
	if err != nil || review.Kind != "repair" {
		t.Fatal(review, err)
	}
	stops, starts := 0, 0
	ops := updateOperations{repair: true, stop: func(context.Context, Config) error { stops++; return nil }, start: func(context.Context, Config, string) error { starts++; return nil }}
	if err = applyStagedWith(context.Background(), c, s, c.Catalog, stage, "", quietUpdate, ops); err != nil {
		t.Fatal(err)
	}
	if stops != 1 || starts != 1 {
		t.Fatal(stops, starts)
	}
	if err = verifyInstalledPrograms(c); err != nil {
		t.Fatal(err)
	}
	after, secrets, err := Load(c.Data)
	if err != nil || after.Owner != c.Owner || !reflect.DeepEqual(s, secrets) {
		t.Fatal("identity changed", err)
	}
	if b, _ := os.ReadFile(filepath.Join(c.Files, "Document.txt")); string(b) != "ordinary original" {
		t.Fatal(string(b))
	}
	backups, err := UpdateBackups(c)
	if err != nil || len(backups) == 0 {
		t.Fatal(backups, err)
	}
	if b, _ := os.ReadFile(filepath.Join(backups[0].Path, "programs", "MyNotes.exe")); string(b) != "damaged program" {
		t.Fatal("previous bytes were not retained", string(b))
	}
}

func TestReinstallRejectsBadReplacementBeforeDowntime(t *testing.T) {
	c, s := updateFixture(t)
	_, stage := fixtureTarget(t, c)
	stops := 0
	err := applyStagedWith(context.Background(), c, s, c.Catalog, stage, "", quietUpdate, updateOperations{repair: true, stop: func(context.Context, Config) error { stops++; return nil }})
	if err == nil || stops != 0 {
		t.Fatal(err, stops)
	}
}

func TestRepairCatalogRetainsIndependentAppsAndExclusions(t *testing.T) {
	c, _ := updateFixture(t)
	incoming := c.Catalog
	incoming.Version = "0.3.0-preview.64"
	incoming.Packages = map[string]catalog.Package{}
	for id, p := range c.Catalog.Packages {
		incoming.Packages[id] = p
	}
	incoming.Packages["manager"] = updatePackage("Workspace.exe", "new manager")
	// App updates are separate from reinstalling an existing app's exact bytes.
	incoming.Apps = append([]catalog.App(nil), c.Catalog.Apps...)
	incoming.Apps[0].Version = "0.1.0-preview.5"
	incoming.Packages["notes"] = updatePackage("MyNotes.exe", "newer notes")
	got, err := repairCatalog(c.Catalog, incoming)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Apps, c.Catalog.Apps) || got.Packages["notes"] != c.Catalog.Packages["notes"] || got.Version != incoming.Version {
		t.Fatal(got)
	}
	same, err := repairCatalog(c.Catalog, c.Catalog)
	if err != nil || !reflect.DeepEqual(same, c.Catalog) {
		t.Fatal(err)
	}
	for _, step := range ReinstallProgress().Snapshot().Steps {
		if strings.Contains(strings.ToLower(step.Label), "create") {
			t.Fatal("repair claims to recreate identity")
		}
	}
}
