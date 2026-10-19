package suite

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Each app still authenticates with its own OIDC client and session. Workspace
// only starts that normal flow; it never supplies a user identity or access token.
type appSignIn struct {
	id, name, route, probe, login, logout string
	cookies                               []string
}

var suiteSignIns = []appSignIn{
	{"tasks", "My Tasks", "/tasks", "/api/me", "/tasks/auth/login", "/tasks/auth/logout", []string{"tasks_session"}},
	{"pdf", "My PDF", "/pdf", "/api/me", "/pdf/auth/login", "/pdf/auth/logout", []string{"pdf_session"}},
	{"monitor", "My Monitor", "/monitor", "/api/me", "/monitor/auth/login", "/monitor/auth/logout", []string{"monitor_session"}},
	{"search", "My Search", "/search", "/api/me", "/search/auth/login", "/search/auth/logout", []string{"search_session"}},

	{"hosting", "My Hosting", "/hosting", "/api/me", "/hosting/auth/login", "/hosting/auth/logout", []string{"myhosting_session"}},
	{"chat", "Chat", "/chat", "/api/me", "/chat/auth/login", "/chat/auth/logout", []string{"chat_session"}},
	{"drive", "Drive", "/drive", "/drive/api/users?username=self", "/drive/api/auth/oidc/login", "/drive/api/auth/logout", []string{"filebrowser_quantum_jwt"}},
	{"gitea", "Git", "/git", "/user/settings", "/git/user/oauth2/myaccounts", "/git/user/logout", []string{"workspace_gitea"}},
	{"music", "Music", "/music", "/auth/user", "/music/auth/myaccounts/login", "/music/auth/logout", []string{"access_token_cookie", "csrf_access_token"}},
	{"mail", "Mail", "/mail", "/api/session", "/mail/auth/login", "/mail/auth/logout", []string{"mail_session"}},
	{"notes", "Notes", "/notes", "/api/me", "/notes/auth/login", "/notes/api/logout", []string{"notes_session"}},
	{"my-tube", "My Tube", "/my-tube", "/api/me", "/my-tube/auth/login", "/my-tube/api/logout", []string{"mytube_session"}},
	{"calendar", "Calendar", "/calendar", "/api/session", "/calendar/auth/login", "/calendar/auth/logout", []string{"calendar_session"}},
	{"contacts", "Contacts", "/contacts", "/api/session", "/contacts/auth/login", "/contacts/auth/logout", []string{"contacts_session"}},
}

func (p *Portal) signInForPath(path string) (appSignIn, bool) {
	if path == "/photos-app" || strings.HasPrefix(path, "/photos-app/") {
		path = "/drive"
	}
	for _, app := range suiteSignIns {
		if _, installed := p.rt.Config.Catalog.Find(app.id); installed && (path == app.route || strings.HasPrefix(path, app.route+"/")) {
			return app, true
		}
	}
	return appSignIn{}, false
}

func browserDocument(r *http.Request) bool {
	return r.Method == http.MethodGet && (r.Header.Get("Sec-Fetch-Dest") == "document" || strings.Contains(r.Header.Get("Accept"), "text/html"))
}

func (p *Portal) signInStamp(r *http.Request) string {
	c, err := r.Cookie("workspace_session")
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", sha256.Sum256([]byte(c.Value)))
}

