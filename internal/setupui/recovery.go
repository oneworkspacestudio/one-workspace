package setupui

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"github.com/oneworkspacestudio/one-workspace/internal/suite"
)

// Restore display state only. Opening an installer never launches, reinstalls or
// resumes writes without an explicit action from this protected local wizard.
func (s *Server) restore(request suite.Request) error {
	if e := request.ResolveStorage(); e != nil {
		return e
	}
	status, active, e := suite.ReadSetupSession(request)
	if e != nil {
		return e
	}
	if active {
		s.recovered = &status
		s.following = true
		s.busy = true
		s.phase = status.Phase
		s.recoveryRequest = request
		return nil
	}
	c, _, loadErr := suite.Load(request.Data)
	if loadErr == nil && c.Owner != "" {
		s.config = c
		s.ready = true
		s.busy = false
		s.cat = c.Catalog
		s.phase = "Installation found. Check its app startup status below."
		s.recoveryRequest = request
		s.recovered = &status
		return nil
	}
	if status.Schema == 1 {
		s.recovered = &status
		s.recoveryRequest = request
		s.phase = status.Phase
		s.problem = status.Error
		return nil
	}
	var legacy map[string]string
	if platform.LoadJSON(filepath.Join(request.Data, "install-state.json"), &legacy) == nil {
		s.recovered = &suite.SetupSession{Schema: 1, Data: request.Data, Programs: request.Programs, Files: request.Files, State: "interrupted", Phase: legacy["phase"], Error: "An older installer stopped here. Its files have been kept. Review recovery to finish safely; this version did not save a continuation checkpoint."}
		if info, e := os.Stat(filepath.Join(request.Data, "install-state.json")); e == nil {
			s.recovered.UpdatedAt = info.ModTime()
		}
		s.recoveryRequest = request
		s.phase = legacy["phase"]
		s.problem = s.recovered.Error
		return nil
	}
	return errors.New("No saved installation status was found here. Choose the Workspace folder to inspect an older installation or use recovery")
}
func (s *Server) recoveryRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/recover", func(w http.ResponseWriter, r *http.Request) {
		var request suite.Request
		if json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&request) != nil {
			http.Error(w, "Choose the installation folder", 400)
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		if s.busy {
			http.Error(w, "An installation is already running. Reconnect to its progress", 409)
			return
		}
		if e := s.restore(request); e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		_ = suite.RememberSetup(s.recoveryRequest)
		w.WriteHeader(202)
	})
	mux.HandleFunc("POST /api/continue", func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		if s.busy || s.ready || s.recovered == nil {
			s.mu.Unlock()
			http.Error(w, "Check the saved installation before continuing", 409)
			return
		}
		request := s.recoveryRequest
		s.busy = true
		s.following = false
		s.problem = ""
		s.recovered = nil
		s.installation = suite.InstallProgress(false)
		s.installation.InstallPhase("Checking saved installation files")
		s.mu.Unlock()
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
			defer cancel()
			c, e := suite.ContinueSetup(ctx, request, s.z, s.cat, func(phase string) {
				s.mu.Lock()
				defer s.mu.Unlock()
				s.phase = phase
				s.installation.InstallPhase(phase)
				_ = suite.SaveSetupSession(request, s.installation.Snapshot())
			})
			s.mu.Lock()
			defer s.mu.Unlock()
			s.busy = false
			if e != nil {
				s.installation.Finish("failed", "Installation needs attention", e)
				s.problem = e.Error()
			} else {
				for _, step := range s.installation.Snapshot().Steps {
					s.installation.Step(step.ID, "done", "")
				}
				s.installation.Finish("installed", "Installation completed. Ready to start.", nil)
				s.config = c
				s.cat = c.Catalog
				s.ready = true
				s.phase = "Installation completed. Ready to start."
			}
			_ = suite.SaveSetupSession(request, s.installation.Snapshot())
			if e != nil {
				p := s.installation.Snapshot()
				s.recovered = &suite.SetupSession{Schema: 1, Data: request.Data, Programs: request.Programs, Files: request.Files, State: p.State, Phase: p.Phase, Error: p.Error, UpdatedAt: p.UpdatedAt, Progress: p}
			}
		}()
		w.WriteHeader(202)
	})
}
func (s *Server) recoveredState(state map[string]any) {
	if s.following || (s.recovered != nil && s.installation == nil && !s.ready) {
		status, active, e := suite.ReadSetupSession(s.recoveryRequest)
		if e == nil {
			s.recovered = &status
			s.phase = status.Phase
			if !active {
				s.following = false
				s.busy = false
				_ = s.restore(s.recoveryRequest)
			}
		}
	}
	if s.recovered != nil {
		state["savedInstallation"] = s.recovered
		state["installation"] = s.recovered.Progress
		state["logPath"] = filepath.Join(s.recoveryRequest.Data, "install-accounts.log")
		if _, e := os.Stat(filepath.Join(suite.SetupSessionDirectory(s.recoveryRequest.Data), "continue.dpapi")); e == nil {
			state["canContinue"] = true
		}
		state["recoveryFolder"] = filepath.Dir(s.recoveryRequest.Data)
	}
	state["busy"], state["ready"], state["phase"], state["error"] = s.busy, s.ready, s.phase, s.problem
	state["origin"] = s.config.Origin
	if s.ready {
		apps := []map[string]any{}
		for _, a := range s.config.Catalog.Apps {
			label := "Included in Workspace"
			if a.Executable != "" {
				label = "Files not found — use recovery"
				if info, e := os.Stat(filepath.Join(s.config.Programs, filepath.FromSlash(a.Executable))); e == nil && !info.IsDir() {
					label = "App files found"
				}
			}
			apps = append(apps, map[string]any{"name": a.Name, "version": a.Version, "status": label})
		}
		state["savedApps"] = apps
	}
}
