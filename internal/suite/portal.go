package suite

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/oneworkspacestudio/one-workspace/internal/branding"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"golang.org/x/oauth2"
	"html/template"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type loginFlow struct {
	Verifier, Nonce, ReturnTo string
	Expires                   time.Time
}
type session struct {
	Subject, Name, Email, IDToken string
	HostingSession                string
	AppSessions                   map[string]string
	Expires                       time.Time
}
type Portal struct {
	maintenanceMu  sync.Mutex
	adminMu        sync.Mutex
	rt             *Runtime
	oauth          oauth2.Config
	verifier       *oidc.IDTokenVerifier
	mu             sync.Mutex
	flows          map[string]loginFlow
	sessions       map[string]session
	secure         bool
	logoutEndpoint string
}

func NewPortal(ctx context.Context, rt *Runtime) (*Portal, error) {
	ctx = oidc.ClientContext(ctx, &http.Client{Timeout: 45 * time.Second})
	// A healthy loopback account service does not guarantee that its public
	// discovery route has propagated through a just-restarted tunnel yet.
	// Retry discovery only; issuer, signatures and owner checks stay unchanged.
	var provider *oidc.Provider
	var e error
	discoveryDeadline := time.Now().Add(2 * time.Minute)
	for {
		provider, e = oidc.NewProvider(ctx, rt.Config.Issuer)
		if e == nil {
			break
		}
		if ctx.Err() != nil || time.Now().After(discoveryDeadline) {
			return nil, fmt.Errorf("Accounts discovery failed after waiting for the public sign-in connection: %w", e)
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
	var metadata struct {
		Logout string `json:"end_session_endpoint"`
	}
	if err := provider.Claims(&metadata); err != nil {
		return nil, err
	}
	if metadata.Logout != "" {
		endpoint, err := url.Parse(metadata.Logout)
		issuer, _ := url.Parse(rt.Config.Issuer)
		if err != nil || endpoint.Scheme != issuer.Scheme || endpoint.Host != issuer.Host {
			return nil, fmt.Errorf("MyAccounts sign-out endpoint has an unexpected origin")
		}
		if err := ensurePortalLogoutRedirect(ctx, rt); err != nil {
			return nil, err
		}
	}
	client := rt.Secrets.Clients["portal"]
	if err := ensureSharedSignIn(ctx, rt); err != nil {
		return nil, err
	}
	return &Portal{rt: rt, oauth: oauth2.Config{ClientID: client.ID, ClientSecret: client.Secret, Endpoint: provider.Endpoint(), RedirectURL: rt.Config.Origin + "/auth/callback", Scopes: []string{oidc.ScopeOpenID, "profile", "email"}}, verifier: provider.Verifier(&oidc.Config{ClientID: client.ID}), flows: map[string]loginFlow{}, sessions: map[string]session{}, secure: strings.HasPrefix(rt.Config.Origin, "https://"), logoutEndpoint: metadata.Logout}, nil
}
func (p *Portal) cookie(w http.ResponseWriter, name, value string, age int) {
	http.SetCookie(w, &http.Cookie{Name: name, Value: value, Path: "/", Secure: p.secure, HttpOnly: true, SameSite: http.SameSiteLaxMode, MaxAge: age})
}
func (p *Portal) user(r *http.Request) (session, bool) {
	cookie, e := r.Cookie("workspace_session")
	if e != nil {
		return session{}, false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	s, ok := p.sessions[cookie.Value]
	if time.Now().After(s.Expires) {
		delete(p.sessions, cookie.Value)
		ok = false
	}
	return s, ok
}
func safeReturnPath(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.IsAbs() || u.Host != "" || !strings.HasPrefix(u.Path, "/") || strings.HasPrefix(u.Path, "//") || strings.ContainsAny(raw+u.Path, "\\\r\n") {
		return "/"
	}
	if strings.HasPrefix(u.Path, "/auth/") || strings.HasPrefix(u.Path, "/_control/") {
		return "/"
	}
	return u.RequestURI()
}
func (p *Portal) login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if _, ok := p.user(r); ok {
		http.Redirect(w, r, safeReturnPath(r.URL.Query().Get("returnTo")), http.StatusSeeOther)
		return
	}
	state, nonce, v := platform.Token(), platform.Token(), oauth2.GenerateVerifier()
	p.mu.Lock()
	for k, f := range p.flows {
		if time.Now().After(f.Expires) {
			delete(p.flows, k)
		}
	}
	if len(p.flows) > 1000 {
		p.mu.Unlock()
		http.Error(w, "Try again later", 429)
		return
	}
	p.flows[state] = loginFlow{Verifier: v, Nonce: nonce, ReturnTo: safeReturnPath(r.URL.Query().Get("returnTo")), Expires: time.Now().Add(5 * time.Minute)}
	p.mu.Unlock()
	p.flowCookie(w, state, state, 300)
	http.Redirect(w, r, p.oauth.AuthCodeURL(state, oidc.Nonce(nonce), oauth2.S256ChallengeOption(v)), http.StatusFound)
}
func (p *Portal) callback(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	f, ok := p.consumeLoginFlow(w, r)
	if !ok {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	token, e := p.oauth.Exchange(ctx, r.URL.Query().Get("code"), oauth2.VerifierOption(f.Verifier))
	if e != nil {
		http.Error(w, "Accounts did not complete login", 401)
		return
	}
	raw, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "Accounts did not return an ID token", 401)
		return
	}
	id, e := p.verifier.Verify(ctx, raw)
	if e != nil || id.Nonce != f.Nonce {
		http.Error(w, "Accounts token did not verify", 401)
		return
	}
	var claims struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		Email       string `json:"email"`
	}
	_ = id.Claims(&claims)
	if strings.TrimSpace(claims.DisplayName) != "" {
		claims.Name = claims.DisplayName
	}
	expiry := time.Now().Add(8 * time.Hour)
	if id.Expiry.Before(expiry) {
		expiry = id.Expiry
	}
	key := platform.Token()
	p.mu.Lock()
	for k, s := range p.sessions {
		if time.Now().After(s.Expires) {
			delete(p.sessions, k)
		}
	}
	p.sessions[key] = session{Subject: id.Subject, Name: claims.Name, Email: claims.Email, IDToken: raw, Expires: expiry}
	p.mu.Unlock()
	clearAppCookies(w, r, p.secure)
	p.cookie(w, "workspace_session", key, int(time.Until(expiry).Seconds()))
	http.Redirect(w, r, safeReturnPath(f.ReturnTo), http.StatusSeeOther)
}
func (p *Portal) Handler(control string) http.Handler {
	mux := http.NewServeMux()
	p.navigationRoutes(mux)
	p.monitorRegistry(mux)
	p.socialStorageRoutes(mux)
	p.myaiStorageRoutes(mux)
	p.clipsStorageRoutes(mux)
	p.chatCatalogRoute(mux)
	mux.HandleFunc("GET /_workspace/open/{id}", p.appOpen)
	p.vaultRoutes(mux)
	p.chatRoutes(mux)
	p.updateRoutes(mux)
	p.backupRoutes(mux)
	p.startupRoutes(mux, control)
	mux.Handle("GET /branding/", branding.Handler())
	mux.HandleFunc("GET /auth/login", p.login)
	mux.HandleFunc("GET /auth/callback", p.callback)
	mux.HandleFunc("GET /auth/signed-out", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/?signedOut=1", http.StatusSeeOther)
	})
	mux.HandleFunc("POST /auth/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Origin") != p.rt.Config.Origin {
			http.Error(w, "Origin rejected", 403)
			return
		}
		current, _ := p.user(r)
		if !p.revokeRemoteAdmin(w, r) {
			return
		}
		p.closeChatSession(r)
		p.closeHostingSession(r)
		p.closeFoundationSessions(r)
		if c, e := r.Cookie("workspace_session"); e == nil {
			p.mu.Lock()
			delete(p.sessions, c.Value)
			p.mu.Unlock()
		}
		p.cookie(w, "workspace_session", "", -1)
		clearAppCookies(w, r, p.secure)
		if p.logoutEndpoint != "" && current.IDToken != "" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Cache-Control", "no-store")
			w.Header().Set("Referrer-Policy", "no-referrer")
			_ = logoutPage.Execute(w, map[string]string{"Endpoint": p.logoutEndpoint, "Hint": current.IDToken, "Client": p.rt.Secrets.Clients["portal"].ID, "Return": p.rt.Config.Origin + "/auth/signed-out"})
			return
		}
		http.Redirect(w, r, "/?signedOut=1", http.StatusSeeOther)
	})
	mux.HandleFunc("POST /_control/stop", func(w http.ResponseWriter, r *http.Request) {
		if r.RemoteAddr == "" || subtle.ConstantTimeCompare([]byte(r.Header.Get("Authorization")), []byte("Bearer "+control)) != 1 {
			http.Error(w, "Forbidden", 403)
			return
		}
		w.WriteHeader(200)
		go func() { time.Sleep(100 * time.Millisecond); p.rt.cancel() }()
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok", "apps": p.rt.Status()})
	})
	if drive, installed := p.rt.Config.Catalog.Find("drive"); installed {
		driveProxy := proxy(fmt.Sprintf("http://127.0.0.1:%d", drive.Port), p.rt.Config.Origin, "", false, true)
		mux.Handle("/drive/", driveProxy)
		mux.HandleFunc("/drive", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/drive/", 302) })
		mux.Handle("/api/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = "/drive" + r.URL.Path
			r.URL.RawPath = ""
			driveProxy.ServeHTTP(w, r)
		}))
		// The bundled editor's converter worker uses a root-relative sdkjs URL.
		// Keep the worker and its Wasm/importScripts on the same trusted origin.
		for _, prefix := range []string{"/sdkjs/", "/fonts/"} {
			mux.Handle(prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r.URL.Path = "/drive" + r.URL.Path
				r.URL.RawPath = ""
				driveProxy.ServeHTTP(w, r)
			}))
		}
		for _, prefix := range []string{"/settings", "/tools", "/files", "/photos"} {
			for _, route := range []string{prefix, prefix + "/"} {
				mux.Handle(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					target := "/drive" + r.URL.Path
					if r.URL.RawQuery != "" {
						target += "?" + r.URL.RawQuery
					}
					http.Redirect(w, r, target, http.StatusTemporaryRedirect)
				}))
			}
		}
	}
	if photos, ok := p.rt.Config.Catalog.Find("photos"); ok {
		local := fmt.Sprintf("http://127.0.0.1:%d", photos.Port)
		photoProxy := proxy(local, local, "/photos-app", true, false)
		mux.Handle("/photos-app/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" && r.Method != "HEAD" {
				http.Error(w, "Method not allowed", 405)
				return
			}
			if r.URL.Path != "/photos-app/" && r.URL.Path != "/photos-app/photos" && r.URL.Path != "/photos-app/photos/" {
				http.NotFound(w, r)
				return
			}
			if r.URL.Path == "/photos-app/" {
				r.URL.Path = "/photos-app/photos"
			}
			photoProxy.ServeHTTP(w, r)
		}))
		mux.Handle("/photos-static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" && r.Method != "HEAD" {
				http.Error(w, "Method not allowed", 405)
				return
			}
			photoProxy.ServeHTTP(w, r)
		}))
		mux.HandleFunc("GET /photos-companion.json", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"filebrowserURL": p.rt.Config.Origin + "/drive"})
		})
	}
	for _, a := range p.rt.Config.Catalog.Apps {
		if a.ID == "accounts" || a.ID == "drive" || a.ID == "photos" || a.ID == "vault" || a.ID == "chat" || a.Kind != "process" {
			continue
		}
		route := strings.TrimRight(a.Route, "/")
		h := proxy(fmt.Sprintf("http://127.0.0.1:%d", a.Port), p.rt.Config.Origin, route, a.StripPrefix, false)
		mux.Handle(route+"/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := p.user(r); !ok && a.ID != "gitea" && a.ID != "music" && !p.hostingWebhook(r) && !p.publicMonitor(r) {
				http.Redirect(w, r, "/auth/login?"+url.Values{"returnTo": {r.URL.RequestURI()}}.Encode(), 302)
				return
			}
			h.ServeHTTP(w, r)
		}))
		mux.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, route+"/", http.StatusTemporaryRedirect)
		})
	}
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		user, ok := p.user(r)
		if !ok && r.URL.Query().Get("signedOut") != "1" {
			p.login(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = home.Execute(w, map[string]any{"SearchInstalled": searchInstalled(p.rt.Config), "LoggedIn": ok, "User": user, "Tiles": p.tiles(ok && user.Subject == p.rt.Config.Owner), "Owner": ok && user.Subject == p.rt.Config.Owner, "Issuer": p.rt.Config.Issuer, "Version": p.rt.Config.Version, "LoginURL": "/auth/login?returnTo=" + url.QueryEscape(safeReturnPath(r.URL.Query().Get("returnTo")))})
	})
	origin, _ := url.Parse(p.rt.Config.Origin)
	local := fmt.Sprintf("localhost:%d", p.rt.Config.Port)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthz" && r.URL.Path != "/_control/stop" && r.URL.Path != "/_integration/monitor" && !strings.HasPrefix(r.URL.Path, "/updates") && !strings.HasPrefix(r.URL.Path, "/auth/") && !strings.HasPrefix(r.URL.Path, "/branding/") && UpdatePending(p.rt.Config) {
			http.Error(w, "Workspace is completing an update. Your data is protected; please wait.", http.StatusServiceUnavailable)
			return
		}
		if r.Host != origin.Host && r.Host != local && r.Host != fmt.Sprintf("127.0.0.1:%d", p.rt.Config.Port) {
			http.Error(w, "Host rejected", 403)
			return
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "same-origin")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		if p.rt.backups.appsPaused() && backupAppPath(p.rt.Config, r.URL.Path) {
			w.Header().Set("Retry-After", "30")
			http.Error(w, "Apps are paused for a backup. Check progress in Workspace → Backups and retry when it finishes.", 503)
			return
		}
		if p.isAppRequest(r) {
			if _, ok := p.user(r); !ok && !nativeGitRequest(r) && !nativeDriveRequest(r) && !p.hostingWebhook(r) && !p.publicMonitor(r) {
				if r.Method == "GET" && strings.Contains(r.Header.Get("Accept"), "text/html") {
					http.Redirect(w, r, "/auth/login?returnTo="+url.QueryEscape(safeReturnPath(r.URL.RequestURI())), http.StatusSeeOther)
				} else {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(401)
					fmt.Fprint(w, `{"error":"Sign in to One Workspace","login":"/auth/login"}`)
				}
				return
			}
		}
		if _, ok := p.user(r); ok && p.seamlessApp(w, r) {
			return
		}
		p.rememberHostingSession(r)
		p.rememberFoundationSessions(r)
		if _, ok := p.user(r); !ok && publicDriveRoute(r.URL.Path) {
			anonymousShareCookies(r)
		}
		mux.ServeHTTP(w, r)
	})
}
func proxy(target, origin, prefix string, strip, cookies bool) http.Handler {
	u, _ := url.Parse(target)
	external, _ := url.Parse(origin)
	p := httputil.NewSingleHostReverseProxy(u)
	director := p.Director
	p.Director = func(r *http.Request) {
		director(r)
		if strings.Contains(r.Header.Get("Accept"), "text/html") || r.Header.Get("Sec-Fetch-Dest") == "document" || (prefix == "/music" && r.URL.Path == "/music/sw.js") {
			r.Header.Set("Accept-Encoding", "identity")
		}
		if strip && strings.HasPrefix(r.URL.Path, prefix) {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, prefix)
			if r.URL.Path == "" {
				r.URL.Path = "/"
			}
			r.URL.RawPath = ""
		}
		r.Host = external.Host
		r.Header.Set("X-Forwarded-Host", external.Host)
		r.Header.Set("X-Forwarded-Proto", external.Scheme)
		r.Header.Del("Forwarded")
	}
	p.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
		if browserDocument(r) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Cache-Control", "no-store")
			w.Header().Set("Retry-After", "3")
			w.WriteHeader(http.StatusBadGateway)
			retry := prefix + "/"
			if retry == "//" || retry == "" {
				retry = "/"
			}
			_ = appUnavailable.Execute(w, map[string]string{"Retry": retry})
			return
		}
		http.Error(w, "Application is starting or unavailable. Check Workspace logs.", 502)
	}

	p.ModifyResponse = func(resp *http.Response) error {
		if err := patchMusicWorker(resp, prefix, strip); err != nil {
			return err
		}
		// Preserve the mounted app path on native redirects (for example Flask's
		// canonical trailing slash), without changing external identity redirects.
		if strip && resp.Header.Get("Location") != "" {
			if location, err := url.Parse(resp.Header.Get("Location")); err == nil &&
				(location.Host == "" || location.Host == external.Host || location.Host == u.Host) &&
				strings.HasPrefix(location.Path, "/") && !strings.HasPrefix(location.Path, "//") &&
				location.Path != prefix && !strings.HasPrefix(location.Path, prefix+"/") {
				location.Path = prefix + location.Path
				location.RawPath = ""
				if location.Host != "" {
					location.Scheme, location.Host = external.Scheme, external.Host
				}
				resp.Header.Set("Location", location.String())
			}
		}
		if cookies {
			values := resp.Cookies()
			resp.Header.Del("Set-Cookie")
			for _, c := range values {
				c.Path = "/"
				resp.Header.Add("Set-Cookie", c.String())
			}
		}
		if resp.StatusCode != 200 || !strings.Contains(resp.Header.Get("Content-Type"), "text/html") || resp.Header.Get("Content-Encoding") != "" {
			return nil
		}
		body, err := io.ReadAll(io.LimitReader(resp.Body, 4*1024*1024+1))
		if err != nil {
			return err
		}
		if len(body) > 4*1024*1024 {
			resp.Body = io.NopCloser(io.MultiReader(bytes.NewReader(body), resp.Body))
			return nil
		}
		resp.Body.Close()
		nonce := ""
		if match := regexp.MustCompile(`nonce=["']([^"']+)["']`).FindSubmatch(body); len(match) == 2 {
			nonce = ` nonce="` + template.HTMLEscapeString(string(match[1])) + `"`
		}
		script := []byte(`<script defer src="/_workspace/navigation.js"` + nonce + `></script>`)
		if i := bytes.LastIndex(bytes.ToLower(body), []byte("</body>")); i >= 0 {
			body = append(append(append([]byte{}, body[:i]...), script...), body[i:]...)
		} else {
			body = append(body, script...)
		}
		resp.Body = io.NopCloser(bytes.NewReader(body))
		resp.ContentLength = int64(len(body))
		resp.Header.Set("Content-Length", strconv.Itoa(len(body)))
		resp.Header.Del("ETag")
		resp.Header.Del("Content-MD5")
		resp.Header.Set("Cache-Control", "no-store")
		return nil
	}

	return p
}