// Probes go only to an installed app's loopback listener. Forward just that
// app's cookies, never the portal/provider session or another app's credentials.
func (p *Portal) appSessionStatus(r *http.Request, app appSignIn) int {
	installed, ok := p.rt.Config.Catalog.Find(app.id)
	if !ok {
		return 0
	}
	ctx, cancel := context.WithTimeout(r.Context(), 6*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d%s", installed.Port, app.probe), nil)
	if err != nil {
		return 0
	}
	external, _ := url.Parse(p.rt.Config.Origin)
	req.Host = external.Host
	req.Header.Set("X-Forwarded-Host", external.Host)
	req.Header.Set("X-Forwarded-Proto", external.Scheme)
	req.Header.Set("Accept", "application/json")
	for _, c := range r.Cookies() {
		for _, name := range app.cookies {
			if c.Name == name {
				req.AddCookie(c)
			}
		}
	}
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	resp, err := client.Do(req)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	// Gitea's API intentionally accepts tokens, not browser sessions. Its
	// authenticated settings page is the authoritative browser-session check.
	if app.id == "gitea" && resp.StatusCode >= 300 && resp.StatusCode < 400 {
		location, err := url.Parse(resp.Header.Get("Location"))
		if err == nil && (location.Host == "" || location.Host == external.Host) && location.Path == "/git/user/login" {
			return http.StatusUnauthorized
		}
		return 0
	}
	if foundationApp(app.id) && resp.StatusCode == http.StatusOK {
		var info struct {
			User struct{ Subject, Issuer string }
		}
		if json.NewDecoder(io.LimitReader(resp.Body, 256*1024)).Decode(&info) != nil {
			return 0
		}
		u, ok := p.user(r)
		if !ok || info.User.Subject != u.Subject || info.User.Issuer != p.rt.Config.Issuer {
			return http.StatusUnauthorized
		}
	}
	if app.id == "my-tube" && resp.StatusCode == http.StatusOK {
		var info struct {
			Mode  string `json:"mode"`
			Setup bool   `json:"setup"`
			User  struct {
				ID string `json:"id"`
			} `json:"user"`
		}
		if json.NewDecoder(io.LimitReader(resp.Body, 256*1024)).Decode(&info) != nil {
			return 0
		}
		if info.Mode == "managed" && !info.Setup && info.User.ID == "" {
			return http.StatusUnauthorized
		}
	}
	return resp.StatusCode
}

// A bookmark or app tile can enter directly, without a second app sign-in
// screen. A failed attempt or explicit app logout pauses automatic sign-in
// until the person deliberately opens the app again from Workspace.
func (p *Portal) seamlessApp(w http.ResponseWriter, r *http.Request) bool {
	app, ok := p.signInForPath(r.URL.Path)
	if !ok {
		return false
	}
	guard, destination := "workspace_sso_"+app.id, "workspace_return_"+app.id
	if r.URL.Path == app.logout && (r.Method == http.MethodPost || app.id == "music") {
		p.cookie(w, "workspace_connect_"+app.id, "", -1)
		p.cookie(w, guard, p.signInStamp(r), 8*3600)
		p.cookie(w, destination, "", -1)
		return false
	}
	if r.URL.Path == app.login {
		// Manual retry uses the native app handler, retaining its state/PKCE checks.
		return false
	}
	if !browserDocument(r) {
		return false
	}
	w.Header().Set("Cache-Control", "no-store")
	// Native cookies are scoped to /app/. Canonicalise before checking them:
	// /app does not carry those cookies and must never start another login.
	if r.URL.Path == app.route {
		target := app.route + "/"
		if r.URL.RawQuery != "" {
			target += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, target, http.StatusSeeOther)
		return true
	}
	relative := strings.TrimPrefix(r.URL.Path, app.route)
	for _, prefix := range []string{"/auth/", "/api/", "/assets/", "/static/", "/public/", "/user/oauth2/"} {
		if strings.HasPrefix(relative, prefix) {
			return false
		}
	}
	if relative == "/setup" || relative == "/setup/" || publicDriveRoute(r.URL.Path) {
		return false
	}
	if p.mailNeedsSetup(r, app) {
		u, ok := p.user(r)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = mailSetupNeeded.Execute(w, map[string]bool{"Owner": ok && u.Subject == p.rt.Config.Owner})
		return true
	}
	status := p.appSessionStatus(r, app)
	if status == http.StatusOK {
		p.cookie(w, guard, "", -1)
		if p.finishConnection(w, r, app) {
			return true
		}
		if c, err := r.Cookie(destination); err == nil {
			p.cookie(w, destination, "", -1)
			if target, err := url.QueryUnescape(c.Value); err == nil && len(target) <= 3000 && safeReturnPath(target) == target {
				u, _ := url.Parse(target)
				if u.Path == app.route {
					u.Path += "/"
					target = u.RequestURI()
				}
				if targetApp, valid := p.signInForPath(u.Path); valid && targetApp.id == app.id && target != r.URL.RequestURI() {
					http.Redirect(w, r, target, http.StatusSeeOther)
					return true
				}
			}
		}
		return false
	}
	// Network, setup, permission and server errors are not expired sessions.
	if status != http.StatusUnauthorized {
		return false
	}
	if c, err := r.Cookie(guard); err == nil && c.Value == p.signInStamp(r) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_ = signInRecovery.Execute(w, map[string]string{"Name": app.name, "Retry": "/_workspace/open/" + app.id + "?" + url.Values{"returnTo": {safeReturnPath(r.URL.RequestURI())}}.Encode()})
		return true
	}
	p.cookie(w, guard, p.signInStamp(r), 8*3600)
	if target := safeReturnPath(r.URL.RequestURI()); len(target) <= 3000 {
		p.cookie(w, destination, url.QueryEscape(target), 300)
	}
	login := app.login
	if app.id == "drive" {
		login += "?" + url.Values{"redirect": {safeReturnPath(r.URL.RequestURI())}}.Encode()
	}
	w.Header().Set("Cache-Control", "no-store")
	http.Redirect(w, r, login, http.StatusSeeOther)
	return true
}

func (p *Portal) appOpen(w http.ResponseWriter, r *http.Request) {
	app, ok := p.rt.Config.Catalog.Find(r.PathValue("id"))
	if !ok || app.ID == "accounts" {
		http.NotFound(w, r)
		return
	}
	if _, ok := p.user(r); !ok {
		http.Redirect(w, r, "/auth/login?"+url.Values{"returnTo": {safeReturnPath(r.URL.RequestURI())}}.Encode(), http.StatusSeeOther)
		return
	}
	target := app.Route
	if candidate := r.URL.Query().Get("returnTo"); candidate != "" && safeReturnPath(candidate) == candidate {
		if signIn, valid := p.signInForPath(strings.SplitN(candidate, "?", 2)[0]); valid && signIn.id == app.ID {
			target = candidate
		}
	}
	if signIn, valid := p.signInForPath(app.Route); valid {
		p.cookie(w, "workspace_sso_"+signIn.id, "", -1)
		p.cookie(w, "workspace_return_"+signIn.id, "", -1)
	}
	w.Header().Set("Cache-Control", "no-store")
	if app.ID == "music" {
		musicOpening(w, target)
		return
	}
	http.Redirect(w, r, target, http.StatusSeeOther)
}

var signInRecovery = template.Must(template.New("app-sign-in").Parse(`<!doctype html><html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Sign in to {{.Name}} Â· One Workspace</title><style>body{max-width:600px;margin:15vh auto;padding:28px;font:17px/1.7 system-ui;color:#202124;background:#f8fafd}a{color:#0b57d0}.button{display:inline-block;padding:10px 24px;border-radius:24px;background:#0b57d0;color:white;text-decoration:none}</style><h1>Continue to {{.Name}}</h1><p>Automatic sign-in has paused because you signed out of this app or its last sign-in did not finish. Your Workspace account is still signed in.</p><p><a class="button" href="{{.Retry}}">Try MyAccounts again</a></p><p>If it keeps failing, check that MyAccounts is available and that this app is connected in Workspace settings. Any account verification required by MyAccounts will appear there.</p><a href="/">Back to Workspace</a></html>`))
