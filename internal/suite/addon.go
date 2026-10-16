package suite

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fahimc/one-workspace/internal/catalog"
	"github.com/fahimc/one-workspace/internal/platform"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Addon struct {
	Schema  int             `json:"schema"`
	App     catalog.App     `json:"app"`
	Package catalog.Package `json:"package"`
}

// AddApp is an explicit local-owner operation. App packages execute with the
// installing user's Windows privileges; hashes establish integrity, not trust.
func AddApp(ctx context.Context, data, manifest, binary string) error {
	c, s, e := Load(data)
	if e != nil {
		return e
	}
	lock, e := platform.Lock(filepath.Join(c.Data, "runtime.lock"))
	if e != nil {
		return errors.New("stop Workspace before adding an app")
	}
	defer lock.Close()
	f, e := os.Open(manifest)
	if e != nil {
		return e
	}
	defer f.Close()
	var add Addon
	decoder := json.NewDecoder(io.LimitReader(f, 1<<20))
	decoder.DisallowUnknownFields()
	if e = decoder.Decode(&add); e != nil {
		return e
	}
	if add.Schema != 1 {
		return errors.New("unsupported addon manifest")
	}
	if add.App.ID == "backups" || add.App.ID == "portal" || add.App.ID == "manager" || add.App.ID == "cloudflared" {
		return errors.New("reserved app ID")
	}
	if _, ok := c.Catalog.Find(add.App.ID); ok {
		return errors.New("app already exists; remove or migrate it explicitly before replacement")
	}
	if add.App.Kind != "process" {
		return errors.New("extensions must be process apps")
	}
	if add.App.OIDC && (add.App.Callback == "" || !strings.HasPrefix(add.App.Callback, strings.TrimRight(add.App.Route, "/")+"/") || strings.ContainsAny(add.App.Callback, "?#\\") || strings.Contains(add.App.Callback, "..")) {
		return errors.New("OIDC callback must be inside the app route")
	}
	for _, prefix := range []string{"/backups", "/updates", "/api", "/auth", "/_", "/branding", "/photos", "/drive", "/sdkjs", "/fonts", "/settings", "/tools", "/files", "/healthz"} {
		if strings.HasPrefix(add.App.Route, prefix) {
			return errors.New("app route conflicts with suite routes")
		}
	}
	key := "addon-" + add.App.ID
	rel := "apps/" + add.App.ID + "/" + add.App.Version + "/app.exe"
	add.Package.File = rel
	add.App.Package = key
	add.App.Executable = rel
	c.Catalog.Packages[key] = add.Package
	c.Catalog.Apps = append(c.Catalog.Apps, add.App)
	if e = c.Catalog.Validate(); e != nil {
		return e
	}
	for _, a := range c.Catalog.Apps {
		if a.Port == c.Port {
			return errors.New("app port conflicts with Workspace")
		}
	}
	// Validate version/path and the supplied digest before copying any executable.
	dest, e := platform.Inside(c.Programs, rel)
	if e != nil {
		return e
	}
	held, e := platform.LockedBinary(binary, add.Package.SHA256)
	if e != nil {
		return e
	}
	defer held.Close()
	if st, e := held.Stat(); e != nil || st.Size() != add.Package.Size {
		return errors.New("app package size mismatch")
	}
	if e = platform.CopyFile(binary, dest); e != nil {
		return e
	}
	if add.App.OIDC {
		if add.App.Callback == "" || !strings.HasPrefix(add.App.Callback, add.App.Route+"/") || strings.ContainsAny(add.App.Callback, "?#\\") || strings.Contains(add.App.Callback, "..") {
			return errors.New("OIDC callback must be inside the app route")
		}
		job, e := platform.NewJob()
		if e != nil {
			return e
		}
		defer job.Close()
		rt := &Runtime{Config: c, Secrets: s, job: job, status: map[string]string{}}
		app, _ := c.Catalog.Find("accounts")
		accountHeld, e := platform.LockedBinary(filepath.Join(c.Programs, app.Executable), c.Catalog.Packages[app.Package].SHA256)
		if e != nil {
			return e
		}
		defer accountHeld.Close()
		cmd := rt.command(app)
		if e = job.Start(cmd); e != nil {
			return e
		}
		go cmd.Wait()
		if e = WaitHTTP(ctx, fmt.Sprintf("http://localhost:%d/healthz", app.Port), 90*time.Second); e != nil {
			return e
		}
		api := newIdentity(c)
		if e = api.login(ctx, s); e != nil {
			return e
		}
		client, e := registerClient(ctx, api, add.App.ID, add.App.Name, c.Origin+add.App.Callback)
		if e != nil {
			return e
		}
		s.Clients[add.App.ID] = client
	}
	return Save(c, s)
}
