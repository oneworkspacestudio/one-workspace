package suite

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

func signInFixture(t *testing.T, app appSignIn, status int) *Portal {
	t.Helper()
	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.id == "mail" && r.URL.Path == "/api/status" {
			w.Write([]byte(`{"mainMailURL":""}`))
			return
		}
		if app.id == "mail" && r.URL.Path == "/healthz" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"configured":true}`))
			return
		}
		if r.URL.RequestURI() != app.probe {
			t.Errorf("unexpected session probe %s", r.URL.RequestURI())
		}
		if r.Header.Get("Authorization") != "" {
			t.Error("probe forwarded bearer credentials")
		}
		for _, c := range r.Cookies() {
			allowed := false
			for _, name := range app.cookies {
				allowed = allowed || c.Name == name
			}
			if !allowed {
				t.Errorf("probe forwarded another app's cookie %s", c.Name)
			}
		}
		if app.id == "gitea" && status == http.StatusUnauthorized {
			http.Redirect(w, r, "/git/user/login", http.StatusSeeOther)
			return
		}
		if app.id == "my-tube" && status == http.StatusUnauthorized {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"mode":"managed","setup":false,"user":{"id":""}}`))
			return
		}
		w.WriteHeader(status)
		if app.id == "my-tube" && status == http.StatusOK {
			w.Write([]byte(`{"mode":"managed","setup":false,"user":{"id":"fixture-user"}}`))
		}
	}))
	t.Cleanup(up.Close)
	u, _ := url.Parse(up.URL)
	port, _ := strconv.Atoi(u.Port())
	return &Portal{rt: &Runtime{Config: Config{Origin: "https://workspace.example.test", Catalog: catalog.Catalog{Apps: []catalog.App{{ID: app.id, Port: port, Route: app.route}}}}}, secure: true}
}

func signInRequest(path string) *http.Request {
	r := httptest.NewRequest("GET", "https://workspace.example.test"+path, nil)
	r.Header.Set("Accept", "text/html")
	r.Header.Set("Authorization", "Bearer must-not-be-forwarded")
	r.AddCookie(&http.Cookie{Name: "workspace_session", Value: "fixture-portal-session"})
	r.AddCookie(&http.Cookie{Name: "unrelated_app_session", Value: "fixture"})
	return r
}

func TestEveryAppStartsNativeSignInWithoutAnotherLoginScreen(t *testing.T) {
	for _, app := range suiteSignIns {
		t.Run(app.id, func(t *testing.T) {
			p := signInFixture(t, app, http.StatusUnauthorized)
			r := signInRequest(app.route + "/")
			w := httptest.NewRecorder()
			if app.id == "chat" {
				if p.seamlessApp(w, r) {
					t.Fatal("Dedicated Chat must not be probed using shared portal cookies")
				}
				return
			}
			if !p.seamlessApp(w, r) || w.Code != http.StatusSeeOther || !strings.HasPrefix(w.Header().Get("Location"), app.login) {
				t.Fatalf("did not start native sign-in: %d %s", w.Code, w.Header().Get("Location"))
			}
		})
	}
}

func TestEveryAuthenticatedAppKeepsItsNativeSession(t *testing.T) {
	for _, app := range suiteSignIns {
		t.Run(app.id, func(t *testing.T) {
			p := signInFixture(t, app, http.StatusOK)
			r := signInRequest(app.route + "/")
			for _, name := range app.cookies {
				r.AddCookie(&http.Cookie{Name: name, Value: "native-fixture"})
			}
			if p.seamlessApp(httptest.NewRecorder(), r) {
				t.Fatal("authenticated app was sent through sign-in again")
			}
		})
	}
}

func TestDriveRelatedAppsShareNativeSignIn(t *testing.T) {
	for _, path := range []string{"/drive/files/My%20Drive/Documents", "/drive/studio", "/photos-app/"} {
		t.Run(path, func(t *testing.T) {
			var app appSignIn
			for _, entry := range suiteSignIns {
				if entry.id == "drive" {
					app = entry
					break
				}
			}
			if app.id == "" {
				t.Fatal("Drive sign-in contract missing")
			}
			p := signInFixture(t, app, http.StatusUnauthorized)
			w := httptest.NewRecorder()
			if !p.seamlessApp(w, signInRequest(path)) || !strings.HasPrefix(w.Header().Get("Location"), app.login+"?") {
				t.Fatalf("Drive-related app lost automatic authentication: %d", w.Code)
			}
		})
	}
}

func TestAppServerFailureDoesNotTriggerSignIn(t *testing.T) {
	for _, app := range suiteSignIns {
		p := signInFixture(t, app, http.StatusServiceUnavailable)
		if p.seamlessApp(httptest.NewRecorder(), signInRequest(app.route+"/")) {
			t.Error("server failure was treated as signed out", app.id)
		}
	}
}

func TestAppCanonicalRouteDoesNotProbeWithoutScopedCookies(t *testing.T) {
	p := &Portal{rt: &Runtime{Config: Config{Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "notes"}}}}}}
	w := httptest.NewRecorder()
	if !p.seamlessApp(w, signInRequest("/notes?note=fixture")) || w.Header().Get("Location") != "/notes/?note=fixture" {
		t.Fatal("missing trailing slash may start a duplicate login")
	}
}
