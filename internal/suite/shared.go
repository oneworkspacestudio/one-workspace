package suite

import (
	"context"
	"errors"
	"github.com/oneworkspacestudio/calendar/shared"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"os"
	"path/filepath"
	"time"
)

func SharedEnabled(c Config) bool {
	_, engine := c.Catalog.Packages["stalwart-shared"]
	_, manager := c.Catalog.Packages["services-shared"]
	return engine && manager
}
func sharedConsumer(id string) bool { return id == "calendar" || id == "contacts" || id == "mail" }

// Published hosting-capable Mail still owns its private engine. Merely adding
// Calendar must not switch Mail's identity, storage or SMTP hosting mode.
func usesShared(c Config, id string) bool {
	if !SharedEnabled(c) || !sharedConsumer(id) {
		return false
	}
	app, ok := c.Catalog.Find(id)
	return ok && !(id == "mail" && app.Env["WORKSPACE_MAIL_BACKEND"] == "private")
}
func sharedOwnsPort(c Config) bool {
	root := filepath.Join(c.Data, "shared")
	config, _, err := shared.Load(root)
	if err != nil || config.EnginePort != c.Port+11 || config.EngineSHA256 != c.Catalog.Packages["stalwart-shared"].SHA256 {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return shared.Healthy(ctx, root)
}
func (rt *Runtime) ensureShared(parent context.Context) error {
	rt.sharedOnce.Do(func() {
		c := rt.Config
		s := rt.Secrets
		root := filepath.Join(c.Data, "shared")
		if _, e := os.Stat(filepath.Join(c.Data, "mail", "engine.json")); e == nil && usesShared(c, "mail") {
			rt.sharedErr = errors.New("Existing Mail private storage requires an administrator migration before enabling shared Stalwart; data was not adopted or replaced")
			return
		}
		credentials := s.Clients["shared"]
		engine := c.Catalog.Packages["stalwart-shared"]
		manager := c.Catalog.Packages["services-shared"]
		if credentials.ID == "" || engine.File == "" || manager.File == "" {
			rt.sharedErr = errors.New("The shared backend package or MyAccounts client is missing")
			return
		}
		enginePath := filepath.Join(c.Programs, filepath.FromSlash(engine.File))
		held, e := platform.LockedBinary(enginePath, engine.SHA256)
		if e != nil {
			rt.sharedErr = e
			return
		}
		defer held.Close()
		managerPath := filepath.Join(c.Programs, filepath.FromSlash(manager.File))
		managerHeld, e := platform.LockedBinary(managerPath, manager.SHA256)
		if e != nil {
			rt.sharedErr = e
			return
		}
		defer managerHeld.Close()
		var initial *shared.Config
		if _, _, e = shared.Load(root); os.IsNotExist(e) {
			initial = &shared.Config{Issuer: c.Issuer, ClientID: credentials.ID, ClientSecret: credentials.Secret, Domain: "workspace.test", EnginePort: c.Port + 11, Version: shared.EngineVersion, EnginePath: enginePath, EngineSHA256: engine.SHA256}
		} else if e != nil {
			rt.sharedErr = e
			return
		}
		var installed shared.Config
		for _, app := range c.Catalog.Apps {
			if usesShared(c, app.ID) {
				if installed, e = shared.Install(root, app.ID, app.Version, initial); e != nil {
					rt.sharedErr = e
					return
				}
				initial = nil
			}
		}
		if e = prepareSharedRuntime(c, installed.EnginePath); e != nil {
			rt.sharedErr = e
			return
		}
		privateManager := filepath.Join(root, "WorkspaceServices.exe")
		if e = platform.CopyFile(managerPath, privateManager); e != nil && !os.IsExist(e) {
			if _, statErr := os.Stat(privateManager); statErr != nil {
				rt.sharedErr = e
				return
			}
		}
		ctx, cancel := context.WithTimeout(parent, 2*time.Minute)
		defer cancel()
		rt.sharedErr = ensureSharedReady(ctx, root, privateManager)
	})
	return rt.sharedErr
}
func registerSharedClient(ctx context.Context, api *identityAPI, c Config) (Client, error) {
	callbacks := []string{c.Origin + "/calendar/auth/callback", c.Origin + "/contacts/auth/callback", c.Origin + "/mail/auth/callback"}
	if len(c.Catalog.Apps) > 0 {
		callbacks = nil
		for _, app := range c.Catalog.Apps {
			if usesShared(c, app.ID) {
				callbacks = append(callbacks, c.Origin+app.Callback)
			}
		}
	}
	return registerClientRedirects(ctx, api, "shared", "Workspace Calendar, Mail and Contacts", callbacks)
}
