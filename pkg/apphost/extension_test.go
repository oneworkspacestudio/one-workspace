package apphost

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEmbeddedUIHandlerRetainsHostBoundaries(t *testing.T) {
	h := testHost(t)
	h.UIHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("embedded custom UI: " + r.URL.Path)) })
	response := call(h, "GET", "/browse/WORK-1", nil, nil, "")
	if response.Code != 200 || !strings.Contains(response.Body.String(), "embedded custom UI") {
		t.Fatal(response.Code, response.Body.String())
	}
	if response.Header().Get("Content-Security-Policy") == "" {
		t.Fatal("security headers missing")
	}
	if response = call(h, "GET", "/healthz", nil, nil, ""); response.Code != 200 || !strings.Contains(response.Body.String(), `"status":"ok"`) {
		t.Fatal(response.Code, response.Body.String())
	}
	if response = call(h, "GET", "/api/me", nil, nil, ""); response.Code != 401 {
		t.Fatal("custom UI bypassed identity")
	}
}

func TestAdditionalSnapshotRequiresAppRestoreAllowlist(t *testing.T) {
	h := testHost(t)
	h.SnapshotFiles = map[string]func(string) error{"core/fixture.db": func(path string) error { return os.WriteFile(path, []byte("consistent snapshot fixture"), 0600) }}
	var archive bytes.Buffer
	if err := h.Snapshot(&archive, "synthetic-backup-password"); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(t.TempDir(), "backup.owbackup")
	if err := os.WriteFile(file, archive.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
	if err := Restore(filepath.Join(t.TempDir(), "denied"), h.Config.ID, file, "synthetic-backup-password"); err == nil {
		t.Fatal("default restore accepted an unknown application file")
	}
	corrupt := filepath.Join(t.TempDir(), "invalid-sqlite")
	if err := RestoreWithDatabaseFiles(corrupt, h.Config.ID, file, "synthetic-backup-password", []string{"core/fixture.db"}, []string{"core/fixture.db"}); err == nil {
		t.Fatal("restore accepted an invalid declared SQLite snapshot")
	}
	if _, err := os.Stat(filepath.Join(corrupt, "state.sqlite")); !os.IsNotExist(err) {
		t.Fatal("failed integrity verification committed restored state", err)
	}
	if err := RestoreWithFiles(filepath.Join(t.TempDir(), "incomplete"), h.Config.ID, file, "synthetic-backup-password", []string{"core/fixture.db", "core/missing.key"}); err == nil {
		t.Fatal("restore accepted a missing required file")
	}
	dest := filepath.Join(t.TempDir(), "restored")
	if err := RestoreWithFiles(dest, h.Config.ID, file, "synthetic-backup-password", []string{"core/fixture.db"}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dest, "core/fixture.db"))
	if err != nil || string(data) != "consistent snapshot fixture" {
		t.Fatal(string(data), err)
	}
	h.SnapshotFiles = map[string]func(string) error{"../escape": func(path string) error { return os.WriteFile(path, []byte("bad"), 0600) }}
	if err := h.Snapshot(&bytes.Buffer{}, "synthetic-backup-password"); err == nil {
		t.Fatal("snapshot accepted traversal")
	}
}
