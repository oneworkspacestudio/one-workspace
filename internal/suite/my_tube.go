package suite

import (
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"os"
	"path/filepath"
	"strings"
)

// Keep configuration outside the application's dedicated ACL-protected data
// directory. OIDC credentials continue to come from the manager's secret store.
func ensureMyTubeConfig(c Config) error {
	config := filepath.Join(c.Data, "my-tube-settings.json")
	if _, e := os.Stat(config); e == nil {
		return nil
	} else if !os.IsNotExist(e) {
		return e
	}
	tube, ok := c.Catalog.Find("my-tube")
	if !ok {
		return fmt.Errorf("My Tube is not in the installed catalog")
	}
	settings := map[string]any{"mode": "managed", "listen": fmt.Sprintf("127.0.0.1:%d", tube.Port), "data_dir": filepath.Join(c.Data, "my-tube"), "media_dir": filepath.Join(c.Data, "my-tube", "media"), "public_url": c.Origin + "/my-tube", "admin_subjects": []string{c.Owner}, "public_access": false, "public_uploads": false}
	// Use the direct loopback Drive mount: a server-to-server request cannot
	// satisfy the portal's separate browser session with only a Drive token.
	if drive, ok := c.Catalog.Find("drive"); ok {
		settings["drive_url"] = fmt.Sprintf("http://127.0.0.1:%d%s", drive.Port, strings.TrimRight(drive.Route, "/"))
	}
	return platform.SaveJSON(config, settings)
}
