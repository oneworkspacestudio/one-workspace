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
			http.Redirect(w, r, strings.TrimRight(u.String(), "/")+"/", http.StatusSeeOther)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_ = vaultLanding.Execute(w, map[string]string{"Local": fmt.Sprintf("https://localhost:%d/", app.Port+1)})
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

var vaultLanding = template.Must(template.New("vault-launcher").Parse(`<!doctype html><html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>My Vault · One Workspace</title><style>body{margin:60px auto;padding:25px;max-width:720px;background:#f5f7fc;color:#17233b;font:16px/1.8 system-ui}h1{font-size:42px;letter-spacing:-1.7px}a{color:#6550cc}code{background:#ede9fa;padding:4px 7px;border-radius:5px}.card{background:white;border:1px solid #e2e6f2;border-radius:20px;padding:32px}</style><div class="card"><h1>Your vault. Independently.</h1><p>My Vault uses its own master password and a separate secure address to protect its browser storage. MyAccounts is not required to unlock or recover it.</p><p><a href="{{.Local}}">Open the local vault on this computer →</a></p><p>Trust the unique local certificate through My Vault’s owner tools before local sign-in. Trust the unique local certificate through My Vault’s owner tools before local sign-in. For access from another device, run <code>MyVault.exe manage</code> on the server and configure a dedicated HTTPS hostname, such as <code>vault.example.com</code>. Point its trusted proxy at My Vault's fixed loopback listener.</p><p>Keep a verified encrypted backup, a password-protected portable export and an offline client copy outside the server. Recovery &amp; Portability explains each step.</p><p>If Workspace or MyAccounts is unavailable, run the installed MyVault.exe directly with the same private data folder. Only one instance can own it.</p><p><a href="/">Return to Workspace</a></p></div></html>`))
