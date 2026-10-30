package suite

import (
	"encoding/json"
	"errors"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// Ticket membership is authoritative in My Tickets. The raw Drive source is disabled
// by default, preventing a broad Drive account from bypassing project permissions.
func ensureTicketsDriveSource(c Config) error {
	if _, ok := c.Catalog.Find("tickets"); !ok {
		return nil
	}
	if _, ok := c.Catalog.Find("drive"); !ok {
		return nil
	}
	root := filepath.Join(c.Data, "tickets", "core", "attachments")
	for current := root; current != filepath.Clean(c.Data); current = filepath.Dir(current) {
		if err := platform.RejectReparse(current); err != nil && !os.IsNotExist(err) {
			return err
		}
		if filepath.Dir(current) == current {
			return errors.New("Tickets storage is outside private data")
		}
	}
	if err := prepareTicketsData(filepath.Join(c.Data, "tickets")); err != nil {
		return err
	}
	if err := os.MkdirAll(root, 0700); err != nil {
		return err
	}
	file := filepath.Join(c.Data, "drive", "config.yaml")
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	var cfg map[string]any
	if json.Unmarshal(b, &cfg) != nil {
		return errors.New("unsupported Drive configuration")
	}
	server, ok := cfg["server"].(map[string]any)
	if !ok {
		return errors.New("Drive server missing")
	}
	sources, ok := server["sources"].([]any)
	if !ok {
		return errors.New("Drive sources missing")
	}
	policy := map[string]any{"defaultEnabled": false, "defaultPermissions": map[string]bool{"view": false, "download": false, "create": false, "modify": false, "delete": false}}
	for _, v := range sources {
		source, ok := v.(map[string]any)
		if !ok {
			return errors.New("invalid Drive source")
		}
		sourcePath, _ := source["path"].(string)
		if source["name"] != "My Tickets" && (sourcePath == "" || containedPath(sourcePath, root) || containedPath(root, sourcePath)) {
			return errors.New("Drive source overlaps private Tickets attachments")
		}
		if source["name"] == "My Tickets" {
			path, _ := source["path"].(string)
			if !strings.EqualFold(filepath.Clean(path), filepath.Clean(root)) {
				return errors.New("My Tickets Drive source path conflict")
			}
			source["config"] = policy
			return platform.SaveJSON(file, cfg)
		}
	}
	server["sources"] = append(sources, map[string]any{"name": "My Tickets", "path": root, "config": policy})
	return platform.SaveJSON(file, cfg)
}

var ticketsAgentToken = regexp.MustCompile(`^Bearer tk_[A-Za-z0-9_-]{20,256}$`)

// Presence only selects the native verifier; the app still validates token scopes
// and current project membership. This never permits session-less browser access.
func nativeTicketsRequest(r *http.Request) bool {
	if !strings.HasPrefix(r.URL.Path, "/tickets/api/v2/") || path.Clean(r.URL.Path) != r.URL.Path || !ticketsAgentToken.MatchString(r.Header.Get("Authorization")) {
		return false
	}
	switch r.Method {
	case http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	}
	return false
}

// prepareTicketsData gives the SDK ownership before Drive creates its empty source.
// Preview105 could create only core/attachments without the SDK marker. Adopt that
// precise empty layout; never adopt files, foreign markers or reparse points.
func prepareTicketsData(root string) error {
	if _, err := os.Lstat(filepath.Join(root, "app.json")); err == nil {
		return apphost.PrepareData(root, "tickets")
	} else if !os.IsNotExist(err) {
		return err
	}
	if _, err := os.Lstat(root); os.IsNotExist(err) {
		return apphost.PrepareData(root, "tickets")
	} else if err != nil {
		return err
	}
	if err := filepath.WalkDir(root, func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if err = platform.RejectReparse(p); err != nil {
			return err
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		if !entry.IsDir() || rel != "." && rel != "core" && rel != filepath.Join("core", "attachments") {
			return errors.New("existing Tickets storage needs explicit migration; no files were changed")
		}
		return nil
	}); err != nil {
		return err
	}
	for _, rel := range []string{filepath.Join("core", "attachments"), "core"} {
		if err := os.Remove(filepath.Join(root, rel)); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return apphost.PrepareData(root, "tickets")
}
