package suite

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/fahimc/one-workspace/internal/platform"
	"golang.org/x/oauth2"
	"html/template"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

type loginFlow struct {
	Verifier, Nonce string
	Expires         time.Time
}
type session struct {
	Subject, Name, Email string
	Expires              time.Time
}
type Portal struct {
	rt       *Runtime
	oauth    oauth2.Config
	verifier *oidc.IDTokenVerifier
	mu       sync.Mutex
	flows    map[string]loginFlow
	sessions map[string]session
	secure   bool
}

func NewPortal(ctx context.Context, rt *Runtime) (*Portal, error) {
	provider, e := oidc.NewProvider(ctx, rt.Config.Issuer)
	if e != nil {
		return nil, fmt.Errorf("Accounts discovery failed: %w", e)
	}
	client := rt.Secrets.Clients["portal"]
	return &Portal{rt: rt, oauth: oauth2.Config{ClientID: client.ID, ClientSecret: client.Secret, Endpoint: provider.Endpoint(), RedirectURL: rt.Config.Origin + "/auth/callback", Scopes: []string{oidc.ScopeOpenID, "profile", "email"}}, verifier: provider.Verifier(&oidc.Config{ClientID: client.ID}), flows: map[string]loginFlow{}, sessions: map[string]session{}, secure: strings.HasPrefix(rt.Config.Origin, "https://")}, nil
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
func (p *Portal) login(w http.ResponseWriter, r *http.Request) {
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
	p.flows[state] = loginFlow{v, nonce, time.Now().Add(5 * time.Minute)}
	p.mu.Unlock()
	p.cookie(w, "workspace_flow", state, 300)
	http.Redirect(w, r, p.oauth.AuthCodeURL(state, oidc.Nonce(nonce), oauth2.S256ChallengeOption(v)), http.StatusFound)
}
func (p *Portal) callback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	cookie, e := r.Cookie("workspace_flow")
	if e != nil || state == "" || subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(state)) != 1 {
		http.Error(w, "Login state did not match", 400)
		return
	}
	p.mu.Lock()
	f, ok := p.flows[state]
	delete(p.flows, state)
	p.mu.Unlock()
	p.cookie(w, "workspace_flow", "", -1)
	if !ok || time.Now().After(f.Expires) {
		http.Error(w, "Login expired", 400)
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
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	_ = id.Claims(&claims)
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
	p.sessions[key] = session{id.Subject, claims.Name, claims.Email, expiry}
	p.mu.Unlock()
	p.cookie(w, "workspace_session", key, int(time.Until(expiry).Seconds()))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (p *Portal) Handler(control string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /auth/login", p.login)
	mux.HandleFunc("GET /auth/callback", p.callback)
	mux.HandleFunc("POST /auth/logout", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Origin") != p.rt.Config.Origin {
			http.Error(w, "Origin rejected", 403)
			return
		}
		if c, e := r.Cookie("workspace_session"); e == nil {
			p.mu.Lock()
			delete(p.sessions, c.Value)
			p.mu.Unlock()
		}
		p.cookie(w, "workspace_session", "", -1)
		http.Redirect(w, r, "/", 303)
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
	drive, _ := p.rt.Config.Catalog.Find("drive")
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
		if a.ID == "accounts" || a.ID == "drive" || a.ID == "photos" || a.Kind != "process" {
			continue
		}
		route := strings.TrimRight(a.Route, "/")
		h := proxy(fmt.Sprintf("http://127.0.0.1:%d", a.Port), p.rt.Config.Origin, route, a.StripPrefix, false)
		mux.Handle(route+"/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := p.user(r); !ok && a.ID != "gitea" && a.ID != "music" {
				http.Redirect(w, r, "/auth/login", 302)
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
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = home.Execute(w, map[string]any{"LoggedIn": ok, "User": user, "Tiles": p.tiles(), "Owner": ok && user.Subject == p.rt.Config.Owner, "Issuer": p.rt.Config.Issuer, "Version": p.rt.Config.Version})
	})
	origin, _ := url.Parse(p.rt.Config.Origin)
	local := fmt.Sprintf("localhost:%d", p.rt.Config.Port)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host != origin.Host && r.Host != local && r.Host != fmt.Sprintf("127.0.0.1:%d", p.rt.Config.Port) {
			http.Error(w, "Host rejected", 403)
			return
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "same-origin")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
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
		http.Error(w, "Application is starting or unavailable. Check Workspace logs.", 502)
	}
	if cookies {
		p.ModifyResponse = func(resp *http.Response) error {
			values := resp.Cookies()
			resp.Header.Del("Set-Cookie")
			for _, c := range values {
				c.Path = "/"
				resp.Header.Add("Set-Cookie", c.String())
			}
			return nil
		}
	}
	return p
}

