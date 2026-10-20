package suite

import (
	"encoding/json"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"os"
	"path/filepath"
	"testing"
)

func TestSocialStorageOutsideSharedRoot(t *testing.T) {
	base := t.TempDir()
	c := Config{Data: filepath.Join(base, "Private"), Files: filepath.Join(base, "Files")}
	rt := &Runtime{Config: c}
	if e := rt.verifySocialStorage(t.Context(), true); e != nil {
		t.Fatal(e)
	}
	if containedPath(c.Files, socialRoot(c)) {
		t.Fatal("shared root")
	}
	if socialScope("a", "issuer") == socialScope("b", "issuer") {
		t.Fatal("identity collision")
	}
	legacy := filepath.Join(c.Files, "My Photos", "My Social")
	os.MkdirAll(legacy, 0700)
	os.WriteFile(filepath.Join(legacy, "original.jpg"), []byte("fixture"), 0600)
	if rt.verifySocialStorage(t.Context(), true) == nil {
		t.Fatal("legacy plaintext silently migrated")
	}
}
func TestSocialSourceRejectsBroadPrivateAlias(t *testing.T) {
	base := t.TempDir()
	c := Config{Data: filepath.Join(base, "Private"), Files: filepath.Join(base, "Files"), Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "social"}}}}
	os.MkdirAll(filepath.Join(c.Data, "drive"), 0700)
	file := filepath.Join(c.Data, "drive", "config.yaml")
	save := func(p string) {
		b, _ := json.Marshal(map[string]any{"server": map[string]any{"sources": []any{map[string]any{"name": "My Drive", "path": p}}}})
		os.WriteFile(file, b, 0600)
	}
	save(c.Files)
	if e := ensureSocialDriveSource(c); e != nil {
		t.Fatal(e)
	}
	save(base)
	if ensureSocialDriveSource(c) == nil {
		t.Fatal("broad source alias accepted")
	}
}
