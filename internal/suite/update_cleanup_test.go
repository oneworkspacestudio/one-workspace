package suite

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

func cleanupCopy(t *testing.T, c Config, purpose string) BackupReceipt {
	t.Helper()
	r, e := makeBackupDetailed(context.Background(), c, "", quietUpdate, nil, purpose)
	if e != nil {
		t.Fatal(e)
	}
	if e = recordBackup(c, r); e != nil {
		t.Fatal(e)
	}
	return r
}
func cleanupReady(t *testing.T, c Config, r BackupReceipt, previous *BackupReceipt) {
	t.Helper()
	j := UpdateJournal{Phase: "complete", Backup: r, PreviousAutomaticBackup: previous}
	if e := saveUpdateJournal(c, &j); e != nil {
		t.Fatal(e)
	}
}
func requireExists(t *testing.T, path string, want bool) {
	t.Helper()
	_, e := os.Stat(path)
	if (e == nil) != want {
		t.Fatalf("exists %s = %v, want %v", path, e == nil, want)
	}
}

func TestUpdateCleanupRetainsRecoverySnapshotsAndUnknownStages(t *testing.T) {
	c, _ := updateFixture(t)
	old := cleanupCopy(t, c, "update")
	manual := cleanupCopy(t, c, "snapshot")
	legacy := cleanupCopy(t, c, "")
	recovery := cleanupCopy(t, c, "snapshot")
	var e error
	recovery, e = markRecoveryOnly(c, recovery)
	if e != nil {
		t.Fatal(e)
	}
	rs, _ := UpdateBackups(c)
	for i := range rs {
		if rs[i].ID == recovery.ID {
			rs[i] = recovery
		}
	}
	if e = platform.SaveJSON(filepath.Join(c.Data, "updates", "backups.json"), rs); e != nil {
		t.Fatal(e)
	}
	keep := cleanupCopy(t, c, "update")
	cleanupReady(t, c, keep, nil)
	stale, e := newUpdateStage(c)
	if e != nil {
		t.Fatal(e)
	}
	active, e := newUpdateStage(c)
	if e != nil {
		t.Fatal(e)
	}
	unknown := filepath.Join(c.Data, "updates", "stage-12345")
	os.MkdirAll(unknown, 0700)
	os.WriteFile(filepath.Join(unknown, "unrecognised.txt"), []byte("keep"), 0600)
	r := cleanupSuccessfulUpdate(context.Background(), c, keep, filepath.Join(active, "packages"), quietUpdate)
	if r.RemovedBackups != 1 || r.RemovedStages != 1 || len(r.Warnings) != 0 || r.BytesFreed <= 0 {
		t.Fatalf("%+v", r)
	}
	requireExists(t, old.Path, false)
	requireExists(t, stale, false)
	for _, p := range []string{manual.Path, legacy.Path, recovery.Path, keep.Path, active, unknown} {
		requireExists(t, p, true)
	}
	if _, e = verifyUpdateBackup(c, keep); e != nil {
		t.Fatal(e)
	}
	rs, _ = UpdateBackups(c)
	if len(rs) != 4 {
		t.Fatal("history still references a removed backup", len(rs))
	}
}

func TestUpdateCleanupRequiresSuccessfulCommitAndIntactRetainedCopy(t *testing.T) {
	for _, mode := range []string{"unfinished", "corrupt"} {
		t.Run(mode, func(t *testing.T) {
			c, _ := updateFixture(t)
			old := cleanupCopy(t, c, "update")
			keep := cleanupCopy(t, c, "update")
			cleanupReady(t, c, keep, nil)
			if mode == "unfinished" {
				j := UpdateJournal{Phase: "testing", Backup: keep}
				saveUpdateJournal(c, &j)
			} else {
				os.WriteFile(filepath.Join(keep.Path, "files", "Document.txt"), []byte("corrupted"), 0600)
			}
			r := cleanupSuccessfulUpdate(context.Background(), c, keep, "", quietUpdate)
			if len(r.Warnings) == 0 || r.RemovedBackups != 0 {
				t.Fatalf("%+v", r)
			}
			requireExists(t, old.Path, true)
		})
	}
}

func TestUpdateCleanupLockedBackupResumesAndCannotBeRestoredWhilePartial(t *testing.T) {
	c, _ := updateFixture(t)
	old := cleanupCopy(t, c, "update")
	keep := cleanupCopy(t, c, "update")
	cleanupReady(t, c, keep, nil)
	lockedPath := filepath.Join(old.Path, "private", "notes.db")
	sum, _, e := fileDigest(lockedPath)
	if e != nil {
		t.Fatal(e)
	}
	// Allow verification reads, but deny deletion until the handle is closed.
	held, e := platform.LockedBinary(lockedPath, sum)
	if e != nil {
		t.Fatal(e)
	}
	r := cleanupSuccessfulUpdate(context.Background(), c, keep, "", quietUpdate)
	if r.RemovedBackups != 0 || len(r.Warnings) == 0 {
		held.Close()
		t.Fatalf("%+v", r)
	}
	rs, _ := UpdateBackups(c)
	pending := false
	for _, v := range rs {
		if v.ID == old.ID {
			pending = v.CleanupPending
		}
	}
	if !pending {
		held.Close()
		t.Fatal("partial deletion was advertised as a healthy backup")
	}
	if e = RollbackUpdate(context.Background(), c.Data, old.ID, quietUpdate); e == nil || !strings.Contains(e.Error(), "being removed") {
		held.Close()
		t.Fatalf("partial backup selected: %v", e)
	}
	held.Close()
	r = cleanupSuccessfulUpdate(context.Background(), c, keep, "", quietUpdate)
	if r.RemovedBackups != 1 || len(r.Warnings) != 0 {
		t.Fatalf("retry %+v", r)
	}
	requireExists(t, old.Path, false)
	if _, e = verifyUpdateBackup(c, keep); e != nil {
		t.Fatal(e)
	}
}

