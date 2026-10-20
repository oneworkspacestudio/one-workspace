package suite

import (
	"context"
	"errors"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
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
	socialGuardMu      sync.Mutex
	socialGuardToken   string
	socialGuardTokenAt time.Time
	sharedOnce         sync.Once
	sharedErr          error
	Config             Config
	Secrets            Secrets
	mu                 sync.RWMutex
	status             map[string]string
	job                *platform.Job
	tunnelJob          *platform.Job
	appCancel          context.CancelFunc
	supervisors        sync.WaitGroup
	children           sync.WaitGroup
	appTasks           sync.WaitGroup
	operations         sync.Mutex
	rootContext        context.Context
	backups            *backupController
	cancel             context.CancelFunc
	progress           *ProgressTracker
	tunnelExit         chan error
	mailProvisionToken string
	chatDirectoryToken string
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
func Run(ctx context.Context, c Config, s Secrets) (runErr error) {
	if e := CheckUpdateRecovery(c); e != nil {
		return e
	}
	lock, e := platform.Lock(filepath.Join(c.Data, "runtime.lock"))
	if e != nil {
		return e
	}
	defer lock.Close()
	progress := NewProgress("starting", StartupSteps(c), filepath.Join(c.Data, "startup-progress.json"))
	progress.Step("checks", "running", "Checking the selected installation")
	defer func() {
		if runErr != nil {
			progress.Finish("failed", "Workspace could not start", runErr)
		} else {
			progress.Finish("stopped", "Workspace has stopped", nil)
		}
	}()
	if e = CheckPorts(c); e != nil {
		return e
	}
	appCount := len(c.Catalog.Apps)
	c, s, e = startupDefaults(ctx, c, s, func(detail string) { progress.Step("checks", "running", detail) })
	if e != nil {
		return e
	}
	// A legacy update can add steps before app startup begins.
	if len(c.Catalog.Apps) != appCount {
		progress = NewProgress("starting", StartupSteps(c), filepath.Join(c.Data, "startup-progress.json"))
	}
	restoreTemp, tempErr := setWorkspaceTemp(c)
	if tempErr != nil {
		return tempErr
	}
	defer restoreTemp()
	progress.Step("checks", "done", "Storage and ports are available")
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	job, e := platform.NewJob()
	if e != nil {
		return e
	}

	rt := &Runtime{Config: c, Secrets: s, status: map[string]string{}, job: job, cancel: cancel, progress: progress, tunnelExit: make(chan error, 1)}
	rt.rootContext = ctx
	rt.backups = newBackupController(c)
	defer func() {
		cancel()
		rt.operations.Lock()
		defer rt.operations.Unlock()
		rt.stopApps()
		rt.tunnelJob.Close()
	}()
	rt.mailProvisionToken = platform.Token()
	rt.chatDirectoryToken = platform.Token()
	token := platform.Token()
	if e = platform.SaveSecret(filepath.Join(c.Data, "control.dpapi"), map[string]string{"token": token}); e != nil {
		return e
	}
	if c.Tunnel {
		rt.tunnelJob, e = platform.NewJob()
		if e != nil {
			return e
		}
		progress.Step("tunnel", "running", "Starting your dedicated Cloudflare connector")
		if e = rt.startTunnel(ctx); e != nil {
			return e
		}
		progress.Step("tunnel", "done", "Connector process started; domain reachability is checked next")
	}
	if e = rt.startApps(ctx, true); e != nil {
		return e
	}
	progress.Step("portal", "running", "Connecting the Workspace portal to MyAccounts")
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
	progress.Step("portal", "done", "Workspace is listening and ready to open")
	progress.Finish("ready", "Your Workspace is running", nil)
	_ = platform.SaveJSON(filepath.Join(c.Data, "running.json"), map[string]any{"pid": os.Getpid(), "origin": c.Origin, "started": time.Now().UTC()})
	go rt.backupSchedule(ctx)
	go startupUpdateCleanup(ctx, c)
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
	if app.ID == "chat" {
		for index := 0; index+1 < len(args); index++ {
			if args[index] == "--origin" {
				args[index+1] = chatOrigin(c)
			}
		}
	}
	cmd := exec.Command(filepath.Join(c.Programs, filepath.FromSlash(app.Executable)), args...)
	cmd.Dir = appWorkingDir(c, app.ID)
	cmd.Env = append(os.Environ(), "TEMP="+runtimeTemp(c, app.ID), "TMP="+runtimeTemp(c, app.ID))
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
			cmd.Env = append(cmd.Env, "TEMP="+musicRuntimeTemp(c), "TMP="+musicRuntimeTemp(c))
		}
		if foundationApp(app.ID) {
			cmd.Env = append(cmd.Env, "WORKSPACE_OWNER_SUBJECT="+c.Owner, "WORKSPACE_HOME_URL="+c.Origin)
			if app.ID == "monitor" {
				cmd.Env = append(cmd.Env, fmt.Sprintf("WORKSPACE_REGISTRY_URL=http://127.0.0.1:%d/_integration/monitor", c.Port), "WORKSPACE_REGISTRY_TOKEN="+monitorCapability(client.Secret))
			}
		}
		if app.ID == "social" {
			cmd.Env = append(cmd.Env, "MYSOCIAL_PHOTOS_ROOT="+socialRoot(c), fmt.Sprintf("MYSOCIAL_STORAGE_GUARD_URL=http://127.0.0.1:%d/_integration/social/storage", c.Port), "MYSOCIAL_STORAGE_GUARD_TOKEN="+socialCapability(client.Secret))
		}
		if app.ID == "hosting" {
			cmd.Env = append(cmd.Env, "WORKSPACE_OWNER_SUBJECT="+c.Owner)
			if c.HostingDomain != "" {
				cmd.Env = append(cmd.Env, hostingSiteEnvironment(c, s)...)
			}
		}
		if app.ID == "chat" {
			cmd.Env = append(cmd.Env, "WORKSPACE_CHAT_DIRECTORY_TOKEN="+rt.chatDirectoryToken)
			cmd.Env = append(cmd.Env, "WORKSPACE_HOME_URL="+c.Origin, fmt.Sprintf("WORKSPACE_CHAT_CATALOG_URL=http://127.0.0.1:%d/_integration/chat/catalog", c.Port))
			cmd.Env = append(cmd.Env, "WORKSPACE_OIDC_REDIRECT_URI="+chatOrigin(c)+"/auth/callback")
			cmd.Env = append(cmd.Env, "WORKSPACE_CHAT_MANAGER_CAPABILITIES=isolated-chat-origin-v1")
		}
		if app.ID == "mail" {
			cmd.Env = append(cmd.Env, "WORKSPACE_BASE_PATH="+strings.TrimRight(app.Route, "/"), "WORKSPACE_OWNER_SUBJECT="+c.Owner)
			cmd.Env = append(cmd.Env, "WORKSPACE_MAIL_PROVISION_TOKEN="+rt.mailProvisionToken)
			if c.HostingDomain != "" && !usesShared(c, "mail") {
				cmd.Env = append(cmd.Env, "WORKSPACE_MAIL_DOMAIN="+c.HostingDomain, "WORKSPACE_MAIL_DNS_ZONE="+c.DNSZoneID, "WORKSPACE_MAIL_DNS_TOKEN="+s.MailDNSToken)
			}
		}
	}
	if usesShared(c, app.ID) {
		cmd.Env = append(cmd.Env, "WORKSPACE_SHARED_DATA="+filepath.Join(c.Data, "shared"))
	}
	log := &platform.Log{Path: filepath.Join(c.Data, "logs", app.ID+".log")}
	cmd.Stdout = log
	cmd.Stderr = log
	return cmd
}
func (rt *Runtime) start(ctx context.Context, app catalog.App) error {
	if e := prepareRuntimeTemp(rt.Config, app.ID); e != nil {
		return e
	}
	if app.ID == "chat" {
		rt.progress.Step(app.ID, "running", "Preparing Chat's separate browser address and existing MyAccounts connection")
		hostingContext, cancelHosting := context.WithTimeout(ctx, 90*time.Second)
		if err := rt.ensureChatHosting(hostingContext); err != nil {
			log := &platform.Log{Path: filepath.Join(rt.Config.Data, "logs", "chat-hosting.log")}
			_, _ = log.Write([]byte(err.Error() + "\n"))
		}
		cancelHosting()
		identityContext, cancel := context.WithTimeout(ctx, 25*time.Second)
		err := rt.prepareChatIdentity(identityContext)
		cancel()
		if err != nil {
			return err
		}
	}
	if app.ID == "vault" {
		rt.progress.Step(app.ID, "running", "Preparing My Vault's permanent secure address")
		hostingContext, cancelHosting := context.WithTimeout(ctx, 90*time.Second)
		if err := rt.ensureVaultHosting(hostingContext); err != nil {
			rt.progress.Step(app.ID, "running", "Remote Vault setup needs attention; local access is retained")
			log := &platform.Log{Path: filepath.Join(rt.Config.Data, "logs", "vault-hosting.log")}
			_, _ = log.Write([]byte(err.Error() + "\n"))
		}
		cancelHosting()
		identityContext, cancel := context.WithTimeout(ctx, 20*time.Second)
		err := rt.prepareVaultIdentity(identityContext)
		cancel()
		if err != nil {
			// An optional identity provider must never block independent vault access.
			rt.progress.Step(app.ID, "running", "MyAccounts connection needs attention; starting independent vault access")
		}
	}
	if app.ID == "music" {
		if e := prepareMusicTemp(rt.Config); e != nil {
			return e
		}
	}
	if app.ID == "notes" {
		if e := ensureNotesConfig(rt.Config); e != nil {
			return e
		}
	}
	if app.ID == "drive" {
		if e := ensureSocialDriveSource(rt.Config); e != nil {
			return e
		}
	}
	if app.ID == "my-tube" {
		if e := ensureMyTubeConfig(rt.Config); e != nil {
			return e
		}
	}
	if usesShared(rt.Config, app.ID) {
		if e := rt.ensureShared(ctx); e != nil {
			return e
		}
	}
	if app.ID == "social" {
		if e := rt.verifySocialStorage(ctx, true); e != nil {
			return e
		}
	}
	rt.progress.Step(app.ID, "running", "Verifying the app and starting its process")
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
	rt.children.Add(1)
	go func() { defer rt.children.Done(); done <- cmd.Wait() }()
	rt.set(app.ID, "starting")
	health := app.Health
	if health == "" {
		health = "/"
	}
	startupTimeout := 90 * time.Second
	if app.ID == "gitea" || app.ID == "mail" {
		startupTimeout = 3 * time.Minute
	}
	if e = WaitHTTPDetailed(ctx, fmt.Sprintf("http://localhost:%d%s", app.Port, health), startupTimeout, done, func(detail string) { rt.progress.Step(app.ID, "running", detail) }); e != nil {
		held.Close()
		return e
	}
	rt.set(app.ID, "running")
	rt.progress.Step(app.ID, "done", "App health check passed")
	if app.ID == "vault" {
		rt.vaultHostingProgress(ctx, app)
	}
	if app.ID == "chat" {
		rt.chatHostingProgress(ctx)
	}
	if app.ID == "mail" && rt.Config.HostingDomain != "" && !usesShared(rt.Config, "mail") {
		rt.mailHostingProgress(ctx, app)
	}
	rt.supervisors.Add(1)
	go func() {
		defer rt.supervisors.Done()
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
			if app.ID == "music" {
				if e := prepareMusicTemp(rt.Config); e != nil {
					rt.set(app.ID, "failed: Music temporary storage needs attention")
					return
				}
			}
			cmd = rt.command(app)
			if e := rt.job.Start(cmd); e != nil {
				rt.set(app.ID, "failed: see logs")
				return
			}
			rt.children.Add(1)
			go func(cmd *exec.Cmd) { defer rt.children.Done(); done <- cmd.Wait() }(cmd)
			if delay < 30*time.Second {
				delay *= 2
			}
			rt.set(app.ID, "running")
		}
	}()
	return nil
}
func (rt *Runtime) startTunnel(ctx context.Context) error {
	if e := prepareRuntimeTemp(rt.Config, "cloudflare"); e != nil {
		return e
	}
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
	cmd.Env = append(os.Environ(), "TUNNEL_TOKEN="+rt.Secrets.TunnelToken, "TEMP="+runtimeTemp(rt.Config, "cloudflare"), "TMP="+runtimeTemp(rt.Config, "cloudflare"))
	cmd.Dir = rt.Config.Data
	log := &platform.Log{Path: filepath.Join(rt.Config.Data, "logs", "cloudflare.log")}
	_, _ = fmt.Fprintln(log, "Starting the dedicated Workspace secure connector")
	cmd.Stdout = log
	cmd.Stderr = log
	if e = rt.tunnelJob.Start(cmd); e != nil {
		held.Close()
		return e
	}
	rt.set("cloudflare", "running")
	go func() {
		defer held.Close()
		e := cmd.Wait()
		_, _ = fmt.Fprintf(log, "Workspace secure connector exited: %v\n", e)
		rt.tunnelExit <- e
		rt.set("cloudflare", "stopped: see logs")
	}()
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

// Only called while holding operations, or during initial startup.
func (rt *Runtime) startApps(ctx context.Context, publicCheck bool) error {
	appContext, cancel := context.WithCancel(ctx)
	rt.appCancel = cancel
	apps, err := rt.Config.Catalog.Order()
	if err != nil {
		return err
	}
	for _, app := range apps {
		if app.Kind != "process" {
			continue
		}
		if err = rt.start(appContext, app); err != nil {
			return fmt.Errorf("start %s: %w", app.ID, err)
		}
		if app.ID == "accounts" && publicCheck && rt.Config.Tunnel {
			rt.progress.Step("public-accounts", "running", "Waiting for the secure sign-in connection")
			if err = WaitHTTPDetailed(appContext, rt.Config.Issuer+"/healthz", 5*time.Minute, rt.tunnelExit, func(detail string) { rt.progress.Step("public-accounts", "running", detail) }); err != nil {
				return err
			}
			rt.progress.Step("public-accounts", "done", "Your secure sign-in address is reachable")
		}
	}
	if _, ok := rt.Config.Catalog.Find("mail"); ok && !usesShared(rt.Config, "mail") {
		rt.appTasks.Add(1)
		go func() { defer rt.appTasks.Done(); rt.maintainMailboxes(appContext) }()
	}
	if _, ok := rt.Config.Catalog.Find("chat"); ok {
		rt.appTasks.Add(1)
		go func() { defer rt.appTasks.Done(); rt.maintainChatDirectory(appContext) }()
	}
	return nil
}
func (rt *Runtime) stopApps() error {
	if rt.appCancel != nil {
		rt.appCancel()
	}
	// Let supervisors leave their restart loops before closing the job handle.
	rt.supervisors.Wait()
	stopContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := rt.job.Stop(stopContext)
	rt.job.Close()
	rt.children.Wait()
	rt.appTasks.Wait()
	return err
}
