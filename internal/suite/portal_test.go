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
	"time"
)

func TestProxyMountedRedirects(t *testing.T) {
	for _, tc := range []struct{ input, want string }{
		{"/nothome/?limit=9", "/music/nothome/?limit=9"},
		{"http://workspace.example.com/nothome/", "https://workspace.example.com/music/nothome/"},
		{"/music/auth/callback", "/music/auth/callback"},
		{"https://accounts.example.com/login", "https://accounts.example.com/login"},
		{"login", "login"},
	} {
		t.Run(tc.input, func(t *testing.T) {
			up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Location", tc.input)
				w.WriteHeader(http.StatusTemporaryRedirect)
			}))
			defer up.Close()
			w := httptest.NewRecorder()
			proxy(up.URL, "https://workspace.example.com", "/music", true, false).ServeHTTP(w, httptest.NewRequest("GET", "https://workspace.example.com/music/nothome", nil))
			if got := w.Header().Get("Location"); got != tc.want {
				t.Fatalf("Location = %q, want %q", got, tc.want)
			}
		})
	}
}

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
	p := &Portal{rt: rt, flows: map[string]loginFlow{}, sessions: map[string]session{"test": {Subject: "qa", Expires: time.Now().Add(time.Hour)}}}
	handler := p.Handler("control")
	for _, tc := range []struct {
		path string
		code int
	}{{"/photos-app/", 200}, {"/photos-app/api/users", 404}} {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "https://workspace.example.com"+tc.path, nil)
		req.AddCookie(&http.Cookie{Name: "workspace_session", Value: "test"})
		handler.ServeHTTP(w, req)
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
		want := 401
		if strings.Contains(path, "evil") {
			want = 403
		}
		if w.Code != want {
			t.Fatalf("%s returned %d", path, w.Code)
		}
	}
}

func TestSafeReturnPath(t *testing.T) {
	for input, want := range map[string]string{"/notes/?note=qa": "/notes/?note=qa", "/drive/files/My%20Drive/test.docx": "/drive/files/My%20Drive/test.docx", "https://evil.test": "/", "//evil.test": "/", "/\\evil": "/", "/auth/logout": "/", "/_control/stop": "/", "": "/"} {
		if got := safeReturnPath(input); got != want {
			t.Errorf("%q = %q, want %q", input, got, want)
		}
	}
}
func TestSuiteSessionRequiredForAppsAndMedia(t *testing.T) {
	p := &Portal{rt: &Runtime{Config: Config{Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "music", Kind: "process", Route: "/music/"}, {ID: "gitea", Kind: "process", Route: "/git/"}}}}}}
	for _, path := range []string{"/drive/api/raw/qa", "/api/resources", "/photos-app/", "/music/file/qa", "/git/api/v1/user"} {
		r := httptest.NewRequest("GET", "https://workspace.example.com"+path, nil)
		if !p.isAppRequest(r) {
			t.Error("unguarded", path)
		}
		if nativeGitRequest(r) {
			t.Error("anonymous native bypass", path)
		}
	}
	r := httptest.NewRequest("GET", "https://workspace.example.com/git/qa/repo.git/info/refs", nil)
	r.Header.Set("Authorization", "Basic synthetic")
	if !nativeGitRequest(r) {
		t.Error("native Git compatibility lost")
	}
	r.URL.Path = "/git/user/settings"
	if nativeGitRequest(r) {
		t.Error("browser route bypass")
	}
}

func TestPublicShareRoutesKeepNativeShareAuthorization(t *testing.T) {
	for _, path := range []string{"/drive/public/share/qa", "/drive/public/api/resources", "/drive/public/api/share/info", "/drive/public/static/assets/qa.js", "/drive/editor", "/sdkjs/word/qa.js"} {
		if !publicDriveRoute(path) {
			t.Error("public route gated", path)
		}
	}
	for _, path := range []string{"/drive/public/api/users", "/drive/api/resources", "/drive/public/api/settings/user-defaults", "/drive/public/api/share/pinned-items"} {
		if publicDriveRoute(path) {
			t.Error("private route exempted", path)
		}
	}
	r := httptest.NewRequest("GET", "https://workspace.example.com/drive/public/api/resources", nil)
	r.AddCookie(&http.Cookie{Name: "filebrowser_quantum_jwt", Value: "old-native-session"})
	r.AddCookie(&http.Cookie{Name: "filebrowser_share_session", Value: "share-grant"})
	anonymousShareCookies(r)
	if _, err := r.Cookie("filebrowser_quantum_jwt"); err == nil {
		t.Error("private session survived anonymous public route")
	}
	if c, err := r.Cookie("filebrowser_share_session"); err != nil || c.Value != "share-grant" {
		t.Error("public share authorization lost")
	}
}

func TestLogoutClearsScopedCookiesAndDomainDriveCookie(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "https://workspace.example.com/auth/logout", nil)
	clearAppCookies(w, r, true)
	found := map[string]bool{}
	for _, cookie := range w.Result().Cookies() {
		if cookie.MaxAge != -1 || !cookie.Secure || !cookie.HttpOnly {
			t.Error("cookie protections lost", cookie.Name)
		}
		if cookie.Name == "filebrowser_quantum_jwt" && cookie.Domain == "workspace.example.com" && cookie.Path == "/" {
			found["drive-domain"] = true
		}
		if cookie.Name == "notes_session" && cookie.Path == "/notes/" {
			found["notes"] = true
		}
		if cookie.Name == "mytube_session" && cookie.Path == "/my-tube/" {
			found["my-tube"] = true
		}
		if cookie.Name == "mail_session" && cookie.Path == "/mail/" {
			found["mail"] = true
		}
	}
	for _, key := range []string{"drive-domain", "notes", "mail", "my-tube"} {
		if !found[key] {
			t.Error("cookie expiry missing", key)
		}
	}
}
