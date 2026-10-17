package suite

import (
	"embed"
	"encoding/json"
	"net"
	"net/http"
	"strings"
)

//go:embed navigation/*
var navigationAssets embed.FS

func (p *Portal) navigationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /_workspace/apps", func(w http.ResponseWriter, r *http.Request) {
		user, ok := p.user(r)
		if !ok {
			http.Error(w, "Sign in to Workspace", 401)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"apps": p.tiles(user.Subject == p.rt.Config.Owner), "name": user.Name, "email": user.Email, "home": "/"})
	})
	for _, name := range []string{"navigation.js", "navigation.css", "logout.js"} {
		mux.HandleFunc("GET /_workspace/"+name, func(w http.ResponseWriter, r *http.Request) {
			body, err := navigationAssets.ReadFile("navigation/" + name)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			if strings.HasSuffix(name, ".js") {
				w.Header().Set("Content-Type", "application/javascript")
			} else {
				w.Header().Set("Content-Type", "text/css")
			}
			w.Header().Set("Cache-Control", "no-cache")
			w.Write(body)
		})
	}
}

func (p *Portal) isAppRequest(r *http.Request) bool {
	if publicDriveRoute(r.URL.Path) {
		return false
	}
	for _, prefix := range []string{"/drive", "/api", "/photos-app", "/photos-static", "/sdkjs", "/fonts"} {
		if r.URL.Path == prefix || strings.HasPrefix(r.URL.Path, prefix+"/") {
			return true
		}
	}
	for _, app := range p.rt.Config.Catalog.Apps {
		route := strings.TrimRight(app.Route, "/")
		if app.ID != "accounts" && app.ID != "vault" && app.Kind == "process" && route != "" && (r.URL.Path == route || strings.HasPrefix(r.URL.Path, route+"/")) {
			return true
		}
	}
	return false
}

func nativeGitRequest(r *http.Request) bool {
	// Native clients still authenticate to Git with their own independently revocable token.
	if !strings.HasPrefix(r.URL.Path, "/git/") || r.Header.Get("Authorization") == "" {
		return false
	}
	return strings.HasPrefix(r.URL.Path, "/git/api/") || strings.HasSuffix(r.URL.Path, "/info/refs") || strings.HasSuffix(r.URL.Path, "/git-upload-pack") || strings.HasSuffix(r.URL.Path, "/git-receive-pack")
}

func clearAppCookies(w http.ResponseWriter, r *http.Request, secure bool) {
	names := map[string][]string{
		"filebrowser_quantum_jwt": {"/", "/drive", "/drive/"},
		"mail_session":            {"/", "/mail/"}, "mail_login": {"/", "/mail/"},
		"calendar_session": {"/", "/calendar/"}, "calendar_login": {"/", "/calendar/"},
		"contacts_session": {"/", "/contacts/"}, "contacts_login": {"/", "/contacts/"},
		"chat_session": {"/", "/chat/"}, "chat_login": {"/", "/chat/"},
		"notes_session": {"/", "/notes/"}, "notes_login": {"/", "/notes/"},
		"mytube_session": {"/", "/my-tube/"}, "mytube_login": {"/", "/my-tube/"},
		"access_token_cookie": {"/", "/music/"}, "csrf_access_token": {"/", "/music/"}, "swingmusic_oidc": {"/", "/music/"},
		"workspace_gitea": {"/", "/git", "/git/"},
	}
	domain := r.Host
	for _, app := range suiteSignIns {
		names["workspace_sso_"+app.id] = []string{"/"}
		names["workspace_return_"+app.id] = []string{"/"}
	}
	if host, _, err := net.SplitHostPort(domain); err == nil {
		domain = host
	}
	for name, paths := range names {
		for _, path := range paths {
			http.SetCookie(w, &http.Cookie{Name: name, Value: "", Path: path, MaxAge: -1, HttpOnly: true, Secure: secure, SameSite: http.SameSiteLaxMode})
			if name == "filebrowser_quantum_jwt" {
				http.SetCookie(w, &http.Cookie{Name: name, Value: "", Domain: domain, Path: path, MaxAge: -1, HttpOnly: true, Secure: secure, SameSite: http.SameSiteLaxMode})
			}
		}
	}
}

func publicDriveRoute(path string) bool {
	for _, prefix := range []string{"/drive/public/static/", "/drive/public/share/", "/drive/assets/", "/drive/sdkjs/", "/drive/web-apps/", "/drive/fonts/", "/drive/icons/", "/drive/libs/", "/drive/ran-fonts/", "/drive/ranui-iife/", "/drive/social/", "/drive/studio/editor/", "/sdkjs/", "/fonts/"} {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	if path == "/drive/editor" {
		return true
	}
	for _, prefix := range []string{"/drive/public/api/resources", "/drive/public/api/office/", "/drive/public/api/media/"} {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return path == "/drive/public/api/raw" || path == "/drive/public/api/share/info" || path == "/drive/public/api/share/image"
}
func nativeDriveRequest(r *http.Request) bool {
	if r.Header.Get("Authorization") == "" {
		return false
	}
	return strings.HasPrefix(r.URL.Path, "/drive/api/") || strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/drive/dav/")
}
func anonymousShareCookies(r *http.Request) {
	cookies := r.Cookies()
	r.Header.Del("Cookie")
	for _, cookie := range cookies {
		if cookie.Name == "filebrowser_share_session" {
			r.AddCookie(cookie)
		}
	}
}
