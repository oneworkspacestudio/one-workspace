package suite

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

// Windows reports a process creation failure before the app runs. Waiting gives
// a growing system page file time to catch up; retries remain strictly bounded.
func memoryStartError(err error) bool {
	return errors.Is(err, syscall.Errno(1455)) || errors.Is(err, syscall.Errno(1450)) || errors.Is(err, syscall.Errno(8)) || errors.Is(err, syscall.Errno(14))
}

func startWithMemoryRetry(ctx context.Context, start func() error, prepare func()) error {
	for attempt := 0; ; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		err := start()
		if err == nil || !memoryStartError(err) || attempt == 2 {
			return err
		}
		prepare()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(attempt+1) * time.Second):
		}
	}
}

func appStartMessage(app catalog.App, err error) string {
	if memoryStartError(err) {
		return app.Name + " couldn't start because Windows is low on available memory. Close unused apps, then retry. Your files are safe."
	}
	return app.Name + " couldn't start. Retry this app; if it still won't open, check for an app update. Your other apps can keep running."
}

func (rt *Runtime) appFailed(app catalog.App, err error) {
	log := &platform.Log{Path: filepath.Join(rt.Config.Data, "logs", app.ID+".log")}
	_, _ = fmt.Fprintf(log, "App startup failed: %v\n", err)
	rt.set(app.ID, "failed")
	rt.progress.Step(app.ID, "failed", appStartMessage(app, err))
	if rt.progress != nil && rt.progress.Snapshot().State == "ready" {
		rt.finishStartup()
	}
}

func (rt *Runtime) dependencyRunning(id string, states map[string]string) bool {
	app, ok := rt.Config.Catalog.Find(id)
	if !ok || app.Kind == "process" {
		return states[id] == "running"
	}
	for _, dep := range app.Depends {
		if !rt.dependencyRunning(dep, states) {
			return false
		}
	}
	return true
}

// Keep independent apps available, while never starting an app without its
// declared dependencies. Accounts remains essential to the portal's sign-in.
func (rt *Runtime) startAvailable(ctx context.Context, apps []catalog.App, start func(context.Context, catalog.App) error) error {
	for _, app := range apps {
		if err := ctx.Err(); err != nil {
			return err
		}
		// Shortcuts follow their dependencies without a fake process status.
		if app.Kind != "process" {
			continue
		}
		rt.mu.RLock()
		active := rt.active[app.ID]
		rt.mu.RUnlock()
		if active {
			continue
		}
		blocked := ""
		states := rt.Status()
		for _, dep := range app.Depends {
			if !rt.dependencyRunning(dep, states) {
				blocked = dep
				break
			}
		}
		if blocked != "" {
			name := blocked
			if dep, ok := rt.Config.Catalog.Find(blocked); ok {
				name = dep.Name
			}
			rt.set(app.ID, "blocked")
			rt.progress.Step(app.ID, "blocked", "Waiting for "+name+". Retry after that app is running.")
			continue
		}
		if err := start(ctx, app); err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			rt.appFailed(app, err)
			if app.ID == "accounts" {
				return errors.New("Sign-in couldn't start. Close unused apps and retry Workspace startup. Your installed apps and files are safe.")
			}
		}
	}
	return nil
}

func (rt *Runtime) finishStartup() {
	phase := "Your Workspace is running"
	for _, state := range rt.Status() {
		if state == "failed" || state == "blocked" {
			phase = "Workspace is running. Some apps need a retry."
			break
		}
	}
	rt.progress.Finish("ready", phase, nil)
}

// The operation lock serializes retries with backups and shutdown. App contexts
// survive an HTTP request but are cancelled when the suite stops or pauses.
func (rt *Runtime) retryApps() error {
	if !rt.operations.TryLock() {
		return errors.New("Workspace is busy. Wait a moment, then retry.")
	}
	if rt.rootContext == nil || rt.rootContext.Err() != nil || rt.appContext == nil || rt.appContext.Err() != nil {
		rt.operations.Unlock()
		return errors.New("Workspace is stopping. Start it again before retrying apps.")
	}
	apps, err := rt.Config.Catalog.Order()
	if err != nil {
		rt.operations.Unlock()
		return err
	}
	rt.progress.Finish("retrying", "Retrying apps that need attention", nil)
	go func() {
		defer rt.operations.Unlock()
		// A failed shared-service preparation may be attempted again. Successful
		// shared services and healthy app processes are retained.
		if rt.sharedErr != nil {
			rt.sharedOnce = sync.Once{}
			rt.sharedErr = nil
		}
		if err := rt.startAvailable(rt.appContext, apps, rt.start); err != nil {
			rt.progress.Finish("failed", "App retry needs attention", err)
			return
		}
		rt.finishStartup()
	}()
	return nil
}

// RetryInstalledApps uses the same protected local capability as suite stop;
// no token is placed in browser URLs or returned to the setup page.
func RetryInstalledApps(ctx context.Context, c Config) error {
	var control map[string]string
	if err := platform.LoadSecret(filepath.Join(c.Data, "control.dpapi"), &control); err != nil {
		return errors.New("Start Workspace before retrying apps")
	}
	r, _ := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://127.0.0.1:%d/_control/retry-apps", c.Port), nil)
	r.Header.Set("Authorization", "Bearer "+control["token"])
	resp, err := (&http.Client{Timeout: 10 * time.Second}).Do(r)
	if err != nil {
		return errors.New("Workspace could not be reached. Start Workspace again.")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return errors.New("Workspace is busy. Wait a moment, then retry.")
	}
	return nil
}
