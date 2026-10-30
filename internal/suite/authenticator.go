package suite

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"
)

var authenticatorAsset = regexp.MustCompile(`^/authenticator/assets/[A-Za-z0-9_-]+\.(js|css|wasm)$`)

// This host has no account APIs, mutable server data or seed endpoints.
func publicAuthenticatorRequest(r *http.Request) bool {
	if (r.Method != "GET" && r.Method != "HEAD") || r.URL.RawQuery != "" || path.Clean(r.URL.Path) != strings.TrimSuffix(r.URL.Path, "/") {
		return false
	}
	switch r.URL.Path {
	case "/authenticator", "/authenticator/", "/authenticator/index.html", "/authenticator/manifest.webmanifest", "/authenticator/icon.svg", "/authenticator/icon-192.png", "/authenticator/icon-512.png", "/authenticator/sw.js", "/authenticator/healthz", "/authenticator/version":
		return true
	}
	return authenticatorAsset.MatchString(r.URL.Path)
}

type accountSessionResult struct {
	key    string
	active bool
}
type accountSessionContextKey struct{}

func (p *Portal) accountSessionActive(r *http.Request, s session, key string) bool {
	if result, ok := r.Context().Value(accountSessionContextKey{}).(accountSessionResult); ok && result.key == key {
		return result.active
	}
	active := false
	defer func() {
		*r = *r.WithContext(context.WithValue(r.Context(), accountSessionContextKey{}, accountSessionResult{key, active}))
	}()
	if s.AccessToken == "" || p.rt == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	body := url.Values{"token": {s.AccessToken}, "token_type_hint": {"access_token"}}
	request, err := http.NewRequestWithContext(ctx, "POST", strings.TrimRight(p.rt.Config.Issuer, "/")+"/api/login/oauth/introspect", strings.NewReader(body.Encode()))
	if err != nil {
		return false
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth(p.oauth.ClientID, p.oauth.ClientSecret)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return false
	}
	defer response.Body.Close()
	var result struct {
		Active  bool   `json:"active"`
		Subject string `json:"sub"`
	}
	if response.StatusCode != 200 || json.NewDecoder(io.LimitReader(response.Body, 8192)).Decode(&result) != nil {
		return false
	}
	active = result.Active && result.Subject == s.Subject
	if !active {
		p.mu.Lock()
		delete(p.sessions, key)
		p.mu.Unlock()
	}
	return active
}

// Cancel authenticated reverse-proxy requests when the provider revokes their
// grant. Go's ReverseProxy closes upgraded WebSocket connections on context
// cancellation, and ordinary streamed responses lose their upstream body too.
func (p *Portal) watchAccountSession(r *http.Request, cancel context.CancelFunc, done <-chan struct{}, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case <-r.Context().Done():
			return
		case <-ticker.C:
			probe := r.Clone(context.WithValue(r.Context(), accountSessionContextKey{}, accountSessionResult{}))
			if _, ok := p.user(probe); !ok {
				cancel()
				return
			}
		}
	}
}
func (p *Portal) serveAccountGuarded(w http.ResponseWriter, r *http.Request, next http.Handler) {
	if _, ok := p.user(r); !ok {
		next.ServeHTTP(w, r)
		return
	}
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()
	r = r.WithContext(ctx)
	done := make(chan struct{})
	defer close(done)
	go p.watchAccountSession(r, cancel, done, 5*time.Second)
	next.ServeHTTP(w, r)
}
