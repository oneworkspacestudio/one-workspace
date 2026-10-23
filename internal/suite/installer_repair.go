package suite

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/distribution"
	"github.com/oneworkspacestudio/one-workspace/internal/operationprogress"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

type InstallerRepairEvent struct {
	Phase    string                   `json:"phase,omitempty"`
	Progress *operationprogress.Event `json:"progress,omitempty"`
}

type InstallerRepairResult struct {
	Version        string            `json:"version"`
	Origin         string            `json:"origin"`
	Checked        time.Time         `json:"checked"`
	Apps           map[string]string `json:"apps"`
	NeedsAttention []string          `json:"needsAttention,omitempty"`
}

// ReinstallInstalled is an explicit CLI maintenance action, never discovery.
// It accepts only a complete, recognised installation and cannot reset accounts.
func ReinstallInstalled(ctx context.Context, data string, z *zip.Reader, cat catalog.Catalog, offline bool, report func(string)) (Config, error) {
	c, _, e := Load(filepath.Clean(data))
	if e != nil {
		return c, e
	}
	if !strings.EqualFold(c.Data, filepath.Clean(data)) {
		return c, errors.New("This installation belongs to another private data folder")
	}
	r := Request{Programs: c.Programs, Data: c.Data, Files: c.Files, Reinstall: true, Offline: offline}
	review, e := InspectSetupStorage(r)
	if e != nil {
		return c, e
	}
	if review.Kind != "repair" {
		return c, errors.New("Only a recognised Workspace can be repaired from Updates; use the setup wizard for incomplete installations")
	}
	r.ReinstallReview = review.Review
	return Install(ctx, r, z, cat, report)
}

// RunInstallerRepair hands execution to the exact reviewed, hash-verified setup.
// The setup runs its own newer catalog/transaction code; this controller does not
// attempt to interpret future application formats or replace itself while running.
func RunInstallerRepair(ctx context.Context, data string, plan distribution.UpdatePlan, report func(string)) (*InstallerRepairResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Hour)
	defer cancel()
	c, e := LoadUpdateConfig(data)
	if e != nil {
		return nil, e
	}
	if e = distribution.ValidateInstallerPlan(c.Catalog, plan); e != nil {
		return nil, e
	}
	guard, e := lockUpdate(c)
	if e != nil {
		return nil, e
	}
	defer func() {
		if guard != nil {
			guard.Close()
		}
	}()
	if UpdatePending(c) {
		return nil, errors.New("Recover the interrupted update before reinstalling")
	}
	root := filepath.Join(c.Data, "updates", "installers")
	if e = os.MkdirAll(root, 0700); e != nil {
		return nil, e
	}
	if e = platform.ProtectMaintenanceDir(root, c.Data); e != nil {
		return nil, e
	}
	stage, e := os.MkdirTemp(root, "verified-")
	if e != nil {
		return nil, e
	}
	defer os.RemoveAll(stage)
	path, e := distribution.DownloadInstaller(ctx, *plan.Installer, stage, report)
	if e != nil {
		return nil, e
	}
	current, e := LoadUpdateConfig(data)
	if e != nil {
		return nil, e
	}
	if e = distribution.ValidateInstallerPlan(current.Catalog, plan); e != nil {
		return nil, e
	}
	held, e := platform.LockedBinary(path, strings.TrimPrefix(plan.Installer.Download.Asset.Digest, "sha256:"))
	if e != nil {
		return nil, e
	}
	defer held.Close()
	log := &platform.Log{Path: filepath.Join(c.Data, "updates", "installer-repair.log")}
	output := &installerRepairOutput{ctx: ctx, report: report}
	cmd := exec.CommandContext(ctx, path, "--repair", c.Data, "--offline")
	cmd.Dir = stage
	cmd.Stdout = io.MultiWriter(log, output)
	cmd.Stderr = log
	operationprogress.Report(ctx, "stage", "Starting the verified replacement installer", 0, 0, "")
	// The new installer must own the transaction lock itself. It independently
	// rechecks paths/state under its setup and update locks before any downtime.
	guard.Close()
	guard = nil
	if e = platform.StartDetached(cmd); e != nil {
		return nil, fmt.Errorf("Windows could not start the verified installer; Workspace has not been stopped: %w", e)
	}
	if e = cmd.Wait(); e != nil {
		return nil, fmt.Errorf("The installer could not finish. Your recovery copies are retained. See %s", log.Path)
	}
	result, e := checkInstallerRepair(ctx, c, plan.Installer.Version)
	if e != nil {
		return nil, e
	}
	if e = platform.SaveJSON(filepath.Join(c.Data, "updates", "last-reinstall.json"), result); e != nil {
		return result, e
	}
	if len(result.NeedsAttention) > 0 {
		return result, fmt.Errorf("Workspace was reinstalled and restarted; these apps still need attention: %s. Open App status in Workspace to retry them", strings.Join(result.NeedsAttention, ", "))
	}
	return result, nil
}

