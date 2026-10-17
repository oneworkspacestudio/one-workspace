package suite

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/fahimc/one-workspace/internal/catalog"
	"github.com/fahimc/one-workspace/internal/platform"
)

func updateFixture(t *testing.T) (Config, Secrets) {
	t.Helper()
	root := t.TempDir()
	c := Config{Schema: 1, Version: "0.3.0-preview.23", Data: filepath.Join(root, "Private"), Programs: filepath.Join(root, "Apps"), Files: filepath.Join(root, "Files"), Origin: "http://localhost:29100", Issuer: "http://localhost:29101", Port: 29100, Owner: "fixture-owner"}
	c.Catalog = catalog.Catalog{Schema: 1, InstallerIntegration: 2, Version: c.Version, Packages: map[string]catalog.Package{"manager": updatePackage("Workspace.exe", "manager-original"), "notes": updatePackage("MyNotes.exe", "notes-original")}, Apps: []catalog.App{{ID: "notes", Name: "Notes", Kind: "process", Package: "notes", Executable: "MyNotes.exe", Port: 29112, Route: "/notes", Version: "0.1.0-preview.4", OIDC: true, Callback: "/notes/auth/callback", Source: "https://github.com/oneworkspacestudio/notes/releases/tag/v0.1.0-preview.4"}}}
	for _, dir := range []string{c.Data, c.Programs, c.Files} {
		if e := os.MkdirAll(dir, 0700); e != nil {
			t.Fatal(e)
		}
	}
	if e := platform.ProtectDir(c.Data); e != nil {
		t.Fatal(e)
	}
	s := Secrets{AdminEmail: "fixture@example.test", AdminPassword: "fixture-recovery", Clients: map[string]Client{"notes": {ID: "existing-client", Secret: "fixture-secret"}}}
	if e := Save(c, s); e != nil {
		t.Fatal(e)
	}
	for path, value := range map[string]string{filepath.Join(c.Data, "Workspace.exe"): "manager-original", filepath.Join(c.Programs, "Workspace.exe"): "manager-original", filepath.Join(c.Programs, "MyNotes.exe"): "notes-original", filepath.Join(c.Data, "notes.db"): "private original", filepath.Join(c.Files, "Document.txt"): "ordinary original"} {
		if e := os.WriteFile(path, []byte(value), 0600); e != nil {
			t.Fatal(e)
		}
	}
	if e := os.MkdirAll(filepath.Join(c.Files, "Empty folder"), 0700); e != nil {
		t.Fatal(e)
	}
	return c, s
}
func updatePackage(name, value string) catalog.Package {
	h := sha256.Sum256([]byte(value))
	return catalog.Package{File: name, Size: int64(len(value)), SHA256: hex.EncodeToString(h[:])}
}
func fixtureTarget(t *testing.T, c Config) (catalog.Catalog, string) {
	t.Helper()
	target := c.Catalog
	target.Version = "0.3.0-preview.24"
	target.Packages = map[string]catalog.Package{}
	for id, p := range c.Catalog.Packages {
		target.Packages[id] = p
	}
	target.Packages["manager"] = updatePackage("Workspace.exe", "manager-updated")
	stage := t.TempDir()
	if e := os.WriteFile(filepath.Join(stage, "Workspace.exe"), []byte("manager-updated"), 0600); e != nil {
		t.Fatal(e)
	}
	return target, stage
}
func quietUpdate(string) {}
func TestUpdateLaunchPreflightFailureLeavesRuntimeAndDataUntouched(t *testing.T) {
	c, s := updateFixture(t)
	target, stage := fixtureTarget(t, c)
	stops, starts := 0, 0
	ops := updateOperations{
		preflight: func(context.Context, Config) error { return errors.New("launch denied by fixture host") },
		stop:      func(context.Context, Config) error { stops++; return nil },
		start:     func(context.Context, Config, string) error { starts++; return nil },
	}
	err := applyStagedWith(context.Background(), c, s, target, stage, "", quietUpdate, ops)
	if err == nil || !strings.Contains(err.Error(), "launch denied") || stops != 0 || starts != 0 {
		t.Fatalf("preflight did not abort safely: %v stops=%d starts=%d", err, stops, starts)
	}
	if UpdatePending(c) {
		t.Fatal("preflight created an unfinished transaction")
	}
	got, _, err := Load(c.Data)
	if err != nil || got.Version != c.Version {
		t.Fatal("preflight changed installation", err)
	}
}