type catalogTile struct{ Name, Description, URL string }

func (p *Portal) tiles() []catalogTile {
	tiles := []catalogTile{}
	for _, a := range p.rt.Config.Catalog.Apps {
		link := a.Route
		if a.ID == "accounts" {
			link = p.rt.Config.Issuer
		}
		tiles = append(tiles, catalogTile{a.Name, a.Description, link})
	}
	sort.SliceStable(tiles, func(i, j int) bool { return tiles[i].Name < tiles[j].Name })
	return tiles
}

var home = template.Must(template.New("home").Parse(`<!doctype html><html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>One Workspace</title><style>body{margin:0;background:#f5f7fc;color:#18233b;font:16px system-ui,sans-serif}header,main{max-width:1040px;margin:auto;padding:32px}header{display:flex;align-items:center;gap:14px}strong.logo{background:#345de5;color:white;padding:12px;border-radius:15px}h1{font-size:44px;letter-spacing:-1.5px;margin-bottom:12px}.muted,p{color:#5e6c84;line-height:1.65}a{color:#345de5;text-decoration:none}button,.button{display:inline-block;background:#345de5;color:white;border:0;padding:13px 22px;border-radius:12px;font:inherit;cursor:pointer}.grid{display:grid;grid-template-columns:repeat(auto-fit,minmax(230px,1fr));gap:18px;margin:32px 0}.tile{background:white;border:1px solid #e3e8f4;border-radius:20px;padding:26px;box-shadow:0 5px 24px #23397605}.tile h2{color:#18233b;margin-top:0;font-size:23px}.tile:hover{border-color:#345de5}.right{margin-left:auto}footer{margin-top:40px;font-size:13px;color:#6e7890}.admin{background:#eaf0ff;padding:22px;border-radius:16px}</style><header><strong class="logo">W</strong><b>One Workspace</b><span class="right muted">Your apps, together</span></header><main>{{if .LoggedIn}}<h1>Your workspace.</h1><p>Welcome, {{if .User.Name}}{{.User.Name}}{{else}}{{.User.Email}}{{end}}. Open an app to continue.</p><div class="grid">{{range .Tiles}}<a class="tile" href="{{.URL}}"><h2>{{.Name}} ↗</h2><p>{{.Description}}</p></a>{{end}}</div>{{if .Owner}}<div class="admin"><b>Workspace administration</b><p>Add people, reset passwords and manage access in MyAccounts.</p><a href="{{.Issuer}}/organizations/myaccounts/users">Manage accounts →</a></div>{{end}}<form action="/auth/logout" method="post" style="margin-top:24px"><button>Sign out of Workspace</button></form>{{else}}<h1>One account.<br>Your whole workspace.</h1><p>Files, documents and photos in one place.<br>Sign in with your MyAccounts account to get started.</p><a class="button" href="/auth/login">Sign in with MyAccounts →</a>{{end}}<footer>One Workspace {{.Version}} · Self-hosted evaluation build. Apps retain their own sessions; use Accounts to revoke access.</footer></main></html>`))
