package suite

import (
	"encoding/json"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"os"
	"path/filepath"
	"testing"
)

func TestClipsStorageOutsideSharedRoot(t *testing.T) {
	base := t.TempDir()
	c := Config{Data: filepath.Join(base, "Private"), Files: filepath.Join(base, "Files")}
	rt := &Runtime{Config: c}
	if e := rt.verifyClipsStorage(t.Context(), true); e != nil {
		t.Fatal(e)
	}
	if containedPath(c.Files, clipsRoot(c)) {
		t.Fatal("shared root")
	}
	if clipsScope("a", "issuer") == clipsScope("b", "issuer") {
		t.Fatal("identity collision")
	}

}
func TestClipsSourceRejectsBroadPrivateAlias(t *testing.T) {
	base := t.TempDir()
	c := Config{Data: filepath.Join(base, "Private"), Files: filepath.Join(base, "Files"), Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "clips"}}}}
	os.MkdirAll(filepath.Join(c.Data, "drive"), 0700)
	file := filepath.Join(c.Data, "drive", "config.yaml")
	save := func(p string) {
		b, _ := json.Marshal(map[string]any{"server": map[string]any{"sources": []any{map[string]any{"name": "My Drive", "path": p}}}})
		os.WriteFile(file, b, 0600)
	}
	save(c.Files)
	if e := ensureClipsDriveSource(c); e != nil {
		t.Fatal(e)
	}
	save(base)
	if ensureClipsDriveSource(c) == nil {
		t.Fatal("broad source alias accepted")
	}
}
