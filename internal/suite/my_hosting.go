package suite

import (
	"context"
	"fmt"
	"github.com/fahimc/one-workspace/internal/platform"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func hostingSiteEnvironment(c Config, s Secrets) []string {
	env := []string{"WORKSPACE_SITES_DOMAIN=" + c.HostingDomain, "WORKSPACE_SITES_STYLE=flat", "WORKSPACE_SITES_TRANSPORT=tunnel"}
	if !c.Tunnel {
		domain := ""
		if c.HostingDomain != "" {
			domain = "sites." + c.HostingDomain
		}
		return []string{"WORKSPACE_SITES_DOMAIN=" + domain, "WORKSPACE_SITES_STYLE=nested", "WORKSPACE_SITES_TRANSPORT=direct"}
	}
	var resources struct {
		AccountID string `json:"accountId"`
		ZoneID    string `json:"zoneId"`
		ZoneName  string `json:"zoneName"`
		TunnelID  string `json:"tunnelId"`
	}
	if platform.LoadJSON(filepath.Join(c.Data, "cloudflare-resources.json"), &resources) != nil {
		return env
	}
	// HostingDomain can be an email subdomain. The flat sites namespace needs
	// the connected root zone, which owns the Universal SSL certificate.
	zone := resources.ZoneName
	if zone == "" {
		// Older connection receipts predate zoneName. Installer-generated portal
		// names are first-level labels (including for an email subdomain).
		if origin, e := url.Parse(c.Origin); e == nil {
			_, zone, _ = strings.Cut(origin.Hostname(), ".")
		}
	}
	if hostingDomainPattern.MatchString(zone) {
		env[0] = "WORKSPACE_SITES_DOMAIN=" + zone
	}
	token := s.HostingToken
	if token == "" {
		token = s.MailDNSToken
	}
	if token != "" {
		env = append(env, "WORKSPACE_SITES_CF_ACCOUNT="+resources.AccountID, "WORKSPACE_SITES_CF_ZONE="+resources.ZoneID, "WORKSPACE_SITES_CF_TUNNEL="+resources.TunnelID, "WORKSPACE_SITES_CF_TOKEN="+token)
	}
	return env
}

var hostingHookPath = regexp.MustCompile(`^/hosting/hooks/github/[a-f0-9]{32}$`)

// Only this exact POST reaches Hosting without a browser session. Hosting still
// requires the project's HMAC and deduplicates its signed payload. No admin API,
// site content, log or general /hosting/api route is made public by this rule.
func (p *Portal) hostingWebhook(r *http.Request) bool {
	app, ok := p.rt.Config.Catalog.Find("hosting")
	return ok && app.Route == "/hosting" && r.Method == http.MethodPost && hostingHookPath.MatchString(r.URL.Path)
}

func (p *Portal) closeHostingSession(r *http.Request) {
	app, ok := p.rt.Config.Catalog.Find("hosting")
	cookie, e := r.Cookie("myhosting_session")
	if !ok {
		return
	}
	if e != nil {
		// The app cookie is scoped to /hosting/, so browsers do not send it to
		// /auth/logout. Use only the current portal browser session's last token.
		current, valid := p.user(r)
		if !valid || current.HostingSession == "" {
			return
		}
		cookie = &http.Cookie{Name: "myhosting_session", Value: current.HostingSession}
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	req, e := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://127.0.0.1:%d/auth/logout", app.Port), nil)
	if e != nil {
		return
	}
	req.AddCookie(cookie)
	req.Header.Set("Origin", p.rt.Config.Origin)
	client := &http.Client{Timeout: 3 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	if response, e := client.Do(req); e == nil {
		response.Body.Close()
	}
}

func (p *Portal) rememberHostingSession(r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/hosting/") {
		return
	}
	app, ok := p.rt.Config.Catalog.Find("hosting")
	if !ok || app.Route != "/hosting" {
		return
	}
	workspace, e := r.Cookie("workspace_session")
	if e != nil {
		return
	}
	appCookie, e := r.Cookie("myhosting_session")
	if e != nil || len(appCookie.Value) > 128 {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	current, ok := p.sessions[workspace.Value]
	if !ok || time.Now().After(current.Expires) {
		return
	}
	current.HostingSession = appCookie.Value
	p.sessions[workspace.Value] = current
}
