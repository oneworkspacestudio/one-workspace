package suite

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

func TestConnectionReturnsOnlyToInstalledEntryPages(t *testing.T) {
	p := &Portal{rt: &Runtime{Config: Config{Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "pdf", Route: "/pdf"}, {ID: "vault", Route: "/vault"}, {ID: "drive", Route: "/drive/"}}}}}}
	for _, tc := range []struct {
		path  string
		valid bool
	}{
		{"/pdf/?drive=open&document=abc", true}, {"/", true}, {"/drive/", true},
		{"//evil.test", false}, {"https://evil.test", false}, {"/%2f%2fevil.test", false},
		{"/pdf/auth/logout", false}, {"/pdf/api/delete", false}, {"/_workspace/connect/drive", false},
		{"/vault", false}, {"/unknown", false}, {"/pdf/#fragment", false}, {"/pdf/\\evil", false},
	} {
		_, valid := p.connectionDestination(tc.path)
		if valid != tc.valid {
			t.Errorf("%q: %v", tc.path, valid)
		}
	}
}

func TestDriveConnectionUsesNativeSessionAndReturns(t *testing.T) {
	var drive appSignIn
	for _, app := range suiteSignIns {
		if app.id == "drive" {
			drive = app
		}
	}
	p := signInFixture(t, drive, http.StatusOK)
	p.rt.Config.Catalog.Apps = append(p.rt.Config.Catalog.Apps, catalog.App{ID: "pdf", Route: "/pdf"})
	p.sessions = map[string]session{"fixture-portal-session": {Subject: "owner", Expires: time.Now().Add(time.Hour)}}
	r := signInRequest("/_workspace/connect/drive?returnTo=%2Fpdf%2F%3Fdrive%3Dopen")
	r.SetPathValue("id", "drive")
	w := httptest.NewRecorder()
	p.connectApp(w, r)
	if w.Code != 303 || w.Header().Get("Location") != "/drive/" {
		t.Fatal(w.Code, w.Header())
	}
	var connect *http.Cookie
	for _, c := range w.Result().Cookies() {
		if c.Name == "workspace_connect_drive" {
			connect = c
		}
	}
	if connect == nil || !connect.HttpOnly || !connect.Secure || connect.MaxAge != 300 {
		t.Fatal("unsafe connection cookie")
	}
	r = signInRequest("/drive/")
	r.AddCookie(connect)
	w = httptest.NewRecorder()
	if !p.seamlessApp(w, r) || w.Header().Get("Location") != "/pdf/?drive=open" {
		t.Fatal("connection did not return", w.Code, w.Header())
	}
	connect.Value = "other-session." + url.QueryEscape("/pdf/")
	r = signInRequest("/drive/")
	r.AddCookie(connect)
	w = httptest.NewRecorder()
	if p.finishConnection(w, r, drive) {
		t.Fatal("different account connection reused")
	}
}

func TestAppUnavailableHasNavigationAndAPIRetainsStatus(t *testing.T) {
	up := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	up.Close()
	handler := proxy(up.URL, "https://workspace.example.test", "/pdf", true, false)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, signInRequest("/pdf/"))
	if w.Code != 502 || !strings.Contains(w.Body.String(), "Try opening the app again") || !strings.Contains(w.Body.String(), `href="/pdf/"`) || w.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("missing recovery page", w.Code, w.Body.String())
	}
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, httptest.NewRequest("GET", "/pdf/api/documents", nil))
	if w.Code != 502 || strings.Contains(w.Body.String(), "<!doctype") {
		t.Fatal("API error became HTML")
	}
}

func TestOwnerOnlyMonitorIsNotOfferedToMembers(t *testing.T) {
	p := &Portal{rt: &Runtime{Config: Config{Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "monitor", Name: "My Monitor"}, {ID: "drive", Name: "Drive"}}}}}}
	if len(p.tiles(false)) != 1 || len(p.tiles(true)) != 3 {
		t.Fatal("owner-only navigation leak")
	}
}

func TestUnconfiguredMailExplainsSetupBeforeLogin(t *testing.T) {
	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthz" {
			t.Error("should not start session probe")
		}
		w.Write([]byte(`{"configured":false}`))
	}))
	defer up.Close()
	u, _ := url.Parse(up.URL)
	var port int
	fmtPort := u.Port()
	for _, c := range fmtPort {
		port = port*10 + int(c-'0')
	}
	p := &Portal{rt: &Runtime{Config: Config{Owner: "owner", Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "mail", Route: "/mail", Port: port}}}}}, sessions: map[string]session{"fixture-portal-session": {Subject: "owner", Expires: time.Now().Add(time.Hour)}}}
	w := httptest.NewRecorder()
	if !p.seamlessApp(w, signInRequest("/mail/")) || !strings.Contains(w.Body.String(), "Set up Mail") {
		t.Fatal("setup guidance missing", w.Body.String())
	}
	p.sessions["fixture-portal-session"] = session{Subject: "member", Expires: time.Now().Add(time.Hour)}
	w = httptest.NewRecorder()
	p.seamlessApp(w, signInRequest("/mail/"))
	if strings.Contains(w.Body.String(), "/_workspace/setup/mail") {
		t.Fatal("owner setup offered to member")
	}
}
