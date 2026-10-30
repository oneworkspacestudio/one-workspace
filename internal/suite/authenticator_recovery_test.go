package suite

import (
	"context"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNewAppRollbackCanReapplyWithoutUnmanagedProgramCollision(t *testing.T) {
	c, s := updateFixture(t)
	target, stage := fixtureTarget(t, c)
	target.Apps = append(append([]catalog.App{}, target.Apps...), catalog.App{ID: "authenticator", Name: "One Authenticator", Version: "0.1.0-preview.2", Kind: "process", Package: "authenticator", Executable: "one-authenticator.exe", Port: 25137, Route: "/authenticator", Source: "https://github.com/oneworkspacestudio/one-authenticator/releases/tag/v0.1.0-preview.2"})
	target.Packages["authenticator"] = updatePackage("one-authenticator.exe", "introduced-program")
	if err := os.WriteFile(filepath.Join(stage, "one-authenticator.exe"), []byte("introduced-program"), 0600); err != nil {
		t.Fatal(err)
	}
	keep := filepath.Join(c.Programs, "unmanaged-document.txt")
	os.WriteFile(keep, []byte("keep-unrelated"), 0600)
	ops := updateOperations{stop: func(context.Context, Config) error { return nil }, start: func(context.Context, Config, string) error { return nil }, connect: func(context.Context, Config, *Secrets) error { return nil }}
	if err := applyStagedWith(context.Background(), c, s, target, stage, "", quietUpdate, ops); err != nil {
		t.Fatal(err)
	}
	next, _, err := Load(c.Data)
	if err != nil {
		t.Fatal(err)
	}
	var journal UpdateJournal
	if err = platform.LoadJSON(updateJournalPath(c), &journal); err != nil {
		t.Fatal(err)
	}
	journal.Phase = "rollback-requested"
	if err = recoverUpdateWith(context.Background(), next, &journal, quietUpdate, ops); err != nil {
		t.Fatal(err)
	}
	if _, err = os.Stat(filepath.Join(c.Programs, "one-authenticator.exe")); !os.IsNotExist(err) {
		t.Fatal("new owned program survived rollback and blocks reapply", err)
	}
	if b, _ := os.ReadFile(keep); string(b) != "keep-unrelated" {
		t.Fatal("unrelated program-folder file changed")
	}
	if journal.Recovery == nil {
		t.Fatal("current state archive missing")
	}
	if b, _ := os.ReadFile(filepath.Join(journal.Recovery.Path, "programs/one-authenticator.exe")); string(b) != "introduced-program" {
		t.Fatal("introduced bytes not retained in recovery archive")
	}
	restored, secret, err := Load(c.Data)
	if err != nil {
		t.Fatal(err)
	}
	stage = t.TempDir()
	os.WriteFile(filepath.Join(stage, "Workspace.exe"), []byte("manager-updated"), 0600)
	os.WriteFile(filepath.Join(stage, "one-authenticator.exe"), []byte("introduced-program"), 0600)
	if err = applyStagedWith(context.Background(), restored, secret, target, stage, "", quietUpdate, ops); err != nil {
		t.Fatal("supported reapplication failed", err)
	}
	if b, _ := os.ReadFile(filepath.Join(c.Files, "Document.txt")); string(b) != "ordinary original" {
		t.Fatal("ordinary document changed")
	}
}

func TestIntroducedProgramChangedAfterSnapshotIsKept(t *testing.T) {
	c, _ := updateFixture(t)
	path := filepath.Join(c.Programs, "one-authenticator.exe")
	if err := os.WriteFile(path, []byte("newer-unarchived-bytes"), 0600); err != nil {
		t.Fatal(err)
	}
	pkg := updatePackage("one-authenticator.exe", "verified-introduced-program")
	before := UpdateBackup{Config: c, Metadata: &UpdateMetadataPolicy{Programs: []string{"one-authenticator.exe"}}}
	current := UpdateBackup{Files: []BackupFile{{Area: "programs", Path: pkg.File, Size: pkg.Size, SHA256: pkg.SHA256}}}
	if err := clearIntroducedRecoveryPrograms(c, before, current); err == nil {
		t.Fatal("unarchived changes were removed")
	}
	if b, _ := os.ReadFile(path); string(b) != "newer-unarchived-bytes" {
		t.Fatal("newer file changed")
	}
}

func TestIntroducedProgramWithUnarchivedNamedStreamIsKept(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("NTFS alternate stream contract")
	}
	c, _ := updateFixture(t)
	path := filepath.Join(c.Programs, "one-authenticator.exe")
	if err := os.WriteFile(path, []byte("verified-introduced-program"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path+":retained-note", []byte("new stream data"), 0600); err != nil {
		t.Fatal(err)
	}
	pkg := updatePackage("one-authenticator.exe", "verified-introduced-program")
	before := UpdateBackup{Config: c, Metadata: &UpdateMetadataPolicy{Programs: []string{"one-authenticator.exe"}}}
	current := UpdateBackup{Files: []BackupFile{{Area: "programs", Path: pkg.File, Size: pkg.Size, SHA256: pkg.SHA256}}}
	if err := clearIntroducedRecoveryPrograms(c, before, current); err == nil {
		t.Fatal("unarchived stream was removed")
	}
	if b, _ := os.ReadFile(path + ":retained-note"); string(b) != "new stream data" {
		t.Fatal("newer stream changed")
	}
}