func TestCompletedRollbackWithLaunchFailureGivesRestartGuidance(t *testing.T) {
	c, s := updateFixture(t)
	target, stage := fixtureTarget(t, c)
	ops := updateOperations{stop: func(context.Context, Config) error { return nil }, start: func(context.Context, Config, string) error { return errors.New("fixture access denied") }}
	err := applyStagedWith(context.Background(), c, s, target, stage, "", quietUpdate, ops)
	if !errors.Is(err, errRestoredNeedsStart) || !strings.Contains(err.Error(), "Recover is not needed") || strings.Contains(err.Error(), "choose Recover") {
		t.Fatalf("completed restore gave incorrect recovery instructions: %v", err)
	}
	if UpdatePending(c) {
		t.Fatal("completed restore should not be pending")
	}
}

func TestUpdateCommitPreservesIdentityDataAndLocations(t *testing.T) {
	c, s := updateFixture(t)
	target, stage := fixtureTarget(t, c)
	starts := 0
	ops := updateOperations{stop: func(context.Context, Config) error { return nil }, start: func(_ context.Context, next Config, token string) error {
		starts++
		if token == "" || !UpdatePending(next) {
			t.Fatal("startup was not isolated until commit")
		}
		got, secret, e := Load(next.Data)
		if e != nil || got.Owner != c.Owner || got.Files != c.Files || got.Programs != c.Programs || secret.Clients["notes"] != s.Clients["notes"] {
			t.Fatal("identity or paths changed", e)
		}
		return nil
	}}
	if e := applyStagedWith(context.Background(), c, s, target, stage, "", quietUpdate, ops); e != nil {
		t.Fatal(e)
	}
	if starts != 1 || UpdatePending(c) {
		t.Fatal("update did not commit")
	}
	now, _, e := Load(c.Data)
	if e != nil || now.Version != target.Version {
		t.Fatal(e, now.Version)
	}
	for path, want := range map[string]string{filepath.Join(c.Data, "Workspace.exe"): "manager-updated", filepath.Join(c.Data, "notes.db"): "private original", filepath.Join(c.Files, "Document.txt"): "ordinary original"} {
		got, _ := os.ReadFile(path)
		if string(got) != want {
			t.Fatal(path, string(got))
		}
	}
	backups, e := UpdateBackups(c)
	if e != nil || len(backups) != 1 {
		t.Fatal(backups, e)
	}
	if _, e = verifyUpdateBackup(c, backups[0]); e != nil {
		t.Fatal(e)
	}
	if _, e = os.Stat(filepath.Join(backups[0].Path, "files", "Empty folder")); e != nil {
		t.Fatal("empty ordinary folder not backed up", e)
	}
}
func TestMaintenanceLockCanCoexistWithOpenWindow(t *testing.T) {
	c, _ := updateFixture(t)
	root := filepath.Join(c.Data, "updates")
	os.MkdirAll(root, 0700)
	if e := platform.ProtectMaintenanceDir(root, c.Data); e != nil {
		t.Fatal(e)
	}
	window, e := platform.Lock(filepath.Join(root, "window.lock"))
	if e != nil {
		t.Fatal(e)
	}
	defer window.Close()
	update, e := lockUpdate(c)
	if e != nil {
		t.Fatal("cannot update from open window", e)
	}
	update.Close()
}
func TestBackupAndRestorePreserveNamedStreamsReadOnlyAndModificationTime(t *testing.T) {
	c, _ := updateFixture(t)
	doc := filepath.Join(c.Files, "Document.txt")
	stream := doc + ":fixture-metadata"
	if e := os.WriteFile(stream, []byte("original named stream"), 0600); e != nil {
		t.Fatal(e)
	}
	originalTime := time.Date(2023, 5, 4, 3, 2, 1, 0, time.UTC)
	if e := os.Chtimes(doc, originalTime, originalTime); e != nil {
		t.Fatal(e)
	}
	if e := os.Chmod(doc, 0400); e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { os.Chmod(doc, 0600) })
	backup, e := makeUpdateBackup(context.Background(), c, "", quietUpdate)
	if e != nil {
		t.Fatal(e)
	}
	os.Chmod(doc, 0600)
	os.WriteFile(doc, []byte("new contents"), 0600)
	os.WriteFile(stream, []byte("new named stream"), 0600)
	recovery, e := makeUpdateBackup(context.Background(), c, "", quietUpdate)
	if e != nil {
		t.Fatal(e)
	}
	if e = restoreUpdateBackup(c, backup, recovery, quietUpdate); e != nil {
		t.Fatal(e)
	}
	value, _ := os.ReadFile(stream)
	if string(value) != "original named stream" {
		t.Fatal("named stream was lost", string(value))
	}
	value, _ = os.ReadFile(filepath.Join(recovery.Path, "files", "Document.txt") + ":fixture-metadata")
	if string(value) != "new named stream" {
		t.Fatal("newer named stream was not archived")
	}
	info, e := os.Stat(doc)
	if e != nil || !info.ModTime().Equal(originalTime) || info.Mode().Perm()&0200 != 0 {
		t.Fatal("file metadata was not restored", info, e)
	}
}
func TestStorageMigrationRefusedBeforeStopping(t *testing.T) {
	c, s := updateFixture(t)
	target, stage := fixtureTarget(t, c)
	target.Apps = append([]catalog.App(nil), c.Catalog.Apps...)
	c.Catalog.Apps[0].Env = map[string]string{"NOTES_DATA_DIR": "${DATA}/notes"}
	target.Apps[0].Env = map[string]string{"NOTES_DATA_DIR": "${DATA}/different-notes"}
	stops := 0
	ops := updateOperations{stop: func(context.Context, Config) error { stops++; return nil }, start: func(context.Context, Config, string) error { return nil }}
	if e := applyStagedWith(context.Background(), c, s, target, stage, "", quietUpdate, ops); e == nil || stops != 0 {
		t.Fatal("persistent path migration was automatic", e, stops)
	}
}
func TestFailedHealthRestoresAllDataAndArchivesMigratedState(t *testing.T) {
	c, s := updateFixture(t)
	target, stage := fixtureTarget(t, c)
	secretBefore, _ := os.ReadFile(filepath.Join(c.Data, "secrets.dpapi"))
	starts := 0
	ops := updateOperations{stop: func(context.Context, Config) error { return nil }, start: func(_ context.Context, next Config, token string) error {
		starts++
		if starts == 1 {
			os.WriteFile(filepath.Join(c.Data, "notes.db"), []byte("migrated private"), 0600)
			os.WriteFile(filepath.Join(c.Data, "new-schema.db"), []byte("new schema"), 0600)
			os.WriteFile(filepath.Join(c.Files, "Document.txt"), []byte("newer ordinary"), 0600)
			os.WriteFile(filepath.Join(c.Files, "Added.txt"), []byte("keep this extra file"), 0600)
			return errors.New("injected unhealthy app")
		}
		if token != "" || next.Version != c.Version {
			t.Fatal("rollback started wrong version")
		}
		return nil
	}}
	e := applyStagedWith(context.Background(), c, s, target, stage, "", quietUpdate, ops)
	if e == nil || !strings.Contains(e.Error(), "was rolled back") || starts != 2 {
		t.Fatal("failure did not roll back", e, starts)
	}
	for path, want := range map[string]string{filepath.Join(c.Data, "Workspace.exe"): "manager-original", filepath.Join(c.Data, "notes.db"): "private original", filepath.Join(c.Files, "Document.txt"): "ordinary original", filepath.Join(c.Files, "Added.txt"): "keep this extra file"} {
		got, _ := os.ReadFile(path)
		if string(got) != want {
			t.Fatal(path, string(got))
		}
	}
	secretAfter, _ := os.ReadFile(filepath.Join(c.Data, "secrets.dpapi"))
	if string(secretBefore) != string(secretAfter) {
		t.Fatal("original credentials not restored byte for byte")
	}
	if _, e = os.Stat(filepath.Join(c.Data, "new-schema.db")); !os.IsNotExist(e) {
		t.Fatal("new schema mixed with old database")
	}
	var j UpdateJournal
	if e = platform.LoadJSON(updateJournalPath(c), &j); e != nil || j.Phase != "rolled-back" || j.Recovery == nil {
		t.Fatal(e, j.Phase)
	}
	for rel, want := range map[string]string{"private/notes.db": "migrated private", "private/new-schema.db": "new schema", "files/Document.txt": "newer ordinary"} {
		got, _ := os.ReadFile(filepath.Join(j.Recovery.Path, filepath.FromSlash(rel)))
		if string(got) != want {
			t.Fatal("newer state was lost", rel, string(got))
		}
	}
}
func TestCorruptBackupRefusedBeforeAnyStopOrRestore(t *testing.T) {
	c, _ := updateFixture(t)
	r, e := makeUpdateBackup(context.Background(), c, "", quietUpdate)
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(filepath.Join(r.Path, "files", "Document.txt"), []byte("corrupt"), 0600); e != nil {
		t.Fatal(e)
	}
	stops := 0
	ops := updateOperations{stop: func(context.Context, Config) error { stops++; return nil }, start: func(context.Context, Config, string) error { return nil }}
	e = recoverUpdateWith(context.Background(), c, &UpdateJournal{Phase: "applying", Backup: r}, quietUpdate, ops)
	if e == nil || stops != 0 {
		t.Fatal("corrupt backup touched installation", e, stops)
	}
	got, _ := os.ReadFile(filepath.Join(c.Files, "Document.txt"))
	if string(got) != "ordinary original" {
		t.Fatal("ordinary file changed")
	}
}
func TestInterruptedRestoreCanResumeWithoutWorkspaceJSON(t *testing.T) {
	c, _ := updateFixture(t)
	if e := cacheUpdateController(c, c.Catalog, t.TempDir()); e != nil {
		t.Fatal(e)
	}
	original, e := makeUpdateBackup(context.Background(), c, "", quietUpdate)
	if e != nil {
		t.Fatal(e)
	}
	os.WriteFile(filepath.Join(c.Data, "notes.db"), []byte("post update"), 0600)
	archive, e := makeUpdateBackup(context.Background(), c, "", quietUpdate)
	if e != nil {
		t.Fatal(e)
	}
	j := UpdateJournal{Phase: "restoring", Backup: original, Recovery: &archive}
	if e = saveUpdateJournal(c, &j); e != nil {
		t.Fatal(e)
	}
	os.Remove(filepath.Join(c.Data, "workspace.json"))
	loaded, e := LoadUpdateConfig(c.Data)
	if e != nil || loaded.Data != c.Data {
		t.Fatal("could not locate recovery state", e)
	}
	ops := updateOperations{stop: func(context.Context, Config) error { return nil }, start: func(context.Context, Config, string) error { return nil }}
	if e = recoverUpdateWith(context.Background(), loaded, &j, quietUpdate, ops); e != nil {
		t.Fatal(e)
	}
	got, _ := os.ReadFile(filepath.Join(c.Data, "notes.db"))
	if string(got) != "private original" {
		t.Fatal("recovery did not resume")
	}
	if _, _, e = Load(c.Data); e != nil {
		t.Fatal(e)
	}
}
func TestUpdateJournalBlocksOrdinaryStartupAndAllowsTestNonce(t *testing.T) {
	c, _ := updateFixture(t)
	j := UpdateJournal{Phase: "testing", Token: "fixture-nonce"}
	if e := saveUpdateJournal(c, &j); e != nil {
		t.Fatal(e)
	}
	t.Setenv("WORKSPACE_UPDATE_TRANSACTION", "")
	if CheckUpdateRecovery(c) == nil {
		t.Fatal("ordinary startup bypassed recovery")
	}
	t.Setenv("WORKSPACE_UPDATE_TRANSACTION", "fixture-nonce")
	if e := CheckUpdateRecovery(c); e != nil {
		t.Fatal(e)
	}
	j.Phase = "applying"
	saveUpdateJournal(c, &j)
	if CheckUpdateRecovery(c) == nil {
		t.Fatal("nonce allowed startup before programs were installed")
	}
	j.Phase = "complete"
	j.Token = ""
	saveUpdateJournal(c, &j)
	if e := CheckUpdateRecovery(c); e != nil {
		t.Fatal(e)
	}
}
func TestUpdateRefusesJunctionsAndExternalInstallationPointers(t *testing.T) {
	c, s := updateFixture(t)
	other := t.TempDir()
	junction := filepath.Join(c.Files, "link")
	// Windows directory symlinks require Developer Mode; junctions do not.
	cmd := execJunction(junction, other)
	if cmd == nil {
		if _, e := makeUpdateBackup(context.Background(), c, "", quietUpdate); e == nil {
			t.Fatal("junction was followed")
		}
		os.Remove(junction)
	} else {
		t.Log("junction creation unavailable", cmd)
	}
	fake := filepath.Join(t.TempDir(), "CopiedPrivate")
	os.MkdirAll(fake, 0700)
	copy := c
	copy.Data = c.Data
	if e := platform.SaveJSON(filepath.Join(fake, "workspace.json"), copy); e != nil {
		t.Fatal(e)
	}
	platform.SaveSecret(filepath.Join(fake, "secrets.dpapi"), s)
	if _, e := LoadUpdateConfig(fake); e == nil {
		t.Fatal("copied pointer could target another installation")
	}
}
func TestUpdateOwnerAndOriginGates(t *testing.T) {
	p := &Portal{rt: &Runtime{Config: Config{Owner: "owner", Origin: "http://localhost:29100"}}, sessions: map[string]session{"owner-token": {Subject: "owner", Expires: time.Now().Add(time.Hour)}, "member-token": {Subject: "member", Expires: time.Now().Add(time.Hour)}}}
	mux := http.NewServeMux()
	p.updateRoutes(mux)
	for _, tc := range []struct {
		method, path, token, origin string
		want                        int
	}{{"GET", "/updates", "owner-token", "", 200}, {"GET", "/updates", "member-token", "", 403}, {"GET", "/updates", "", "", 403}, {"POST", "/updates/open", "owner-token", "https://other.example", 403}} {
		r := httptest.NewRequest(tc.method, tc.path, nil)
		r.AddCookie(&http.Cookie{Name: "workspace_session", Value: tc.token})
		r.Header.Set("Origin", tc.origin)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != tc.want {
			t.Fatal(tc, w.Code)
		}
	}
}
