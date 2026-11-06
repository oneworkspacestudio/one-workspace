package websearch

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func controlFixture(t *testing.T) (*harness, string) {
	f := fixture(t)
	ctx, cancel := context.WithCancel(context.Background())
	f.s.parent = ctx
	t.Cleanup(func() {
		cancel()
		f.s.mu.Lock()
		if f.s.desktop != nil {
			f.s.desktop.close()
		}
		f.s.mu.Unlock()
	})
	w := f.call("POST", "/api/desktop/start", map[string]any{}, "", true)
	if w.Code != 202 {
		t.Fatal(w.Code, w.Body.String())
	}
	var data struct {
		ID string `json:"id"`
	}
	_ = json.Unmarshal(w.Body.Bytes(), &data)
	return f, data.ID
}

func TestDesktopSessionCSRFIsolationAndRevocation(t *testing.T) {
	f, id := controlFixture(t)
	saved := f.csrf
	f.csrf = ""
	if w := f.call("POST", "/api/desktop/heartbeat", nil, "", true, id); w.Code != 403 {
		t.Fatal("CSRF bypass", w.Code)
	}
	f.csrf = saved
	if w := f.call("GET", "/api/desktop/frame", nil, "", true, "wrong"); w.Code != 409 {
		t.Fatal("lease bypass", w.Code)
	}
	if w := f.call("GET", "/api/desktop/frame?token=x", nil, "", true, id); w.Code != 409 {
		t.Fatal("URL handle accepted", w.Code)
	}
	key, _ := f.grant(t, "Synthetic server client")
	cookie := f.cookie
	f.cookie = nil
	for _, path := range []string{"/api/desktop/frame", "/api/status"} {
		if w := f.call("GET", path, nil, key, false, id); w.Code != 401 {
			t.Fatal("server key gained desktop", path, w.Code)
		}
	}
	f.cookie = cookie
	login := f.call("POST", "/api/login", map[string]string{"email": "synthetic-owner@example.test", "password": "synthetic-fixture-password"}, "", true)
	f.cookie = login.Result().Cookies()[0]
	if w := f.call("GET", "/api/desktop/frame", nil, "", true, id); w.Code != 409 {
		t.Fatal("different owner cookie gained lease", w.Code)
	}
	f.cookie = cookie
	d := f.s.desktop
	if w := f.call("POST", "/auth/logout", map[string]any{}, "", true); w.Code != 200 && w.Code != 303 {
		t.Fatal(w.Code)
	}
	select {
	case <-d.ctx.Done():
	case <-time.After(2500 * time.Millisecond):
		t.Fatal("logout retained remote control")
	}
	if w := f.call("GET", "/api/desktop/frame", nil, "", true, id); w.Code != 401 {
		t.Fatal("logged out owner read frame", w.Code)
	}
}

func TestDesktopExpiryAndRoleRevocation(t *testing.T) {
	for _, kind := range []string{"expiry", "role"} {
		t.Run(kind, func(t *testing.T) {
			f, id := controlFixture(t)
			d := f.s.desktop
			if kind == "expiry" {
				f.s.mu.Lock()
				d.expires = time.Now().Add(-time.Second)
				f.s.mu.Unlock()
			} else {
				_, _ = f.h.DB.Exec("UPDATE host_users SET role='member'")
			}
			select {
			case <-d.ctx.Done():
			case <-time.After(2500 * time.Millisecond):
				t.Fatal("expired/demoted session retained control")
			}
			w := f.call("POST", "/api/desktop/input", DesktopInput{Type: "text", Text: "not sent"}, "", true, id)
			if w.Code != 409 && w.Code != 403 {
				t.Fatal("input after expiry", w.Code)
			}
		})
	}
}

