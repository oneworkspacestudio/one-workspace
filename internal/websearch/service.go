package websearch

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

const Version = "0.1.0-preview.1"

func randomID() string  { return apphost.ID() }
func itoa(n int) string { return strconv.Itoa(n) }

type Job struct {
	ID      string  `json:"id"`
	Query   string  `json:"query"`
	Status  string  `json:"status"`
	Created int64   `json:"created"`
	Result  *Result `json:"result,omitempty"`
	Error   string  `json:"error,omitempty"`
}
type Service struct {
	h             *apphost.Host
	mu            sync.Mutex
	enabled       bool
	state         string
	agent         *Agent
	active        string
	runtimeCancel context.CancelFunc
	cancel        context.CancelFunc
	wake          chan struct{}
	runtime       string
	models        string
	setupRunning  bool
	setupBytes    int64
	setupError    string
	parent        context.Context
}

func New(h *apphost.Host) (*Service, error) {
	_, e := h.DB.Exec(`CREATE TABLE IF NOT EXISTS search_settings(id INTEGER PRIMARY KEY CHECK(id=1),enabled INTEGER NOT NULL);INSERT OR IGNORE INTO search_settings VALUES(1,0);
CREATE TABLE IF NOT EXISTS search_grants(hash TEXT PRIMARY KEY,name TEXT NOT NULL,created INTEGER NOT NULL);
CREATE TABLE IF NOT EXISTS search_jobs(id TEXT PRIMARY KEY,principal TEXT NOT NULL,session TEXT NOT NULL DEFAULT '',query TEXT NOT NULL,status TEXT NOT NULL,created INTEGER NOT NULL,result TEXT NOT NULL DEFAULT '',error TEXT NOT NULL DEFAULT '');
CREATE INDEX IF NOT EXISTS search_queue ON search_jobs(status,created);
UPDATE search_jobs SET status='interrupted',error='Service restarted; submit again' WHERE status IN ('queued','running');`)
	if e != nil {
		return nil, e
	}
	exe, e := os.Executable()
	if e != nil {
		return nil, e
	}
	s := &Service{h: h, state: "off", wake: make(chan struct{}, 1), runtime: filepath.Join(filepath.Dir(exe), "websearch-runtime")}
	s.models = filepath.Join(h.Config.Data, "models")
	if root := os.Getenv("WEBSEARCH_MODELS_ROOT"); root != "" {
		s.models = root
	}
	if root := os.Getenv("WEBSEARCH_RUNTIME_ROOT"); root != "" {
		s.runtime = root
	}
	var enabled int
	if e = h.DB.QueryRow("SELECT enabled FROM search_settings WHERE id=1").Scan(&enabled); e != nil {
		return nil, e
	}
	s.enabled = enabled == 1
	h.BackupPaths = append(h.BackupPaths, "browser-profile")
	h.Mux.HandleFunc("GET /api/status", h.Owner(s.status))
	h.Mux.HandleFunc("POST /api/settings", h.Owner(s.settings))
	h.Mux.HandleFunc("POST /api/model/install", h.Owner(s.installModel))
	h.Mux.HandleFunc("GET /api/grants", h.Owner(s.grants))
	h.Mux.HandleFunc("POST /api/grants", h.Owner(s.createGrant))
	h.Mux.HandleFunc("DELETE /api/grants/{hash}", h.Owner(s.revokeGrant))
	h.Mux.HandleFunc("POST /api/search", h.Owner(s.ownerSubmit))
	h.Mux.HandleFunc("GET /api/jobs", h.Owner(s.ownerJobs))
	h.Mux.HandleFunc("DELETE /api/jobs/{id}", h.Owner(s.ownerCancel))
	h.Mux.HandleFunc("GET /api/browser/screenshot", h.Owner(s.screenshot))
	h.Mux.HandleFunc("POST /api/browser/action", h.Owner(s.manualAction))
	h.Mux.HandleFunc("POST /api/client/search", s.clientSubmit)
	h.Mux.HandleFunc("GET /api/client/jobs/{id}", s.clientJob)
	h.Mux.HandleFunc("DELETE /api/client/jobs/{id}", s.clientCancel)
	return s, nil
}
func (s *Service) signal() {
	select {
	case s.wake <- struct{}{}:
	default:
	}
}
func (s *Service) status(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var queued int
	_ = s.h.DB.QueryRow("SELECT COUNT(*) FROM search_jobs WHERE status='queued'").Scan(&queued)
	apphost.JSON(w, 200, map[string]any{"enabled": s.enabled, "state": s.state, "active": s.active, "queued": queued, "model": "Fara 1.5 4B", "version": Version, "profile": "Private Web Search browser", "control": "Screenshots, mouse and keyboard", "retentionHours": 24, "installingModel": s.setupRunning, "modelInstalled": modelsPresent(s.models), "downloadedBytes": s.setupBytes, "setupError": s.setupError})
}
func (s *Service) installModel(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.enabled || s.setupRunning || s.parent == nil {
		apphost.Problem(w, 409, "Turn off the browser service before installing its model")
		return
	}
	s.setupRunning = true
	s.setupBytes = 0
	s.setupError = ""
	ctx := s.parent
	go func() {
		e := installModels(ctx, s.models, func(n int64) { s.mu.Lock(); s.setupBytes += n; s.mu.Unlock() })
		s.mu.Lock()
		s.setupRunning = false
		if e != nil {
			s.setupError = "Model installation could not finish; check free space and retry."
		}
		s.mu.Unlock()
		s.signal()
	}()
	apphost.JSON(w, 202, map[string]string{"status": "installing"})
}
func (s *Service) settings(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Enabled bool `json:"enabled"`
	}
	if !apphost.Decode(w, r, &body) {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, e := s.h.DB.Exec("UPDATE search_settings SET enabled=? WHERE id=1", body.Enabled); e != nil {
		apphost.Problem(w, 500, "Could not save settings")
		return
	}
	s.enabled = body.Enabled
	if !body.Enabled {
		if s.cancel != nil {
			s.cancel()
		}
		if s.runtimeCancel != nil {
			s.runtimeCancel()
		}
		_, _ = s.h.DB.Exec("UPDATE search_jobs SET status='cancelled',query='',result='',error='Service turned off' WHERE status IN ('queued','running')")
	}
	s.signal()
	apphost.JSON(w, 200, map[string]bool{"enabled": s.enabled})
}
func (s *Service) grants(w http.ResponseWriter, r *http.Request) {
	rows, e := s.h.DB.Query("SELECT hash,name,created FROM search_grants ORDER BY created DESC")
	if e != nil {
		apphost.Problem(w, 500, "Could not read credentials")
		return
	}
	defer rows.Close()
	items := []map[string]any{}
	for rows.Next() {
		var hash, name string
		var created int64
		if rows.Scan(&hash, &name, &created) == nil {
			items = append(items, map[string]any{"hash": hash, "name": name, "created": created})
		}
	}
	apphost.JSON(w, 200, items)
}
func (s *Service) createGrant(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if !apphost.Decode(w, r, &body) {
		return
	}
	body.Name = strings.TrimSpace(body.Name)
	if body.Name == "" || len(body.Name) > 80 {
		apphost.Problem(w, 400, "Choose a credential name (up to 80 characters)")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	var n int
	_ = s.h.DB.QueryRow("SELECT COUNT(*) FROM search_grants").Scan(&n)
	if n >= 32 {
		apphost.Problem(w, 409, "Revoke an unused credential first")
		return
	}
	token := randomID() + randomID()
	hash := apphost.Hash(token)
	if _, e := s.h.DB.Exec("INSERT INTO search_grants VALUES(?,?,?)", hash, body.Name, time.Now().Unix()); e != nil {
		apphost.Problem(w, 500, "Could not create credential")
		return
	}
	apphost.JSON(w, 201, map[string]string{"token": token, "hash": hash, "name": body.Name})
}
func (s *Service) revokeGrant(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	hash := r.PathValue("hash")
	_, e := s.h.DB.Exec("DELETE FROM search_grants WHERE hash=?", hash)
	if e != nil {
		apphost.Problem(w, 500, "Could not revoke credential")
		return
	}
	_, _ = s.h.DB.Exec("UPDATE search_jobs SET status='cancelled',query='',result='',error='Credential revoked' WHERE principal=? AND status IN ('queued','running')", "grant:"+hash)
	s.signal()
	apphost.JSON(w, 200, map[string]bool{"revoked": true})
}
func (s *Service) credential(w http.ResponseWriter, r *http.Request) (string, bool) {
	if r.Header.Get("Origin") != "" || r.URL.RawQuery != "" {
		apphost.Problem(w, 403, "Use a server credential without browser origins or URL tokens")
		return "", false
	}
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if token == r.Header.Get("Authorization") || len(token) != 96 {
		apphost.Problem(w, 401, "Server credential required")
		return "", false
	}
	hash := apphost.Hash(token)
	var n int
	if s.h.DB.QueryRow("SELECT COUNT(*) FROM search_grants WHERE hash=?", hash).Scan(&n) != nil || n != 1 {
		apphost.Problem(w, 401, "Server credential invalid or revoked")
		return "", false
	}
	return "grant:" + hash, true
}
func (s *Service) ownerSubmit(w http.ResponseWriter, r *http.Request) {
	c, e := r.Cookie("websearch_session")
	if e != nil {
		apphost.Problem(w, 401, "Reconnect your session")
		return
	}
	s.submit(w, r, "owner:"+apphost.Current(r).User.ID, apphost.Hash(c.Value))
}
func (s *Service) clientSubmit(w http.ResponseWriter, r *http.Request) {
	if p, ok := s.credential(w, r); ok {
		s.submit(w, r, p, "")
	}
}
func (s *Service) submit(w http.ResponseWriter, r *http.Request, principal, session string) {
	var body struct {
		Query string `json:"query"`
	}
	if !apphost.Decode(w, r, &body) {
		return
	}
	body.Query = strings.TrimSpace(body.Query)
	if body.Query == "" || len(body.Query) > 4000 {
		apphost.Problem(w, 400, "Enter a research question (up to 4,000 characters)")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.enabled || s.state != "ready" {
		apphost.Problem(w, 503, "Web Search is off or its browser and model are starting")
		return
	}
	var total, per int
	e := s.h.DB.QueryRow("SELECT COUNT(*),COALESCE(SUM(CASE WHEN principal=? THEN 1 ELSE 0 END),0) FROM search_jobs WHERE status IN ('queued','running')", principal).Scan(&total, &per)
	if e != nil {
		apphost.Problem(w, 500, "Queue unavailable")
		return
	}
	if total >= 64 || per >= 8 {
		w.Header().Set("Retry-After", "10")
		apphost.Problem(w, 429, "Research queue is full; retry shortly")
		return
	}
	id := randomID()
	created := time.Now().Unix()
	if _, e = s.h.DB.Exec("INSERT INTO search_jobs(id,principal,session,query,status,created) VALUES(?,?,?,?,'queued',?)", id, principal, session, body.Query, created); e != nil {
		apphost.Problem(w, 500, "Could not queue research")
		return
	}
	s.signal()
	apphost.JSON(w, 202, Job{ID: id, Query: body.Query, Status: "queued", Created: created})
}
func (s *Service) readJob(id, principal string) (Job, error) {
	var j Job
	var result string
	e := s.h.DB.QueryRow("SELECT id,query,status,created,result,error FROM search_jobs WHERE id=? AND principal=?", id, principal).Scan(&j.ID, &j.Query, &j.Status, &j.Created, &result, &j.Error)
	if e == nil && result != "" {
		_ = json.Unmarshal([]byte(result), &j.Result)
	}
	return j, e
}
func (s *Service) clientJob(w http.ResponseWriter, r *http.Request) {
	p, ok := s.credential(w, r)
	if !ok {
		return
	}
	j, e := s.readJob(r.PathValue("id"), p)
	if e != nil {
		apphost.Problem(w, 404, "Research request not found")
		return
	}
	apphost.JSON(w, 200, j)
}
func (s *Service) ownerJobs(w http.ResponseWriter, r *http.Request) {
	rows, e := s.h.DB.Query("SELECT id,query,status,created,result,error FROM search_jobs ORDER BY created DESC,rowid DESC LIMIT 50")
	if e != nil {
		apphost.Problem(w, 500, "Could not read queue")
		return
	}
	defer rows.Close()
	jobs := []Job{}
	for rows.Next() {
		var j Job
		var result string
		if rows.Scan(&j.ID, &j.Query, &j.Status, &j.Created, &result, &j.Error) == nil {
			if result != "" {
				_ = json.Unmarshal([]byte(result), &j.Result)
			}
			jobs = append(jobs, j)
		}
	}
	apphost.JSON(w, 200, jobs)
}
func (s *Service) cancelJob(id, principal string, owner bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	query := "UPDATE search_jobs SET status='cancelled',query='',result='',error='Cancelled' WHERE id=? AND status IN ('queued','running')"
	args := []any{id}
	if !owner {
		query += " AND principal=?"
		args = append(args, principal)
	}
	res, e := s.h.DB.Exec(query, args...)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	if s.active == id && s.cancel != nil {
		s.cancel()
	}
	return nil
}
func (s *Service) ownerCancel(w http.ResponseWriter, r *http.Request) {
	if e := s.cancelJob(r.PathValue("id"), "", true); e != nil {
		apphost.Problem(w, 404, "No waiting or active request found")
		return
	}
	apphost.JSON(w, 200, map[string]bool{"cancelled": true})
}
func (s *Service) clientCancel(w http.ResponseWriter, r *http.Request) {
	p, ok := s.credential(w, r)
	if !ok {
		return
	}
	if e := s.cancelJob(r.PathValue("id"), p, false); e != nil {
		apphost.Problem(w, 404, "No waiting or active request found")
		return
	}
	apphost.JSON(w, 200, map[string]bool{"cancelled": true})
}
func (s *Service) screenshot(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	a := s.agent
	s.mu.Unlock()
	if a == nil {
		apphost.Problem(w, 503, "Browser is off")
		return
	}
	image, e := a.browser.Screenshot(r.Context())
	if e != nil {
		apphost.Problem(w, 503, "Browser unavailable")
		return
	}
	apphost.JSON(w, 200, map[string]string{"image": "data:image/png;base64," + image})
}
func (s *Service) manualAction(w http.ResponseWriter, r *http.Request) {
	var action Action
	if !apphost.Decode(w, r, &action) {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.agent == nil || s.active != "" {
		apphost.Problem(w, 409, "Wait until research completes before using the browser")
		return
	}
	if e := s.agent.browser.Act(r.Context(), action); e != nil {
		apphost.Problem(w, 400, "Unsupported browser action or invalid destination")
		return
	}
	apphost.JSON(w, 200, map[string]bool{"ok": true})
}
func (s *Service) validJob(id, principal, session string) bool {
	var status string
	if s.h.DB.QueryRow("SELECT status FROM search_jobs WHERE id=?", id).Scan(&status) != nil || status != "running" {
		return false
	}
	if strings.HasPrefix(principal, "grant:") {
		var n int
		_ = s.h.DB.QueryRow("SELECT COUNT(*) FROM search_grants WHERE hash=?", strings.TrimPrefix(principal, "grant:")).Scan(&n)
		return n == 1
	}
	var expires string
	var role string
	if s.h.DB.QueryRow("SELECT s.expires,u.role FROM host_sessions s JOIN host_users u ON u.id=s.user_id WHERE s.id=?", session).Scan(&expires, &role) != nil || role == "disabled" {
		return false
	}
	until, e := time.Parse(time.RFC3339Nano, expires)
	return e == nil && time.Now().Before(until)
}
func (s *Service) Run(ctx context.Context) error {
	s.mu.Lock()
	s.parent = ctx
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		if s.agent != nil {
			s.agent.Close()
		}
		if s.runtimeCancel != nil {
			s.runtimeCancel()
		}
	}()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		case <-s.wake:
		}
		_, _ = s.h.DB.Exec("DELETE FROM search_jobs WHERE created<?", time.Now().Add(-24*time.Hour).Unix())
		_, _ = s.h.DB.Exec("DELETE FROM search_jobs WHERE status NOT IN ('queued','running') AND id NOT IN (SELECT id FROM search_jobs ORDER BY created DESC,rowid DESC LIMIT 1000)")
		s.mu.Lock()
		enabled := s.enabled
		a := s.agent
		if !enabled {
			if a != nil {
				a.Close()
				s.agent = nil
			}
			s.state = "off"
			s.mu.Unlock()
			continue
		}
		if a == nil {
			s.state = "starting"
			runtimeCtx, cancel := context.WithCancel(ctx)
			s.runtimeCancel = cancel
			s.mu.Unlock()
			next, e := StartAgent(runtimeCtx, s.runtime, s.models, s.h.Config.Data)
			s.mu.Lock()
			if e != nil {
				s.state = "runtime unavailable"
				s.setupError = e.Error()
				s.mu.Unlock()
				if delay(ctx, 5*time.Second) != nil {
					return nil
				}
				continue
			}
			if !s.enabled {
				next.Close()
				s.state = "off"
				s.mu.Unlock()
				continue
			}
			s.agent = next
			s.state = "ready"
			s.setupError = ""
			s.mu.Unlock()
			s.signal()
			continue
		}
		s.mu.Unlock()
		_, _ = s.h.DB.Exec("UPDATE search_jobs SET status='expired',query='',error='Queue time limit reached' WHERE status='queued' AND created<?", time.Now().Add(-10*time.Minute).Unix())
		var id, principal, session, query string
		e := s.h.DB.QueryRow("SELECT id,principal,session,query FROM search_jobs WHERE status='queued' ORDER BY created,rowid LIMIT 1").Scan(&id, &principal, &session, &query)
		if errors.Is(e, sql.ErrNoRows) {
			continue
		}
		if e != nil {
			return e
		}
		s.mu.Lock()
		if !s.enabled {
			s.mu.Unlock()
			continue
		}
		changed, updateErr := s.h.DB.Exec("UPDATE search_jobs SET status='running' WHERE id=? AND status='queued'", id)
		e = updateErr
		if e != nil {
			s.mu.Unlock()
			return e
		}
		n, _ := changed.RowsAffected()
		if n != 1 {
			s.mu.Unlock()
			continue
		}
		jobCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		s.cancel = cancel
		s.active = id
		s.mu.Unlock()
		if !s.validJob(id, principal, session) {
			cancel()
			_, _ = s.h.DB.Exec("UPDATE search_jobs SET status='cancelled',query='',error='Permission expired' WHERE id=?", id)
			s.mu.Lock()
			s.active = ""
			s.cancel = nil
			s.mu.Unlock()
			continue
		}
		done := make(chan struct{})
		go func() {
			defer close(done)
			t := time.NewTicker(time.Second)
			defer t.Stop()
			for {
				if !s.validJob(id, principal, session) {
					cancel()
					return
				}
				select {
				case <-jobCtx.Done():
					return
				case <-t.C:
				}
			}
		}()
		result, err := a.Search(jobCtx, query, 32)
		cancel()
		<-done
		status := "completed"
		message := ""
		raw := ""
		if err != nil {
			status = "failed"
			message = "Browser research could not finish. Retry or check the browser and local model."
		} else {
			status = result.Status
			b, _ := json.Marshal(result)
			raw = string(b)
		}
		if !s.validJob(id, principal, session) {
			status = "cancelled"
			raw = ""
			message = "Request permission expired or was revoked"
		}
		_, _ = s.h.DB.Exec("UPDATE search_jobs SET status=?,result=?,error=?,query=CASE WHEN ?='cancelled' THEN '' ELSE query END WHERE id=? AND status='running'", status, raw, message, status, id)
		s.mu.Lock()
		s.cancel = nil
		s.active = ""
		if err != nil {
			a.Close()
			s.agent = nil
			s.state = "restarting"
		}
		s.mu.Unlock()
		s.signal()
	}
}