type installerRepairOutput struct {
	ctx     context.Context
	report  func(string)
	pending []byte
}

func (p *installerRepairOutput) Write(b []byte) (int, error) {
	p.pending = append(p.pending, b...)
	for {
		i := bytes.IndexByte(p.pending, '\n')
		if i < 0 {
			break
		}
		line := p.pending[:i]
		var e InstallerRepairEvent
		if json.Unmarshal(line, &e) == nil {
			if e.Phase != "" && p.report != nil {
				p.report(e.Phase)
			}
			if e.Progress != nil {
				v := e.Progress
				operationprogress.Report(p.ctx, v.Stage, v.Detail, v.Done, v.Total, v.Unit)
			}
		}
		p.pending = p.pending[i+1:]
	}
	if len(p.pending) > 64<<10 {
		p.pending = nil
	}
	return len(b), nil
}

func checkInstallerRepair(ctx context.Context, c Config, version string) (*InstallerRepairResult, error) {
	// Read only the stable location header, so an older controller can report
	// completion even when it deliberately cannot read the new catalog format.
	var actual Config
	if e := platform.LoadJSON(filepath.Join(c.Data, "workspace.json"), &actual); e != nil {
		return nil, e
	}
	if actual.Version != version || actual.Owner != c.Owner || actual.Origin != c.Origin || actual.Port != c.Port || !strings.EqualFold(actual.Data, c.Data) || !strings.EqualFold(actual.Programs, c.Programs) || !strings.EqualFold(actual.Files, c.Files) {
		return nil, errors.New("Reinstall did not confirm the expected version and original account/storage settings")
	}
	r, e := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://127.0.0.1:%d/healthz", c.Port), nil)
	if e != nil {
		return nil, e
	}
	response, e := (&http.Client{Timeout: 15 * time.Second}).Do(r)
	if e != nil {
		return nil, errors.New("Reinstall finished, but Workspace could not be reached. Open its startup status before retrying")
	}
	defer response.Body.Close()
	return readInstallerHealth(response, c, version)
}
func readInstallerHealth(response *http.Response, c Config, version string) (*InstallerRepairResult, error) {
	var health struct {
		Status string            `json:"status"`
		Apps   map[string]string `json:"apps"`
	}
	if response.StatusCode != 200 || json.NewDecoder(io.LimitReader(response.Body, 64<<10)).Decode(&health) != nil || health.Status != "ok" || len(health.Apps) == 0 {
		return nil, errors.New("Workspace did not confirm its restarted app status")
	}
	for _, a := range c.Catalog.Apps {
		if a.Kind == "process" {
			if _, ok := health.Apps[a.ID]; !ok {
				health.Apps[a.ID] = "status missing"
			}
		}
	}
	if c.Tunnel {
		if _, ok := health.Apps["cloudflare"]; !ok {
			health.Apps["cloudflare"] = "status missing"
		}
	}
	result := &InstallerRepairResult{Version: version, Origin: c.Origin, Checked: time.Now().UTC(), Apps: health.Apps}
	for id, status := range health.Apps {
		if status != "running" {
			result.NeedsAttention = append(result.NeedsAttention, id)
		}
	}
	sort.Strings(result.NeedsAttention)
	return result, nil
}
