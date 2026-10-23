package suite

import (
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func flowTestPortal() *Portal {
	return &Portal{rt: &Runtime{Config: Config{Origin: "https://workspace.example.test"}}, secure: true,
		oauth: oauth2.Config{ClientID: "synthetic", RedirectURL: "https://workspace.example.test/auth/callback", Endpoint: oauth2.Endpoint{AuthURL: "https://accounts.example.test/authorize"}},
		flows: map[string]loginFlow{}, sessions: map[string]session{}}
}

func beginFlow(t *testing.T, p *Portal, target string) (string, *http.Cookie) {
	t.Helper()
	w := httptest.NewRecorder()
	p.login(w, httptest.NewRequest("GET", "https://workspace.example.test/auth/login?returnTo="+url.QueryEscape(target), nil))
	u, err := url.Parse(w.Header().Get("Location"))
	if err != nil {
		t.Fatal(err)
	}
	state := u.Query().Get("state")
	cs := w.Result().Cookies()
	if state == "" || len(cs) != 1 {
		t.Fatal("missing isolated flow")
	}
	c := cs[0]
	if c.Name != portalFlowCookieName(state) || c.Value != state || c.Path != "/auth/callback" || !c.Secure || !c.HttpOnly || c.SameSite != http.SameSiteLaxMode || c.MaxAge != 300 {
		t.Fatal("flow cookie lost protections")
	}
	if u.Query().Get("code_challenge_method") != "S256" || u.Query().Get("nonce") == "" {
		t.Fatal("PKCE/nonce missing")
	}
	return state, c
}

func TestPortalConcurrentLoginFlows(t *testing.T) {
	p := flowTestPortal()
	s1, c1 := beginFlow(t, p, "/drive/")
	s2, c2 := beginFlow(t, p, "/mail/")
	if s1 == s2 || c1.Name == c2.Name {
		t.Fatal("tabs share flow state")
	}
	for _, tc := range []struct{ state, target string }{{s1, "/drive/"}, {s2, "/mail/"}} {
		r := httptest.NewRequest("GET", "https://workspace.example.test/auth/callback?state="+tc.state+"&code=synthetic", nil)
		r.AddCookie(c1)
		r.AddCookie(c2)
		r.AddCookie(&http.Cookie{Name: "workspace_flow", Value: s2}) // former global overwrite
		w := httptest.NewRecorder()
		f, ok := p.consumeLoginFlow(w, r)
		if !ok || f.ReturnTo != tc.target || f.Verifier == "" || f.Nonce == "" {
			t.Fatal("concurrent flow was lost")
		}
		for _, c := range w.Result().Cookies() {
			if c.Name != portalFlowCookieName(tc.state) || c.MaxAge != -1 {
				t.Fatal("another tab's cookie was cleared")
			}
		}
	}
	if len(p.flows) != 0 {
		t.Fatal("flows were not single use")
	}
}

func TestPortalFlowRejectsUnboundCallback(t *testing.T) {
	for _, kind := range []string{"missing", "wrong-value", "other-tab", "oversized"} {
		t.Run(kind, func(t *testing.T) {
			p := flowTestPortal()
			state, c := beginFlow(t, p, "/mail/")
			callbackState := state
			if kind == "oversized" {
				callbackState = strings.Repeat("x", 257)
			}
			r := httptest.NewRequest("GET", "https://workspace.example.test/auth/callback?state="+callbackState+"&code=do-not-exchange", nil)
			switch kind {
			case "wrong-value":
				copy := *c
				copy.Value = "wrong"
				r.AddCookie(&copy)
			case "other-tab":
				_, other := beginFlow(t, p, "/drive/")
				r.AddCookie(other)
			}
			w := httptest.NewRecorder()
			if _, ok := p.consumeLoginFlow(w, r); ok || w.Code != http.StatusBadRequest {
				t.Fatal("unbound callback accepted")
			}
			if _, ok := p.flows[state]; !ok {
				t.Fatal("invalid callback consumed valid flow")
			}
			if strings.Contains(w.Body.String(), "do-not-exchange") || !strings.Contains(w.Body.String(), "Continue sign-in") || w.Header().Get("Referrer-Policy") != "no-referrer" {
				t.Fatal("unsafe or unhelpful recovery")
			}
		})
	}
}

func TestPortalFlowExpiryRestartAndReplayStartFresh(t *testing.T) {
	for _, kind := range []string{"expired", "restart", "replay"} {
		t.Run(kind, func(t *testing.T) {
			p := flowTestPortal()
			state, c := beginFlow(t, p, "/drive/?view=grid")
			r := httptest.NewRequest("GET", "https://workspace.example.test/auth/callback?state="+state+"&code=never-reuse", nil)
			r.AddCookie(c)
			if kind == "expired" {
				f := p.flows[state]
				f.Expires = time.Now().Add(-time.Second)
				p.flows[state] = f
			} else if kind == "restart" {
				p.flows = map[string]loginFlow{}
			} else {
				if _, ok := p.consumeLoginFlow(httptest.NewRecorder(), r); !ok {
					t.Fatal("initial flow failed")
				}
			}
			w := httptest.NewRecorder()
			if _, ok := p.consumeLoginFlow(w, r); ok || w.Code != http.StatusSeeOther {
				t.Fatal("stale flow was exchanged")
			}
			u, _ := url.Parse(w.Header().Get("Location"))
			if u.Host != "" || u.Path != "/auth/login" || strings.Contains(u.RawQuery, "never-reuse") {
				t.Fatal("unsafe recovery target")
			}
			if kind == "expired" && u.Query().Get("returnTo") != "/drive/?view=grid" {
				t.Fatal("known safe destination lost")
			}
		})
	}
}

func TestPortalLegacyFlowRemainsBound(t *testing.T) {
	p := flowTestPortal()
	state, _ := beginFlow(t, p, "/mail/")
	r := httptest.NewRequest("GET", "https://workspace.example.test/auth/callback?state="+state, nil)
	r.AddCookie(&http.Cookie{Name: "workspace_flow", Value: state})
	w := httptest.NewRecorder()
	if f, ok := p.consumeLoginFlow(w, r); !ok || f.ReturnTo != "/mail/" {
		t.Fatal("matching legacy flow rejected")
	}
	cs := w.Result().Cookies()
	if len(cs) != 1 || cs[0].Name != "workspace_flow" || cs[0].MaxAge != -1 || cs[0].Path != "/" {
		t.Fatal("legacy cookie not consumed")
	}
}
