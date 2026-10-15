package suite

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMusicTempCleanupIsOwnedAndScoped(t *testing.T) {
	c := Config{Data: t.TempDir()}
	if err := prepareMusicTemp(c); err != nil {
		t.Fatal(err)
	}
	root := musicRuntimeTemp(c)
	old := filepath.Join(root, "_MEI-test")
	os.MkdirAll(old, 0700)
	os.WriteFile(filepath.Join(old, "runtime.bin"), []byte("runtime"), 0600)
	keep := filepath.Join(root, "keep.txt")
	os.WriteFile(keep, []byte("preserve"), 0600)
	if err := prepareMusicTemp(c); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Fatal("stale extraction retained")
	}
	if _, err := os.Stat(keep); err != nil {
		t.Fatal("unrelated file removed", err)
	}
	other := Config{Data: t.TempDir()}
	os.MkdirAll(filepath.Join(musicRuntimeTemp(other), "_MEI-unknown"), 0700)
	if err := prepareMusicTemp(other); err == nil {
		t.Fatal("unowned folder accepted")
	}
}
