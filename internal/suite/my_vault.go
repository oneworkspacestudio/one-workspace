package suite

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/fahimc/one-workspace/internal/platform"
)

// Keep the vault on a distinct origin. No Workspace scripts, portal session
// gate or MyAccounts credentials enter the compatible vault client.
func (p *Portal) vaultRoutes(mux *http.ServeMux) {
	app, ok := p.rt.Config.Catalog.Find("vault")
	if !ok {
		return
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/vault" && r.URL.Path != "/vault/" {
			http.NotFound(w, r)
			return
		}
		var settings struct {
			Origin string `json:"origin"`
		}
		_ = platform.LoadJSON(filepath.Join(p.rt.Config.Data, "my-vault", "myvault.json"), &settings)
		u, e := url.Parse(settings.Origin)
		workspace, _ := url.Parse(p.rt.Config.Origin)
		if e == nil && u.Scheme == "https" && u.Host != "" && u.User == nil && u.RawQuery == "" && u.Fragment == "" && !localVaultHostname(u.Hostname()) && !sameVaultOrigin(u, workspace) {
			http.Redirect(w, r, strings.TrimRight(u.String(), "/")+"/?workspace=1", http.StatusSeeOther)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		localBase := ""
		if e == nil && u.Scheme == "https" && u.User == nil && u.RawQuery == "" && u.Fragment == "" {
			localBase = strings.TrimRight(u.EscapedPath(), "/")
		}
		_ = vaultLanding.Execute(w, map[string]string{"Local": fmt.Sprintf("https://localhost:%d%s/?workspace=1", app.Port+1, localBase)})
	}
	mux.HandleFunc("GET /vault", handler)
	mux.HandleFunc("GET /vault/", handler)
}
func localVaultHostname(host string) bool {
	return strings.EqualFold(host, "localhost") || (net.ParseIP(host) != nil && net.ParseIP(host).IsLoopback())
}
func sameVaultOrigin(a, b *url.URL) bool {
	port := func(u *url.URL) string {
		if u.Port() != "" {
			return u.Port()
		}
		if strings.EqualFold(u.Scheme, "https") {
			return "443"
		}
		return "80"
	}
	return strings.EqualFold(a.Scheme, b.Scheme) && strings.EqualFold(strings.TrimSuffix(a.Hostname(), "."), strings.TrimSuffix(b.Hostname(), ".")) && port(a) == port(b)
}

var vaultLanding = template.Must(template.New("vault-launcher").Parse(`<!doctype html><html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>My Vault · One Workspace</title><link rel="icon" href="/branding/icons/my_vault.svg"><style>body{margin:60px auto;padding:25px;max-width:720px;background:#f8fafd;color:#202124;font:16px/1.8 Arial,system-ui}h1{font-size:36px;font-weight:400;letter-spacing:-1px}a{color:#0b57d0}code{background:#edf2fa;padding:4px 7px;border-radius:5px}.card{background:white;border:1px solid #e2e6f2;border-radius:28px;padding:36px}.logo{width:48px}.button{display:inline-block;border-radius:24px;padding:10px 22px;background:#0b57d0;color:#fff;text-decoration:none}</style><div class="card"><img class="logo" src="/branding/icons/my_vault.svg" alt="My Vault"><h1>Your passwords, together.</h1><p>My Vault reuses its active session when you open it. When connected, MyAccounts handles sign-in and your vault master password unlocks the encrypted items.</p><p><a class="button" href="{{.Local}}">Open My Vault on this computer →</a></p><p>For local access, trust the unique certificate through My Vault’s owner tools. To use another device, configure a dedicated HTTPS address such as <code>vault.example.com</code> in those tools and point its trusted proxy at the vault’s loopback listener.</p><p>Independent vault sign-in remains available if MyAccounts is down. Keep a verified encrypted backup, a portable encrypted export and an offline client copy outside the server.</p><p><a href="/">Return to Workspace</a></p></div></html>`))
