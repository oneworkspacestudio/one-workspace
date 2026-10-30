package suite

import (
	"encoding/json"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTicketsUsesManagedNativeContract(t *testing.T) {
	c := Config{Data: filepath.Join(t.TempDir(), "private")}
	if !foundationApp("tickets") || !knownUpdateStorageProfile("tickets") || appWorkingDir(c, "tickets") != filepath.Join(c.Data, "tickets") {
		t.Fatal("missing native runtime or storage profile")
	}
	if paths := appUpdatePrivatePaths("tickets"); len(paths) != 1 || paths[0] != "tickets" {
		t.Fatal(paths)
	}
	if !preservedUpdateData("tickets/core/attachments/WORK/WORK-1/file") || !preservedUpdateData("tickets/tmp/partial") || preservedUpdateData("tickets/core/mytickets.db") {
		t.Fatal("update recovery mixed attachment data and schema metadata")
	}
	if !catalog.DefaultInstall("tickets") || catalog.CategoryFor("tickets").ID != "general" {
		t.Fatal("Tickets should be a default general app")
	}
	found := false
	for _, p := range suiteSignIns {
		if p.id == "tickets" {
			found = true
		}
	}
	if !found {
		t.Fatal("Tickets lacks native seamless login")
	}
}
func TestTicketsDriveSourceIsIdempotentAndPrivate(t *testing.T) {
	c := Config{Data: t.TempDir(), Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "tickets"}, {ID: "drive"}}}}
	if err := os.MkdirAll(filepath.Join(c.Data, "drive"), 0700); err != nil {
		t.Fatal(err)
	}
	cfg := map[string]any{"server": map[string]any{"sources": []any{map[string]any{"name": "My Drive", "path": filepath.Join(t.TempDir(), "files")}}}}
	b, _ := json.Marshal(cfg)
	file := filepath.Join(c.Data, "drive", "config.yaml")
	if err := os.WriteFile(file, b, 0600); err != nil {
		t.Fatal(err)
	}
	for n := 0; n < 2; n++ {
		if err := ensureTicketsDriveSource(c); err != nil {
			t.Fatal(err)
		}
	}
	b, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(b, &cfg); err != nil {
		t.Fatal(err)
	}
	sources := cfg["server"].(map[string]any)["sources"].([]any)
	if len(sources) != 2 {
		t.Fatal("duplicated source")
	}
	source := sources[1].(map[string]any)
	if source["name"] != "My Tickets" || source["path"] != filepath.Join(c.Data, "tickets", "core", "attachments") || source["config"].(map[string]any)["defaultEnabled"] != false {
		t.Fatal("unsafe Drive source")
	}
}

func TestTicketsAgentRoutingRequiresNarrowPathAndCredentialShape(t *testing.T) {
	for _, test := range []struct {
		method, path, auth string
		want               bool
	}{
		{"POST", "/tickets/api/v2/mcp", "Bearer tk_" + strings.Repeat("a", 48), true},
		{"POST", "/tickets/api/v2/tickets/projects/1/git/webhook", "Bearer tk_" + strings.Repeat("a", 48), true},
		{"GET", "/tickets/", "Bearer tk_" + strings.Repeat("a", 48), false},
		{"GET", "/tickets/api/v2/tickets/issues/WORK-1", "", false},
		{"GET", "/tickets/api/v2/../../_integration/monitor", "Bearer tk_" + strings.Repeat("a", 48), false},
		{"POST", "/_integration/tickets/drive", "Bearer tk_" + strings.Repeat("a", 48), false},
		{"CONNECT", "/tickets/api/v2/mcp", "Bearer tk_" + strings.Repeat("a", 48), false},
	} {
		r := httptest.NewRequest(test.method, test.path, nil)
		r.Header.Set("Authorization", test.auth)
		if got := nativeTicketsRequest(r); got != test.want {
			t.Errorf("%s %s: %v", test.method, test.path, got)
		}
	}
}
