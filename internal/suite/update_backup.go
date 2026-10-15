package suite

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fahimc/one-workspace/internal/platform"
)

type BackupFile struct {
	Area       string         `json:"area"`
	Path       string         `json:"path"`
	Directory  bool           `json:"directory,omitempty"`
	Size       int64          `json:"size"`
	SHA256     string         `json:"sha256"`
	ModTime    int64          `json:"modTime"`
	Attributes uint32         `json:"attributes,omitempty"`
	DACL       string         `json:"dacl,omitempty"`
	Streams    []BackupStream `json:"streams,omitempty"`
}
type BackupStream struct {
	Name   string `json:"name"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
}
type UpdateBackup struct {
	Schema       int          `json:"schema"`
	ID           string       `json:"id"`
	Created      time.Time    `json:"created"`
	Config       Config       `json:"config"`
	Files        []BackupFile `json:"files"`
	RecoveryOnly bool         `json:"recoveryOnly,omitempty"`
}
type BackupReceipt struct {
	ID           string    `json:"id"`
	Path         string    `json:"path"`
	SHA256       string    `json:"sha256"`
	Version      string    `json:"version"`
	Created      time.Time `json:"created"`
	RecoveryOnly bool      `json:"recoveryOnly,omitempty"`
}

var privateUpdateExclusions = map[string]bool{".owner-sid": true, "backups": true, "updates": true, "maintenance": true, "logs": true, "music-runtime-temp": true, "runtime.lock": true, "install.lock": true, "running.json": true, "startup-progress.json": true, "control.dpapi": true}

func updateRoots(c Config) error {
	roots := []string{c.Data, c.Programs, c.Files}
	for _, p := range roots {
		if !filepath.IsAbs(p) || filepath.Clean(p) != p || filepath.Dir(p) == p {
			return errors.New("update requires dedicated absolute installation folders")
		}
		current := p
		for {
			if e := platform.RejectReparse(current); e != nil {
				return e
			}
			next := filepath.Dir(current)
			if next == current {
				break
			}
			current = next
		}
	}
	for i, a := range roots {
		for j, b := range roots {
			if i != j && (strings.EqualFold(a, b) || strings.HasPrefix(strings.ToLower(a), strings.ToLower(b)+string(os.PathSeparator))) {
				return errors.New("installation folders must not overlap")
			}
		}
	}
	_, e := platform.CheckDataVolume(c.Data)
	return e
}
func walkBackup(c Config, visit func(area, root, path string, entry fs.DirEntry) error) error {
	for _, area := range []struct{ name, root string }{{"private", c.Data}, {"programs", c.Programs}, {"files", c.Files}} {
		e := filepath.WalkDir(area.root, func(path string, d fs.DirEntry, e error) error {
			if e != nil {
				return e
			}
			if path == area.root {
				return platform.RejectReparse(path)
			}
			rel, e := filepath.Rel(area.root, path)
			if e != nil {
				return e
			}
			first := strings.SplitN(filepath.ToSlash(rel), "/", 2)[0]
			if area.name == "private" && (privateUpdateExclusions[first] || strings.HasPrefix(first, ".apps-stage-")) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if e = platform.RejectReparse(path); e != nil {
				return e
			}
			if !d.IsDir() && !d.Type().IsRegular() {
				return errors.New("unsupported special file in backup")
			}
			return visit(area.name, area.root, path, d)
		})
		if e != nil {
			return e
		}
	}
	return nil
}
func fileDigest(path string) (string, int64, error) {
	f, e := os.Open(path)
	if e != nil {
		return "", 0, e
	}
	defer f.Close()
	h := sha256.New()
	n, e := io.Copy(h, f)
	return hex.EncodeToString(h.Sum(nil)), n, e
}
func snapshotSize(c Config) (int64, error) {
	var total int64
	e := walkBackup(c, func(_, _, path string, d fs.DirEntry) error {
		if d.IsDir() {
			return nil
		}
		st, e := d.Info()
		if e == nil {
			total += st.Size()
		}
		streams, err := platform.NamedStreams(path)
		if err != nil {
			return err
		}
		for _, stream := range streams {
			total += stream.Size
		}
		return e
	})
	return total, e
}
func safeCopy(ctx context.Context, src, dst string) (BackupFile, error) {
	var result BackupFile
	if e := ctx.Err(); e != nil {
		return result, e
	}
	if e := platform.CopyFile(src, dst); e != nil {
		return result, e
	}
	sum, size, e := fileDigest(dst)
	if e != nil {
		return result, e
	}
	sourceSum, sourceSize, e := fileDigest(src)
	if e != nil {
		return result, e
	}
	if sum != sourceSum || size != sourceSize {
		return result, errors.New("a file changed while its backup was being copied; retry after closing editors")
	}
	result.Size, result.SHA256 = size, sum
	st, e := os.Stat(src)
	if e != nil {
		return result, e
	}
	result.ModTime = st.ModTime().UnixNano()
	stamp := time.Unix(0, result.ModTime)
	if e = os.Chtimes(dst, stamp, stamp); e != nil {
		return result, e
	}
	return result, nil
}
func makeUpdateBackup(ctx context.Context, c Config, parent string, report func(string)) (BackupReceipt, error) {
	var receipt BackupReceipt
	partial := false
	if actual, _, err := Load(c.Data); err == nil {
		if actual.Data != c.Data || actual.Programs != c.Programs || actual.Files != c.Files {
			return receipt, errors.New("installation locations changed during update")
		}
		c = actual
	} else {
		partial = true
	}
	if e := updateRoots(c); e != nil {
		return receipt, e
	}
	if parent == "" {
		parent = filepath.Join(c.Data, "backups")
	}
	parent = filepath.Clean(parent)
	if !filepath.IsAbs(parent) {
		return receipt, errors.New("backup location must be absolute")
	}
	if strings.HasPrefix(parent, `\\`) {
		return receipt, errors.New("choose a local NTFS backup folder, not a network share")
	}
	for _, root := range []string{c.Programs, c.Files} {
		if strings.EqualFold(parent, root) || strings.HasPrefix(strings.ToLower(parent), strings.ToLower(root)+string(os.PathSeparator)) {
			return receipt, errors.New("backup cannot be inside programs or ordinary files")
		}
	}
	if strings.EqualFold(parent, c.Data) || strings.HasPrefix(strings.ToLower(parent), strings.ToLower(c.Data)+string(os.PathSeparator)) && !strings.EqualFold(parent, filepath.Join(c.Data, "backups")) {
		return receipt, errors.New("private backups must use the installation's backups folder")
	}
	if _, e := platform.CheckDataVolume(parent); e != nil {
		return receipt, e
	}
	probe := parent
	for {
		if _, e := os.Stat(probe); e == nil {
			if e = platform.RejectReparse(probe); e != nil {
				return receipt, e
			}
			if filepath.Dir(probe) == probe {
				break
			}
		} else if !os.IsNotExist(e) {
			return receipt, e
		}
		next := filepath.Dir(probe)
		if next == probe {
			break
		}
		probe = next
	}
	if e := os.MkdirAll(parent, 0700); e != nil {
		return receipt, e
	}
	size, e := snapshotSize(c)
	if e != nil {
		return receipt, e
	}
	free, e := platform.FreeBytes(parent)
	if e != nil {
		return receipt, e
	}
	if free < uint64(size)*2+(512<<20) {
		return receipt, fmt.Errorf("backup volume needs at least %d MB free for a verified backup and rollback archive; choose a larger local NTFS backup location", (size*2+(512<<20))>>20)
	}
	id := "update-" + time.Now().UTC().Format("20060102T150405") + "-" + platform.Token()[:12]
	path := filepath.Join(parent, id)
	if e = platform.RequireNewDir(path); e != nil {
		return receipt, e
	}
	if e = platform.ProtectMaintenanceDir(path, c.Data); e != nil {
		return receipt, e
	}
	backup := UpdateBackup{Schema: 1, ID: id, Created: time.Now().UTC(), Config: c, RecoveryOnly: partial}
	report("Backing up accounts, app data, programs and ordinary files")
	e = walkBackup(c, func(area, root, src string, d fs.DirEntry) error {
		rel, e := filepath.Rel(root, src)
		if e != nil {
			return e
		}
		dest, e := platform.Inside(filepath.Join(path, area), filepath.ToSlash(rel))
		if e != nil {
			return e
		}
		var file BackupFile
		attrs, e := platform.FileAttributes(src)
		if e != nil {
			return e
		}
		if platform.EncryptedFileAttributes(attrs) {
			return errors.New("EFS-encrypted files require a Windows backup that preserves their encryption; no update files were changed")
		}
		streams, e := platform.NamedStreams(src)
		if e != nil {
			return e
		}
		if d.IsDir() {
			if len(streams) > 0 {
				return errors.New("a directory contains alternate data streams and needs a Windows directory-stream backup; no update files were changed")
			}
			e = os.MkdirAll(dest, 0700)
			file.Directory = true
			if info, err := d.Info(); err == nil {
				file.ModTime = info.ModTime().UnixNano()
			} else {
				return err
			}
		} else {
			file, e = safeCopy(ctx, src, dest)
		}
		if e != nil {
			return e
		}
		file.Attributes = attrs
		if area == "files" {
			file.DACL, e = platform.FileDACL(src)
			if e != nil {
				return e
			}
		}
		for _, stream := range streams {
			copy, e := safeCopy(ctx, src+":"+stream.Name, dest+":"+stream.Name)
			if e != nil {
				return e
			}
			file.Streams = append(file.Streams, BackupStream{Name: stream.Name, Size: copy.Size, SHA256: copy.SHA256})
		}
		after, e := platform.NamedStreams(src)
		if e != nil || len(after) != len(streams) {
			return errors.New("alternate file streams changed during backup; close editors and retry")
		}
		file.Area, file.Path = area, filepath.ToSlash(rel)
		backup.Files = append(backup.Files, file)
		return nil
	})
	if e != nil {
		return receipt, e
	}
	if e = platform.SaveJSON(filepath.Join(path, "backup.json"), backup); e != nil {
		return receipt, e
	}
	sum, _, e := fileDigest(filepath.Join(path, "backup.json"))
	if e != nil {
		return receipt, e
	}
	receipt = BackupReceipt{ID: id, Path: path, SHA256: sum, Version: c.Version, Created: backup.Created, RecoveryOnly: partial}
	if _, e = verifyUpdateBackup(c, receipt); e != nil {
		return receipt, e
	}
	return receipt, nil
}
func verifyUpdateBackup(c Config, r BackupReceipt) (UpdateBackup, error) {
	var b UpdateBackup
	if !filepath.IsAbs(r.Path) || r.ID == "" || filepath.Base(r.Path) != r.ID {
		return b, errors.New("invalid backup receipt")
	}
	if _, e := platform.CheckDataVolume(r.Path); e != nil {
		return b, e
	}
	if e := platform.RejectReparse(r.Path); e != nil {
		return b, e
	}
	sum, _, e := fileDigest(filepath.Join(r.Path, "backup.json"))
	if e != nil {
		return b, e
	}
	if sum != r.SHA256 {
		return b, errors.New("backup manifest checksum mismatch")
	}
	if e = platform.LoadJSON(filepath.Join(r.Path, "backup.json"), &b); e != nil {
		return b, e
	}
	if b.Schema != 1 || b.ID != r.ID || b.Config.Data != c.Data || b.Config.Programs != c.Programs || b.Config.Files != c.Files {
		return b, errors.New("backup belongs to a different installation")
	}
	if b.RecoveryOnly != r.RecoveryOnly {
		return b, errors.New("backup recovery type mismatch")
	}
	seen := map[string]bool{}
	for _, f := range b.Files {
		if f.Area != "private" && f.Area != "programs" && f.Area != "files" {
			return b, errors.New("invalid backup area")
		}
		key := strings.ToLower(f.Area + "/" + f.Path)
		if seen[key] {
			return b, errors.New("duplicate backup path")
		}
		seen[key] = true
		path, e := platform.Inside(filepath.Join(r.Path, f.Area), f.Path)
		if e != nil {
			return b, e
		}
		if e = platform.RejectReparse(path); e != nil {
			return b, e
		}
		actual, e := platform.NamedStreams(path)
		if e != nil || len(actual) != len(f.Streams) {
			return b, errors.New("backup alternate stream list differs from manifest")
		}
		streamNames := map[string]bool{}
		for _, stream := range f.Streams {
			if !platform.ValidStreamName(stream.Name) || streamNames[strings.ToLower(stream.Name)] {
				return b, errors.New("invalid or duplicate backup alternate stream")
			}
			streamNames[strings.ToLower(stream.Name)] = true
			sum, size, e := fileDigest(path + ":" + stream.Name)
			if e != nil || size != stream.Size || sum != stream.SHA256 {
				return b, errors.New("backup alternate stream checksum mismatch")
			}
		}
		if f.Directory {
			st, e := os.Stat(path)
			if e != nil || !st.IsDir() {
				return b, errors.New("backup directory missing")
			}
			continue
		}
		sum, size, e := fileDigest(path)
		if e != nil {
			return b, e
		}
		if sum != f.SHA256 || size != f.Size {
			return b, errors.New("backup file checksum mismatch")
		}
	}
	for _, required := range []string{"private/workspace.json", "private/secrets.dpapi", "private/Workspace.exe"} {
		if b.RecoveryOnly {
			continue
		}
		if !seen[strings.ToLower(required)] {
			return b, errors.New("backup is missing required installation state")
		}
	}
	return b, nil
}
func markRecoveryOnly(c Config, r BackupReceipt) (BackupReceipt, error) {
	b, e := verifyUpdateBackup(c, r)
	if e != nil {
		return r, e
	}
	b.RecoveryOnly = true
	if e = platform.SaveJSON(filepath.Join(r.Path, "backup.json"), b); e != nil {
		return r, e
	}
	sum, _, e := fileDigest(filepath.Join(r.Path, "backup.json"))
	if e != nil {
		return r, e
	}
	r.SHA256 = sum
	r.RecoveryOnly = true
	_, e = verifyUpdateBackup(c, r)
	return r, e
}
func writeReplace(src, dst string) error {
	if e := os.MkdirAll(filepath.Dir(dst), 0700); e != nil {
		return e
	}
	if e := platform.RejectReparse(filepath.Dir(dst)); e != nil {
		return e
	}
	if _, e := os.Lstat(dst); e == nil {
		if e = platform.RejectReparse(dst); e != nil {
			return e
		}
		if e = platform.ClearFileReadOnly(dst); e != nil {
			return e
		}
	} else if !os.IsNotExist(e) {
		return e
	}
	temp := filepath.Join(filepath.Dir(dst), ".workspace-replace-"+platform.Token())
	if e := platform.CopyFile(src, temp); e != nil {
		return e
	}
	defer os.Remove(temp)
	return os.Rename(temp, dst)
}

// Before restoration the caller creates a second verified snapshot of the
// current state. Extra ordinary files remain; private schema artifacts remain in
// the protected recovery snapshot so an old database cannot mix with them.
func restoreUpdateBackup(c Config, r BackupReceipt, recovery BackupReceipt, report func(string)) error {
	if e := updateRoots(c); e != nil {
		return e
	}
	b, e := verifyUpdateBackup(c, r)
	if e != nil {
		return e
	}
	if _, e = verifyUpdateBackup(c, recovery); e != nil {
		return e
	}
	entries, e := os.ReadDir(c.Data)
	if e != nil {
		return e
	}
	for _, entry := range entries {
		if privateUpdateExclusions[entry.Name()] || strings.EqualFold(entry.Name(), "Workspace.exe") {
			continue
		}
		src := filepath.Join(c.Data, entry.Name())
		if e = platform.RejectReparse(src); e != nil {
			return e
		}
		// Same-volume moves are ideal, but external protected backup drives are
		// supported too. Current bytes are already in the verified recovery snapshot.
		if e = filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if err = platform.RejectReparse(path); err != nil {
				return err
			}
			return platform.ClearFileReadOnly(path)
		}); e != nil {
			return e
		}
		if e = os.RemoveAll(src); e != nil {
			return e
		}
	}
	report("Restoring the verified pre-update state; newer data remains in the recovery backup")
	var managers []BackupFile
	for _, f := range b.Files {
		if (f.Area == "private" || f.Area == "programs") && strings.EqualFold(f.Path, "Workspace.exe") {
			managers = append(managers, f)
			continue
		}
		root := c.Data
		if f.Area == "programs" {
			root = c.Programs
		} else if f.Area == "files" {
			root = c.Files
		}
		dst, e := platform.Inside(root, f.Path)
		if e != nil {
			return e
		}
		src, e := platform.Inside(filepath.Join(r.Path, f.Area), f.Path)
		if e != nil {
			return e
		}
		if current, e := os.Lstat(dst); e == nil && current.IsDir() != f.Directory {
			if e = platform.RejectReparse(dst); e != nil {
				return e
			}
			// A verified full recovery snapshot already preserves this conflicting
			// file/directory before its original type is restored.
			if e = os.RemoveAll(dst); e != nil {
				return e
			}
		} else if e != nil && !os.IsNotExist(e) {
			return e
		}
		if f.Directory {
			e = os.MkdirAll(dst, 0700)
		} else {
			e = writeReplace(src, dst)
			if e == nil {
				for _, stream := range f.Streams {
					if e = platform.CopyFile(src+":"+stream.Name, dst+":"+stream.Name); e != nil {
						break
					}
				}
			}
			if e == nil {
				stamp := time.Unix(0, f.ModTime)
				e = os.Chtimes(dst, stamp, stamp)
			}
			if e == nil {
				e = platform.RestoreFileAttributes(dst, f.Attributes)
			}
		}
		if e != nil {
			return e
		}
	}
	for i := len(b.Files) - 1; i >= 0; i-- {
		f := b.Files[i]
		if !f.Directory {
			continue
		}
		root := c.Data
		if f.Area == "programs" {
			root = c.Programs
		} else if f.Area == "files" {
			root = c.Files
		}
		dst, e := platform.Inside(root, f.Path)
		if e != nil {
			return e
		}
		stamp := time.Unix(0, f.ModTime)
		if e = os.Chtimes(dst, stamp, stamp); e != nil {
			return e
		}
		if e = platform.RestoreFileAttributes(dst, f.Attributes); e != nil {
			return e
		}
	}
	for i := len(b.Files) - 1; i >= 0; i-- {
		f := b.Files[i]
		if f.Area != "files" {
			continue
		}
		dst, e := platform.Inside(c.Files, f.Path)
		if e != nil {
			return e
		}
		if e = platform.RestoreFileDACL(dst, f.DACL); e != nil {
			return e
		}
	}
	// Legacy managers do not understand the recovery journal. Restore them only
	// after every database, file and program is back, so a reboot cannot start a
	// legacy launcher against a partially restored installation.
	for _, f := range managers {
		root := c.Data
		if f.Area == "programs" {
			root = c.Programs
		}
		dst, e := platform.Inside(root, f.Path)
		if e != nil {
			return e
		}
		src, e := platform.Inside(filepath.Join(r.Path, f.Area), f.Path)
		if e != nil {
			return e
		}
		if e = writeReplace(src, dst); e != nil {
			return e
		}
		for _, stream := range f.Streams {
			if e = platform.CopyFile(src+":"+stream.Name, dst+":"+stream.Name); e != nil {
				return e
			}
		}
		stamp := time.Unix(0, f.ModTime)
		if e = os.Chtimes(dst, stamp, stamp); e != nil {
			return e
		}
		if e = platform.RestoreFileAttributes(dst, f.Attributes); e != nil {
			return e
		}
	}
	return nil
}
