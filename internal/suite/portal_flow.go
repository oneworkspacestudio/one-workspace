package suite

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"time"
)

// Each browser authorization has its own short-lived cookie. One shared cookie
// lets another tab overwrite the first tab's state while MyAccounts is opening.
func portalFlowCookieName(state string) string {
	sum := sha256.Sum256([]byte(state))
	return fmt.Sprintf("workspace_flow_%x", sum[:16])
}

func (p *Portal) flowCookie(w http.ResponseWriter, state, value string, age int) {
	http.SetCookie(w, &http.Cookie{Name: portalFlowCookieName(state), Value: value,
		Path: "/auth/callback", Secure: p.secure, HttpOnly: true,
		SameSite: http.SameSiteLaxMode, MaxAge: age})
}

func (p *Portal) consumeLoginFlow(w http.ResponseWriter, r *http.Request) (loginFlow, bool) {
	state := r.URL.Query().Get("state")
	if len(state) < 32 || len(state) > 256 {
		portalSignInProblem(w)
		return loginFlow{}, false
	}
	cookie, err := r.Cookie(portalFlowCookieName(state))
	legacy := false
	if err != nil {
		// Retain a correctly bound in-flight flow from a prior manager version. New
		// authorizations never create or overwrite this shared legacy cookie.
		cookie, err = r.Cookie("workspace_flow")
		legacy = err == nil
	}
	if err != nil || subtle.ConstantTimeCompare([]byte(cookie.Value), []byte(state)) != 1 {
		portalSignInProblem(w)
		return loginFlow{}, false
	}
	p.mu.Lock()
	flow, ok := p.flows[state]
	delete(p.flows, state) // single use, before any token exchange
	p.mu.Unlock()
	if legacy {
		p.cookie(w, "workspace_flow", "", -1)
	} else {
		p.flowCookie(w, state, "", -1)
	}
	if !ok || time.Now().After(flow.Expires) {
		// A restart can discard the server-side flow. Discard this callback's code
		// and begin normal PKCE/nonce authorization again; never exchange an unbound
		// code, trust the callback's identity or relax any token verification.
		target := "/auth/login?returnTo=" + url.QueryEscape(safeReturnPath(flow.ReturnTo))
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Referrer-Policy", "no-referrer")
		http.Redirect(w, r, target, http.StatusSeeOther)
		return loginFlow{}, false
	}
	return flow, true
}

func portalSignInProblem(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.WriteHeader(http.StatusBadRequest)
	_ = portalSignInProblemPage.Execute(w, nil)
}

var portalSignInProblemPage = template.Must(template.New("sign-in-retry").Parse(`<!doctype html><html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Continue sign-in · One Workspace</title><style>body{margin:0;background:#f5f7fc;color:#18233b;font:16px/1.65 "Segoe UI",system-ui,sans-serif}main{max-width:520px;margin:12vh auto;padding:32px;background:#fff;border:1px solid #e3e8f4;border-radius:20px}p{color:#5e6c84}a{display:inline-block;padding:12px 22px;background:#345de5;color:#fff;border-radius:24px;text-decoration:none}</style><main><h1>Let’s start sign-in again</h1><p>This sign-in request expired or no longer matches this browser. Your password has not been changed.</p><p>Continue to start a fresh request. If MyAccounts is already signed in, you can return to Workspace automatically.</p><a href="/auth/login">Continue sign-in</a></main></html>`))
