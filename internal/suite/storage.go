package suite

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

// Program payloads are replaceable; working data and ordinary user files must
// never be children of that directory (or contain it). Enforce this on startup
// and before update recovery is planned, including legacy/headless configs.
func validateProgramDataSeparation(c Config) error {
	roots := []string{c.Programs, c.Data, c.Files}
	for i, path := range roots {
		if !filepath.IsAbs(path) || filepath.Dir(filepath.Clean(path)) == filepath.Clean(path) {
			return errors.New("Workspace programs, private data and files require dedicated absolute folders")
		}
		for j, other := range roots {
			if i == j {
				continue
			}
			a, b := strings.ToLower(filepath.Clean(path)), strings.ToLower(filepath.Clean(other))
			if a == b || strings.HasPrefix(a, b+string(os.PathSeparator)) {
				return errors.New("Workspace program, private-data and user-file folders must be separate")
			}
		}
	}
	dataFlags := map[string]bool{"--data": true, "-data": true, "--data-dir": true, "--config": true, "-c": true, "--media-dir": true, "--cache-dir": true, "--shared-data": true}
	vars := strings.NewReplacer("${DATA}", c.Data, "${FILES}", c.Files, "${PROGRAMS}", c.Programs)
	insidePrograms := func(base, raw string) bool {
		path := vars.Replace(raw)
		if !filepath.IsAbs(path) {
			path = filepath.Join(base, path)
		}
		rel, err := filepath.Rel(c.Programs, filepath.Clean(path))
		return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
	}
	for _, app := range c.Catalog.Apps {
		for i, arg := range app.Args {
			flag, assigned, equal := strings.Cut(arg, "=")
			if !dataFlags[flag] {
				continue
			}
			if !equal && i+1 < len(app.Args) {
				assigned = app.Args[i+1]
			}
			if insidePrograms(appWorkingDir(c, app.ID), assigned) {
				return fmt.Errorf("%s stores mutable data in the program folder; choose a separate data location before starting or updating it", app.Name)
			}
		}
		for key, value := range app.Env {
			key = strings.ToUpper(key)
			for _, suffix := range []string{"_DATA", "_DATA_DIR", "_DATA_ROOT", "_MEDIA_DIR", "_MEDIA_ROOT", "_MODELS", "_MODELS_DIR", "_CACHE_DIR", "_CONFIG"} {
				if strings.HasSuffix(key, suffix) && insidePrograms(appWorkingDir(c, app.ID), value) {
					return fmt.Errorf("%s has a mutable storage setting inside the program folder", app.Name)
				}
			}
		}
	}
	return nil
}

// Read live settings only while planning/startup, never during snapshot
// verification: recovery must also work when current settings are damaged.
func validateManagedStorageSettings(c Config, recovery bool) error {
	for _, id := range []string{"notes", "my-tube"} {
		app, installed := c.Catalog.Find(id)
		if !installed {
			continue
		}
		var settings struct {
			DataDir  string `json:"data_dir"`
			MediaDir string `json:"media_dir"`
		}
		if err := platform.LoadJSON(filepath.Join(c.Data, id+"-settings.json"), &settings); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("cannot check %s storage settings: %w", app.Name, err)
		}
		for _, path := range []string{settings.DataDir, settings.MediaDir} {
			if path == "" {
				continue
			}
			rel, err := filepath.Rel(c.Programs, filepath.Clean(path))
			if !filepath.IsAbs(path) || (err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))) {
				return fmt.Errorf("%s requires an absolute data location outside the Apps folder", app.Name)
			}
		}
		if !recovery {
			continue
		}
		if settings.DataDir != "" && !strings.EqualFold(filepath.Clean(settings.DataDir), filepath.Join(c.Data, id)) {
			return fmt.Errorf("%s uses a custom database location; migrate it to the managed private folder before using automatic update recovery", app.Name)
		}
		if settings.MediaDir != "" && !strings.EqualFold(filepath.Clean(settings.MediaDir), filepath.Join(c.Data, id, "media")) {
			rel, err := filepath.Rel(c.Files, filepath.Clean(settings.MediaDir))
			if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
				return fmt.Errorf("%s uses a custom media location; choose its managed media folder or a folder under Files before updating", app.Name)
			}
		}
	}
	return nil
}

// ResolveStorage keeps existing headless three-path requests compatible while
// the normal installer asks for a single parent folder.
func (r *Request) ResolveStorage() error {
	if r.Folder == "" {
		return nil
	}
	if !filepath.IsAbs(r.Folder) || filepath.Dir(filepath.Clean(r.Folder)) == filepath.Clean(r.Folder) {
		return errors.New("choose a dedicated Workspace folder, not a drive root")
	}
	r.Folder = filepath.Clean(r.Folder)
	r.Programs = filepath.Join(r.Folder, "Apps")
	r.Data = filepath.Join(r.Folder, "Private")
	r.Files = filepath.Join(r.Folder, "Files")
	return nil
}

// PreflightStorage does not create folders or change ACLs. The UI uses the same
// checks on selection, before Continue, and before accepting an installation.
func (r Request) PreflightStorage() (platform.DataVolume, error) {
	var empty platform.DataVolume
	if e := r.ResolveStorage(); e != nil {
		return empty, e
	}
	if r.Data == "" || r.Programs == "" || r.Files == "" {
		return empty, errors.New("choose your Workspace folder")
	}
	volume, e := platform.CheckDataVolume(r.Data)
	if e != nil {
		return volume, e
	}
	for _, path := range []string{r.Data, r.Programs} {
		info, e := os.Lstat(path)
		if e == nil && (!info.IsDir() || info.Mode()&os.ModeSymlink != 0) {
			return volume, fmt.Errorf("%s must be an ordinary folder", path)
		}
		if e != nil && !os.IsNotExist(e) {
			return volume, fmt.Errorf("cannot inspect the selected Workspace folder: %w", e)
		}
		entries, e := os.ReadDir(path)
		if e == nil && len(entries) > 0 {
			return volume, errors.New("this folder already contains Workspace apps or private data. Choose another parent folder for a new installation")
		}
		if e != nil && !os.IsNotExist(e) {
			return volume, fmt.Errorf("cannot inspect the selected Workspace folder: %w", e)
		}
	}
	if info, e := os.Stat(r.Files); e == nil && !info.IsDir() {
		return volume, errors.New("the Workspace files location is already used by a file")
	} else if e != nil && !os.IsNotExist(e) {
		return volume, e
	}
	return volume, nil
}
