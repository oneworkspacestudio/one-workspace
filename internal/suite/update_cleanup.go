package suite

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/operationprogress"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

const cleanupMarker = ".workspace-update-cleanup.json"

type cleanupOwner struct {
	Schema int    `json:"schema"`
	Kind   string `json:"kind"`
	Data   string `json:"data"`
	ID     string `json:"id"`
	Digest string `json:"digest,omitempty"`
}

// Only receipts explicitly labelled by the automatic updater are retired.
// Legacy snapshots cannot be distinguished from scheduled/manual backups.
type UpdateCleanupReport struct {
	Finished         time.Time `json:"finished"`
	KeptBackup       string    `json:"keptBackup"`
	RetainedVerified bool      `json:"retainedVerified"`
	RemovedBackups   int       `json:"removedBackups"`
	RemovedStages    int       `json:"removedStages"`
	BytesFreed       int64     `json:"bytesFreed"`
	Warnings         []string  `json:"warnings,omitempty"`
}

func (r UpdateCleanupReport) Summary() string {
	if !r.RetainedVerified {
		return "Cleanup kept older recovery copies because the completed update or retained recovery copy could not be verified. See the cleanup details."
	}
	s := fmt.Sprintf("Cleanup removed %d older update backup(s) and %d temporary update folder(s), freeing approximately %.1f MB. The latest verified recovery copy is kept.", r.RemovedBackups, r.RemovedStages, float64(r.BytesFreed)/1e6)
	if len(r.Warnings) > 0 {
		s += " Some files were kept because cleanup needs attention; see the cleanup details."
	}
	return s
}

func LoadUpdateCleanup(c Config) (*UpdateCleanupReport, error) {
	var r UpdateCleanupReport
	e := platform.LoadJSON(filepath.Join(c.Data, "updates", "cleanup.json"), &r)
	if os.IsNotExist(e) {
		return nil, nil
	}
	return &r, e
}

func cleanupPath(path string) error {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path || filepath.Dir(path) == path {
		return errors.New("cleanup requires a dedicated absolute folder")
	}
	for p := path; ; p = filepath.Dir(p) {
		if e := platform.RejectReparse(p); e != nil {
			return e
		}
		if filepath.Dir(p) == p {
			break
		}
	}
	return nil
}

func withinCleanup(root, path string) bool {
	r, e := filepath.Rel(root, path)
	return e == nil && r != "." && r != ".." && !strings.HasPrefix(r, ".."+string(os.PathSeparator)) && !filepath.IsAbs(r)
}

// Preflight the entire tree before deleting anything; never follow junctions.
// Remove individual entries with fresh ancestor checks. Metadata is removed last
// so a failed deletion can be retried without mistaking a partial backup for a
// valid rollback candidate. Restricted backup/stage directories serialize writers.
func removeCleanupTree(ctx context.Context, root string, allowed map[string]bool) (int64, error) {
	if e := cleanupPath(root); e != nil {
		return 0, e
	}
	busy, e := platform.DirectoryHasProcesses(root)
	if e != nil || busy {
		return 0, errors.New("cleanup folder is in use or process ownership could not be checked")
	}
	var paths []string
	e = filepath.WalkDir(root, func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if e = ctx.Err(); e != nil {
			return e
		}
		if e = platform.RejectReparse(path); e != nil {
			return e
		}
		if path != root && !withinCleanup(root, path) {
			return errors.New("cleanup path escaped its folder")
		}
		if !d.IsDir() && !d.Type().IsRegular() {
			return errors.New("cleanup found a special file")
		}
		rel, e := filepath.Rel(root, path)
		if e != nil {
			return e
		}
		if allowed != nil && !allowed[strings.ToLower(filepath.ToSlash(rel))] {
			return errors.New("backup contains an unrecorded file; it was kept")
		}
		paths = append(paths, path)
		return nil
	})
	if e != nil {
		return 0, e
	}
	sort.Slice(paths, func(i, j int) bool {
		// Children first. Root metadata (especially the marker) outlives data.
		rank := func(p string) int {
			if p == root {
				return -3
			}
			if p == filepath.Join(root, cleanupMarker) {
				return -2
			}
			if p == filepath.Join(root, "backup.json") {
				return -1
			}
			return strings.Count(p, string(os.PathSeparator))
		}
		return rank(paths[i]) > rank(paths[j])
	})
	var freed int64
	for _, path := range paths {
		if e = ctx.Err(); e != nil {
			return freed, e
		}
		if e = cleanupPath(path); e != nil {
			return freed, e
		}
		st, e := os.Lstat(path)
		if e != nil {
			return freed, e
		}
		if e = os.Remove(path); e != nil {
			return freed, e
		}
		if !st.IsDir() {
			freed += st.Size()
		}
	}
	return freed, nil
}

