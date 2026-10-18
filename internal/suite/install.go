package suite

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/distribution"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func Install(ctx context.Context, r Request, z *zip.Reader, cat catalog.Catalog, progress func(string)) (Config, error) {
	var empty Config
	var overrides map[string]string
	if e := r.ResolveStorage(); e != nil {
		return empty, e
	}
	if _, e := r.PreflightStorage(); e != nil {
		return empty, e
	}
	if r.Mode == "automatic" {
		if !r.Authorize {
			return empty, errors.New("authorize automatic domain hosting before installation")
		}
		if r.CFCertFile == "" {
			var e error
			r.CFCertFile, e = CloudflareLoginPath()
			if e != nil {
				return empty, e
			}
		}
		progress("Discovering your domain and Cloudflare settings")
		h, e := ResolveHosting(ctx, r.Domain, r.CFCertFile)
		if e != nil {
			return empty, e
		}
		r.Mode, r.Domain, r.Origin, r.Issuer = h.Mode, h.Domain, h.Origin, h.Issuer
		r.CFToken, r.AccountID, r.ZoneID, r.Authorize = h.CFToken, h.AccountID, h.ZoneID, true
	}
	if e := r.Validate(); e != nil {
		return empty, e
	}
	// Record exclusions against the apps actually offered by this installer.
	// A newer online bundle can contain apps absent from its embedded app picker.
	offered, selectionErr := SelectInstallApps(cat, r)
	if selectionErr != nil {
		return empty, selectionErr
	}
	if !r.Offline {
		downloadRoot := filepath.Dir(r.Programs)
		if err := os.MkdirAll(downloadRoot, 0700); err != nil {
			return empty, err
		}
		bundle, err := distribution.Latest(ctx, progress, downloadRoot)
		if err != nil {
			return empty, err
		}
		defer bundle.Close()
		if err = distribution.CheckInstallerBundle(cat, bundle.Catalog); err != nil {
			return empty, err
		}
		z, cat = bundle.ZIP, bundle.Catalog
		overrides = bundle.Overrides
	} else {
		progress("Using this installer’s bundled versions (offline mode)")
	}
	r.SelectedApps = []string{}
	for _, app := range cat.DefaultApps(offered.Selection) {
		r.SelectedApps = append(r.SelectedApps, app.ID)
	}
	selected, selectionErr := SelectInstallApps(cat, r)
	if selectionErr != nil {
		return empty, selectionErr
	}
	selected.Selection = offered.Selection
	c := NewConfig(r, selected)
	if e := CheckPorts(c); e != nil {
		return empty, e
	}
	if e := platform.RequireNewDir(c.Data); e != nil {
		return empty, e
	}
	if e := platform.ProtectDir(c.Data); e != nil {
		return empty, e
	}
	lock, e := platform.Lock(filepath.Join(c.Data, "install.lock"))
	if e != nil {
		return empty, e
	}
	defer lock.Close()
	if e = platform.RequireNewDir(c.Programs); e != nil {
		return empty, e
	}
	if e = os.MkdirAll(c.Files, 0700); e != nil {
		return empty, e
	}
	for _, folder := range []string{"Documents", "Photos", "Music"} {
		if e = os.MkdirAll(filepath.Join(c.Files, folder), 0700); e != nil {
			return empty, e
		}
	}
	state := func(phase string) {
		progress(phase)
		_ = platform.SaveJSON(filepath.Join(c.Data, "install-state.json"), map[string]string{"phase": phase, "version": c.Version})
	}
	state("Verifying and extracting Workspace " + cat.Version + " applications")
	if e = catalog.Extract(z, c.Catalog, c.Programs, overrides); e != nil {
		return empty, e
	}
	// Keep the trusted manager and catalog on the protected data volume. Large
	// component packages remain at the selected program location, hash-locked at launch.
	manager := filepath.Join(c.Data, "Workspace.exe")
	if e = platform.CopyFile(filepath.Join(c.Programs, "Workspace.exe"), manager); e != nil {
		return empty, e
	}
	s := Secrets{AdminEmail: r.Email, AdminPassword: platform.Token(), RecoveryPassword: platform.Token(), Clients: map[string]Client{}}
	if e = platform.SaveSecret(filepath.Join(c.Data, "secrets.dpapi"), s); e != nil {
		return empty, e
	}
	if e = os.MkdirAll(filepath.Join(c.Data, "accounts", "secrets"), 0700); e != nil {
		return empty, e
	}
	a, _ := c.Catalog.Find("accounts")
	if e = platform.SaveJSON(filepath.Join(c.Data, "accounts", "config.json"), map[string]any{"version": 1, "issuer": c.Issuer, "port": a.Port, "adminEmail": r.Email, "adminName": "Workspace server recovery", "tunnelMode": "manual", "sharedEmail": false}); e != nil {
		return empty, e
	}
	if e = platform.SaveSecret(filepath.Join(c.Data, "accounts", "secrets", "bootstrap.dpapi"), map[string]string{"name": "Workspace server recovery", "email": r.Email, "password": s.AdminPassword}); e != nil {
		return empty, e
	}
	state("Creating Accounts and your Workspace owner")
	job, e := platform.NewJob()
	if e != nil {
		return empty, e
	}
	defer job.Close()
	exe := filepath.Join(c.Programs, filepath.FromSlash(a.Executable))
	held, e := platform.LockedBinary(exe, c.Catalog.Packages[a.Package].SHA256)
	if e != nil {
		return empty, e
	}
	defer held.Close()
	cmd := exec.Command(exe, "-data", filepath.Join(c.Data, "accounts"), "-no-browser")
	cmd.Dir = c.Data
	if e = prepareRuntimeTemp(c, "accounts"); e != nil {
		return empty, e
	}
	cmd.Env = append(os.Environ(), "TEMP="+runtimeTemp(c, "accounts"), "TMP="+runtimeTemp(c, "accounts"))
	log := &platform.Log{Path: filepath.Join(c.Data, "install-accounts.log")}
	cmd.Stdout = log
	cmd.Stderr = log
	if e = job.Start(cmd); e != nil {
		return empty, e
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	defer func() {
		_ = cmd.Process.Kill()
		select {
		case <-done:
		case <-time.After(3 * time.Second):
		}
	}()
	if e = WaitHTTP(ctx, fmt.Sprintf("http://localhost:%d/healthz", a.Port), 90*time.Second); e != nil {
		return empty, e
	}
	if e = ConfigureIdentity(ctx, &c, &s, r); e != nil {
		return empty, e
	}
	stop := exec.Command(exe, "-data", filepath.Join(c.Data, "accounts"), "stop")
	stop.Dir = c.Data
	_ = stop.Run()
	select {
	case <-done:
		done <- nil
	case <-time.After(10 * time.Second):
		return empty, errors.New("Accounts did not stop after configuration")
	}
	state("Connecting your selected Workspace apps")
	if e = writeDrive(c, s); e != nil {
		return empty, e
	}
	if r.Mode == "cloudflare" {
		state("Creating the authorized Cloudflare tunnel and exact DNS records")
		s.TunnelToken, e = ProvisionCloudflare(ctx, r, c)
		if e != nil {
			return empty, e
		}
		s.HostingToken = r.CFToken
		if c.HostingDomain != "" {
			s.MailDNSToken = r.CFToken
		}
	}
	if e = Save(c, s); e != nil {
		return empty, e
	}
	if e = platform.Autostart(manager, c.Data, r.Autostart); e != nil {
		return empty, e
	}
	if e = platform.Shortcut(manager, c.Data); e != nil {
		return empty, fmt.Errorf("apps installed but Start Menu shortcut failed: %w", e)
	}
	// The program-side pointer contains no credentials and enables a normal launch.
	if e = platform.SaveJSON(filepath.Join(c.Programs, "workspace-location.json"), map[string]string{"data": c.Data}); e != nil {
		return empty, e
	}
	state("Ready to launch Workspace")
	return c, nil
}
func writeDrive(c Config, s Secrets) error {
	app, ok := c.Catalog.Find("drive")
	if !ok {
		return nil
	}
	if e := os.MkdirAll(filepath.Join(c.Data, "drive"), 0700); e != nil {
		return e
	}
	config := map[string]any{"http": map[string]any{"listen": "127.0.0.1", "port": app.Port, "baseURL": "/drive", "trustProxyHeaders": true}, "server": map[string]any{"database": map[string]string{"path": filepath.Join(c.Data, "drive", "settings.sqlite")}, "cacheDir": filepath.Join(c.Data, "drive", "cache"), "disableUpdateCheck": true, "sources": []any{map[string]any{"path": c.Files, "name": "My Drive", "config": map[string]any{"defaultEnabled": true, "defaultPermissions": map[string]bool{"view": true, "download": true, "modify": true, "create": true, "delete": true}}}}}, "frontend": map[string]any{"name": "Workspace Drive", "oidcLoginButtonText": "Sign in with MyAccounts"}, "auth": map[string]any{"methods": map[string]any{"noauth": false, "password": map[string]any{"enabled": true, "signup": false, "adminUsername": "workspace-recovery"}, "oidc": map[string]any{"enabled": true, "issuerUrl": c.Issuer, "clientId": s.Clients["drive"].ID, "scopes": "openid profile email", "userIdentifier": "sub", "disableVerifyTLS": false, "adminGroup": ""}}}, "userDefaults": map[string]any{"account": map[string]any{"permissions": map[string]bool{"admin": false, "api": false}}}}
	return platform.SaveJSON(filepath.Join(c.Data, "drive", "config.yaml"), config)
}
func WaitHTTP(ctx context.Context, target string, timeout time.Duration) error {
	return WaitHTTPDetailed(ctx, target, timeout, nil, nil)
}

