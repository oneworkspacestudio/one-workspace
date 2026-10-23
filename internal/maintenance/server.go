// Package maintenance hosts the Windows-local update window. Its capability
// token never passes through the remotely accessible Workspace portal.
package maintenance

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/branding"
	"github.com/oneworkspacestudio/one-workspace/internal/distribution"
	"github.com/oneworkspacestudio/one-workspace/internal/operationprogress"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"github.com/oneworkspacestudio/one-workspace/internal/suite"
)

type Server struct {
	mu                  sync.Mutex
	operationMu         sync.Mutex
	Data, Token, Origin string
	Busy                bool
	Phase, Problem      string
	Operation           *operationState
	Plan                *distribution.UpdatePlan
	Context             context.Context
	Close               context.CancelFunc
	Initial             suite.Config
	Repair              *suite.InstallerRepairResult
}

func Serve(ctx context.Context, data string, open bool) error {
	c, e := suite.LoadUpdateConfig(data)
	if e != nil {
		return e
	}
	if open && suite.OpenExistingUpdates(c) {
		return nil
	}
	root := filepath.Join(c.Data, "updates")
	if e = os.MkdirAll(root, 0700); e != nil {
		return e
	}
	if e = platform.ProtectMaintenanceDir(root, c.Data); e != nil {
		return e
	}
	lock, e := platform.Lock(filepath.Join(root, "window.lock"))
	if e != nil {
		return errors.New("the update window is already open on this PC")
	}
	defer lock.Close()
	listener, e := net.Listen("tcp4", "127.0.0.1:0")
	if e != nil {
		return e
	}
	defer listener.Close()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	s := &Server{Data: c.Data, Initial: c, Token: platform.Token(), Origin: fmt.Sprintf("http://localhost:%d", listener.Addr().(*net.TCPAddr).Port), Context: ctx, Close: cancel, Phase: "Check for updates when you are ready."}
	exe, _ := os.Executable()
	link := s.Origin + "/#" + s.Token
	marker := filepath.Join(root, "window.json")
	if e = platform.SaveJSON(marker, map[string]any{"pid": os.Getpid(), "executable": exe, "url": link}); e != nil {
		return e
	}
	defer os.Remove(marker)
	fmt.Println("Workspace updates:", link)
	if open {
		platform.OpenBrowser(link)
	}
	server := &http.Server{Handler: s.Handler(), ReadHeaderTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second}
	go func() { <-ctx.Done(); _ = server.Close() }()
	e = server.Serve(listener)
	if errors.Is(e, http.ErrServerClosed) {
		return nil
	}
	return e
}
func (s *Server) progress(phase string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Phase = phase
	if p := s.Operation; p != nil {
		if len(p.History) == 0 || p.History[len(p.History)-1].Text != phase {
			p.History = append(p.History, operationLog{At: time.Now().UTC(), Text: phase})
			if len(p.History) > 40 {
				p.History = p.History[len(p.History)-40:]
			}
		}
	}
}
func (s *Server) work(kind string, fn func(context.Context) error) {
	s.mu.Lock()
	s.Busy = true
	s.Problem = ""
	s.beginLocked(kind)
	s.mu.Unlock()
	go func() {
		e := fn(operationprogress.Observe(s.Context, s.detail))
		s.mu.Lock()
		defer s.mu.Unlock()
		s.Busy = false
		if e != nil {
			s.Problem = e.Error()
			s.Phase = "Update needs attention. Your verified backups are retained."
		}
		p := s.Operation
		p.FinishedAt = time.Now().UTC()
		p.Status = "complete"
		p.Detail = s.Phase
		if e != nil {
			p.Status = "failed"
		} else {
			p.Done, p.Total, p.Unit = 0, 0, ""
		}
		for i := range p.Steps {
			if e == nil {
				p.Steps[i].State = "done"
			} else if p.Steps[i].State == "running" {
				p.Steps[i].State = "failed"
			}
		}
	}()
}
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /branding/", branding.Handler())
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, branding.VersionHTML(page))
	})
	mux.HandleFunc("GET /api/state", func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		c, e := suite.LoadUpdateConfig(s.Data)
		newController := false
		if e != nil && s.Repair != nil && s.Initial.Data == s.Data {
			c = s.Initial
			c.Version = s.Repair.Version
			e = nil
			newController = true
		}
		if e != nil {
			http.Error(w, e.Error(), 500)
			return
		}
		backups, e := suite.UpdateBackups(c)
		if e != nil {
			http.Error(w, e.Error(), 500)
			return
		}
		cleanup, _ := suite.LoadUpdateCleanup(c)
		json.NewEncoder(w).Encode(map[string]any{"busy": s.Busy, "phase": s.Phase, "error": s.Problem, "plan": s.Plan, "backups": backups, "cleanup": cleanup, "pending": suite.UpdatePending(c), "version": c.Version, "origin": c.Origin, "progress": s.snapshotLocked(), "repair": s.Repair, "newController": newController})
	})
	mux.HandleFunc("POST /api/check", func(w http.ResponseWriter, r *http.Request) {
		s.work("check", func(ctx context.Context) error {
			s.mu.Lock()
			s.Plan = nil
			s.mu.Unlock()
			c, e := suite.LoadUpdateConfig(s.Data)
			if e != nil {
				return e
			}
			s.progress("Checking Workspace and each installed app repository")
			operationprogress.Report(ctx, "check", "Checking published releases for installed apps", 0, 0, "")
			p, e := distribution.CheckUpdates(ctx, c.Catalog)
			if e != nil {
				return e
			}
			s.mu.Lock()
			s.Plan = &p
			s.Phase = "Release check complete. Choose the updates to install."
			if p.Installer != nil && p.Installer.Required {
				s.Phase = "A newer installer is needed. Choose Back up and reinstall Workspace below."
			}
			s.mu.Unlock()
			return nil
		})
		w.WriteHeader(202)
	})
	mux.HandleFunc("POST /api/reinstall", func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		if s.Plan == nil || s.Plan.Installer == nil {
			s.mu.Unlock()
			http.Error(w, "Check for updates to find the verified installer first", 409)
			return
		}
		plan := *s.Plan
		s.mu.Unlock()
		c, e := suite.LoadUpdateConfig(s.Data)
		if e != nil {
			http.Error(w, e.Error(), 409)
			return
		}
		if suite.UpdatePending(c) {
			http.Error(w, "Recover the interrupted update before reinstalling", 409)
			return
		}
		s.work("reinstall", func(ctx context.Context) error {
			plan, e := s.preparePlan(ctx, plan, nil, true)
			if e != nil {
				return e
			}
			result, e := suite.RunInstallerRepair(ctx, s.Data, plan, s.progress)
			s.mu.Lock()
			defer s.mu.Unlock()
			s.Repair = result
			if result != nil {
				s.Plan = nil
			}
			if e == nil {
				s.Phase = "Workspace reinstalled and restarted. All apps are running."
			}
			return e
		})
		w.WriteHeader(202)
	})
	mux.HandleFunc("POST /api/apply", func(w http.ResponseWriter, r *http.Request) {
		var in suite.UpdateRequest
		if e := json.NewDecoder(http.MaxBytesReader(w, r.Body, 32<<10)).Decode(&in); e != nil {
			http.Error(w, "Invalid update request", 400)
			return
		}
		s.mu.Lock()
		if s.Plan == nil {
			s.mu.Unlock()
			http.Error(w, "Check for updates first", 409)
			return
		}
		plan := *s.Plan
		s.mu.Unlock()
		s.work("apply", func(ctx context.Context) error {
			plan, e := s.preparePlan(ctx, plan, in.Selected, false)
			if e != nil {
				return e
			}
			e = suite.ApplyUpdates(ctx, s.Data, plan, in, s.progress)
			if e == nil {
				s.mu.Lock()
				s.Plan = nil
				s.mu.Unlock()
			}
			return e
		})
		w.WriteHeader(202)
	})
	mux.HandleFunc("POST /api/recover", func(w http.ResponseWriter, r *http.Request) {
		s.work("recover", func(ctx context.Context) error { return suite.RecoverUpdate(ctx, s.Data, s.progress) })
		w.WriteHeader(202)
	})
	mux.HandleFunc("POST /api/rollback", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			ID string `json:"id"`
		}
		if json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<10)).Decode(&in) != nil {
			http.Error(w, "Invalid backup", 400)
			return
		}
		s.work("rollback", func(ctx context.Context) error { return suite.RollbackUpdate(ctx, s.Data, in.ID, s.progress) })
		w.WriteHeader(202)
	})
	mux.HandleFunc("POST /api/close", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(202)
		go func() { time.Sleep(100 * time.Millisecond); s.Close() }()
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'unsafe-inline'; style-src 'unsafe-inline'; connect-src 'self'; frame-ancestors 'none'; base-uri 'none'")
		if r.Host != strings.TrimPrefix(s.Origin, "http://") {
			http.Error(w, "Host rejected", 403)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/") {
			if subtle.ConstantTimeCompare([]byte(r.Header.Get("X-Workspace-Update")), []byte(s.Token)) != 1 {
				http.Error(w, "Open the update link on the server PC", 403)
				return
			}
			if r.Method != "GET" && r.Header.Get("Origin") != s.Origin {
				http.Error(w, "Origin rejected", 403)
				return
			}
			if r.Method != "GET" {
				if !s.operationMu.TryLock() {
					http.Error(w, "Another update request is running", 409)
					return
				}
				defer s.operationMu.Unlock()
				s.mu.Lock()
				busy := s.Busy
				s.mu.Unlock()
				if busy {
					http.Error(w, "An update operation is already running", 409)
					return
				}
			}
		}
		mux.ServeHTTP(w, r)
	})
}
