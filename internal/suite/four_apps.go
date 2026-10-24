package suite

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func foundationApp(id string) bool {
	return id == "cinema" || id == "myai" || id == "my3d" || id == "create" || id == "maps" || id == "meet" || id == "clips" || id == "social" || id == "slicer" || id == "tasks" || id == "pdf" || id == "monitor" || id == "search"
}
func monitorCapability(secret string) string {
	if secret == "" {
		return ""
	}
	sum := sha256.Sum256([]byte("workspace-monitor-registry-v1\x00" + secret))
	return hex.EncodeToString(sum[:])
}
func (p *Portal) monitorRegistry(mux *http.ServeMux) {
	mux.HandleFunc("GET /_integration/monitor", func(w http.ResponseWriter, r *http.Request) {
		host, _, e := net.SplitHostPort(r.RemoteAddr)
		ip := net.ParseIP(host)
		want := monitorCapability(p.rt.Secrets.Clients["monitor"].Secret)
		got := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if e != nil || ip == nil || !ip.IsLoopback() || want == "" || subtle.ConstantTimeCompare([]byte(want), []byte(got)) != 1 {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		states := p.rt.Status()
		apps := []map[string]string{}
		for _, a := range p.rt.Config.Catalog.Apps {
			if a.Kind != "process" || a.Health == "" {
				continue
			}
			apps = append(apps, map[string]string{"ID": a.ID, "Name": a.Name, "URL": fmt.Sprintf("http://127.0.0.1:%d%s", a.Port, a.Health), "State": states[a.ID]})
		}
		update := "Idle"
		var journal UpdateJournal
		if platform.LoadJSON(updateJournalPath(p.rt.Config), &journal) == nil {
			update = journal.Phase
		}
		backup := "Not scheduled"
		if p.rt.backups != nil {
			settings, state := p.rt.backups.snapshot()
			if settings.Enabled {
				backup = "Scheduled"
			}
			if state.Phase != "" {
				backup = state.Phase
			}
			if !state.LastSuccess.IsZero() {
				backup += " · last verified " + state.LastSuccess.UTC().Format(time.RFC3339)
			}
		}
		json.NewEncoder(w).Encode(map[string]any{"apps": apps, "update": update, "backup": backup})
	})
}
func (p *Portal) publicMonitor(r *http.Request) bool {
	a, ok := p.rt.Config.Catalog.Find("monitor")
	return ok && a.Route == "/monitor" && r.Method == "GET" && r.URL.Path == "/monitor/status"
}
func (p *Portal) rememberFoundationSessions(r *http.Request) {
	workspace, e := r.Cookie("workspace_session")
	if e != nil {
		return
	}
	for _, a := range p.rt.Config.Catalog.Apps {
		if !foundationApp(a.ID) || !strings.HasPrefix(r.URL.Path, a.Route+"/") {
			continue
		}
		c, e := r.Cookie(a.ID + "_session")
		if e != nil || len(c.Value) != 96 {
			continue
		}
		p.mu.Lock()
		s, ok := p.sessions[workspace.Value]
		if ok && time.Now().Before(s.Expires) {
			if s.AppSessions == nil {
				s.AppSessions = map[string]string{}
			}
			s.AppSessions[a.ID] = c.Value
			p.sessions[workspace.Value] = s
		}
		p.mu.Unlock()
	}
}
func (p *Portal) closeFoundationSessions(r *http.Request) {
	s, ok := p.user(r)
	if !ok {
		return
	}
	for id, token := range s.AppSessions {
		a, ok := p.rt.Config.Catalog.Find(id)
		if !ok || !foundationApp(id) {
			continue
		}
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		req, e := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://127.0.0.1:%d/auth/logout", a.Port), nil)
		if e == nil {
			req.Header.Set("Origin", p.rt.Config.Origin)
			req.AddCookie(&http.Cookie{Name: id + "_session", Value: token})
			client := http.Client{Timeout: 3 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
			if resp, e := client.Do(req); e == nil {
				resp.Body.Close()
			}
		}
		cancel()
	}
}

// These private folders are included in the existing stop/checkpoint/verified
// backup transaction, rather than introducing a second central backup path.
func foundationData(c Config, id string) string { return filepath.Join(c.Data, id) }

func searchInstalled(c Config) bool { _, ok := c.Catalog.Find("search"); return ok }
