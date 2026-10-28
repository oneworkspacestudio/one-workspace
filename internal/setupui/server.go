package setupui

import (
	"archive/zip"
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/branding"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/maintenance"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"github.com/oneworkspacestudio/one-workspace/internal/suite"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	recovered                  *suite.SetupSession
	following                  bool
	recoveryRequest            suite.Request
	offline                    bool
	z                          *zip.Reader
	cat                        catalog.Catalog
	token, origin              string
	mu                         sync.Mutex
	pickerMu                   sync.Mutex
	busy, ready                bool
	phase, problem             string
	config                     suite.Config
	hosting                    *suite.Request
	connecting                 bool
	hostingPhase, hostingError string
	installation               *suite.ProgressTracker
	launching                  bool
	launchStarted              time.Time
	launchError                string
}

func Serve(ctx context.Context, z *zip.Reader, c catalog.Catalog, open bool, offline ...bool) error {
	return serve(ctx, z, c, open, nil, len(offline) > 0 && offline[0])
}
func ServeInstalled(ctx context.Context, z *zip.Reader, c catalog.Catalog, installed suite.Config, open bool) error {
	return serve(ctx, z, c, open, &installed, false)
}
func serve(ctx context.Context, z *zip.Reader, c catalog.Catalog, open bool, installed *suite.Config, offline bool) error {
	listener, e := net.Listen("tcp4", "127.0.0.1:0")
	if e != nil {
		return e
	}
	defer listener.Close()
	s := &Server{offline: offline, z: z, cat: c, token: platform.Token(), origin: "http://localhost:" + fmt.Sprint(listener.Addr().(*net.TCPAddr).Port), phase: "Ready to install"}
	if installed != nil {
		s.config = *installed
		s.ready = true
		s.phase = "Apps installed. Check startup status below."
	}
	if installed == nil {
		if previous, e := suite.LastSetup(); e == nil {
			_ = s.restore(previous)
		}
	}
	server := &http.Server{Handler: s.Handler(), ReadHeaderTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second}
	go func() { <-ctx.Done(); _ = server.Close() }()
	link := s.origin + "/#" + s.token
	fmt.Println("Setup wizard:", link)
	if open {
		platform.OpenBrowser(link)
	}
	e = server.Serve(listener)
	if errors.Is(e, http.ErrServerClosed) {
		return nil
	}
	return e
}
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	s.recoveryRoutes(mux)
	mux.Handle("GET /branding/", branding.Handler())
	mux.HandleFunc("GET /api/startup-report", func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		c, ready := s.config, s.ready
		s.mu.Unlock()
		if !ready {
			http.Error(w, "Install Workspace first", 409)
			return
		}
		var progress suite.Progress
		if err := platform.LoadJSON(filepath.Join(c.Data, "startup-progress.json"), &progress); err != nil {
			http.Error(w, "Start Workspace to create a diagnostic report", 409)
			return
		}
		apps := []map[string]string{}
		for _, app := range c.Catalog.Apps {
			apps = append(apps, map[string]string{"id": app.ID, "name": app.Name, "version": app.Version})
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", `attachment; filename="Workspace-startup-report.json"`)
		// Same bounded status already displayed in the capability-protected wizard;
		// no account details, secrets, environment, raw logs or command lines.
		_ = json.NewEncoder(w).Encode(map[string]any{"schema": 1, "workspaceVersion": c.Catalog.Version, "createdAt": time.Now().UTC(), "startup": progress, "apps": apps, "logsFolder": filepath.Join(c.Data, "logs"), "sharedLogsFolder": filepath.Join(c.Data, "shared"), "note": "This report contains startup status, not raw diagnostic logs or credentials."})
	})
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, branding.VersionHTML(page))
	})
	mux.HandleFunc("GET /api/state", func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		state := map[string]any{"busy": s.busy, "ready": s.ready, "phase": s.phase, "error": s.problem, "version": s.cat.Version, "origin": s.config.Origin, "connecting": s.connecting, "hostingPhase": s.hostingPhase, "hostingError": s.hostingError}
		s.recoveredState(state)
		available := []map[string]any{}
		for _, a := range s.cat.Apps {
			category := catalog.CategoryFor(a.ID)
			available = append(available, map[string]any{"id": a.ID, "name": a.Name, "description": a.Description, "depends": a.Depends, "category": category.ID, "categoryName": category.Name, "defaultInstall": category.DefaultInstall})
		}
		state["availableApps"] = available
		state["offline"] = s.offline
		state["launching"] = s.launching
		state["launchError"] = s.launchError
		state["launchStarted"] = s.launchStarted
		if s.installation != nil {
			state["installation"] = s.installation.Snapshot()
		}
		if s.ready {
			state["diagnosticFolder"] = s.config.Data
			state["workspaceFolder"] = filepath.Dir(s.config.Data)
			var startup suite.Progress
			if platform.LoadJSON(filepath.Join(s.config.Data, "startup-progress.json"), &startup) == nil && (s.launchStarted.IsZero() || !startup.StartedAt.Before(s.launchStarted.Add(-time.Second))) {
				if (startup.State == "starting" || startup.State == "ready" || startup.State == "retrying") && !platform.ProcessRunning(startup.PID, filepath.Join(s.config.Data, "Workspace.exe")) {
					startup.State = "failed"
					startup.Phase = "Workspace is no longer running"
					startup.Error = "The Workspace launcher stopped. Retry startup or inspect the startup log."
				}
				state["startup"] = startup
			}
			state["logPath"] = filepath.Join(s.config.Data, "launcher.log")
		}
		if s.hosting != nil {
			state["hosting"] = map[string]string{"domain": s.hosting.Domain, "origin": s.hosting.Origin, "issuer": s.hosting.Issuer}
		}
		_ = json.NewEncoder(w).Encode(state)
	})
	mux.HandleFunc("POST /api/hosting/connect", s.connectHosting)
	mux.HandleFunc("POST /api/hosting/review", s.reviewHosting)
	mux.HandleFunc("POST /api/updates/open", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Folder string `json:"folder"`
		}
		if json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&in) != nil {
			http.Error(w, "Choose the existing Workspace parent folder", 400)
			return
		}
		request := suite.Request{Folder: in.Folder}
		if e := request.ResolveStorage(); e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		if _, e := suite.LoadUpdateConfig(request.Data); e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		go func() { _ = maintenance.Serve(context.Background(), request.Data, true) }()
		w.WriteHeader(202)
	})
	mux.HandleFunc("POST /api/browse", func(w http.ResponseWriter, r *http.Request) {
		if !s.pickerMu.TryLock() {
			http.Error(w, "The folder picker is already open", 409)
			return
		}
		defer s.pickerMu.Unlock()
		var in struct {
			Folder string `json:"folder"`
		}
		if json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&in) != nil {
			http.Error(w, "Invalid folder request", 400)
			return
		}
		path, e := platform.ChooseDataDir(in.Folder)
		if e != nil {
			http.Error(w, "Folder selection cancelled", 400)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"path": path})
	})
	mux.HandleFunc("POST /api/storage", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			Folder string `json:"folder"`
		}
		if json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&in) != nil {
			http.Error(w, "Invalid folder request", 400)
			return
		}
		if in.Folder == "" {
			http.Error(w, "Choose your Workspace folder", 400)
			return
		}
		request := suite.Request{Folder: in.Folder}
		request.Offline = s.offline
		if e := request.ResolveStorage(); e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		review, e := suite.InspectSetupStorage(request)
		if e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"folder": request.Folder, "volume": review.Volume, "reinstall": review})
	})
	mux.HandleFunc("POST /api/install", func(w http.ResponseWriter, r *http.Request) {
		var request suite.Request
		decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 32<<10))
		decoder.DisallowUnknownFields()
		if e := decoder.Decode(&request); e != nil {
			http.Error(w, "Invalid setup request", 400)
			return
		}
		request.Offline = s.offline
		if e := request.ResolveStorage(); e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		review, storageErr := suite.InspectSetupStorage(request)
		if storageErr != nil {
			http.Error(w, storageErr.Error(), 400)
			return
		}
		if review.Kind != "new" && (!request.Reinstall || request.ReinstallReview != review.Review) {
			http.Error(w, "Review this folder and choose Reinstall in this folder before continuing", 409)
			return
		}
		if request.Reinstall && (review.Kind == "new" || request.ReinstallReview != review.Review) {
			http.Error(w, "This folder changed. Check the reinstall option again.", 409)
			return
		}
		if !request.Reinstall {
			if _, e := request.PreflightStorage(); e != nil {
				http.Error(w, e.Error(), 400)
				return
			}
		}
		s.mu.Lock()
		if request.Mode == "automatic" && review.Kind != "repair" {
			domain, e := suite.NormalizeDomain(request.Domain)
			if e != nil || s.hosting == nil || s.connecting || s.hosting.Domain != domain {
				s.mu.Unlock()
				http.Error(w, "Connect Cloudflare for this domain first", 400)
				return
			}
			h := *s.hosting
			if e := applyHostingSelection(&request, h); e != nil {
				s.mu.Unlock()
				http.Error(w, e.Error(), 400)
				return
			}
		}
		s.mu.Unlock()
		if e := request.Validate(); e != nil && review.Kind != "repair" {
			http.Error(w, e.Error(), 400)
			return
		}
		if _, e := suite.SelectInstallApps(s.cat, request); e != nil && review.Kind != "repair" {
			http.Error(w, e.Error(), 400)
			return
		}
		s.mu.Lock()
		if s.busy || s.connecting || (s.ready && !(request.Reinstall && review.Kind == "repair")) {
			s.mu.Unlock()
			http.Error(w, "Setup is already running or complete", 409)
			return
		}
		s.busy = true
		s.ready = false
		s.recovered = nil
		s.recoveryRequest = request
		_ = suite.RememberSetup(request)
		s.installation = suite.InstallProgress(request.Mode == "cloudflare")
		if review.Kind == "repair" {
			s.installation = suite.ReinstallProgress()
		}
		s.installation.InstallPhase("Preparing storage and ports")
		_ = suite.SaveSetupSession(request, s.installation.Snapshot())
		s.problem = ""
		s.mu.Unlock()
		go func() {
			limit := 2 * time.Hour
			if request.Reinstall {
				limit = 2 * time.Hour
			}
			ctx, cancel := context.WithTimeout(context.Background(), limit)
			defer cancel()
			c, e := suite.Install(ctx, request, s.z, s.cat, func(phase string) {
				s.mu.Lock()
				s.phase = phase
				if review.Kind == "repair" {
					s.installation.ReinstallPhase(phase)
				} else {
					s.installation.InstallPhase(phase)
				}
				_ = suite.SaveSetupSession(request, s.installation.Snapshot())
				s.mu.Unlock()
			})
			s.mu.Lock()
			defer s.mu.Unlock()
			s.busy = false
			defer func() {
				_ = suite.SaveSetupSession(request, s.installation.Snapshot())
				if !s.ready {
					p := s.installation.Snapshot()
					s.recovered = &suite.SetupSession{Schema: 1, Data: request.Data, Programs: request.Programs, Files: request.Files, State: p.State, Phase: p.Phase, Error: p.Error, UpdatedAt: p.UpdatedAt, Progress: p}
				}
			}()
			if e != nil {
				s.installation.Finish("failed", "Installation stopped", e)
				s.problem = e.Error()
				s.phase = "Setup stopped. Return to the folder step and choose Reinstall in this folder to try again."
			} else {
				for _, step := range s.installation.Snapshot().Steps {
					s.installation.Step(step.ID, "done", "")
				}
				s.installation.Finish("installed", "Apps installed. Ready to start.", nil)
				s.ready = true
				s.config = c
			}
		}()
		w.WriteHeader(http.StatusAccepted)
	})
	mux.HandleFunc("POST /api/retry-apps", func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		ready, c := s.ready, s.config
		s.mu.Unlock()
		if !ready {
			http.Error(w, "Install Workspace first", 409)
			return
		}
		if err := suite.RetryInstalledApps(r.Context(), c); err != nil {
			http.Error(w, err.Error(), 409)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	})
	mux.HandleFunc("POST /api/launch", func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		ready, c := s.ready, s.config
		if s.launching {
			s.mu.Unlock()
			http.Error(w, "Workspace startup is already being monitored", 409)
			return
		}
		s.mu.Unlock()
		if !ready {
			http.Error(w, "Installation is not ready", 409)
			return
		}
		s.mu.Lock()
		s.launching = true
		s.launchStarted = time.Now().UTC()
		s.launchError = ""
		s.mu.Unlock()
		launched := false
		defer func() {
			if !launched {
				s.mu.Lock()
				s.launching = false
				s.mu.Unlock()
			}
		}()
		cmd := exec.Command(filepath.Join(c.Data, "Workspace.exe"), "--data", c.Data, "run")
		cmd.Dir = c.Data
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
		log, e := os.OpenFile(filepath.Join(c.Data, "launcher.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if e != nil {
			http.Error(w, "Cannot open launcher log", 500)
			return
		}
		cmd.Stdout = log
		cmd.Stderr = log
		if e = cmd.Start(); e != nil {
			log.Close()
			http.Error(w, "Could not launch Workspace", 500)
			return
		}
		launched = true
		go func() {
			e := cmd.Wait()
			_ = log.Close()
			s.mu.Lock()
			defer s.mu.Unlock()
			s.launching = false
			if e != nil {
				s.launchError = "Workspace exited: " + e.Error() + ". See " + filepath.Join(c.Data, "launcher.log")
			}
		}()
		w.WriteHeader(202)
	})
	host := strings.TrimPrefix(s.origin, "http://")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'unsafe-inline'; style-src 'unsafe-inline'; connect-src 'self'; frame-ancestors 'none'; base-uri 'none'")
		if r.Host != host {
			http.Error(w, "Host rejected", 403)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/") {
			if subtle.ConstantTimeCompare([]byte(r.Header.Get("X-Setup-Token")), []byte(s.token)) != 1 {
				http.Error(w, "Open the setup link from the installer", 403)
				return
			}
			if r.Method != "GET" && r.Header.Get("Origin") != s.origin {
				http.Error(w, "Origin rejected", 403)
				return
			}
		}
		mux.ServeHTTP(w, r)
	})
}
