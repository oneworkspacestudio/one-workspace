package suite

import (
	"encoding/json"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMyTubeManagedSettingsPreserveUserConfiguration(t *testing.T) {
	cat := catalog.Catalog{Apps: []catalog.App{{ID: "drive", Route: "/drive", Port: 18002}, {ID: "my-tube", Route: "/my-tube", Port: 18013}}}
	c := Config{Data: t.TempDir(), Origin: "http://localhost:18000", Owner: "owner-subject", Catalog: cat}
	if e := ensureMyTubeConfig(c); e != nil {
		t.Fatal(e)
	}
	p := filepath.Join(c.Data, "my-tube-settings.json")
	b, e := os.ReadFile(p)
	if e != nil {
		t.Fatal(e)
	}
	var settings map[string]any
	if e = json.Unmarshal(b, &settings); e != nil {
		t.Fatal(e)
	}
	if settings["mode"] != "managed" || settings["public_url"] != "http://localhost:18000/my-tube" || settings["drive_url"] != "http://127.0.0.1:18002/drive" || settings["listen"] != "127.0.0.1:18013" {
		t.Fatal(settings)
	}
	if settings["admin_subjects"].([]any)[0] != "owner-subject" || settings["public_access"] != false {
		t.Fatal(settings)
	}
	if strings.Contains(string(b), "client_secret") {
		t.Fatal("Credentials written to ordinary settings")
	}
	if _, e = os.Stat(filepath.Join(c.Data, "my-tube", "my-tube-settings.json")); !os.IsNotExist(e) {
		t.Fatal("Dedicated data folder polluted")
	}
	custom := []byte(`{"mode":"managed","workers":2}`)
	if e = os.WriteFile(p, custom, 0600); e != nil {
		t.Fatal(e)
	}
	if e = ensureMyTubeConfig(c); e != nil {
		t.Fatal(e)
	}
	b, _ = os.ReadFile(p)
	if string(b) != string(custom) {
		t.Fatal("Existing settings replaced")
	}
}
func TestMyTubePortsFollowSelectedWorkspacePort(t *testing.T) {
	cat := catalog.Catalog{Apps: []catalog.App{{ID: "drive", Port: 8992}, {ID: "notes", Port: 9002}, {ID: "my-tube", Port: 9003}}}
	c := NewConfig(Request{Port: 19000, Mode: "local"}, cat)
	tube, _ := c.Catalog.Find("my-tube")
	notes, _ := c.Catalog.Find("notes")
	if tube.Port != 19013 || notes.Port != 19012 {
		t.Fatal("Ports were not remapped")
	}
}
