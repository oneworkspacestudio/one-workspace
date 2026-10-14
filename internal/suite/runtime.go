package suite

import (
	"context"
	"errors"
	"fmt"
	"github.com/fahimc/one-workspace/internal/catalog"
	"github.com/fahimc/one-workspace/internal/platform"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Runtime struct {
	Config  Config
	Secrets Secrets
	mu      sync.RWMutex
	status  map[string]string
	job     *platform.Job
	cancel  context.CancelFunc
}

func (rt *Runtime) set(id, state string) { rt.mu.Lock(); rt.status[id] = state; rt.mu.Unlock() }
func (rt *Runtime) Status() map[string]string {
	rt.mu.RLock()
	defer rt.mu.RUnlock()
	out := map[string]string{}
	for k, v := range rt.status {
		out[k] = v
	}
	return out
}
func Run(ctx context.Context, c Config, s Secrets) error {
	lock, e := platform.Lock(filepath.Join(c.Data, "runtime.lock"))
	if e != nil {
		return e
	}
	defer lock.Close()
	if e = CheckPorts(c); e != nil {
		return e
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	job, e := platform.NewJob()
	if e != nil {
		return e
	}
	defer job.Close()
	rt := &Runtime{Config: c, Secrets: s, status: map[string]string{}, job: job, cancel: cancel}
	token := platform.Token()
	if e = platform.SaveSecret(filepath.Join(c.Data, "control.dpapi"), map[string]string{"token": token}); e != nil {
		return e
	}
	if c.Tunnel {
		if e = rt.startTunnel(ctx); e != nil {
			return e
		}
	}
	apps, e := c.Catalog.Order()
	if e != nil {
		return e
	}
	for _, app := range apps {
		if app.Kind != "process" {
			continue
		}
		if e = rt.start(ctx, app); e != nil {
			return fmt.Errorf("start %s: %w", app.ID, e)
		}
		if app.ID == "accounts" && c.Tunnel {
			if e = WaitHTTP(ctx, c.Issuer+"/healthz", 2*time.Minute); e != nil {
				return e
			}
		}
	}
	portal, e := NewPortal(ctx, rt)
	if e != nil {
		return e
	}
	handler := portal.Handler(token)
	listener, e := net.Listen("tcp4", fmt.Sprintf("127.0.0.1:%d", c.Port))
	if e != nil {
		return e
	}
	server := &http.Server{Handler: handler, ReadHeaderTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20}
	errc := make(chan error, 1)
	go func() { errc <- server.Serve(listener) }()
	_ = platform.SaveJSON(filepath.Join(c.Data, "running.json"), map[string]any{"pid": os.Getpid(), "origin": c.Origin, "started": time.Now().UTC()})
	select {
	case <-ctx.Done():
	case e = <-errc:
		if !errors.Is(e, http.ErrServerClosed) {
			return e
		}
	}
	stop, cancelStop := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelStop()
	_ = server.Shutdown(stop)
	return nil
}
func (rt *Runtime) command(app catalog.App) *exec.Cmd {
	c, s := rt.Config, rt.Secrets
	vars := strings.NewReplacer("${DATA}", c.Data, "${FILES}", c.Files, "${ORIGIN}", c.Origin, "${ISSUER}", c.Issuer, "${PORT}", strconv.Itoa(app.Port), "${ROUTE}", app.Route)
	args := []string{}
	for _, v := range app.Args {
		args = append(args, vars.Replace(v))
	}
	cmd := exec.Command(filepath.Join(c.Programs, filepath.FromSlash(app.Executable)), args...)
	cmd.Dir = c.Data
	cmd.Env = os.Environ()
	for k, v := range app.Env {
		cmd.Env = append(cmd.Env, k+"="+vars.Replace(v))
	}
	if app.ID == "drive" {
		cmd.Env = append(cmd.Env, "FILEBROWSER_ADMIN_PASSWORD="+s.RecoveryPassword, "FILEBROWSER_OIDC_CLIENT_SECRET="+s.Clients["drive"].Secret)
	} else if app.OIDC {
		client := s.Clients[app.ID]
		cmd.Env = append(cmd.Env, "WORKSPACE_OIDC_ISSUER="+c.Issuer, "WORKSPACE_OIDC_CLIENT_ID="+client.ID, "WORKSPACE_OIDC_CLIENT_SECRET="+client.Secret, "WORKSPACE_OIDC_REDIRECT_URI="+c.Origin+app.Callback)
		if app.ID == "music" {
			cmd.Env = append(cmd.Env, "WORKSPACE_MUSIC_DIR="+filepath.Join(c.Files, "Music"), "WORKSPACE_BASE_PATH="+strings.TrimRight(app.Route, "/"), "WORKSPACE_OWNER_SUBJECT="+c.Owner)
		}
		if app.ID == "mail" {
			cmd.Env = append(cmd.Env, "WORKSPACE_BASE_PATH="+strings.TrimRight(app.Route, "/"), "WORKSPACE_OWNER_SUBJECT="+c.Owner)
		}
	}
	log := &platform.Log{Path: filepath.Join(c.Data, "logs", app.ID+".log")}
	cmd.Stdout = log
	cmd.Stderr = log
	return cmd
}
func (rt *Runtime) start(ctx context.Context, app catalog.App) error {
	held, e := platform.LockedBinary(filepath.Join(rt.Config.Programs, filepath.FromSlash(app.Executable)), rt.Config.Catalog.Packages[app.Package].SHA256)
	if e != nil {
		return e
	}
	cmd := rt.command(app)
	if e = rt.job.Start(cmd); e != nil {
		held.Close()
		return e
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	rt.set(app.ID, "starting")
	health := app.Health
	if health == "" {
		health = "/"
	}
	startupTimeout := 90 * time.Second
	if app.ID == "gitea" {
		startupTimeout = 3 * time.Minute
	}
	if e = WaitHTTP(ctx, fmt.Sprintf("http://localhost:%d%s", app.Port, health), startupTimeout); e != nil {
		held.Close()
		return e
	}
	rt.set(app.ID, "running")
	go func() {
		defer held.Close()
		delay := time.Second
		for {
			select {
			case <-ctx.Done():
				return
			case <-done:
				rt.set(app.ID, "restarting")
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
			}
			if ctx.Err() != nil {
				return
			}
			cmd = rt.command(app)
			if e := rt.job.Start(cmd); e != nil {
				rt.set(app.ID, "failed: see logs")
				return
			}
			go func(cmd *exec.Cmd) { done <- cmd.Wait() }(cmd)
			if delay < 30*time.Second {
				delay *= 2
			}
			rt.set(app.ID, "running")
		}
	}()
	return nil
}
func (rt *Runtime) startTunnel(ctx context.Context) error {
	p, ok := rt.Config.Catalog.Packages["cloudflared"]
	if !ok {
		return errors.New("connector is absent from this bundle")
	}
	held, e := platform.LockedBinary(filepath.Join(rt.Config.Programs, p.File), p.SHA256)
	if e != nil {
		return e
	}
	config := filepath.Join(rt.Config.Data, "cloudflared.yml")
	if e = platform.AtomicWrite(config, []byte("loglevel: error\n")); e != nil {
		held.Close()
		return e
	}
	cmd := exec.Command(filepath.Join(rt.Config.Programs, p.File), "tunnel", "--config", config, "--no-autoupdate", "run")
	cmd.Env = append(os.Environ(), "TUNNEL_TOKEN="+rt.Secrets.TunnelToken)
	cmd.Dir = rt.Config.Data
	log := &platform.Log{Path: filepath.Join(rt.Config.Data, "logs", "cloudflare.log")}
	cmd.Stdout = log
	cmd.Stderr = log
	if e = rt.job.Start(cmd); e != nil {
		held.Close()
		return e
	}
	rt.set("cloudflare", "running")
	go func() { defer held.Close(); _ = cmd.Wait(); rt.set("cloudflare", "stopped: see logs") }()
	return nil
}
func Stop(ctx context.Context, c Config) error {
	var control map[string]string
	if e := platform.LoadSecret(filepath.Join(c.Data, "control.dpapi"), &control); e != nil {
		return e
	}
	req, _ := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://localhost:%d/_control/stop", c.Port), nil)
	req.Header.Set("Authorization", "Bearer "+control["token"])
	resp, e := (&http.Client{Timeout: 10 * time.Second}).Do(req)
	if e != nil {
		return e
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("runtime refused stop")
	}
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		lock, e := platform.Lock(filepath.Join(c.Data, "runtime.lock"))
		if e == nil {
			lock.Close()
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
	return errors.New("suite is still stopping; inspect launcher and app logs")
}
