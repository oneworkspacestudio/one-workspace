package suite

import (
	"context"
	"encoding/json"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"os"
	"path/filepath"
	"testing"
)

func TestMyAINativeContract(t *testing.T) {
	c := NewConfig(Request{Port: 29000, Mode: "local", Programs: "C:/Synthetic/Apps", Data: "C:/Synthetic/Private", Files: "C:/Synthetic/Files"}, catalog.Catalog{Apps: []catalog.App{{ID: "myai", Kind: "process", Port: 9019, Route: "/myai"}}})
	a, ok := c.Catalog.Find("myai")
	if !ok || a.Port != 29029 || !foundationApp(a.ID) || appWorkingDir(c, "myai") == c.Data {
		t.Fatal(a)
	}
	if containedPath(c.Files, myaiRoot(c)) || containedPath(myaiRoot(c), c.Files) {
		t.Fatal("private source overlaps shared files")
	}
	found := false
	for _, s := range suiteSignIns {
		if s.id == "myai" {
			found = s.probe == "/api/me" && s.route == "/myai"
		}
	}
	if !found {
		t.Fatal("missing seamless MyAccounts session contract")
	}
}
func TestMyAIDefaultDenySourceAndNoBroadOverlap(t *testing.T) {
	root := t.TempDir()
	c := Config{Data: filepath.Join(root, "Private"), Files: filepath.Join(root, "Files"), Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "myai"}}}}
	os.MkdirAll(filepath.Join(c.Data, "drive"), 0700)
	config := filepath.Join(c.Data, "drive", "config.yaml")
	save := func(path string) {
		b, _ := json.Marshal(map[string]any{"server": map[string]any{"sources": []any{map[string]any{"name": "My Drive", "path": path}}}})
		os.WriteFile(config, b, 0600)
	}
	save(c.Files)
	if e := ensureMyAIDriveSource(c); e != nil {
		t.Fatal(e)
	}
	var cfg map[string]any
	b, _ := os.ReadFile(config)
	json.Unmarshal(b, &cfg)
	sources := cfg["server"].(map[string]any)["sources"].([]any)
	scope := sources[1].(map[string]any)["config"].(map[string]any)
	if scope["defaultEnabled"] != false {
		t.Fatal("enabled by default")
	}
	for k, v := range scope["defaultPermissions"].(map[string]any) {
		if v != false {
			t.Fatal("default grant", k)
		}
	}
	rt := Runtime{Config: c}
	if e := rt.verifyMyAIStorage(context.Background(), true); e != nil {
		t.Fatal(e)
	}
	if e := rt.verifyMyAIStorage(context.Background(), false); e != nil {
		t.Fatal(e)
	}
	save(root)
	if ensureMyAIDriveSource(c) == nil {
		t.Fatal("broad source permitted")
	}
	if myaiScope("alice", "issuer") == myaiScope("bob", "issuer") || myaiScope("alice", "issuer") == myaiScope("alice", "other") {
		t.Fatal("identity collision")
	}
}