func TestUpdateCleanupOnlyRetiresLegacyBackupProvenByCompletedJournal(t *testing.T) {
	c, _ := updateFixture(t)
	unknown := cleanupCopy(t, c, "")
	proven := cleanupCopy(t, c, "")
	keep := cleanupCopy(t, c, "update")
	cleanupReady(t, c, keep, &proven)
	r := cleanupSuccessfulUpdate(context.Background(), c, keep, "", quietUpdate)
	if r.RemovedBackups != 1 || len(r.Warnings) != 0 {
		t.Fatalf("%+v", r)
	}
	requireExists(t, unknown.Path, true)
	requireExists(t, proven.Path, false)
}

func TestUpdateCleanupPreservesExtraFilesAndForeignPaths(t *testing.T) {
	c, _ := updateFixture(t)
	old := cleanupCopy(t, c, "update")
	extra := filepath.Join(old.Path, "extra-personal-file.txt")
	os.WriteFile(extra, []byte("never delete unrecorded files"), 0600)
	keep := cleanupCopy(t, c, "update")
	cleanupReady(t, c, keep, nil)
	r := cleanupSuccessfulUpdate(context.Background(), c, keep, "", quietUpdate)
	if r.RemovedBackups != 0 || len(r.Warnings) == 0 {
		t.Fatalf("%+v", r)
	}
	requireExists(t, extra, true)
	forged := old
	forged.Path = filepath.Join(c.Files, old.ID)
	os.MkdirAll(forged.Path, 0700)
	os.WriteFile(filepath.Join(forged.Path, "document.txt"), []byte("ordinary file"), 0600)
	if e := cleanupBackupBoundary(c, forged, keep); e == nil {
		t.Fatal("accepted cleanup inside ordinary files")
	}
	requireExists(t, filepath.Join(forged.Path, "document.txt"), true)
}

func TestUpdateCleanupRefusesJunctionsWithoutTouchingTarget(t *testing.T) {
	c, _ := updateFixture(t)
	guard, e := lockUpdate(c)
	if e != nil {
		t.Fatal(e)
	}
	defer guard.Close()
	root, e := newUpdateStage(c)
	if e != nil {
		t.Fatal(e)
	}
	outside := t.TempDir()
	ordinary := filepath.Join(outside, "ordinary.txt")
	os.WriteFile(ordinary, []byte("must survive"), 0600)
	link := filepath.Join(root, "link")
	if e = execJunction(link, outside); e != nil {
		t.Skip(e)
	}
	defer os.Remove(link)
	if _, e = removeCleanupTree(context.Background(), root, nil); e == nil {
		t.Fatal("cleanup followed a junction")
	}
	requireExists(t, ordinary, true)
	requireExists(t, root, true)
}

func TestUpdateCleanupTransactionRunsOnlyAfterHealthCheckedCommit(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "failed startup"}[fail], func(t *testing.T) {
			c, s := updateFixture(t)
			old := cleanupCopy(t, c, "update")
			target, stage := fixtureTarget(t, c)
			starts := 0
			ops := updateOperations{stop: func(context.Context, Config) error { return nil }, start: func(context.Context, Config, string) error {
				starts++
				if fail && starts == 1 {
					return errors.New("synthetic startup failure")
				}
				return nil
			}}
			e := applyStagedWith(context.Background(), c, s, target, stage, "", quietUpdate, ops)
			if fail {
				if e == nil {
					t.Fatal("failed update succeeded")
				}
				requireExists(t, old.Path, true)
				if r, _ := LoadUpdateCleanup(c); r != nil {
					t.Fatal("cleanup ran on rollback")
				}
			} else {
				if e != nil {
					t.Fatal(e)
				}
				requireExists(t, old.Path, false)
				rs, _ := UpdateBackups(c)
				if len(rs) != 1 || rs[0].Purpose != "update" {
					t.Fatalf("%+v", rs)
				}
				if _, e = verifyUpdateBackup(c, rs[0]); e != nil {
					t.Fatal(e)
				}
			}
		})
	}
}

func TestUpdateCleanupStartupCompletesLegacyHandoffAndIsIdempotent(t *testing.T) {
	c, _ := updateFixture(t)
	old := cleanupCopy(t, c, "update")
	keep := cleanupCopy(t, c, "")
	cleanupReady(t, c, keep, nil)
	startupUpdateCleanup(context.Background(), c)
	requireExists(t, old.Path, false)
	requireExists(t, keep.Path, true)
	r, e := LoadUpdateCleanup(c)
	if e != nil || r == nil || r.RemovedBackups != 1 || len(r.Warnings) != 0 {
		t.Fatalf("%+v %v", r, e)
	}
	before := r.Finished
	startupUpdateCleanup(context.Background(), c)
	r, _ = LoadUpdateCleanup(c)
	if !r.Finished.Equal(before) {
		t.Fatal("startup repeated a completed cleanup")
	}
}
