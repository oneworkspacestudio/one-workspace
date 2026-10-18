package suite

import (
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"os"
	"path/filepath"
)

func ensureNotesConfig(c Config) error {
	path := filepath.Join(c.Data, "notes-settings.json")
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	return platform.SaveJSON(path, map[string]any{"mode": "managed", "data_dir": filepath.Join(c.Data, "notes"), "public_url": c.Origin + "/notes", "admin_subjects": []string{c.Owner}, "session_hours": 8, "upload_mb": 25, "quota_mb": 1024})
}
