package suite

import (
	"net/http"
	"net/url"
	"strings"
)

// A connection uses the app's existing OIDC flow. Only a destination and a
// session fingerprint are remembered; no provider/app credential is shared.
func (p *Portal) connectionDestination(raw string) (string, bool) {
	if raw == "" || len(raw) > 2000 || safeReturnPath(raw) != raw {
		return "", false
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", false
	}
	if u.Path == "/" {
		return raw, true
	}
	for _, app := range p.rt.Config.Catalog.Apps {
		if app.ID != "accounts" && app.ID != "vault" && app.Route != "" && strings.TrimRight(u.Path, "/") == strings.TrimRight(app.Route, "/") {
			return raw, true
		}
	}
	return "", false
}

func (p *Portal) connectApp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if _, ok := p.user(r); !ok {
		http.Redirect(w, r, "/auth/login?"+url.Values{"returnTo": {safeReturnPath(r.URL.RequestURI())}}.Encode(), http.StatusSeeOther)
		return
	}
	installed, ok := p.rt.Config.Catalog.Find(r.PathValue("id"))
	app, supported := p.signInForPath(installed.Route)
	target, valid := p.connectionDestination(r.URL.Query().Get("returnTo"))
	if !ok || !supported || app.id != installed.ID || !valid || strings.TrimRight(strings.SplitN(target, "?", 2)[0], "/") == strings.TrimRight(app.route, "/") {
		http.Error(w, "Choose an installed app and a Workspace return page.", http.StatusBadRequest)
		return
	}
	p.cookie(w, "workspace_connect_"+app.id, p.signInStamp(r)+"."+url.QueryEscape(target), 300)
	p.cookie(w, "workspace_sso_"+app.id, "", -1)
	p.cookie(w, "workspace_return_"+app.id, "", -1)
	http.Redirect(w, r, app.route+"/", http.StatusSeeOther)
}

func (p *Portal) finishConnection(w http.ResponseWriter, r *http.Request, app appSignIn) bool {
	c, err := r.Cookie("workspace_connect_" + app.id)
	if err != nil {
		return false
	}
	p.cookie(w, c.Name, "", -1)
	stamp, value, ok := strings.Cut(c.Value, ".")
	if !ok || stamp != p.signInStamp(r) {
		return false
	}
	raw, err := url.QueryUnescape(value)
	target, valid := p.connectionDestination(raw)
	if err != nil || !valid {
		return false
	}
	p.cookie(w, "workspace_return_"+app.id, "", -1)
	http.Redirect(w, r, target, http.StatusSeeOther)
	return true
}
