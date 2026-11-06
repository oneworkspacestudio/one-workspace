package websearch

import (
	"bytes"
	"encoding/json"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

type harness struct {
	h      *apphost.Host
	s      *Service
	cookie *http.Cookie
	csrf   string
}

func fixture(t *testing.T) *harness {
	t.Helper()
	root := "C:/Projects/.cache/web-search-20261106"
	if e := os.MkdirAll(root, 0700); e != nil {
		t.Fatal(e)
	}
	dir, e := os.MkdirTemp(root, "security-")
	if e != nil {
		t.Fatal(e)
	}
	h, e := apphost.New(apphost.Config{ID: "websearch", Name: "Web Search", Version: "test", Data: filepath.Join(dir, "data"), Listen: "127.0.0.1:29443", Origin: "http://localhost:29443"}, UI())
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { h.Close() })
	s, e := New(h)
	if e != nil {
		t.Fatal(e)
	}
	if e = h.Bootstrap("synthetic-owner@example.test", "synthetic-fixture-password"); e != nil {
		t.Fatal(e)
	}
	f := &harness{h: h, s: s}
	w := f.call("POST", "/api/login", map[string]string{"email": "synthetic-owner@example.test", "password": "synthetic-fixture-password"}, "", true)
	if w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	f.cookie = w.Result().Cookies()[0]
	w = f.call("GET", "/api/me", nil, "", true)
	var me apphost.Session
	_ = json.Unmarshal(w.Body.Bytes(), &me)
	f.csrf = me.CSRF
	return f
}
func (f *harness) call(method, path string, body any, token string, origin bool, control ...string) *httptest.ResponseRecorder {
	raw, _ := json.Marshal(body)
	r := httptest.NewRequest(method, f.h.Config.Origin+path, bytes.NewReader(raw))
	r.RemoteAddr = "127.0.0.1:19001"
	if origin {
		r.Header.Set("Origin", f.h.Config.Origin)
	}
	if token != "" {
		r.Header.Set("Authorization", "Bearer "+token)
	}
	if f.cookie != nil {
		r.AddCookie(f.cookie)
	}
	r.Header.Set("X-CSRF-Token", f.csrf)
	if len(control) > 0 {
		r.Header.Set("X-WebSearch-Control", control[0])
	}
	w := httptest.NewRecorder()
	f.h.Handler().ServeHTTP(w, r)
	return w
}
func (f *harness) grant(t *testing.T, name string) (string, string) {
	t.Helper()
	w := f.call("POST", "/api/grants", map[string]string{"name": name}, "", true)
	if w.Code != 201 {
		t.Fatal(w.Code, w.Body.String())
	}
	var v map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &v)
	return v["token"], v["hash"]
}
func TestCredentialIsolationBoundsRevocationAndOff(t *testing.T) {
	f := fixture(t)
	a, ah := f.grant(t, "Synthetic agent A")
	b, _ := f.grant(t, "Synthetic agent B")
	if w := f.call("POST", "/api/client/search", map[string]string{"query": "synthetic public research"}, a, false); w.Code != 503 {
		t.Fatal("off accepted", w.Code)
	}
	f.s.enabled = true
	f.s.state = "ready"
	if w := f.call("POST", "/api/client/search", map[string]string{"query": "synthetic public research"}, a, true); w.Code != 403 {
		t.Fatal("browser accepted server token", w.Code)
	}
	if w := f.call("POST", "/api/client/search?token=secret", map[string]string{"query": "synthetic public research"}, a, false); w.Code != 403 {
		t.Fatal("query credential accepted", w.Code)
	}
	var id string
	for n := 0; n < 8; n++ {
		w := f.call("POST", "/api/client/search", map[string]string{"query": "synthetic public research"}, a, false)
		if w.Code != 202 {
			t.Fatal(w.Code, w.Body.String())
		}
		var j Job
		_ = json.Unmarshal(w.Body.Bytes(), &j)
		id = j.ID
	}
	if w := f.call("POST", "/api/client/search", map[string]string{"query": "synthetic overflow"}, a, false); w.Code != 429 {
		t.Fatal("principal bound", w.Code)
	}
	if w := f.call("GET", "/api/client/jobs/"+id, nil, b, false); w.Code != 404 {
		t.Fatal("cross credential read", w.Code)
	}
	if w := f.call("DELETE", "/api/client/jobs/"+id, nil, b, false); w.Code != 404 {
		t.Fatal("cross credential cancellation", w.Code)
	}
	if w := f.call("GET", "/api/client/jobs/"+id, nil, a, false); w.Code != 200 {
		t.Fatal("own job unavailable", w.Code)
	}
	if w := f.call("DELETE", "/api/grants/"+ah, nil, "", true); w.Code != 200 {
		t.Fatal(w.Code)
	}
	if w := f.call("GET", "/api/client/jobs/"+id, nil, a, false); w.Code != 401 {
		t.Fatal("revoked credential accepted", w.Code)
	}
	var pending int
	_ = f.h.DB.QueryRow("SELECT COUNT(*) FROM search_jobs WHERE principal=? AND status IN ('running','queued')", "grant:"+ah).Scan(&pending)
	if pending != 0 {
		t.Fatal("revocation kept jobs alive")
	}
	var leaked int
	_ = f.h.DB.QueryRow("SELECT COUNT(*) FROM search_grants WHERE hash=?", a).Scan(&leaked)
	if leaked != 0 {
		t.Fatal("plaintext credential persisted")
	}
}
func TestOwnerCSRFAndExactSessionLogout(t *testing.T) {
	f := fixture(t)
	f.s.enabled = true
	f.s.state = "ready"
	saved := f.csrf
	f.csrf = ""
	if w := f.call("POST", "/api/settings", map[string]bool{"enabled": false}, "", true); w.Code != 403 {
		t.Fatal("CSRF bypass", w.Code)
	}
	f.csrf = saved
	w := f.call("POST", "/api/search", map[string]string{"query": "Synthetic session research"}, "", true)
	if w.Code != 202 {
		t.Fatal(w.Code, w.Body.String())
	}
	var j Job
	_ = json.Unmarshal(w.Body.Bytes(), &j)
	_, _ = f.h.DB.Exec("UPDATE search_jobs SET status='running' WHERE id=?", j.ID)
	var principal, session string
	_ = f.h.DB.QueryRow("SELECT principal,session FROM search_jobs WHERE id=?", j.ID).Scan(&principal, &session)
	if !f.s.validJob(j.ID, principal, session) {
		t.Fatal("live session rejected")
	}
	if w = f.call("POST", "/auth/logout", map[string]string{}, "", true); w.Code != 200 && w.Code != 303 {
		t.Fatal(w.Code, w.Body.String())
	}
	if f.s.validJob(j.ID, principal, session) {
		t.Fatal("logged out session kept browser authority")
	}
}