func newUpdateStage(c Config) (string, error) {
	root := filepath.Join(c.Data, "updates")
	if e := cleanupPath(root); e != nil {
		return "", e
	}
	stage, e := os.MkdirTemp(root, "stage-")
	if e != nil {
		return "", e
	}
	owner := cleanupOwner{Schema: 1, Kind: "stage", Data: c.Data, ID: filepath.Base(stage)}
	if e = platform.SaveJSON(filepath.Join(stage, cleanupMarker), owner); e != nil {
		_ = os.Remove(stage)
		return "", e
	}
	return stage, nil
}

func validStage(c Config, path string) bool {
	if filepath.Dir(path) != filepath.Join(c.Data, "updates") || !regexp.MustCompile(`^stage-[0-9]+$`).MatchString(filepath.Base(path)) {
		return false
	}
	if cleanupPath(path) != nil {
		return false
	}
	var m cleanupOwner
	return platform.LoadJSON(filepath.Join(path, cleanupMarker), &m) == nil && m.Schema == 1 && m.Kind == "stage" && m.Data == c.Data && m.ID == filepath.Base(path)
}

func removeUpdateStage(c Config, stage string, report func(string)) {
	if !validStage(c, stage) {
		return
	}
	if _, e := removeCleanupTree(context.Background(), stage, nil); e != nil && report != nil {
		report("Some temporary update files could not be removed. They will be retried after a successful update.")
	}
}

func cleanupBackupBoundary(c Config, r BackupReceipt, keep BackupReceipt) error {
	if !regexp.MustCompile(`^update-[0-9]{8}T[0-9]{6}-[a-f0-9]{12}$`).MatchString(r.ID) || filepath.Base(r.Path) != r.ID {
		return errors.New("unrecognised backup location")
	}
	if e := cleanupPath(r.Path); e != nil {
		return e
	}
	for _, protected := range []string{c.Data, c.Programs, c.Files, keep.Path} {
		if strings.EqualFold(protected, r.Path) || withinCleanup(r.Path, protected) {
			return errors.New("backup overlaps protected installation data")
		}
	}
	if withinCleanup(c.Programs, r.Path) || withinCleanup(c.Files, r.Path) {
		return errors.New("backup is inside ordinary files or programs")
	}
	if withinCleanup(c.Data, r.Path) && filepath.Dir(r.Path) != filepath.Join(c.Data, "backups") {
		return errors.New("unexpected private backup location")
	}
	return nil
}

func backupCleanupFiles(c Config, r BackupReceipt, legacy bool) (map[string]bool, error) {
	var b UpdateBackup
	if e := platform.LoadJSON(filepath.Join(r.Path, "backup.json"), &b); e != nil {
		return nil, e
	}
	sum, _, e := fileDigest(filepath.Join(r.Path, "backup.json"))
	if e != nil || sum != r.SHA256 || b.Schema != 1 || b.ID != r.ID || !b.Created.Equal(r.Created) || b.Config.Version != r.Version || (b.Purpose != "update" && !(legacy && b.Purpose == "")) || b.Purpose != r.Purpose || b.RecoveryOnly || b.Config.Data != c.Data || b.Config.Programs != c.Programs || b.Config.Files != c.Files || b.Config.Owner != c.Owner {
		return nil, errors.New("backup ownership could not be verified")
	}
	allowed := map[string]bool{".": true, "backup.json": true, cleanupMarker: true}
	for _, f := range b.Files {
		if f.Area != "private" && f.Area != "programs" && f.Area != "files" {
			return nil, errors.New("invalid backup area")
		}
		path, e := platform.Inside(filepath.Join(r.Path, f.Area), f.Path)
		if e != nil {
			return nil, e
		}
		for path != r.Path {
			rel, e := filepath.Rel(r.Path, path)
			if e != nil {
				return nil, e
			}
			allowed[strings.ToLower(filepath.ToSlash(rel))] = true
			path = filepath.Dir(path)
		}
	}
	return allowed, nil
}

