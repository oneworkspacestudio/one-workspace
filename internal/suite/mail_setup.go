package suite

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

// An unconfigured Mail service is healthy enough to show setup but cannot
// create a mailbox. Do not start OIDC and strand the user on a JSON callback.
func (p *Portal) mailNeedsSetup(r *http.Request, app appSignIn) bool {
	if app.id != "mail" {
		return false
	}
	installed, ok := p.rt.Config.Catalog.Find("mail")
	if !ok {
		return false
	}
	req, err := http.NewRequestWithContext(r.Context(), "GET", fmt.Sprintf("http://127.0.0.1:%d/healthz", installed.Port), nil)
	if err != nil {
		return false
	}
	client := &http.Client{Timeout: 3 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	var state struct {
		Configured *bool `json:"configured"`
	}
	return resp.StatusCode == http.StatusOK && json.NewDecoder(io.LimitReader(resp.Body, 16384)).Decode(&state) == nil && state.Configured != nil && !*state.Configured
}

func (p *Portal) mailSetup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	u, ok := p.user(r)
	if !ok || u.Subject != p.rt.Config.Owner {
		http.Error(w, "Only the Workspace owner can configure Mail.", http.StatusForbidden)
		return
	}
	if _, ok := p.rt.Config.Catalog.Find("mail"); !ok {
		http.NotFound(w, r)
		return
	}
	var saved map[string]string
	err := platform.LoadSecret(filepath.Join(p.rt.Config.Data, "mail", "setup.dpapi"), &saved)
	setup, parseErr := url.Parse(saved["url"])
	if err != nil || parseErr != nil || len(setup.Fragment) < 32 || len(setup.Fragment) > 256 {
		http.Error(w, "Mail setup is not ready. Start Mail, then try again.", http.StatusServiceUnavailable)
		return
	}
	// Never redirect to a URL supplied by a file: only reuse its setup fragment.
	http.Redirect(w, r, "/mail/setup#"+url.PathEscape(setup.Fragment), http.StatusSeeOther)
}

var mailSetupNeeded = template.Must(template.New("mail-setup-needed").Parse(`<!doctype html><html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Connect Mail · One Workspace</title><style>body{margin:0;background:#f5f7fc;color:#18233b;font:16px/1.65 "Segoe UI",system-ui,sans-serif}main{max-width:560px;margin:12vh auto;padding:32px;background:#fff;border:1px solid #e3e8f4;border-radius:20px}a{color:#345de5}.button{display:inline-block;padding:12px 22px;background:#345de5;color:white;border-radius:12px;text-decoration:none}p{color:#5e6c84}</style><main><h1>Finish connecting Mail</h1><p>Your Workspace account is signed in. Mail still needs its server and email-domain setup before it can create your mailbox.</p>{{if .Owner}}<p><a class="button" href="/_workspace/setup/mail">Set up Mail</a></p><p>Workspace opens the protected setup for you. Existing mail and DNS settings must be checked before Internet delivery is enabled.</p>{{else}}<p>Ask your Workspace owner to finish Mail setup. You can continue using your other apps.</p>{{end}}<a href="/">Back to Workspace</a></main><script defer src="/_workspace/navigation.js"></script></html>`))
