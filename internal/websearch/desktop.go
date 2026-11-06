package websearch

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
)

const controlTTL = 90 * time.Second

type desktopLease struct {
	id, session, state, problem string
	expires                     time.Time
	browser                     *Browser
	agent                       *Agent
	ctx                         context.Context
	cancel                      context.CancelFunc
	browserCancel               context.CancelFunc
	io                          sync.Mutex
}

func (d *desktopLease) close() {
	if d.cancel != nil {
		d.cancel()
	}
	if d.agent != nil {
		d.agent.Close()
	} else if d.browser != nil {
		d.browser.Close()
	}
	if d.browserCancel != nil {
		d.browserCancel()
	}
}
func ownerSession(r *http.Request) string {
	c, e := r.Cookie("websearch_session")
	if e != nil {
		return ""
	}
	return apphost.Hash(c.Value)
}
func (s *Service) validOwnerSession(id string) bool {
	var expires, role string
	if s.h.DB.QueryRow("SELECT s.expires,u.role FROM host_sessions s JOIN host_users u ON u.id=s.user_id WHERE s.id=?", id).Scan(&expires, &role) != nil || role != "owner" {
		return false
	}
	end, e := time.Parse(time.RFC3339Nano, expires)
	return e == nil && time.Now().Before(end)
}
func (s *Service) controlValid(d *desktopLease) bool {
	return d != nil && time.Now().Before(d.expires) && s.validOwnerSession(d.session)
}
func (s *Service) registerDesktop() {
	s.h.Mux.HandleFunc("POST /api/desktop/start", s.h.Owner(s.desktopStart))
	s.h.Mux.HandleFunc("POST /api/desktop/heartbeat", s.h.Owner(s.desktopHeartbeat))
	s.h.Mux.HandleFunc("POST /api/desktop/stop", s.h.Owner(s.desktopStop))
	s.h.Mux.HandleFunc("POST /api/desktop/reconnect", s.h.Owner(s.reconnectDesktop))
	s.h.Mux.HandleFunc("GET /api/desktop/frame", s.h.Owner(s.desktopFrame))
	s.h.Mux.HandleFunc("POST /api/desktop/input", s.h.Owner(s.desktopInput))
}
func (s *Service) desktopStart(w http.ResponseWriter, r *http.Request) {
	session := ownerSession(r)
	if !s.validOwnerSession(session) {
		apphost.Problem(w, 401, "Reconnect your owner session")
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.parent == nil {
		apphost.Problem(w, 503, "Web Search is starting")
		return
	}
	if d := s.desktop; d != nil {
		if s.controlValid(d) && d.session == session {
			d.expires = time.Now().Add(controlTTL)
			apphost.JSON(w, 200, map[string]any{"id": d.id, "state": d.state, "expires": d.expires})
			return
		}
		apphost.Problem(w, 409, "Another remote session is closing or in control. Close it or wait 90 seconds, then reopen")
		return
	}
	if _, e := s.h.DB.Exec("UPDATE search_jobs SET status='interrupted',result='',error='Owner took browser control; submit again when finished' WHERE status='running'"); e != nil {
		apphost.Problem(w, 500, "Could not pause the research queue")
		return
	}
	d := &desktopLease{id: randomID(), session: session, state: "waiting", expires: time.Now().Add(controlTTL)}
	d.ctx, d.cancel = context.WithCancel(s.parent)
	s.desktop = d
	go s.watchControl(d)
	if s.cancel != nil {
		s.cancel()
	}
	// Interrupt a model startup that has not yet produced a browser. A healthy
	// browser is handed over intact by Run, including forms and challenge pages.
	if s.agent == nil && s.runtimeCancel != nil {
		s.runtimeCancel()
	}
	s.signal()
	apphost.JSON(w, 202, map[string]any{"id": d.id, "state": d.state, "expires": d.expires})
}
func (s *Service) watchControl(d *desktopLease) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
		}
		s.mu.Lock()
		valid := s.desktop == d && s.controlValid(d)
		if !valid {
			d.expires = time.Time{}
			d.cancel()
			s.signal()
		}
		s.mu.Unlock()
		if !valid {
			return
		}
	}
}
func (s *Service) controlRequest(w http.ResponseWriter, r *http.Request) *desktopLease {
	// Called with s.mu held. The handle alone grants no authority: the exact
	// native owner cookie, live owner role and expiry must also match every time.
	d := s.desktop
	if r.URL.RawQuery != "" || d == nil || r.Header.Get("X-WebSearch-Control") != d.id || ownerSession(r) != d.session || !s.controlValid(d) {
		apphost.Problem(w, 409, "Remote browser control ended. Reopen it to continue")
		return nil
	}
	return d
}
func (s *Service) desktopHeartbeat(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	d := s.controlRequest(w, r)
	if d == nil {
		return
	}
	d.expires = time.Now().Add(controlTTL)
	apphost.JSON(w, 200, map[string]any{"state": d.state, "error": d.problem, "expires": d.expires})
}
func (s *Service) desktopStop(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	d := s.controlRequest(w, r)
	if d == nil {
		return
	}
	d.expires = time.Time{}
	if d.cancel != nil {
		d.cancel()
	}
	s.signal()
	apphost.JSON(w, 200, map[string]bool{"closed": true})
}
func (s *Service) controlOperation(w http.ResponseWriter, r *http.Request, fn func(context.Context, *Browser) (any, error)) {
	s.mu.Lock()
	d := s.controlRequest(w, r)
	if d == nil {
		s.mu.Unlock()
		return
	}
	if d.state != "ready" || d.browser == nil {
		state := d.state
		problem := d.problem
		s.mu.Unlock()
		apphost.JSON(w, 202, map[string]string{"state": state, "error": problem})
		return
	}
	leaseCtx := d.ctx
	s.mu.Unlock()
	d.io.Lock()
	defer d.io.Unlock()
	s.mu.Lock()
	valid := s.desktop == d && s.controlValid(d)
	state := d.state
	s.mu.Unlock()
	if !valid {
		apphost.Problem(w, 409, "Remote browser control ended")
		return
	}
	if state != "ready" {
		apphost.JSON(w, 202, map[string]string{"state": state})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	stop := context.AfterFunc(leaseCtx, cancel)
	defer stop()
	result, e := fn(ctx, d.browser)
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.desktop != d || !s.controlValid(d) {
		apphost.Problem(w, 409, "Remote browser control ended")
		return
	}
	if e != nil {
		// An interrupted CDP read invalidates its connection. Reopening the lease
		// reconnects the same live tab rather than losing a sign-in form.
		d.state = "disconnected"
		d.problem = "Browser connection interrupted. Choose Reconnect to continue."
		apphost.Problem(w, 503, d.problem)
		return
	}
	apphost.JSON(w, 200, result)
}
func (s *Service) desktopFrame(w http.ResponseWriter, r *http.Request) {
	s.controlOperation(w, r, func(ctx context.Context, b *Browser) (any, error) {
		image, loc, tabs, e := b.DesktopFrame(ctx)
		return map[string]any{"state": "ready", "image": image, "url": loc.URL, "title": loc.Title, "tabs": tabs, "width": 1280, "height": 800}, e
	})
}
func (s *Service) desktopInput(w http.ResponseWriter, r *http.Request) {
	var input DesktopInput
	if !apphost.Decode(w, r, &input) {
		return
	}
	if !input.valid() {
		apphost.Problem(w, 400, "Unsupported browser input or destination")
		return
	}
	s.controlOperation(w, r, func(ctx context.Context, b *Browser) (any, error) {
		e := b.DesktopInput(ctx, input)
		return map[string]bool{"ok": e == nil}, e
	})
}

// Called by the sole queue worker with s.mu held. It returns with the mutex
// released when control owns this tick, preventing LLM/manual input overlap.
func (s *Service) runDesktop(ctx context.Context) bool {
	d := s.desktop
	if d == nil {
		return false
	}
	if !s.controlValid(d) {
		d.close()
		s.desktop = nil
		s.state = "off"
		if s.enabled {
			s.state = "starting"
		}
		s.mu.Unlock()
		s.signal()
		return true
	}
	s.state = "owner control"
	if d.browser == nil && s.agent != nil {
		d.agent = s.agent
		d.browser = s.agent.browser
		d.agent.StopModel()
		s.agent = nil
		s.runtimeCancel = nil
	}
	if d.state != "waiting" {
		s.mu.Unlock()
		return true
	}
	d.state = "starting"
	b := d.browser
	if b != nil {
		select {
		case <-b.done:
			if d.agent != nil {
				d.agent.Close()
				d.agent = nil
			} else {
				b.Close()
			}
			d.browser = nil
			b = nil
		default:
		}
	}
	var browserCtx context.Context
	if b == nil {
		browserCtx, d.browserCancel = context.WithCancel(ctx)
	}
	s.mu.Unlock()
	var e error
	d.io.Lock()
	if b != nil {
		e = b.Reconnect(d.ctx)
	} else {
		stop := context.AfterFunc(d.ctx, d.browserCancel)
		b, e = startPortableBrowser(browserCtx, s.runtime, s.h.Config.Data)
		stop()
	}
	d.io.Unlock()
	s.mu.Lock()
	if s.desktop != d || !s.controlValid(d) {
		if b != nil {
			b.Close()
		}
		if d.browserCancel != nil {
			d.browserCancel()
		}
		s.mu.Unlock()
		s.signal()
		return true
	}
	if e != nil {
		if d.browserCancel != nil {
			d.browserCancel()
			d.browserCancel = nil
		}
		d.state = "disconnected"
		d.problem = "Browser could not connect. Choose Reconnect or repair Web Search in Workspace."
	} else {
		d.browser = b
		d.state = "ready"
		d.problem = ""
	}
	s.mu.Unlock()
	return true
}

func (s *Service) reconnectDesktop(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	d := s.controlRequest(w, r)
	if d == nil {
		return
	}
	if d.state == "disconnected" || d.state == "ready" {
		d.state = "waiting"
		d.problem = ""
		s.signal()
	}
	apphost.JSON(w, 200, map[string]string{"state": d.state})
}