func cleanupSuccessfulUpdate(ctx context.Context, c Config, keep BackupReceipt, active string, report func(string)) UpdateCleanupReport {
	r := UpdateCleanupReport{KeptBackup: keep.ID}
	if report == nil {
		report = func(string) {}
	}
	operationprogress.Report(ctx, "cleanup", "Checking recovery copy before cleanup", 0, 0, "")
	report("Checking the recovery copy before removing older update files")
	finish := func() UpdateCleanupReport {
		r.Finished = time.Now().UTC()
		if e := platform.SaveJSON(filepath.Join(c.Data, "updates", "cleanup.json"), r); e != nil {
			r.Warnings = append(r.Warnings, "Could not save the cleanup receipt")
		}
		return r
	}
	var j UpdateJournal
	if e := platform.LoadJSON(updateJournalPath(c), &j); e != nil || j.Phase != "complete" || j.Backup.ID != keep.ID || j.Backup.SHA256 != keep.SHA256 || j.Backup.Purpose != keep.Purpose || (keep.Purpose != "update" && keep.Purpose != "") || keep.RecoveryOnly || keep.CleanupPending {
		r.Warnings = append(r.Warnings, "Cleanup skipped: a completed healthy update and its recovery copy are required")
		return finish()
	}
	verifyCtx := operationprogress.Observe(ctx, func(e operationprogress.Event) {
		operationprogress.Report(ctx, "cleanup", "Checking the retained recovery copy", e.Done, e.Total, e.Unit)
	})
	if e := cleanupPath(keep.Path); e != nil {
		r.Warnings = append(r.Warnings, "Cleanup skipped: recovery copy location could not be verified")
		return finish()
	}
	if _, e := verifyUpdateBackupProgress(verifyCtx, c, keep); e != nil {
		r.Warnings = append(r.Warnings, "Cleanup skipped: retained recovery copy could not be verified")
		return finish()
	}
	rs, e := UpdateBackups(c)
	if e != nil {
		r.Warnings = append(r.Warnings, "Could not read backup history")
		return finish()
	}
	kept := false
	for _, v := range rs {
		if v.ID == keep.ID && v.Path == keep.Path && v.SHA256 == keep.SHA256 && !v.CleanupPending {
			kept = true
		}
	}
	if !kept {
		r.Warnings = append(r.Warnings, "Cleanup skipped: recovery copy is missing from backup history")
		return finish()
	}
	r.RetainedVerified = true
	operationprogress.Report(ctx, "cleanup", "Removing older automatic update backups", 0, int64(len(rs)), "backups")
	for i := 0; i < len(rs); {
		v := rs[i]
		legacy := v.Purpose == "" && (v.CleanupPending || (j.PreviousAutomaticBackup != nil && j.PreviousAutomaticBackup.ID == v.ID && j.PreviousAutomaticBackup.Path == v.Path && j.PreviousAutomaticBackup.SHA256 == v.SHA256))
		if v.ID == keep.ID || (v.Purpose != "update" && !legacy) || v.RecoveryOnly || !v.Created.Before(keep.Created) || (j.Recovery != nil && v.ID == j.Recovery.ID) {
			i++
			continue
		}
		if e = ctx.Err(); e != nil {
			r.Warnings = append(r.Warnings, "Cleanup interrupted; remaining copies were kept")
			break
		}
		if _, e = os.Lstat(v.Path); os.IsNotExist(e) && v.CleanupPending {
			rs = append(rs[:i], rs[i+1:]...)
			if e = platform.SaveJSON(filepath.Join(c.Data, "updates", "backups.json"), rs); e != nil {
				r.Warnings = append(r.Warnings, "Could not refresh backup history")
				return finish()
			}
			continue
		}
		if e = cleanupBackupBoundary(c, v, keep); e != nil {
			r.Warnings = append(r.Warnings, "Kept older backup "+v.ID+": "+e.Error())
			i++
			continue
		}
		allowed, e := backupCleanupFiles(c, v, legacy)
		if e != nil {
			r.Warnings = append(r.Warnings, "Kept older backup "+v.ID+": metadata could not be verified")
			i++
			continue
		}
		marker := cleanupOwner{Schema: 1, Kind: "backup", Data: c.Data, ID: v.ID, Digest: v.SHA256}
		if !v.CleanupPending {
			if _, e = verifyUpdateBackupProgress(verifyCtx, c, v); e != nil {
				r.Warnings = append(r.Warnings, "Kept older backup "+v.ID+": verification failed")
				i++
				continue
			}
			// Reject extra files before marking the copy unavailable for rollback.
			if e = validateCleanupContents(v.Path, allowed); e != nil {
				r.Warnings = append(r.Warnings, "Kept older backup "+v.ID+": "+e.Error())
				i++
				continue
			}
			if e = platform.SaveJSON(filepath.Join(v.Path, cleanupMarker), marker); e != nil {
				r.Warnings = append(r.Warnings, "Could not mark older backup for cleanup")
				i++
				continue
			}
			rs[i].CleanupPending = true
			if e = platform.SaveJSON(filepath.Join(c.Data, "updates", "backups.json"), rs); e != nil {
				r.Warnings = append(r.Warnings, "Could not record cleanup; older backup was kept")
				return finish()
			}
		} else {
			var m cleanupOwner
			if platform.LoadJSON(filepath.Join(v.Path, cleanupMarker), &m) != nil || m != marker {
				r.Warnings = append(r.Warnings, "Kept incomplete cleanup with unverified ownership")
				i++
				continue
			}
		}
		bytes, e := removeCleanupTree(ctx, v.Path, allowed)
		r.BytesFreed += bytes
		if e != nil {
			r.Warnings = append(r.Warnings, "Could not finish removing older backup "+v.ID+": "+e.Error())
			i++
			continue
		}
		rs = append(rs[:i], rs[i+1:]...)
		r.RemovedBackups++
		if e = platform.SaveJSON(filepath.Join(c.Data, "updates", "backups.json"), rs); e != nil {
			r.Warnings = append(r.Warnings, "Could not refresh backup history")
			return finish()
		}
		operationprogress.Report(ctx, "cleanup", "Removing older automatic update backups", int64(r.RemovedBackups), int64(len(rs)+r.RemovedBackups), "backups")
	}
	entries, e := os.ReadDir(filepath.Join(c.Data, "updates"))
	if e == nil {
		for _, entry := range entries {
			path := filepath.Join(c.Data, "updates", entry.Name())
			if path == active || withinCleanup(path, active) || !validStage(c, path) {
				continue
			}
			bytes, e := removeCleanupTree(ctx, path, nil)
			r.BytesFreed += bytes
			if e != nil {
				r.Warnings = append(r.Warnings, "Some temporary update files were kept because they are in use or need attention")
			} else {
				r.RemovedStages++
			}
		}
	}
	operationprogress.Report(ctx, "cleanup", "Cleanup finished; latest recovery copy retained", int64(r.RemovedBackups+r.RemovedStages), int64(r.RemovedBackups+r.RemovedStages), "folders")
	return finish()
}

