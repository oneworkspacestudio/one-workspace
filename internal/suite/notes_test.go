package suite

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNotesSettingsAreSeparateAndNeverReset(t *testing.T) {
	c := Config{Data: t.TempDir(), Origin: "http://localhost:18890", Owner: "owner-subject"}
	if err := ensureNotesConfig(c); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(c.Data, "notes-settings.json")
	if _, err := os.Stat(filepath.Join(c.Data, "notes", "notes-settings.json")); !os.IsNotExist(err) {
		t.Fatal("configuration polluted dedicated Notes data")
	}
	if err := os.WriteFile(path, []byte(`{"preserve":true}`), 0600); err != nil {
		t.Fatal(err)
	}
	if err := ensureNotesConfig(c); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(path)
	if string(b) != `{"preserve":true}` {
		t.Fatal("existing Notes configuration reset")
	}
}
