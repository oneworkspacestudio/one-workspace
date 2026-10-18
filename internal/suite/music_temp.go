package suite

import (
	"errors"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"os"
	"path/filepath"
	"strings"
)

func musicRuntimeTemp(c Config) string { return filepath.Join(c.Data, "music-runtime-temp") }

// PyInstaller's extraction parent cannot clean up after the Windows job closes
// all children. Own a dedicated private location and reclaim only its stale
// extraction directories before starting the next Music process.
func prepareMusicTemp(c Config) error {
	root := musicRuntimeTemp(c)
	if err := os.MkdirAll(root, 0700); err != nil {
		return err
	}
	info, err := os.Lstat(root)
	if err != nil {
		return err
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return errors.New("Music temporary storage must be an ordinary private folder")
	}
	markerPath := filepath.Join(root, ".workspace-owner.json")
	var marker struct {
		Schema int    `json:"schema"`
		Data   string `json:"data"`
	}
	err = platform.LoadJSON(markerPath, &marker)
	if os.IsNotExist(err) {
		entries, e := os.ReadDir(root)
		if e != nil {
			return e
		}
		if len(entries) != 0 {
			return errors.New("Music temporary folder is not empty and has no Workspace ownership marker")
		}
		marker.Schema = 1
		marker.Data = filepath.Clean(c.Data)
		return platform.SaveJSON(markerPath, marker)
	}
	if err != nil {
		return err
	}
	if marker.Schema != 1 || !strings.EqualFold(filepath.Clean(marker.Data), filepath.Clean(c.Data)) {
		return errors.New("Music temporary folder belongs to a different installation")
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "_MEI") {
			continue
		}
		path, e := platform.Inside(root, entry.Name())
		if e != nil {
			return e
		}
		if e = os.RemoveAll(path); e != nil {
			return e
		}
	}
	return nil
}
