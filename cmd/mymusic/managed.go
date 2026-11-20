package main

import (
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

type engineSettings struct {
	Engine      string
	ACEKey      string
	DeepSeekKey string
}

func runManaged() error {
	// The no-argument launch remains a loopback-only desktop app. Workspace
	// launches the same binary with its supervised apphost arguments.
	ui, err := fs.Sub(assets, "web")
	if err != nil {
		return err
	}
	cfg := apphost.Config{ID: "mymusic", Name: "Vuno", Version: version, Listen: "127.0.0.1:9049", RestoreFiles: []string{"audio", "engine-settings"}}
	return apphost.Main(cfg, ui, func(h *apphost.Host) (apphost.Application, error) {
		if h.Config.Issuer == "" {
			return nil, errors.New("configure the MyAccounts client before starting managed Vuno")
		}
		return nil, registerManaged(h)
	})
}

func registerManaged(h *apphost.Host) error {
	dir := filepath.Join(h.Config.Data, "audio")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	h.BackupPaths = append(h.BackupPaths, "audio", "engine-settings")

	path := filepath.Join(h.Config.Data, "engine-settings")
	settings := engineSettings{Engine: defaultEngine}
	if err := platform.LoadSecret(path, &settings); err != nil && !errors.Is(err, os.ErrNotExist) {
		return errors.New("could not read Vuno engine settings")
	}
	if !validEngine(settings.Engine) {
		return errors.New("invalid saved ACE-Step origin")
	}
	var mu sync.Mutex
	snapshot := func() engineSettings {
		mu.Lock()
		defer mu.Unlock()
		return settings
	}
	app := application{dir: dir, config: snapshot, freeze: &h.Freeze}
	app.register(h.Mux, h.Owner)
	h.Mux.HandleFunc("PUT /api/settings", h.Owner(func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Engine              string
			ACEKey, DeepSeekKey *string
		}
		if !apphost.Decode(w, r, &input) {
			return
		}
		engine := strings.TrimSpace(input.Engine)
		if !validEngine(engine) || !validKey(input.ACEKey) || !validKey(input.DeepSeekKey) {
			apphost.Problem(w, 400, "Use a loopback HTTP origin and valid API keys")
			return
		}
		mu.Lock()
		defer mu.Unlock()
		next := settings
		next.Engine = engine
		if input.ACEKey != nil {
			next.ACEKey = strings.TrimSpace(*input.ACEKey)
		}
		if input.DeepSeekKey != nil {
			next.DeepSeekKey = strings.TrimSpace(*input.DeepSeekKey)
		}
		if err := platform.SaveSecret(path, next); err != nil {
			apphost.Problem(w, 500, "Could not save engine settings")
			return
		}
		settings = next
		h.Audit(apphost.Current(r).User.ID, "settings-updated", "mymusic-engine")
		apphost.JSON(w, 200, map[string]bool{"saved": true})
	}))
	return nil
}

func validKey(value *string) bool {
	return value == nil || len(*value) <= 1024 && !strings.ContainsAny(*value, "\r\n\x00")
}
