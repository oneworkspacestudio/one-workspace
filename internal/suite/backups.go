package suite

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"github.com/oneworkspacestudio/calendar/shared"
)

type BackupSettings struct {
	Enabled       bool   `json:"enabled"`
	Destination   string `json:"destination"`
	IntervalHours int    `json:"intervalHours"`
}
type BackupState struct {
	AppsPaused  bool           `json:"appsPaused"`
	Busy        bool           `json:"busy"`
	Phase       string         `json:"phase"`
	Error       string         `json:"error,omitempty"`
	Started     time.Time      `json:"started,omitempty"`
	Finished    time.Time      `json:"finished,omitempty"`
	LastSuccess time.Time      `json:"lastSuccess,omitempty"`
	NextRun     time.Time      `json:"nextRun,omitempty"`
	BytesDone   int64          `json:"bytesDone"`
	BytesTotal  int64          `json:"bytesTotal"`
	Receipt     *BackupReceipt `json:"receipt,omitempty"`
}
type backupController struct {
	mu          sync.Mutex
	picker      sync.Mutex
	settings    BackupSettings
	state       BackupState
	nextAttempt time.Time
	config      Config
}

func newBackupController(c Config) *backupController {
	b := &backupController{config: c, settings: BackupSettings{IntervalHours: 24}}
	settingsProblem := ""
	if e := platform.LoadJSON(filepath.Join(c.Data, "backup-settings.json"), &b.settings); e != nil && !os.IsNotExist(e) {
		b.settings.Enabled = false
		settingsProblem = "Backup settings could not be read. Save the setup again."
	}
	var old BackupState
	if platform.LoadJSON(b.statePath(), &old) == nil {
		b.state = old
		if old.Busy {
			b.state.Busy = false
			b.state.Phase = "Interrupted"
			b.state.Error = "Workspace stopped during a backup. The incomplete copy is not a verified backup. Run again when ready."
		}
	}
	if settingsProblem != "" {
		b.state.Error = settingsProblem
		b.state.Phase = "Setup needs attention"
	}
	if b.settings.IntervalHours < 1 || b.settings.IntervalHours > 168 {
		b.settings.Enabled = false
		b.settings.IntervalHours = 24
	}
	return b
}
func (b *backupController) statePath() string {
	return filepath.Join(b.config.Data, "maintenance", "backups", "status.json")
}
func (b *backupController) snapshot() (BackupSettings, BackupState) {
	b.mu.Lock()
	defer b.mu.Unlock()
	s := b.state
	if b.settings.Enabled {
		s.NextRun = s.LastSuccess.Add(time.Duration(b.settings.IntervalHours) * time.Hour)
		if s.LastSuccess.IsZero() {
			s.NextRun = time.Now().UTC()
		}
		if b.nextAttempt.After(s.NextRun) {
			s.NextRun = b.nextAttempt
		}
	}
	return b.settings, s
}
func (b *backupController) busy() bool {
	if b == nil {
		return false
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state.Busy
}
func (b *backupController) appsPaused() bool {
	if b == nil {
		return false
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state.Busy && b.state.AppsPaused
}
func (b *backupController) change(fn func(*BackupState)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	fn(&b.state)
	// This folder is excluded from snapshots, so progress never changes a copied file.
	if e := os.MkdirAll(filepath.Dir(b.statePath()), 0700); e == nil {
		_ = platform.SaveJSON(b.statePath(), b.state)
	}
}

// CheckBackupDestination is read-only and used before enabling a schedule and
// again at run time (a removable disk may have been unplugged or replaced).
func CheckBackupDestination(c Config, destination string) (string, error) {
	destination = strings.TrimSpace(destination)
	if !filepath.IsAbs(destination) || strings.HasPrefix(destination, `\\`) {
		return "", errors.New("choose a folder on a local NTFS drive, for example E:\\Workspace backups")
	}
	destination = filepath.Clean(destination)
	if filepath.Dir(destination) == destination {
		return "", errors.New("choose a backup folder rather than the drive root")
	}
	for _, root := range []string{c.Programs, c.Data, c.Files} {
		a, b := strings.ToLower(destination), strings.ToLower(filepath.Clean(root))
		if a == b || strings.HasPrefix(a, b+string(os.PathSeparator)) || strings.HasPrefix(b, a+string(os.PathSeparator)) {
			return "", errors.New("choose a separate backup folder outside Apps, Private and Files; it must not contain those folders either")
		}
	}
	if _, e := platform.CheckDataVolume(destination); e != nil {
		return "", fmt.Errorf("backups contain account secrets and need local NTFS permissions: %w", e)
	}
	for p := destination; ; p = filepath.Dir(p) {
		if _, e := os.Lstat(p); e == nil {
			if e = platform.RejectReparse(p); e != nil {
				return "", e
			}
			info, e := os.Stat(p)
			if e != nil || !info.IsDir() {
				return "", errors.New("the backup location must be an ordinary folder")
			}
		} else if !os.IsNotExist(e) {
			return "", e
		}
		if filepath.Dir(p) == p {
			break
		}
	}
	return destination, nil
}
func (b *backupController) save(in BackupSettings) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.state.Busy {
		return errors.New("wait for the current backup to finish before changing its setup")
	}
	if in.IntervalHours < 1 || in.IntervalHours > 168 {
		return errors.New("choose an interval from 1 to 168 hours")
	}
	if in.Enabled || (strings.TrimSpace(in.Destination) != "" && in.Destination != b.settings.Destination) {
		var err error
		in.Destination, err = CheckBackupDestination(b.config, in.Destination)
		if err != nil {
			return err
		}
	}
	if e := platform.SaveJSON(filepath.Join(b.config.Data, "backup-settings.json"), in); e != nil {
		return e
	}
	b.settings = in
	b.nextAttempt = time.Time{}
	return nil
}
func (rt *Runtime) backupSchedule(ctx context.Context) {
	timer := time.NewTicker(30 * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
		}
		settings, state := rt.backups.snapshot()
		if settings.Enabled && !state.Busy && !time.Now().Before(state.NextRun) {
			if e := rt.beginBackup(ctx); e != nil && ctx.Err() == nil {
				rt.backups.mu.Lock()
				rt.backups.nextAttempt = time.Now().Add(30 * time.Minute)
				rt.backups.mu.Unlock()
				rt.backups.change(func(s *BackupState) { s.Phase = "Scheduled backup is waiting"; s.Error = e.Error() })
			}
		}
	}
}
func (rt *Runtime) beginBackup(ctx context.Context) error {
	if e := ctx.Err(); e != nil {
		return e
	}
	if !rt.operations.TryLock() {
		return errors.New("Workspace is already performing maintenance")
	}
	guard, e := lockUpdate(rt.Config)
	if e != nil {
		rt.operations.Unlock()
		return errors.New("an update or recovery operation is using this installation; retry after it finishes")
	}
	if UpdatePending(rt.Config) {
		guard.Close()
		rt.operations.Unlock()
		return errors.New("finish update recovery before backing up")
	}
	b := rt.backups
	b.mu.Lock()
	if b.state.Busy {
		b.mu.Unlock()
		guard.Close()
		rt.operations.Unlock()
		return errors.New("a backup is already running")
	}
	destination, e := CheckBackupDestination(rt.Config, b.settings.Destination)
	if e != nil {
		b.nextAttempt = time.Now().Add(30 * time.Minute)
		b.mu.Unlock()
		guard.Close()
		rt.operations.Unlock()
		b.change(func(s *BackupState) { s.Phase = "Destination needs attention"; s.Error = e.Error() })
		return e
	}
	b.state.Busy = true
	b.mu.Unlock()
	go func() { defer rt.operations.Unlock(); defer guard.Close(); rt.runBackup(ctx, destination) }()
	return nil
}
func (rt *Runtime) runBackup(ctx context.Context, destination string) {
	b, c := rt.backups, rt.Config
	b.change(func(s *BackupState) {
		last := s.LastSuccess
		*s = BackupState{Busy: true, Phase: "Checking storage and available space", Started: time.Now().UTC(), LastSuccess: last}
	})
	var err error
	paused := false
	var receipt BackupReceipt
	defer func() {
		if paused && ctx.Err() == nil {
			b.change(func(s *BackupState) { s.Phase = "Restarting apps and checking their health" })
			rt.sharedOnce = sync.Once{}
			rt.sharedErr = nil
			var errJob error
			rt.job, errJob = platform.NewJob()
			if errJob == nil {
				errJob = rt.startApps(ctx, false)
			}
			if errJob != nil {
				rt.progress.Finish("failed", "App restart after backup needs attention", errJob)
				err = errors.Join(err, fmt.Errorf("apps could not restart: %w; restart Workspace and inspect startup progress", errJob))
			}
		}
		b.mu.Lock()
		b.nextAttempt = time.Now().Add(30 * time.Minute)
		b.mu.Unlock()
		b.change(func(s *BackupState) {
			s.Busy = false
			s.AppsPaused = false
			s.Finished = time.Now().UTC()
			s.Phase = "Backup complete"
			if receipt.ID != "" {
				s.LastSuccess = receipt.Created
				s.Receipt = &receipt
			}
			if err != nil {
				s.Phase = "Backup needs attention"
				s.Error = err.Error()
			}
		})
	}()
	// Capacity, configuration and program integrity checks precede downtime.
	actual, _, loadErr := Load(c.Data)
	if loadErr != nil {
		err = fmt.Errorf("cannot read the installation configuration: %w", loadErr)
		return
	}
	if actual.Data != c.Data || actual.Programs != c.Programs || actual.Files != c.Files || actual.Version != c.Version || actual.Owner != c.Owner {
		err = errors.New("installation settings changed; restart Workspace before backing up")
		return
	}
	if err = verifyInstalledPrograms(c); err != nil {
		return
	}
	var size int64
	size, err = snapshotSize(c)
	// Online caches may rotate while estimating. The stopped tree is sized again before copying.
	if os.IsNotExist(err) {
		size = 0
		err = nil
	}
	if err != nil {
		return
	}
	if err = os.MkdirAll(destination, 0700); err != nil {
		return
	}
	var free uint64
	free, err = platform.FreeBytes(destination)
	if err != nil {
		return
	}
	if free < uint64(size)*2+(512<<20) {
		err = errors.New("the backup drive needs more free space: allow twice the installation size plus 512 MB for verified backup and recovery")
		return
	}
	b.change(func(s *BackupState) {
		s.Phase = "Pausing apps; the Backups page stays available"
		s.BytesTotal = size
		s.AppsPaused = true
	})
	if err = cacheUpdateController(c, c.Catalog, ""); err != nil {
		return
	}
	paused = true
	if err = rt.stopApps(); err != nil {
		return
	}
	for _, app := range c.Catalog.Apps {
		if app.Kind == "process" {
			rt.set(app.ID, "paused for backup")
		}
	}
	if err = stopBackupShared(ctx, c); err != nil {
		return
	}
	if err = waitBackupPorts(ctx, c); err != nil {
		return
	}
	if err = ctx.Err(); err != nil {
		return
	}
	receipt, err = makeUpdateBackupDetailed(ctx, c, destination, func(phase string) { b.change(func(s *BackupState) { s.Phase = phase }) }, func(done, total int64) { b.change(func(s *BackupState) { s.BytesDone = done; s.BytesTotal = total }) })
	if err != nil {
		receipt = BackupReceipt{}
		return
	}
	if err = recordBackup(c, receipt); err != nil {
		return
	}
	// Keep a receipt with the external snapshot too, for recovery if the source disk is lost.
	err = platform.SaveJSON(filepath.Join(receipt.Path, "receipt.json"), receipt)
}
func stopBackupShared(ctx context.Context, c Config) error {
	root := filepath.Join(c.Data, "shared")
	if run, e := shared.ReadRunning(root); e == nil && platform.ProcessRunning(run.PID, filepath.Join(root, "WorkspaceServices.exe")) {
		cfg, _, e := shared.Load(root)
		if e != nil || cfg.EnginePort != c.Port+11 {
			return errors.New("shared service ownership does not match this installation")
		}
		return shared.Stop(ctx, root)
	}
	return nil
}
func waitBackupPorts(ctx context.Context, c Config) error {
	ports := []int{}
	for _, a := range c.Catalog.Apps {
		if a.Kind != "process" {
			continue
		}
		ports = append(ports, a.Port)
		if a.ID == "vault" {
			ports = append(ports, a.Port+1)
		}
		if a.ID == "mail" && !usesShared(c, "mail") {
			ports = append(ports, a.Port+1, a.Port+2)
		}
	}
	if SharedEnabled(c) {
		ports = append(ports, c.Port+11)
	}
	deadline := time.Now().Add(30 * time.Second)
	for {
		quiet := true
		for _, port := range ports {
			l, e := net.Listen("tcp4", fmt.Sprintf("127.0.0.1:%d", port))
			if e != nil {
				quiet = false
				break
			}
			l.Close()
		}
		if quiet {
			return nil
		}
		if time.Now().After(deadline) {
			return errors.New("apps have not fully stopped; no snapshot was marked complete")
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(200 * time.Millisecond):
		}
	}
}
