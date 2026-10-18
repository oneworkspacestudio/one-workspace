package suite

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

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