type catalogTile struct{ Name, Description, URL, Icon string }

var appUnavailable = template.Must(template.New("app-unavailable").Parse(`<!doctype html><html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>App unavailable · One Workspace</title><style>body{margin:0;background:#f5f7fc;color:#18233b;font:16px/1.65 "Segoe UI",system-ui,sans-serif}main{max-width:560px;margin:12vh auto;padding:32px;border:1px solid #e3e8f4;border-radius:20px;background:#fff}img{width:48px;height:48px}h1{line-height:1.2}p{color:#5e6c84}a{color:#345de5}nav{display:flex;gap:20px;align-items:center;flex-wrap:wrap}.button{padding:12px 22px;background:#345de5;color:#fff;border-radius:12px;text-decoration:none}a:focus-visible{outline:3px solid #345de5;outline-offset:4px}@media(max-width:640px){main{margin:32px 16px;padding:24px}}</style><main><img src="/branding/icons/one_workspace.svg" alt="One Workspace"><h1>This app isn’t ready yet</h1><p>It may be starting, restarting after an update, or temporarily stopped. Wait a moment, then try opening it again.</p><nav><a class="button" href="{{.Retry}}">Try opening the app again</a><a href="/">Back to Workspace</a></nav><p>If it still won’t open, ask your Workspace owner to check <a href="/_workspace/startup">App status and retry</a>. This page does not mean your files have been deleted.</p></main><script defer src="/_workspace/navigation.js"></script></html>`))