func TestDesktopTakeoverInterruptsActiveWorkAndKeepsQueue(t *testing.T) {
	f := fixture(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	f.s.parent = ctx
	f.s.enabled = true
	f.s.state = "ready"
	key, hash := f.grant(t, "Synthetic queue")
	var ids []string
	for n := 0; n < 2; n++ {
		w := f.call("POST", "/api/client/search", map[string]string{"query": "synthetic research"}, key, false)
		var j Job
		_ = json.Unmarshal(w.Body.Bytes(), &j)
		ids = append(ids, j.ID)
	}
	_, _ = f.h.DB.Exec("UPDATE search_jobs SET status='running' WHERE id=?", ids[0])
	cancelled := false
	f.s.cancel = func() { cancelled = true }
	w := f.call("POST", "/api/desktop/start", nil, "", true)
	if w.Code != 202 || !cancelled {
		t.Fatal("takeover did not stop active request", w.Code)
	}
	defer f.s.desktop.close()
	j, e := f.s.readJob(ids[0], "grant:"+hash)
	if e != nil || j.Status != "interrupted" {
		t.Fatal("active job replayed", j, e)
	}
	j, e = f.s.readJob(ids[1], "grant:"+hash)
	if e != nil || j.Status != "queued" {
		t.Fatal("waiting work lost", j, e)
	}
	if w = f.call("POST", "/api/client/search", map[string]string{"query": "waiting during takeover"}, key, false); w.Code != 202 {
		t.Fatal("paused queue rejected work", w.Code)
	}
}

// The CDP stand-in is limited to unit tests. Native acceptance uses actual Chrome.
func fakeDesktopBrowser(t *testing.T, block <-chan struct{}, started chan<- struct{}) *Browser {
	t.Helper()
	up := websocket.Upgrader{}
	var once sync.Once
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, e := up.Upgrade(w, r, nil)
		if e != nil {
			return
		}
		defer conn.Close()
		for {
			var call struct {
				ID     int    `json:"id"`
				Method string `json:"method"`
			}
			if conn.ReadJSON(&call) != nil {
				return
			}
			result := any(map[string]any{})
			if call.Method == "Page.captureScreenshot" {
				if started != nil {
					once.Do(func() { started <- struct{}{} })
				}
				if block != nil {
					<-block
				}
				result = map[string]string{"data": "cHJpdmF0ZS1waXhlbHM="}
			}
			if call.Method == "Target.getTargetInfo" {
				result = map[string]any{"targetInfo": map[string]string{"url": "https://example.com/", "title": "Synthetic page"}}
			}
			if conn.WriteJSON(map[string]any{"id": call.ID, "result": result}) != nil {
				return
			}
		}
	}))
	t.Cleanup(server.Close)
	conn, _, e := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http"), nil)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { conn.Close() })
	return &Browser{conn: conn, target: "synthetic"}
}

func TestDesktopInFlightFrameIsRevokedOnLogout(t *testing.T) {
	f, id := controlFixture(t)
	release := make(chan struct{})
	started := make(chan struct{}, 1)
	defer close(release)
	b := fakeDesktopBrowser(t, release, started)
	f.s.mu.Lock()
	f.s.desktop.browser = b
	f.s.desktop.state = "ready"
	f.s.mu.Unlock()
	response := make(chan *httptest.ResponseRecorder, 1)
	go func() { response <- f.call("GET", "/api/desktop/frame", nil, "", true, id) }()
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("frame never started")
	}
	if w := f.call("POST", "/auth/logout", nil, "", true); w.Code != 200 && w.Code != 303 {
		t.Fatal(w.Code)
	}
	select {
	case w := <-response:
		if strings.Contains(w.Body.String(), "cHJpdmF0ZS1waXhlbHM=") || w.Code == 200 {
			t.Fatal("pixels returned after logout", w.Code)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("frame retained revoked authority")
	}
}

func TestDesktopInputSurfaceRejectsCodeAndLocalNavigation(t *testing.T) {
	for _, in := range []DesktopInput{{Type: "evaluate", Text: "alert(1)"}, {Type: "navigate", URL: "file:///C:/"}, {Type: "navigate", URL: "javascript:alert(1)"}, {Type: "navigate", URL: "https://user:pass@example.com"}, {Type: "navigate", URL: "https://example.com:8443/"}, {Type: "pointerDown", X: 1280, Y: 50}, {Type: "text", Text: strings.Repeat("a", 8001)}, {Type: "wheel", DeltaY: 1601}} {
		if in.valid() {
			t.Fatal("unsafe input", in.Type)
		}
	}
	for _, in := range []DesktopInput{{Type: "navigate", URL: "https://example.com/"}, {Type: "text", Text: "synthetic input"}, {Type: "pointerDown", X: 320, Y: 200, Buttons: 1, Count: 2}, {Type: "wheel", X: 640, Y: 400, DeltaY: 550}, {Type: "key", Keys: []string{"CTRL", "a"}}} {
		if !in.valid() {
			t.Fatal("valid input rejected", in.Type)
		}
	}
}
