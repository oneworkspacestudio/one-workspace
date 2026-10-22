package suite

import (
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSetupSessionDetectsLiveLockAndInterruptedCheckpoint(t *testing.T) {
	root := filepath.Clean(t.TempDir())
	r := Request{Folder: root, Password: "synthetic-sensitive-value", Confirm: "synthetic-sensitive-value"}
	if e := r.ResolveStorage(); e != nil {
		t.Fatal(e)
	}
	if e := prepareSetupSession(r); e != nil {
		t.Fatal(e)
	}
	if e := saveSetupCheckpoint(r, Config{}, "extract"); e != nil {
		t.Fatal(e)
	}
	p := InstallProgress(false)
	p.InstallPhase("Verifying app file 1 of 3")
	if e := SaveSetupSession(r, p.Snapshot()); e != nil {
		t.Fatal(e)
	}
	lock, e := platform.Lock(setupLockPath(r))
	if e != nil {
		t.Fatal(e)
	}
	status, active, e := ReadSetupSession(r)
	if e != nil || !active || status.State != "installing" {
		t.Fatal(status, active, e)
	}
	lock.Close()
	status, active, e = ReadSetupSession(r)
	if e != nil || active || status.State != "interrupted" || status.Phase == "" {
		t.Fatal(status, active, e)
	}
	for _, name := range []string{"continue.dpapi", "status.json"} {
		raw, e := os.ReadFile(filepath.Join(setupSessionDir(r.Data), name))
		if e != nil || strings.Contains(string(raw), r.Password) {
			t.Fatal("sensitive setup value exposed", name, e)
		}
	}
}
