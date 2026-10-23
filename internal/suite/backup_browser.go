package suite

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

type backupFolder struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
type backupFolderListing struct {
	Path      string         `json:"path"`
	Parent    string         `json:"parent"`
	Folders   []backupFolder `json:"folders"`
	Truncated bool           `json:"truncated"`
}

// Owner-only directory names, never file contents or a desktop OS picker.
// Limit work and reject reparse points before opening any directory.
func browseBackupFolders(path string) (backupFolderListing, error) {
	result := backupFolderListing{Folders: []backupFolder{}}
	path = strings.TrimSpace(path)
	if path == "" {
		for drive := 'A'; drive <= 'Z'; drive++ {
			root := string(drive) + ":\\"
			if info, e := os.Stat(root); e == nil && info.IsDir() {
				result.Folders = append(result.Folders, backupFolder{Name: root, Path: root})
			}
		}
		return result, nil
	}
	if !filepath.IsAbs(path) || strings.HasPrefix(path, "\\\\") {
		return result, errors.New("Choose an absolute folder on a local server drive")
	}
	path = filepath.Clean(path)
	for parent := path; ; parent = filepath.Dir(parent) {
		if e := platform.RejectReparse(parent); e != nil {
			return result, errors.New("Linked folders are not supported; choose a normal server folder")
		}
		if filepath.Dir(parent) == parent {
			break
		}
	}
	dir, e := os.Open(path)
	if e != nil {
		return result, errors.New("Could not open that server folder. Choose a parent folder or enter a destination path directly")
	}
	defer dir.Close()
	entries, e := dir.ReadDir(512)
	if e != nil && !errors.Is(e, io.EOF) {
		return result, errors.New("Could not list that server folder")
	}
	result.Path = path
	result.Parent = filepath.Dir(path)
	if result.Parent == path {
		result.Parent = ""
	}
	result.Truncated = len(entries) == 512
	for _, entry := range entries {
		child := filepath.Join(path, entry.Name())
		if !entry.IsDir() || entry.Type()&os.ModeSymlink != 0 || platform.RejectReparse(child) != nil {
			continue
		}
		result.Folders = append(result.Folders, backupFolder{Name: entry.Name(), Path: child})
	}
	sort.Slice(result.Folders, func(i, j int) bool {
		return strings.ToLower(result.Folders[i].Name) < strings.ToLower(result.Folders[j].Name)
	})
	return result, nil
}
