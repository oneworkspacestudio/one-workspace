package suite

import (
	"bytes"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestProviderLogoutUsesSupportedOIDCNavigation(t *testing.T) {
	var html bytes.Buffer
	err := logoutPage.Execute(&html, map[string]string{"Endpoint": "https://accounts.example.test/api/logout", "Hint": "synthetic-id-token", "Client": "workspace-portal", "Return": "https://workspace.example.test/auth/signed-out"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html.String(), `method="get"`) || !strings.Contains(html.String(), `name="id_token_hint"`) || !strings.Contains(html.String(), `name="post_logout_redirect_uri"`) {
		t.Fatal("Provider logout is no longer an OIDC navigation")
	}
	// MyAccounts accepts RP-initiated GET logout with a token bound to a
	// registered client. Cross-origin browser POSTs are rejected by its CSRF
	// boundary. The response keeps Referrer-Policy:no-referrer and no-store.
}

func TestHostingLogoutRevokesItsPathScopedSession(t *testing.T) {
	called := false
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, e := r.Cookie("myhosting_session")
		if r.Method != "POST" || r.URL.Path != "/auth/logout" || e != nil || cookie.Value != "synthetic-app-session" || r.Header.Get("Origin") != "https://workspace.example.test" {
			t.Error("Unexpected app logout request")
		}
		called = true
		w.WriteHeader(200)
	}))
	defer backend.Close()
	u, _ := url.Parse(backend.URL)
	port, _ := strconv.Atoi(u.Port())
	p := &Portal{rt: &Runtime{Config: Config{Origin: "https://workspace.example.test", Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "hosting", Route: "/hosting", Port: port}}}}}, sessions: map[string]session{"synthetic-portal-session": {Subject: "fixture-owner", Expires: time.Now().Add(time.Hour)}}}
	appRequest := httptest.NewRequest("GET", "https://workspace.example.test/hosting/api/me", nil)
	appRequest.AddCookie(&http.Cookie{Name: "workspace_session", Value: "synthetic-portal-session"})
	appRequest.AddCookie(&http.Cookie{Name: "myhosting_session", Value: "synthetic-app-session"})
	p.rememberHostingSession(appRequest)
	rootRequest := httptest.NewRequest("POST", "https://workspace.example.test/auth/logout", nil)
	rootRequest.AddCookie(&http.Cookie{Name: "workspace_session", Value: "synthetic-portal-session"})
	p.closeHostingSession(rootRequest)
	if !called {
		t.Fatal("The browser-scoped native session was not revoked")
	}
	foreign := httptest.NewRequest("GET", "https://workspace.example.test/hosting/api/me", nil)
	foreign.AddCookie(&http.Cookie{Name: "myhosting_session", Value: "untrusted-cookie"})
	p.rememberHostingSession(foreign)
	if p.sessions["synthetic-portal-session"].HostingSession != "synthetic-app-session" {
		t.Fatal("Unauthenticated request changed the remembered session")
	}
}
