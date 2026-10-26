package suite

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

type backupReuseFile struct {
	path string
	file BackupFile
}

// Only immutable copies already in protected, recorded backups are linked.
// Never link the live app data: databases/models may be overwritten by apps.
// Each new name owns its bytes even when the older snapshot is later removed.
type backupReuser map[string][]backupReuseFile

func reusableBackupFiles(c Config) backupReuser {
	index := backupReuser{}
	receipts, err := UpdateBackups(c)
	if err != nil {
		return index
	}
	for _, r := range receipts {
		if r.RecoveryOnly || r.CleanupPending || filepath.Base(r.Path) != r.ID || cleanupPath(r.Path) != nil {
			continue
		}
		sum, _, err := fileDigest(filepath.Join(r.Path, "backup.json"))
		if err != nil || sum != r.SHA256 {
			continue
		}
		var b UpdateBackup
		if platform.LoadJSON(filepath.Join(r.Path, "backup.json"), &b) != nil || b.ID != r.ID || b.Config.Data != c.Data || b.Config.Programs != c.Programs || b.Config.Files != c.Files || b.Config.Owner != c.Owner || b.RecoveryOnly || b.Purpose != r.Purpose || validateBackupScope(b, r) != nil {
			continue
		}
		for _, f := range b.Files {
			if f.Directory || f.Size < 1<<20 || len(f.Streams) != 0 || len(f.SHA256) != 64 || (f.Area != "private" && f.Area != "programs") {
				continue
			}
			path, err := platform.Inside(filepath.Join(r.Path, f.Area), f.Path)
			if err != nil {
				continue
			}
			key := strings.ToLower(f.Area + "/" + f.Path)
			index[key] = append(index[key], backupReuseFile{path, f})
		}
	}
	return index
}

func (index backupReuser) copy(ctx context.Context, area, rel, src, dst string) (BackupFile, bool, error) {
	if err := ctx.Err(); err != nil {
		return BackupFile{}, false, err
	}
	candidates := index[strings.ToLower(area+"/"+filepath.ToSlash(rel))]
	info, err := os.Stat(src)
	if err != nil {
		return BackupFile{}, false, err
	}
	if len(candidates) > 0 && info.Size() >= 1<<20 {
		streams, err := platform.NamedStreams(src)
		if err == nil && len(streams) == 0 {
			liveSum, size, err := fileDigest(src)
			if err != nil {
				return BackupFile{}, false, err
			}
			for _, candidate := range candidates {
				if candidate.file.Size != size || candidate.file.SHA256 != liveSum || cleanupPath(candidate.path) != nil {
					continue
				}
				prior, err := os.Stat(candidate.path)
				if err != nil || !prior.Mode().IsRegular() || os.SameFile(info, prior) {
					continue
				}
				streams, err := platform.NamedStreams(candidate.path)
				if err != nil || len(streams) != 0 {
					continue
				}
				priorSum, priorSize, err := fileDigest(candidate.path)
				if err != nil || priorSum != liveSum || priorSize != size {
					continue
				}
				if err = os.MkdirAll(filepath.Dir(dst), 0700); err != nil {
					return BackupFile{}, false, err
				}
				if err = cleanupPath(filepath.Dir(dst)); err != nil {
					return BackupFile{}, false, err
				}
				if err = os.Link(candidate.path, dst); err != nil {
					// Cross-volume or unsupported link operations use a normal copy.
					continue
				}
				// Do not set attributes/timestamps on a shared backup inode. Each
				// manifest records its own restore metadata independently.
				liveAfter, afterSize, err := fileDigest(src)
				if err != nil || liveAfter != liveSum || afterSize != size {
					_ = os.Remove(dst)
					return BackupFile{}, false, errors.New("a file changed while its backup was being checked; close editors and retry")
				}
				return BackupFile{Size: size, SHA256: liveSum, ModTime: info.ModTime().UnixNano()}, true, nil
			}
		}
	}
	f, err := safeCopy(ctx, src, dst)
	return f, false, err
}
