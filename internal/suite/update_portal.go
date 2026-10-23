package suite

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

type updateWindow struct {
	PID            int    `json:"pid"`
	Executable     string `json:"executable"`
	URL            string `json:"url"`
	RemoteProtocol int    `json:"remoteProtocol"`
}

// Refresh idle older controllers so the owner receives the current reconnect UI.
// Active operations still use the existing busy check before any retirement.
const RemoteUpdateProtocol = 2

func (p *Portal) updateWindow() (updateWindow, *url.URL, error) {
	var window updateWindow
	if e := platform.LoadJSON(filepath.Join(p.rt.Config.Data, "updates", "window.json"), &window); e != nil {
		return window, nil, e
	}
	expected := filepath.Join(p.rt.Config.Data, "maintenance", "WorkspaceUpdate.exe")
	if !strings.EqualFold(filepath.Clean(window.Executable), expected) || !platform.ProcessRunning(window.PID, expected) {
		return window, nil, errors.New("update controller is not running")
	}
	u, e := url.Parse(window.URL)
	if e != nil || u.Scheme != "http" || u.Hostname() != "localhost" || u.Port() == "" || u.Path != "/" || u.RawQuery != "" || u.User != nil || len(u.Fragment) != 64 {
		return window, nil, errors.New("invalid update controller address")
	}
	return window, u, nil
}

func (p *Portal) ensureUpdateWindow(ctx context.Context) (*url.URL, error) {
	p.maintenanceMu.Lock()
	defer p.maintenanceMu.Unlock()
	old, u, e := p.updateWindow()
	if e == nil && old.RemoteProtocol == RemoteUpdateProtocol {
		return u, nil
	}
	if e == nil {
		// An idle legacy controller can be retired without requiring access to the
		// server desktop. Never interrupt an operation already running there.
		client := &http.Client{Timeout: 5 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
		origin := "http://" + u.Host
		req, _ := http.NewRequestWithContext(ctx, "GET", origin+"/api/state", nil)
		req.Header.Set("X-Workspace-Update", u.Fragment)
		response, err := client.Do(req)
		if err != nil {
			return nil, errors.New("An existing update controller is not responding. Retry after it finishes.")
		}
		var state struct {
			Busy bool `json:"busy"`
		}
		err = json.NewDecoder(io.LimitReader(response.Body, 16<<20)).Decode(&state)
		response.Body.Close()
		if err != nil || response.StatusCode != 200 || state.Busy {
			return nil, errors.New("An update operation is already running. Wait for it to finish, then refresh this page.")
		}
		req, _ = http.NewRequestWithContext(ctx, "POST", origin+"/api/close", nil)
		req.Header.Set("X-Workspace-Update", u.Fragment)
		req.Header.Set("Origin", origin)
		response, err = client.Do(req)
		if err != nil {
			return nil, errors.New("Could not restart the idle update controller. Retry shortly.")
		}
		response.Body.Close()
		if response.StatusCode != 202 {
			return nil, errors.New("The update controller is busy. Retry after it finishes.")
		}
		for i := 0; i < 40 && platform.ProcessRunning(old.PID, old.Executable); i++ {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(100 * time.Millisecond):
			}
		}
		if platform.ProcessRunning(old.PID, old.Executable) {
			return nil, errors.New("The previous update controller is closing. Refresh to continue.")
		}
	}
	if e = LaunchUpdatesBackground(p.rt.Config); e != nil {
		return nil, e
	}
	for i := 0; i < 60; i++ {
		window, target, err := p.updateWindow()
		if err == nil && window.RemoteProtocol == RemoteUpdateProtocol {
			return target, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(250 * time.Millisecond):
		}
	}
	return nil, errors.New("The update controller is starting slowly. Refresh this page to reconnect; no update has been started.")
}

func (p *Portal) updateRoutes(mux *http.ServeMux) {
	serve := func(w http.ResponseWriter, r *http.Request) {
		if !p.remoteUpdateOwner(w, r) {
			return
		}
		endpoint := "/"
		if r.URL.Path != "/updates" && r.URL.Path != "/updates/" && r.URL.Path != "/updates/open" {
			action := strings.TrimPrefix(r.URL.Path, "/updates/api/")
			methods := map[string]string{"state": "GET", "check": "POST", "apply": "POST", "reinstall": "POST", "recover": "POST", "rollback": "POST", "close": "POST"}
			if methods[action] != r.Method {
				http.NotFound(w, r)
				return
			}
			endpoint = "/api/" + action
		}
		target, e := p.ensureUpdateWindow(r.Context())
		if e != nil {
			http.Error(w, e.Error(), http.StatusServiceUnavailable)
			return
		}
		if r.URL.Path == "/updates/open" {
			http.Redirect(w, r, "/updates", http.StatusSeeOther)
			return
		}
		origin := "http://" + target.Host
		req, e := http.NewRequestWithContext(r.Context(), r.Method, origin+endpoint, http.MaxBytesReader(w, r.Body, 64<<10))
		if e != nil {
			http.Error(w, "Could not contact update controller", 500)
			return
		}
		req.Host = target.Host
		req.Header.Set("Origin", origin)
		req.Header.Set("X-Workspace-Update", target.Fragment)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{Timeout: 30 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
		resp, e := client.Do(req)
		if e != nil {
			http.Error(w, "Update controller is reconnecting. Refresh to continue.", 503)
			return
		}
		defer resp.Body.Close()
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'unsafe-inline'; style-src 'unsafe-inline'; connect-src 'self'; frame-ancestors 'none'; base-uri 'none'")
		w.WriteHeader(resp.StatusCode)
		_, _ = io.Copy(w, io.LimitReader(resp.Body, 16<<20))
	}
	mux.HandleFunc("GET /updates", serve)
	mux.HandleFunc("GET /updates/{$}", serve)
	mux.HandleFunc("POST /updates/open", serve)
	mux.HandleFunc("GET /updates/api/{action}", serve)
	mux.HandleFunc("POST /updates/api/{action}", serve)
}
