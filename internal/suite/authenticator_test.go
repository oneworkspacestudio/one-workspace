package suite

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestAuthenticatorPublicBootstrapIsNarrow(t *testing.T) {
	for _, v := range []struct {
		method, path string
		allowed      bool
	}{{"GET", "/authenticator/", true}, {"GET", "/authenticator/sw.js", true}, {"GET", "/authenticator/assets/index-Ab12.js", true}, {"POST", "/authenticator/", false}, {"GET", "/authenticator/api/secrets", false}, {"GET", "/authenticator/assets/../key.js", false}, {"GET", "/authenticator/sw.js?seed=x", false}, {"GET", "/other/sw.js", false}} {
		r := httptest.NewRequest(v.method, v.path, nil)
		if got := publicAuthenticatorRequest(r); got != v.allowed {
			t.Fatalf("%s %s: %v", v.method, v.path, got)
		}
	}
}
func TestPortalMfaTransitionRejectsOldSessions(t *testing.T) {
	active := true
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		id, secret, ok := r.BasicAuth()
		if !ok || id != "client" || secret != "synthetic" {
			t.Error("missing client binding")
		}
		if active {
			w.Write([]byte(`{"active":true,"sub":"synthetic-member"}`))
		} else {
			w.Write([]byte(`{"active":false}`))
		}
	}))
	defer server.Close()
	p := &Portal{mu: sync.Mutex{}, rt: &Runtime{Config: Config{Issuer: server.URL}}, oauth: oauth2.Config{ClientID: "client", ClientSecret: "synthetic"}, sessions: map[string]session{"session": {Subject: "synthetic-member", AccessToken: "synthetic-access", Expires: time.Now().Add(time.Hour)}}}
	request := func() *http.Request {
		r := httptest.NewRequest("GET", "http://workspace.test/drive/", nil)
		r.AddCookie(&http.Cookie{Name: "workspace_session", Value: "session"})
		return r
	}
	r := request()
	if _, ok := p.user(r); !ok {
		t.Fatal("MFA-complete session denied")
	}
	if _, ok := p.user(r); !ok || calls != 1 {
		t.Fatal("per-request verification was not reused")
	}
	active = false
	if _, ok := p.user(request()); ok {
		t.Fatal("revoked pre-MFA session retained access")
	}
	if _, ok := p.sessions["session"]; ok {
		t.Fatal("invalid session retained")
	}
}