// Older maintenance binaries cannot run a newly installed manager's cleanup
// code. Once Run has passed all app checks, wait for their committed journal and
// maintenance lock. The new manager can then finish cleanup without racing the
// old updater or treating its still-testing transaction as successful.
func startupUpdateCleanup(ctx context.Context, c Config) {
	deadline := time.NewTimer(2 * time.Minute)
	defer deadline.Stop()
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()
	for {
		var j UpdateJournal
		if platform.LoadJSON(updateJournalPath(c), &j) != nil {
			return
		}
		if j.Phase == "rolled-back" {
			return
		}
		if j.Phase == "complete" {
			last, e := LoadUpdateCleanup(c)
			if e == nil && last != nil && last.KeptBackup == j.Backup.ID && !last.Finished.Before(j.Updated) && len(last.Warnings) == 0 {
				return
			}
			guard, e := lockUpdate(c)
			if e == nil {
				defer guard.Close()
				// Configuration/journal may have changed while the prior updater held
				// the lock. Only this installation's current healthy manager may act.
				actual, _, e := Load(c.Data)
				if e != nil || actual.Version != c.Version || actual.Owner != c.Owner {
					return
				}
				if platform.LoadJSON(updateJournalPath(c), &j) != nil || j.Phase != "complete" {
					return
				}
				last, e = LoadUpdateCleanup(c)
				if e == nil && last != nil && last.KeptBackup == j.Backup.ID && !last.Finished.Before(j.Updated) && len(last.Warnings) == 0 {
					return
				}
				log := &platform.Log{Path: filepath.Join(c.Data, "updates", "cleanup.log")}
				result := cleanupSuccessfulUpdate(ctx, actual, j.Backup, "", func(s string) { fmt.Fprintln(log, s) })
				fmt.Fprintln(log, result.Summary())
				return
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-deadline.C:
			return
		case <-tick.C:
		}
	}
}

func validateCleanupContents(root string, allowed map[string]bool) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if e = platform.RejectReparse(path); e != nil {
			return e
		}
		rel, e := filepath.Rel(root, path)
		if e != nil {
			return e
		}
		if !allowed[strings.ToLower(filepath.ToSlash(rel))] {
			return errors.New("unrecorded files were found")
		}
		return nil
	})
}
