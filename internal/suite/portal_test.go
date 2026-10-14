package suite

import (
	"fmt"
	"github.com/fahimc/one-workspace/internal/catalog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestProxyPreservesSecureCookieAndExternalOrigin(t *testing.T) {
	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Forwarded-Proto") != "https" || r.Host != "workspace.example.com" {
			t.Error("external origin lost")
		}
		http.SetCookie(w, &http.Cookie{Name: "filebrowser_quantum_jwt", Value: "test", Path: "/drive", Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode})
		w.WriteHeader(200)
	}))
	defer up.Close()
	w := httptest.NewRecorder()
	proxy(up.URL, "https://workspace.example.com", "", false, true).ServeHTTP(w, httptest.NewRequest("GET", "http://localhost/drive/", nil))
	cookies := w.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Path != "/" || !cookies[0].Secure || !cookies[0].HttpOnly {
		t.Fatal("shared cookie protections changed", cookies)
	}
}

func TestPhotosLocalHostAndAPIAreSeparated(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Host, "127.0.0.1:") {
			t.Error("Photos received external host")
		}
		if r.URL.Path != "/photos" {
			t.Error("wrong Photos UI path", r.URL.Path)
		}
		fmt.Fprint(w, "Photos UI")
	}))
	defer upstream.Close()
	u, _ := url.Parse(upstream.URL)
	port, _ := strconv.Atoi(u.Port())
	rt := &Runtime{Config: Config{Origin: "https://workspace.example.com", Port: 8990, Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "drive", Port: 8992}, {ID: "photos", Port: port}}}}, status: map[string]string{}}
	p := &Portal{rt: rt, flows: map[string]loginFlow{}, sessions: map[string]session{}}
	handler := p.Handler("control")
	for _, tc := range []struct {
		path string
		code int
	}{{"/photos-app/", 200}, {"/photos-app/api/users", 404}} {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, httptest.NewRequest("GET", "https://workspace.example.com"+tc.path, nil))
		if w.Code != tc.code {
			t.Fatalf("%s = %d", tc.path, w.Code)
		}
	}
}
func TestPortalRejectsHostAndAnonymousApp(t *testing.T) {
	rt := &Runtime{Config: Config{Origin: "http://localhost:8990", Port: 8990, Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "drive", Port: 8992}, {ID: "addon", Kind: "process", Route: "/extra", Port: 9000}}}}, status: map[string]string{}}
	p := &Portal{rt: rt, flows: map[string]loginFlow{}, sessions: map[string]session{}}
	h := p.Handler("control")
	for _, path := range []string{"http://localhost:8990/extra/", "http://evil.test/"} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest("GET", path, nil))
		want := 302
		if strings.Contains(path, "evil") {
			want = 403
		}
		if w.Code != want {
			t.Fatalf("%s returned %d", path, w.Code)
		}
	}
}
