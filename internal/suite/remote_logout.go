package suite

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Owner sign-out must revoke existing desktop streams, not only HTTP cookies.
func (p *Portal) revokeRemoteDesktop(w http.ResponseWriter, r *http.Request) bool {
	app, installed := p.rt.Config.Catalog.Find("remote")
	user, signedIn := p.user(r)
	if !installed || !signedIn || user.Subject != p.rt.Config.Owner || p.rt.Status()["remote"] != "running" {
		return true
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("http://127.0.0.1:%d/_integration/revoke-owner-sessions", app.Port), nil)
	if err == nil {
		req.Header.Set("Origin", p.rt.Config.Origin)
		req.Header.Set("Authorization", "Bearer "+p.rt.Secrets.Clients["remote"].Secret)
		transport := &http.Transport{Proxy: nil}
		defer transport.CloseIdleConnections()
		client := http.Client{Transport: transport, Timeout: 5 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
		if resp, e := client.Do(req); e == nil {
			defer resp.Body.Close()
			var ack struct {
				OK bool `json:"ok"`
			}
			if resp.StatusCode == 200 && json.NewDecoder(io.LimitReader(resp.Body, 4096)).Decode(&ack) == nil && ack.OK {
				return true
			}
		}
	}
	w.Header().Set("Cache-Control", "no-store")
	http.Error(w, "Could not confirm My Remote desktop connections stopped. Retry sign out, or stop My Remote first.", http.StatusServiceUnavailable)
	return false
}
