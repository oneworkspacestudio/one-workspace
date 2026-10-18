package suite

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/distribution"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

// InstallApps adds the verified suite apps to an existing stopped installation.
// Existing process binaries must match: changing a working app or migrating a
// mailbox is a separate operation. The manager, notices and new apps are backed
// up/copied, while identity stores and ordinary user files are retained.
func InstallApps(ctx context.Context, data, bundle string, report func(string)) (err error) {
	return InstallRequestedApps(ctx, data, bundle, nil, report)
}

// Explicit requested IDs re-enable a previous opt-out, preserving its data and
// existing OIDC client. All downloads still obey the verified bundle contract.
func InstallRequestedApps(ctx context.Context, data, bundle string, requested []string, report func(string)) (err error) {
	c, s, err := Load(data)
	if err != nil {
		return err
	}
	lock, err := platform.Lock(filepath.Join(c.Data, "runtime.lock"))
	if err != nil {
		return errors.New("stop Workspace before installing its additional apps")
	}
	defer lock.Close()
	z, err := zip.OpenReader(bundle)
	if err != nil {
		return err
	}
	defer z.Close()
	var incoming catalog.Catalog
	for _, f := range z.File {
		if f.Name == "catalog.json" {
			if incoming.Schema != 0 || f.UncompressedSize64 > 1<<20 {
				return errors.New("invalid app bundle catalog")
			}
			r, e := f.Open()
			if e != nil {
				return e
			}
			incoming, e = readAppCatalog(r)
			r.Close()
			if e != nil {
				return e
			}
		}
	}
	if err = distribution.CheckInstallerBundle(c.Catalog, incoming); err != nil {
		return err
	}
	selection := c.Catalog.EffectiveSelection()
	for _, id := range requested {
		if !foundationApp(id) {
			return errors.New("explicit reinstallation supports tasks, pdf, monitor and search")
		}
		if _, ok := incoming.Find(id); !ok {
			return fmt.Errorf("%s is not in this verified bundle", id)
		}
		kept := selection.Excluded[:0]
		for _, excluded := range selection.Excluded {
			if excluded != id {
				kept = append(kept, excluded)
			}
		}
		selection.Excluded = kept
	}
	request := Request{SelectedApps: []string{}}
	for _, app := range incoming.DefaultApps(selection) {
		request.SelectedApps = append(request.SelectedApps, app.ID)
	}
	incoming, err = SelectInstallApps(incoming, request)
	if err != nil {
		return err
	}
	incoming.Selection = selection
	if _, hadDrive := c.Catalog.Find("drive"); !hadDrive {
		if _, addingDrive := incoming.Find("drive"); addingDrive {
			if _, e := os.Lstat(filepath.Join(c.Data, "drive", "config.yaml")); !os.IsNotExist(e) {
				return errors.New("Drive already has private settings; migrate that data explicitly before adding Drive")
			}
		}
	}
	if err = incoming.Validate(); err != nil {
		return err
	}
	for _, app := range incoming.Apps {
		if !organizationSource(app.Source) {
			return fmt.Errorf("%s does not declare an organization repository source", app.ID)
		}
		if old, ok := c.Catalog.Find(app.ID); ok && old.Kind == "process" {
			a, b := c.Catalog.Packages[old.Package], incoming.Packages[app.Package]
			if old.Executable != app.Executable || a.SHA256 != b.SHA256 || a.Size != b.Size {
				return fmt.Errorf("%s differs from the installed version; preserve it and prepare an explicit app upgrade before installing this bundle", app.ID)
			}
		}
	}
	if err = catalog.Verify(&z.Reader, incoming); err != nil {
		return err
	}
	stage, err := os.MkdirTemp(c.Data, ".apps-stage-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(stage)
	// Stage only files that will be copied. Re-extracting an unchanged 570 MB
	// Drive on the protected data volume can exhaust a small Windows system disk.
	partial := catalog.Catalog{Schema: incoming.Schema, Version: incoming.Version, Packages: map[string]catalog.Package{}}
	for id, pkg := range incoming.Packages {
		old, ok := c.Catalog.Packages[id]
		if id == "manager" || !ok || old.SHA256 != pkg.SHA256 || old.Size != pkg.Size {
			partial.Packages[id] = pkg
		}
	}
	if err = catalog.Extract(&z.Reader, partial, stage); err != nil {
		return err
	}
	payloadPackages := make(map[string]catalog.Package, len(incoming.Packages))
	for id, p := range incoming.Packages {
		payloadPackages[id] = p
	}
	// Retain extensions in the configuration without pretending their files were
	// supplied by this bundle or changing their executables.
	for _, old := range c.Catalog.Apps {
		if _, ok := incoming.Find(old.ID); !ok {
			incoming.Apps = append(incoming.Apps, old)
			if old.Kind == "process" {
				incoming.Packages[old.Package] = c.Catalog.Packages[old.Package]
			}
		}
	}
	if err = incoming.Validate(); err != nil {
		return err
	}
	oldConfig := c
	c.Catalog, c.Version = incoming, incoming.Version
	ports := NewConfig(Request{Port: c.Port, Photos: true}, incoming).Catalog.Apps
	for i, app := range ports {
		if old, installed := oldConfig.Catalog.Find(app.ID); installed {
			ports[i].Port = old.Port
		}
	}
	c.Catalog.Apps = ports
	if err = CheckPorts(c); err != nil {
		return err
	}
	backup := filepath.Join(c.Data, "backups", "install-apps-"+time.Now().UTC().Format("20060102T150405.000000000"))
	if err = os.MkdirAll(backup, 0700); err != nil {
		return err
	}
	for _, name := range []string{"workspace.json", "secrets.dpapi"} {
		if err = platform.CopyFile(filepath.Join(c.Data, name), filepath.Join(backup, name)); err != nil {
			return err
		}
	}
	savedConfig, err := os.ReadFile(filepath.Join(backup, "workspace.json"))
	if err != nil {
		return err
	}
	savedSecrets, err := os.ReadFile(filepath.Join(backup, "secrets.dpapi"))
	if err != nil {
		return err
	}
	if report == nil {
		report = func(string) {}
	}
	report("Connecting the additional apps to MyAccounts")
	if err = connectAddedApps(ctx, c, &s); err != nil {
		return err
	}
	type change struct {
		target, saved string
		existed       bool
	}
	var changes []change
	committed := false
	defer func() {
		if committed {
			return
		}
		for i := len(changes) - 1; i >= 0; i-- {
			item := changes[i]
			_ = os.Remove(item.target)
			if item.existed {
				_ = platform.CopyFile(item.saved, item.target)
			}
		}
		_ = platform.AtomicWrite(filepath.Join(c.Data, "workspace.json"), savedConfig)
		_ = platform.AtomicWrite(filepath.Join(c.Data, "secrets.dpapi"), savedSecrets)
	}()
	for id, pkg := range payloadPackages {
		target := filepath.Join(c.Programs, filepath.FromSlash(pkg.File))
		if old, ok := oldConfig.Catalog.Packages[id]; ok && old.SHA256 == pkg.SHA256 && old.Size == pkg.Size {
			held, e := platform.LockedBinary(target, pkg.SHA256)
			if e != nil {
				return e
			}
			held.Close()
			continue
		}
		saved := filepath.Join(backup, "programs", filepath.FromSlash(pkg.File))
		_, e := os.Stat(target)
		exists := e == nil
		if e != nil && !os.IsNotExist(e) {
			return e
		}
		if exists {
			if err = platform.CopyFile(target, saved); err != nil {
				return err
			}
		}
		changes = append(changes, change{target, saved, exists})
		if exists {
			if err = os.Remove(target); err != nil {
				return err
			}
		}
		if err = platform.CopyFile(filepath.Join(stage, filepath.FromSlash(pkg.File)), target); err != nil {
			return err
		}
	}
	manager := filepath.Join(c.Data, "Workspace.exe")
	saved := filepath.Join(backup, "Workspace.exe")
	if err = platform.CopyFile(manager, saved); err != nil {
		return err
	}
	changes = append(changes, change{manager, saved, true})
	if err = os.Remove(manager); err != nil {
		return err
	}
	if err = platform.CopyFile(filepath.Join(stage, "Workspace.exe"), manager); err != nil {
		return err
	}
	if _, hadDrive := oldConfig.Catalog.Find("drive"); !hadDrive {
		if _, hasDrive := c.Catalog.Find("drive"); hasDrive {
			changes = append(changes, change{target: filepath.Join(c.Data, "drive", "config.yaml")})
			if err = writeDrive(c, s); err != nil {
				return err
			}
		}
	}
	if err = Save(c, s); err != nil {
		return err
	}
	committed = true
	report("Installed all available suite apps. Existing files, accounts and Mail storage were preserved. Backup: " + backup)
	return nil
}

func readAppCatalog(r io.Reader) (catalog.Catalog, error) {
	data, err := io.ReadAll(io.LimitReader(r, 1<<20))
	if err != nil {
		return catalog.Catalog{}, err
	}
	return catalog.Read(data)
}
func organizationSource(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && u.Scheme == "https" && u.Host == "github.com" && u.User == nil && u.RawQuery == "" && u.Fragment == "" && strings.HasPrefix(u.Path, "/oneworkspacestudio/") && len(strings.Split(strings.Trim(u.Path, "/"), "/")) >= 2
}

func connectAddedApps(ctx context.Context, c Config, s *Secrets) (result error) {
	job, err := platform.NewJob()
	if err != nil {
		return err
	}
	defer func() {
		stop, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if e := job.Stop(stop); e != nil && result == nil {
			result = fmt.Errorf("stop temporary account provisioning: %w", e)
		}
		job.Close()
	}()
	rt := &Runtime{Config: c, Secrets: *s, job: job, status: map[string]string{}}
	a, _ := c.Catalog.Find("accounts")
	held, err := platform.LockedBinary(filepath.Join(c.Programs, a.Executable), c.Catalog.Packages[a.Package].SHA256)
	if err != nil {
		return err
	}
	defer held.Close()
	if err = prepareRuntimeTemp(c, a.ID); err != nil {
		return err
	}
	cmd := rt.command(a)
	if err = job.Start(cmd); err != nil {
		return err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	defer func() {
		_ = cmd.Process.Kill()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
		}
	}()
	if err = WaitHTTPDetailed(ctx, fmt.Sprintf("http://localhost:%d/healthz", a.Port), 90*time.Second, done, nil); err != nil {
		return err
	}
	api := newIdentity(c)
	if err = api.login(ctx, *s); err != nil {
		return err
	}
	if s.Clients == nil {
		s.Clients = map[string]Client{}
	}
	if _, hasDrive := c.Catalog.Find("drive"); hasDrive && s.Clients["drive"].ID == "" {
		client, e := registerClient(ctx, api, "drive", "Workspace Drive and Photos", c.Origin+"/drive/api/auth/oidc/callback")
		if e != nil {
			return e
		}
		s.Clients["drive"] = client
	}
	for _, app := range c.Catalog.Apps {
		if !app.OIDC || usesShared(c, app.ID) || s.Clients[app.ID].ID != "" {
			continue
		}
		client, e := registerClient(ctx, api, app.ID, app.Name, c.Origin+app.Callback)
		if e != nil {
			return fmt.Errorf("register %s: %w", app.ID, e)
		}
		s.Clients[app.ID] = client
	}
	if SharedEnabled(c) && s.Clients["shared"].ID == "" {
		client, e := registerSharedClient(ctx, api, c)
		if e != nil {
			return e
		}
		s.Clients["shared"] = client
	}
	for _, app := range c.Catalog.Apps {
		if usesShared(c, app.ID) {
			s.Clients[app.ID] = s.Clients["shared"]
		}
	}
	return nil
}
