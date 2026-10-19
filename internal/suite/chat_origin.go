package suite

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

type chatHostingRecord struct {
	Origin      string `json:"origin"`
	DNSRecordID string `json:"dnsRecordId,omitempty"`
	State       string `json:"state"`
	Message     string `json:"message,omitempty"`
}

func chatOrigin(c Config) string {
	var record chatHostingRecord
	if platform.LoadJSON(filepath.Join(c.Data, "chat", "workspace-hosting.json"), &record) == nil && validChatOrigin(record.Origin, c.Origin) {
		return record.Origin
	}
	app, ok := c.Catalog.Find("chat")
	if !ok {
		return ""
	}
	host := "127.0.0.1"
	home, _ := url.Parse(c.Origin)
	if home.Hostname() == host {
		host = "localhost"
	}
	return fmt.Sprintf("http://%s:%d/chat", host, app.Port)
}

func validChatOrigin(raw, workspace string) bool {
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Path != "/chat" || strings.ContainsAny(raw, "\\\r\n") {
		return false
	}
	home, err := url.Parse(workspace)
	if err != nil || strings.EqualFold(home.Hostname(), u.Hostname()) {
		return false
	}
	return u.Scheme == "https" || u.Scheme == "http" && localVaultHostname(u.Hostname())
}

// Update the existing OIDC application, never allocate a second identity client
// or rotate its secret during an origin migration or a normal restart.
func (rt *Runtime) prepareChatIdentity(ctx context.Context) error {
	origin := chatOrigin(rt.Config)
	if !validChatOrigin(origin, rt.Config.Origin) {
		return errors.New("Chat needs a separate secure hostname")
	}
	api := newIdentity(rt.Config)
	if err := api.login(ctx, rt.Secrets); err != nil {
		return err
	}
	var app map[string]any
	const id = "admin/workspace-chat"
	if err := api.call(ctx, "/api/get-application?id="+id, nil, &app); err != nil {
		return err
	}
	credentials := rt.Secrets.Clients["chat"]
	if app == nil || credentials.ID == "" || app["clientId"] != credentials.ID || app["clientSecret"] != credentials.Secret {
		return errors.New("Chat's existing MyAccounts client needs repair; its credentials were preserved")
	}
	app["redirectUris"] = []string{origin + "/auth/callback"}
	app["enableAutoSignin"], app["enableSigninSession"] = true, true
	return api.call(ctx, "/api/update-application?id="+id, app, nil)
}

func (p *Portal) chatRoutes(mux *http.ServeMux) {
	if _, ok := p.rt.Config.Catalog.Find("chat"); !ok {
		return
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		if r.Method != "GET" && r.Method != "HEAD" {
			http.Error(w, "Reopen Chat from Workspace. Its address changed; retry this action there.", 409)
			return
		}
		// OAuth transactions started before migration cannot be transplanted into a
		// new host-only cookie. Restart the normal PKCE flow; never forward codes.
		target := chatOrigin(p.rt.Config) + strings.TrimPrefix(r.URL.Path, "/chat")
		if strings.HasPrefix(r.URL.Path, "/chat/auth/") {
			target = chatOrigin(p.rt.Config) + "/auth/login"
		} else if r.URL.RawQuery != "" {
			query := r.URL.Query()
			// Old native media links may contain an engine credential. The new
			// host authenticates with its own cookie; never transplant that URL
			// credential into another origin or the browser's new history entry.
			for key := range query {
				switch strings.ToLower(key) {
				case "auth", "secret", "token", "sid", "apikey":
					query.Del(key)
				}
			}
			if encoded := query.Encode(); encoded != "" {
				target += "?" + encoded
			}
		}
		origin, _ := url.Parse(target)
		if !localVaultHostname(strings.Split(r.Host, ":")[0]) && localVaultHostname(origin.Hostname()) {
			http.Error(w, "Chat is running locally. Its dedicated remote address needs attention; ask the Workspace owner to restart hosting. Your chats are preserved.", 503)
			return
		}
		http.Redirect(w, r, target, http.StatusTemporaryRedirect)
	}
	mux.HandleFunc("/chat", handler)
	mux.HandleFunc("/chat/", handler)
}

func (rt *Runtime) ensureChatHosting(ctx context.Context) error {
	if !rt.Config.Tunnel {
		return nil
	}
	path := filepath.Join(rt.Config.Data, "chat", "workspace-hosting.json")
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	if err := platform.ProtectDir(filepath.Dir(path)); err != nil {
		return err
	}
	var record chatHostingRecord
	if err := platform.LoadJSON(path, &record); err != nil && !os.IsNotExist(err) {
		return err
	}
	origin, dnsID, err := rt.ensureDedicatedTunnelHost(ctx, "chat", record.Origin, func(origin string) error {
		record.Origin, record.State = origin, "pending"
		return platform.SaveJSON(path, record)
	})
	if err != nil {
		record.State, record.Message = "attention", "Chat's dedicated remote address needs attention. "+err.Error()
		_ = platform.SaveJSON(path, record)
		return err
	}
	record.Origin, record.DNSRecordID, record.State, record.Message = origin, dnsID, "ready", "Chat has a permanent dedicated address; updates retain it."
	return platform.SaveJSON(path, record)
}

// Health runs through the permanent public endpoint too; process start is not
// presented as proof that a tunnel hostname is reachable.
func (rt *Runtime) chatHostingProgress(ctx context.Context) {
	if !rt.Config.Tunnel {
		return
	}
	origin := chatOrigin(rt.Config)
	u, _ := url.Parse(origin)
	if localVaultHostname(u.Hostname()) {
		return
	}
	check, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := WaitHTTPDetailed(check, origin+"/healthz", 30*time.Second, nil, nil); err != nil {
		rt.progress.Step("chat", "done", "Chat is running; its dedicated remote address is still connecting. Restart hosting if it remains unavailable.")
	}
}