func (p *Portal) tiles(owner ...bool) []catalogTile {
	tiles := []catalogTile{}
	for _, a := range p.rt.Config.Catalog.Apps {
		if a.ID == "monitor" && (len(owner) == 0 || !owner[0]) {
			continue // Infrastructure controls are only available to the owner.
		}
		link := "/_workspace/open/" + a.ID
		if a.ID == "accounts" {
			link = p.rt.Config.Issuer
		}
		icon := branding.Icon(a.ID)
		tiles = append(tiles, catalogTile{a.Name, a.Description, link, icon})
	}
	if len(owner) > 0 && owner[0] {
		tiles = append(tiles, catalogTile{"Backups", "Protect your workspace with optional scheduled backups", "/backups", branding.Icon("backups")})
	}
	sort.SliceStable(tiles, func(i, j int) bool { return tiles[i].Name < tiles[j].Name })
	return tiles
}

var home = template.Must(template.New("home").Parse(branding.VersionHTML(`<!doctype html><html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>One Workspace</title><link rel="icon" type="image/svg+xml" href="/branding/icons/one_workspace.svg"><style>body{margin:0;background:#f5f7fc;color:#18233b;font:16px system-ui,sans-serif}header,main{max-width:1040px;margin:auto;padding:32px}header{display:flex;align-items:center;gap:14px}.logo{width:48px;height:48px;flex:none}.app-logo{display:block;width:56px;height:56px;margin-bottom:18px}h1{font-size:44px;letter-spacing:-1.5px;margin-bottom:12px}.muted,p{color:#5e6c84;line-height:1.65}a{color:#345de5;text-decoration:none}button,.button{display:inline-block;background:#345de5;color:white;border:0;padding:13px 22px;border-radius:12px;font:inherit;cursor:pointer}.grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(230px,1fr));gap:18px;margin:32px 0}.tile{background:white;border:1px solid #e3e8f4;border-radius:20px;padding:26px;box-shadow:0 5px 24px #23397605}.tile h2{color:#18233b;margin-top:0;font-size:23px}.tile:hover{border-color:#345de5}.right{margin-left:auto}footer{margin-top:40px;font-size:13px;color:#6e7890}.admin{background:#eaf0ff;padding:22px;border-radius:16px}</style><header><img class="logo" src="/branding/icons/one_workspace.svg" alt="" width="48" height="48"><b>One Workspace</b><span class="right muted">Your apps, together</span></header><main>{{if .LoggedIn}}<h1>Your workspace.</h1>{{if .SearchInstalled}}<form action="/search/" method="get" role="search"><label for="workspace-global-search">Search your Workspace</label><input id="workspace-global-search" name="q" type="search" maxlength="200" placeholder="Files, notes, events and tasks" style="font:inherit;padding:13px;border:1px solid #e3e8f4;border-radius:12px;margin:12px"><button>Search</button></form>{{end}}<p>Welcome, {{if .User.Name}}{{.User.Name}}{{else}}{{.User.Email}}{{end}}. Open an app to continue.</p><div class="grid">{{range .Tiles}}<a class="tile" href="{{.URL}}">{{if .Icon}}<img class="app-logo" src="{{.Icon}}" alt="" width="56" height="56">{{end}}<h2>{{.Name}} ↗</h2><p>{{.Description}}</p></a>{{end}}</div>{{if .Owner}}<div class="admin"><b>Workspace administration</b><p>Add people, reset passwords and manage access in MyAccounts.</p><a href="{{.Issuer}}/organizations/myaccounts/users">Manage accounts →</a> · <a href="/updates">Updates and backups →</a> · <a href="/_workspace/startup">App status and retry →</a></div>{{end}}<form action="/auth/logout" method="post" style="margin-top:24px"><button>Sign out of Workspace</button></form>{{else}}<h1>One account.<br>Your whole workspace.</h1><p>Files, documents and photos in one place.<br>Sign in with your MyAccounts account to get started.</p><a class="button" href="{{.LoginURL}}">Sign in with MyAccounts →</a>{{end}}<footer>One Workspace {{.Version}} · Self-hosted evaluation build. Signing out closes Workspace app access and opens MyAccounts to end its browser session. If your browser blocks that step, sign out of MyAccounts separately. App access tokens and public sharing links are managed separately.</footer></main><script defer src="/_workspace/navigation.js"></script></html>`)))

var logoutPage = template.Must(template.New("logout").Parse(`<!doctype html><html lang="en"><meta charset="utf-8"><title>Signing out of One Workspace</title><h1>Signing out</h1><p>Workspace app access is closed. Opening MyAccounts to end its browser session. If your browser blocks that step, sign out of MyAccounts separately.</p><form id="workspace-provider-logout" method="get" action="{{.Endpoint}}"><input type="hidden" name="id_token_hint" value="{{.Hint}}"><input type="hidden" name="client_id" value="{{.Client}}"><input type="hidden" name="post_logout_redirect_uri" value="{{.Return}}"><button>Continue sign-out</button></form><script src="/_workspace/logout.js" defer></script></html>`))