// Report observed reachability, and stop waiting immediately if the child exits.
func WaitHTTPDetailed(ctx context.Context, target string, timeout time.Duration, exit <-chan error, progress func(string)) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	client := &http.Client{Timeout: 2 * time.Second}
	last := "No successful response yet"
	reported := ""
	lastReport := time.Time{}
	for {
		select {
		case e := <-exit:
			if e == nil {
				return errors.New("application exited before it became ready")
			}
			return fmt.Errorf("application exited before it became ready: %w", e)
		default:
		}
		req, e := http.NewRequestWithContext(ctx, "GET", target, nil)
		if e != nil {
			return e
		}
		if resp, e := client.Do(req); e == nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				return nil
			}
			last = fmt.Sprintf("HTTP %d from %s", resp.StatusCode, target)
		} else if ctx.Err() == nil {
			last = "Not reachable yet: " + target
		}
		if progress != nil && (last != reported || time.Since(lastReport) > 3*time.Second) {
			progress(last)
			reported = last
			lastReport = time.Now()
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("health check timed out after %s; last result: %s", timeout, last)
		case e := <-exit:
			if e == nil {
				return errors.New("application exited before it became ready")
			}
			return fmt.Errorf("application exited before it became ready: %w", e)
		case <-time.After(500 * time.Millisecond):
		}
	}
}
