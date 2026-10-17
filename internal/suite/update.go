package suite

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fahimc/one-workspace/internal/catalog"
	"github.com/fahimc/one-workspace/internal/distribution"
	"github.com/fahimc/one-workspace/internal/platform"
	"github.com/oneworkspacestudio/calendar/shared"
)

type UpdateJournal struct {
	Schema   int            `json:"schema"`
	Phase    string         `json:"phase"`
	Backup   BackupReceipt  `json:"backup"`
	Recovery *BackupReceipt `json:"recovery,omitempty"`
	Token    string         `json:"token,omitempty"`
	Error    string         `json:"error,omitempty"`
	Updated  time.Time      `json:"updated"`
}

var errRestoredNeedsStart = errors.New("original data and programs restored, but Workspace could not start")

func updateJournalPath(c Config) string { return filepath.Join(c.Data, "updates", "transaction.json") }
func LoadUpdateConfig(data string) (Config, error) {
	data = filepath.Clean(data)
	c, _, e := Load(data)
	if e != nil {
		var j UpdateJournal
		if err := platform.LoadJSON(filepath.Join(data, "updates", "transaction.json"), &j); err != nil {
			return c, e
		}
		sum, _, err := fileDigest(filepath.Join(j.Backup.Path, "backup.json"))
		if err != nil || sum != j.Backup.SHA256 {
			return c, errors.New("recovery backup manifest could not be verified")
		}
		var b UpdateBackup
		if err = platform.LoadJSON(filepath.Join(j.Backup.Path, "backup.json"), &b); err != nil {
			return c, err
		}
		c = b.Config
	}
	if !strings.EqualFold(c.Data, data) {
		return c, errors.New("this configuration belongs to another private data folder; open the original installation")
	}
	return c, updateRoots(c)
}
func saveUpdateJournal(c Config, j *UpdateJournal) error {
	j.Schema = 1
	j.Updated = time.Now().UTC()
	return platform.SaveJSON(updateJournalPath(c), j)
}
func UpdatePending(c Config) bool {
	var j UpdateJournal
	if e := platform.LoadJSON(updateJournalPath(c), &j); os.IsNotExist(e) {
		return false
	} else if e != nil {
		return true
	}
	return j.Phase != "complete" && j.Phase != "rolled-back"
}
func CheckUpdateRecovery(c Config) error {
	var j UpdateJournal
	e := platform.LoadJSON(updateJournalPath(c), &j)
	if os.IsNotExist(e) {
		return nil
	}
	if e != nil {
		return errors.New("update journal is unreadable; open Workspace updates to recover")
	}
	if j.Schema != 1 {
		return errors.New("unsupported update journal")
	}
	if j.Phase == "complete" || j.Phase == "rolled-back" {
		return nil
	}
	if j.Phase == "testing" && j.Token != "" && os.Getenv("WORKSPACE_UPDATE_TRANSACTION") == j.Token {
		return nil
	}
	return errors.New("an unfinished update requires recovery; run WorkspaceSetup --update with this private data folder")
}
func UpdateBackups(c Config) ([]BackupReceipt, error) {
	var rs []BackupReceipt
	e := platform.LoadJSON(filepath.Join(c.Data, "updates", "backups.json"), &rs)
	if os.IsNotExist(e) {
		e = nil
	}
	return rs, e
}
func recordBackup(c Config, r BackupReceipt) error {
	rs, e := UpdateBackups(c)
	if e != nil {
		return e
	}
	rs = append([]BackupReceipt{r}, rs...)
	return platform.SaveJSON(filepath.Join(c.Data, "updates", "backups.json"), rs)
}
func lockUpdate(c Config) (*os.File, error) {
	root := filepath.Join(c.Data, "updates")
	if e := os.MkdirAll(root, 0700); e != nil {
		return nil, e
	}
	if e := platform.ProtectMaintenanceDir(root, c.Data); e != nil {
		return nil, e
	}
	return platform.Lock(filepath.Join(root, "update.lock"))
}

