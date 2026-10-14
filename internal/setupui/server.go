package setupui

import (
	"archive/zip"
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fahimc/one-workspace/internal/catalog"
	"github.com/fahimc/one-workspace/internal/platform"
	"github.com/fahimc/one-workspace/internal/suite"
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
	z              *zip.Reader
	cat            catalog.Catalog
	token, origin  string
	mu             sync.Mutex
	busy, ready    bool
	phase, problem string
	config         suite.Config
}

func Serve(ctx context.Context, z *zip.Reader, c catalog.Catalog, open bool) error {
	listener, e := net.Listen("tcp4", "127.0.0.1:0")
	if e != nil {
		return e
	}
	defer listener.Close()
	s := &Server{z: z, cat: c, token: platform.Token(), origin: "http://localhost:" + fmt.Sprint(listener.Addr().(*net.TCPAddr).Port), phase: "Ready to install"}
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
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, page)
	})
	mux.HandleFunc("GET /api/state", func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"busy": s.busy, "ready": s.ready, "phase": s.phase, "error": s.problem, "version": s.cat.Version, "origin": s.config.Origin})
	})
	mux.HandleFunc("POST /api/browse", func(w http.ResponseWriter, r *http.Request) {
		path, e := platform.ChooseDataDir()
		if e != nil {
			http.Error(w, "Folder selection cancelled", 400)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"path": path})
	})
	mux.HandleFunc("POST /api/install", func(w http.ResponseWriter, r *http.Request) {
		var request suite.Request
		decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 32<<10))
		decoder.DisallowUnknownFields()
		if e := decoder.Decode(&request); e != nil {
			http.Error(w, "Invalid setup request", 400)
			return
		}
		if e := request.Validate(); e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		s.mu.Lock()
		if s.busy || s.ready {
			s.mu.Unlock()
			http.Error(w, "Setup is already running or complete", 409)
			return
		}
		s.busy = true
		s.problem = ""
		s.mu.Unlock()
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
			defer cancel()
			c, e := suite.Install(ctx, request, s.z, s.cat, func(phase string) { s.mu.Lock(); s.phase = phase; s.mu.Unlock() })
			s.mu.Lock()
			defer s.mu.Unlock()
			s.busy = false
			if e != nil {
				s.problem = e.Error()
				s.phase = "Setup stopped; inspect the error before retrying in fresh folders"
			} else {
				s.ready = true
				s.config = c
			}
		}()
		w.WriteHeader(http.StatusAccepted)
	})
	mux.HandleFunc("POST /api/launch", func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		ready, c := s.ready, s.config
		s.mu.Unlock()
		if !ready {
			http.Error(w, "Installation is not ready", 409)
			return
		}
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
		go func() { _ = cmd.Wait(); _ = log.Close() }()
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
