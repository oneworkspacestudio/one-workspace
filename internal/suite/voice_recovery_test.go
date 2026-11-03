package suite

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"os"
	"path/filepath"
	"testing"
)

func TestNewVoiceRollbackRetainsCompletePrivateIdentity(t *testing.T) {
	before := UpdateBackup{Config: Config{Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "notes"}}}}}
	after := before
	after.Config.Catalog.Apps = append([]catalog.App{{ID: "voice"}}, before.Config.Catalog.Apps...)
	roots := introducedAppRecoveryRoots(before, after)
	if len(roots) != 1 || roots[0] != "voice" {
		t.Fatal("new app private scope missing", roots)
	}
	root := t.TempDir()
	for _, rel := range []string{"voice/app.json", "voice/app.db", "voice/projects/recording.wav", "notes/settings.json"} {
		p := filepath.Join(root, filepath.FromSlash(rel))
		if e := os.MkdirAll(filepath.Dir(p), 0700); e != nil {
			t.Fatal(e)
		}
		if e := os.WriteFile(p, []byte("synthetic"), 0600); e != nil {
			t.Fatal(e)
		}
	}
	if e := clearRecoveryMetadata(root, root, roots...); e != nil {
		t.Fatal(e)
	}
	for _, rel := range []string{"voice/app.json", "voice/app.db", "voice/projects/recording.wav"} {
		if _, e := os.Stat(filepath.Join(root, filepath.FromSlash(rel))); e != nil {
			t.Fatal("new-app identity/state removed", rel, e)
		}
	}
	if _, e := os.Stat(filepath.Join(root, "notes", "settings.json")); !os.IsNotExist(e) {
		t.Fatal("existing app metadata was not cleared for restoration")
	}
	before.Files = []BackupFile{{Area: "private", Path: "voice/app.json"}}
	if roots := introducedAppRecoveryRoots(before, after); len(roots) != 0 {
		t.Fatal("pre-existing metadata must restore from its snapshot", roots)
	}
	before.Files = nil
	before.Config.Catalog = after.Config.Catalog
	if roots := introducedAppRecoveryRoots(before, after); len(roots) != 0 {
		t.Fatal("existing apps must use normal schema recovery", roots)
	}
}