// Stop only processes proved to belong to this installation. Port ownership is
// checked again after both the manager job and the persistent shared service stop.
func stopUpdateRuntime(ctx context.Context, c Config) error {
	lock, e := platform.Lock(filepath.Join(c.Data, "runtime.lock"))
	if e != nil {
		var run struct {
			PID int `json:"pid"`
		}
		if e = platform.LoadJSON(filepath.Join(c.Data, "running.json"), &run); e != nil || !platform.ProcessRunning(run.PID, filepath.Join(c.Data, "Workspace.exe")) {
			return errors.New("Workspace is starting or another operation holds this installation; retry when it is ready")
		}
		if e = Stop(ctx, c); e != nil {
			return e
		}
		deadline := time.Now().Add(45 * time.Second)
		for time.Now().Before(deadline) {
			lock, e = platform.Lock(filepath.Join(c.Data, "runtime.lock"))
			if e == nil && !platform.ProcessRunning(run.PID, filepath.Join(c.Data, "Workspace.exe")) {
				break
			}
			if lock != nil {
				lock.Close()
				lock = nil
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(200 * time.Millisecond):
			}
		}
		if lock == nil {
			return errors.New("Workspace processes have not stopped; no update files were changed")
		}
	}
	lock.Close()
	root := filepath.Join(c.Data, "shared")
	if run, e := shared.ReadRunning(root); e == nil && platform.ProcessRunning(run.PID, filepath.Join(root, "WorkspaceServices.exe")) {
		cfg, _, e := shared.Load(root)
		if e != nil || cfg.EnginePort != c.Port+11 {
			return errors.New("shared service does not match this installation")
		}
		if e = shared.Stop(ctx, root); e != nil {
			return e
		}
	}
	deadline := time.Now().Add(15 * time.Second)
	for {
		e = CheckPorts(c)
		if e == nil {
			return nil
		}
		if time.Now().After(deadline) {
			return e
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(200 * time.Millisecond):
		}
	}
}
func startUpdatedRuntime(ctx context.Context, c Config, token string) (result error) {
	log := &platform.Log{Path: filepath.Join(c.Data, "updates", "restart.log")}
	cmd := exec.Command(filepath.Join(c.Data, "Workspace.exe"), "--data", c.Data, "run")
	cmd.Dir = c.Data
	cmd.Stdout = log
	cmd.Stderr = log
	cmd.Env = os.Environ()
	if token != "" {
		cmd.Env = append(cmd.Env, "WORKSPACE_UPDATE_TRANSACTION="+token)
	}
	if e := platform.StartDetached(cmd); e != nil {
		return e
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	defer func() {
		if result != nil {
			_ = cmd.Process.Kill()
			select {
			case <-done:
			case <-time.After(10 * time.Second):
			}
		}
	}()
	if e := WaitHTTPDetailed(ctx, fmt.Sprintf("http://localhost:%d/healthz", c.Port), 8*time.Minute, done, nil); e != nil {
		return e
	}
	// Run writes the ready marker after all app checks and the portal have passed.
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		var progress Progress
		if platform.LoadJSON(filepath.Join(c.Data, "startup-progress.json"), &progress) == nil && progress.PID == cmd.Process.Pid && progress.State == "ready" {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return errors.New("Workspace did not confirm all startup checks")
}

type UpdateRequest struct {
	Selected   []string `json:"selected"`
	BackupRoot string   `json:"backupRoot"`
}

type updateOperations struct {
	preflight func(context.Context, Config) error
	stop      func(context.Context, Config) error
	start     func(context.Context, Config, string) error
}

// Check process creation from the updater itself before taking working apps
// offline. "version" reads no private data and starts no app processes.
func preflightUpdateRestart(ctx context.Context, c Config) error {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, filepath.Join(c.Data, "Workspace.exe"), "version")
	cmd.Dir = c.Data
	if e := platform.StartDetached(cmd); e != nil {
		return fmt.Errorf("Windows could not launch Workspace; your apps have not been stopped. Check executable permissions or Windows security protection: %w", e)
	}
	if e := cmd.Wait(); e != nil {
		return fmt.Errorf("Workspace launch check failed; your apps have not been stopped: %w", e)
	}
	return nil
}

var nativeUpdateOperations = updateOperations{preflight: preflightUpdateRestart, stop: stopUpdateRuntime, start: startUpdatedRuntime}

func verifyInstalledPrograms(c Config) error {
	if p, ok := c.Catalog.Packages["manager"]; !ok || p.File != "Workspace.exe" {
		return errors.New("unsupported manager location")
	}
	for _, p := range c.Catalog.Packages {
		path, e := platform.Inside(c.Programs, p.File)
		if e != nil {
			return e
		}
		held, e := platform.LockedBinary(path, p.SHA256)
		if e != nil {
			return fmt.Errorf("installed package %s is not verified: %w", p.File, e)
		}
		st, e := held.Stat()
		held.Close()
		if e != nil || st.Size() != p.Size {
			return errors.New("installed package size mismatch")
		}
	}
	p, ok := c.Catalog.Packages["manager"]
	if !ok {
		return errors.New("installed manager package missing")
	}
	held, e := platform.LockedBinary(filepath.Join(c.Data, "Workspace.exe"), p.SHA256)
	if e != nil {
		return e
	}
	held.Close()
	return nil
}
func cacheUpdateController(c Config, target catalog.Catalog, stage string) error {
	source := filepath.Join(c.Data, "Workspace.exe")
	p := c.Catalog.Packages["manager"]
	if !distribution.UpdateControllerSupported(c.Version) {
		p = target.Packages["manager"]
		source = filepath.Join(stage, p.File)
	}
	held, e := platform.LockedBinary(source, p.SHA256)
	if e != nil {
		return e
	}
	held.Close()
	dir := filepath.Join(c.Data, "updates", "controller")
	if e = os.MkdirAll(dir, 0700); e != nil {
		return e
	}
	if e = writeReplace(source, filepath.Join(dir, "Workspace.exe")); e != nil {
		return e
	}
	return platform.SaveJSON(filepath.Join(dir, "package.json"), p)
}
func installRecoveryGuard(c Config) error {
	dir := filepath.Join(c.Data, "updates", "controller")
	var p catalog.Package
	if e := platform.LoadJSON(filepath.Join(dir, "package.json"), &p); e != nil {
		return fmt.Errorf("verified update controller is missing: %w", e)
	}
	if p.File != "Workspace.exe" {
		return errors.New("invalid recovery controller package")
	}
	source := filepath.Join(dir, "Workspace.exe")
	held, e := platform.LockedBinary(source, p.SHA256)
	if e != nil {
		return e
	}
	held.Close()
	if e = writeReplace(source, filepath.Join(c.Data, "Workspace.exe")); e != nil {
		return e
	}
	return writeReplace(source, filepath.Join(c.Programs, "Workspace.exe"))
}

// ApplyUpdates stages and verifies every selected byte before downtime, then
// snapshots the entire installation. Its durable journal blocks ordinary startup
// until either health-checked commit or a verified rollback completes.
func ApplyUpdates(ctx context.Context, data string, plan distribution.UpdatePlan, request UpdateRequest, report func(string)) error {
	c, e := LoadUpdateConfig(data)
	if e != nil {
		return e
	}
	_, s, e := Load(data)
	if e != nil {
		return e
	}
	if report == nil {
		report = func(string) {}
	}
	if e = updateRoots(c); e != nil {
		return e
	}
	guard, e := lockUpdate(c)
	if e != nil {
		return e
	}
	defer guard.Close()
	if UpdatePending(c) {
		return errors.New("recover the unfinished update before applying another")
	}
	target, sources, e := distribution.SelectUpdates(c.Catalog, plan, request.Selected)
	if e != nil {
		return e
	}
	stage, e := os.MkdirTemp(filepath.Join(c.Data, "updates"), "stage-")
	if e != nil {
		return e
	}
	defer os.RemoveAll(stage)
	if e = distribution.StageUpdateFiles(ctx, target, sources, stage, report); e != nil {
		return e
	}
	return applyStagedUpdate(ctx, c, s, target, filepath.Join(stage, "packages"), request.BackupRoot, report)
}

// UpdateFromBundle is the explicit offline path, and bootstraps an installation
// made by older managers. It never adopts an unverified archive or removes apps.
func UpdateFromBundle(ctx context.Context, data, bundle, backupRoot string, selected []string, report func(string)) error {
	c, e := LoadUpdateConfig(data)
	if e != nil {
		return e
	}
	_, s, e := Load(data)
	if e != nil {
		return e
	}
	if report == nil {
		report = func(string) {}
	}
	if e = updateRoots(c); e != nil {
		return e
	}
	guard, e := lockUpdate(c)
	if e != nil {
		return e
	}
	defer guard.Close()
	if UpdatePending(c) {
		return errors.New("recover the unfinished update first")
	}
	z, incoming, e := distribution.OpenUpdateBundle(bundle)
	if e != nil {
		return e
	}
	defer z.Close()
	target, e := distribution.SelectBundleUpdates(c.Catalog, incoming, selected)
	if e != nil {
		return e
	}
	partial := catalog.Catalog{Schema: 1, Version: target.Version, Packages: map[string]catalog.Package{}}
	for id, p := range target.Packages {
		if c.Catalog.Packages[id] != p {
			partial.Packages[id] = p
		}
	}
	stage, e := os.MkdirTemp(filepath.Join(c.Data, "updates"), "stage-")
	if e != nil {
		return e
	}
	defer os.RemoveAll(stage)
	if e = catalog.Extract(&z.Reader, partial, stage); e != nil {
		return e
	}
	return applyStagedUpdate(ctx, c, s, target, stage, backupRoot, report)
}
func applyStagedUpdate(ctx context.Context, c Config, s Secrets, target catalog.Catalog, stage, backupRoot string, report func(string)) (result error) {
	return applyStagedWith(ctx, c, s, target, stage, backupRoot, report, nativeUpdateOperations)
}
func applyStagedWith(ctx context.Context, c Config, s Secrets, target catalog.Catalog, stage, backupRoot string, report func(string), ops updateOperations) (result error) {
	if e := verifyInstalledPrograms(c); e != nil {
		return e
	}
	for _, app := range target.Apps {
		old, ok := c.Catalog.Find(app.ID)
		if !ok || old.Callback != app.Callback || old.OIDC != app.OIDC {
			return errors.New("update changes an identity callback and needs an explicit migration")
		}
		for key, value := range old.Env {
			if (strings.Contains(value, "${DATA}") || strings.Contains(value, "${FILES}") || key == "WORKSPACE_MAIL_BACKEND") && app.Env[key] != value {
				return fmt.Errorf("%s changes a persistent storage binding and needs an explicit migration", app.Name)
			}
		}
		for _, arg := range old.Args {
			if strings.Contains(arg, "${DATA}") || strings.Contains(arg, "${FILES}") {
				retained := false
				for _, nextArg := range app.Args {
					if arg == nextArg {
						retained = true
						break
					}
				}
				if !retained {
					return fmt.Errorf("%s changes an existing data path and needs an explicit migration", app.Name)
				}
			}
		}
	}
	if !distribution.UpdateControllerSupported(c.Version) && target.Version == c.Version {
		return errors.New("update One Workspace first, or select it together with app updates")
	}
	if target.Packages["manager"].File != "Workspace.exe" {
		return errors.New("manager update changes its executable location")
	}
	if ops.preflight != nil {
		report("Checking Windows can restart Workspace before stopping your apps")
		if e := ops.preflight(ctx, c); e != nil {
			return e
		}
	}
	if e := cacheUpdateController(c, target, stage); e != nil {
		return e
	}
	before, _, e := fileDigest(filepath.Join(c.Data, "workspace.json"))
	if e != nil {
		return e
	}
	report("Stopping this Workspace and its apps for a consistent backup")
	if e = ops.stop(ctx, c); e != nil {
		return e
	}
	transactionStarted := false
	defer func() {
		if !transactionStarted && result != nil {
			_ = ops.start(context.Background(), c, "")
		}
	}()
	lock, e := platform.Lock(filepath.Join(c.Data, "runtime.lock"))
	if e != nil {
		return e
	}
	defer func() {
		if lock != nil {
			lock.Close()
		}
	}()
	after, _, e := fileDigest(filepath.Join(c.Data, "workspace.json"))
	if e != nil || before != after {
		return errors.New("installation changed while preparing the update; check again")
	}
	backup, e := makeUpdateBackup(ctx, c, backupRoot, report)
	if e != nil {
		lock.Close()
		lock = nil
		return e
	}
	if e = recordBackup(c, backup); e != nil {
		return e
	}
	j := UpdateJournal{Phase: "applying", Backup: backup, Token: platform.Token()}
	if e = saveUpdateJournal(c, &j); e != nil {
		return e
	}
	transactionStarted = true
	committed := false
	defer func() {
		if committed {
			return
		}
		if lock != nil {
			lock.Close()
			lock = nil
		}
		report("Update did not finish; restoring the verified backup")
		if rollbackErr := recoverUpdateWith(context.Background(), c, &j, report, ops); rollbackErr != nil {
			if errors.Is(rollbackErr, errRestoredNeedsStart) {
				result = fmt.Errorf("update failed: %v; %w", result, rollbackErr)
			} else {
				result = fmt.Errorf("%v; automatic recovery needs attention: %w. Open Workspace updates and choose Recover", result, rollbackErr)
			}
		} else {
			result = fmt.Errorf("update was rolled back; original data and programs restored: %w", result)
		}
	}()
	next := c
	next.Catalog = target
	next.Version = target.Version
	// Callbacks are part of identity configuration. Changing them requires a
	// publisher migration rather than silently invalidating existing OIDC clients.
	report("Installing verified program files; retaining existing locations and identities")
	ids := make([]string, 0, len(target.Packages))
	for id := range target.Packages {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool {
		if ids[i] == "manager" {
			return true
		}
		if ids[j] == "manager" {
			return false
		}
		return ids[i] < ids[j]
	})
	for _, id := range ids {
		p := target.Packages[id]
		if c.Catalog.Packages[id] == p {
			continue
		}
		src, e := platform.Inside(stage, p.File)
		if e != nil {
			return e
		}
		held, e := platform.LockedBinary(src, p.SHA256)
		if e != nil {
			return e
		}
		held.Close()
		dst, e := platform.Inside(c.Programs, p.File)
		if e != nil {
			return e
		}
		if e = writeReplace(src, dst); e != nil {
			return e
		}
		if id == "manager" {
			if e = writeReplace(src, filepath.Join(c.Data, "Workspace.exe")); e != nil {
				return e
			}
		}
	}
	if e = Save(next, s); e != nil {
		return e
	}
	j.Phase = "testing"
	if e = saveUpdateJournal(c, &j); e != nil {
		return e
	}
	lock.Close()
	lock = nil
	report("Restarting and checking every app before completing the update")
	if e = ops.start(ctx, next, j.Token); e != nil {
		return e
	}
	j.Phase = "complete"
	j.Token = ""
	if e = saveUpdateJournal(c, &j); e != nil {
		return e
	}
	committed = true
	report("Update complete. All apps are healthy. Verified backup: " + backup.Path)
	return nil
}
func RecoverUpdate(ctx context.Context, data string, report func(string)) error {
	c, e := LoadUpdateConfig(data)
	if e != nil {
		return e
	}
	guard, e := lockUpdate(c)
	if e != nil {
		return e
	}
	defer guard.Close()
	var j UpdateJournal
	if e = platform.LoadJSON(updateJournalPath(c), &j); e != nil {
		return e
	}
	if j.Phase == "complete" || j.Phase == "rolled-back" {
		return errors.New("there is no unfinished update to recover")
	}
	return recoverUpdate(ctx, c, &j, report)
}
func recoverUpdate(ctx context.Context, c Config, j *UpdateJournal, report func(string)) error {
	return recoverUpdateWith(ctx, c, j, report, nativeUpdateOperations)
}
func recoverUpdateWith(ctx context.Context, c Config, j *UpdateJournal, report func(string), ops updateOperations) error {
	if report == nil {
		report = func(string) {}
	}
	if _, e := verifyUpdateBackup(c, j.Backup); e != nil {
		return e
	}
	if e := ops.stop(ctx, c); e != nil {
		return e
	}
	lock, e := platform.Lock(filepath.Join(c.Data, "runtime.lock"))
	if e != nil {
		return e
	}
	defer func() {
		if lock != nil {
			lock.Close()
		}
	}()
	if j.Recovery == nil {
		recovery, e := makeUpdateBackup(ctx, c, filepath.Dir(j.Backup.Path), report)
		if e != nil {
			return e
		}
		if j.Phase != "rollback-requested" {
			recovery, e = markRecoveryOnly(c, recovery)
			if e != nil {
				return e
			}
		}
		if e = recordBackup(c, recovery); e != nil {
			return e
		}
		j.Recovery = &recovery
	}
	j.Phase = "restoring"
	if e = saveUpdateJournal(c, j); e != nil {
		return e
	}
	if e = installRecoveryGuard(c); e != nil {
		return e
	}
	if e = restoreUpdateBackup(c, j.Backup, *j.Recovery, report); e != nil {
		return e
	}
	restored, _, e := Load(c.Data)
	if e != nil {
		return e
	}
	j.Phase = "rolled-back"
	j.Token = ""
	if e = saveUpdateJournal(c, j); e != nil {
		return e
	}
	lock.Close()
	lock = nil
	if e = ops.start(ctx, restored, ""); e != nil {
		return fmt.Errorf("%w: %w. Open %s to start the restored installation, then use a corrected updater for the next attempt. Keep your backups; Recover is not needed after this completed rollback", errRestoredNeedsStart, e, filepath.Join(c.Data, "Workspace.exe"))
	}
	report("Original installation restored. Newer state is retained in " + j.Recovery.Path)
	return nil
}
func RollbackUpdate(ctx context.Context, data, id string, report func(string)) error {
	c, e := LoadUpdateConfig(data)
	if e != nil {
		return e
	}
	if e = updateRoots(c); e != nil {
		return e
	}
	guard, e := lockUpdate(c)
	if e != nil {
		return e
	}
	defer guard.Close()
	if UpdatePending(c) {
		return errors.New("recover the unfinished update before choosing another backup")
	}
	rs, e := UpdateBackups(c)
	if e != nil {
		return e
	}
	var backup *BackupReceipt
	for i := range rs {
		if rs[i].ID == id {
			backup = &rs[i]
			break
		}
	}
	if backup == nil {
		return errors.New("choose a recorded backup for this installation")
	}
	if backup.RecoveryOnly {
		return errors.New("this archive contains an interrupted state; choose a complete installation backup")
	}
	if _, e = verifyUpdateBackup(c, *backup); e != nil {
		return e
	}
	j := UpdateJournal{Phase: "rollback-requested", Backup: *backup}
	if e = saveUpdateJournal(c, &j); e != nil {
		return e
	}
	return recoverUpdate(ctx, c, &j, report)
}

func LaunchUpdates(c Config) error {
	root := filepath.Join(c.Data, "maintenance")
	if e := os.MkdirAll(root, 0700); e != nil {
		return e
	}
	if e := platform.ProtectMaintenanceDir(root, c.Data); e != nil {
		return e
	}
	if OpenExistingUpdates(c) {
		return nil
	}
	helper := filepath.Join(root, "WorkspaceUpdate.exe")
	current, e := os.Executable()
	if e != nil {
		return e
	}
	if !strings.EqualFold(filepath.Clean(current), filepath.Join(c.Data, "Workspace.exe")) {
		return errors.New("open updates from the protected installed manager")
	}
	if e = writeReplace(current, helper); e != nil {
		return e
	}
	cmd := exec.Command(helper, "--data", c.Data, "updates")
	cmd.Dir = root
	cmd.Stdout = &platform.Log{Path: filepath.Join(root, "window.log")}
	cmd.Stderr = cmd.Stdout
	if e = platform.StartDetached(cmd); e != nil {
		return e
	}
	go cmd.Wait()
	return nil
}

func OpenExistingUpdates(c Config) bool {
	var window struct {
		PID        int    `json:"pid"`
		Executable string `json:"executable"`
		URL        string `json:"url"`
	}
	if platform.LoadJSON(filepath.Join(c.Data, "updates", "window.json"), &window) == nil && platform.ProcessRunning(window.PID, window.Executable) {
		u, e := url.Parse(window.URL)
		if e == nil && u.Scheme == "http" && u.Hostname() == "localhost" && u.Port() != "" && u.Path == "/" && len(u.Fragment) == 64 {
			platform.OpenBrowser(window.URL)
			return true
		}
	}
	return false
}
