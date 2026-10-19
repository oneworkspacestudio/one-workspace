package suite

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/branding"
)

// Clearing a portal cookie cannot close a dedicated-origin WebSocket. Revoke
// the verified account's Chat sessions through a private backchannel; no
// browser, provider or portal credential is forwarded to Chat.
func (p *Portal) closeChatSession(r *http.Request) {
	app, installed := p.rt.Config.Catalog.Find("chat")
	user, authenticated := p.user(r)
	if !installed || !authenticated || len(p.rt.chatDirectoryToken) < 32 {
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	body, _ := json.Marshal(map[string]string{"issuer": p.rt.Config.Issuer, "subject": user.Subject})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://127.0.0.1:%d/api/session/revoke", app.Port), bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+p.rt.chatDirectoryToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 3 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := client.Do(req)
	if err == nil {
		response.Body.Close()
	}
}

func (p *Portal) chatCatalogRoute(mux *http.ServeMux) {
	mux.HandleFunc("POST /_integration/chat/catalog", func(w http.ResponseWriter, r *http.Request) {
		host, _, _ := net.SplitHostPort(r.RemoteAddr)
		if net.ParseIP(host) == nil || !net.ParseIP(host).IsLoopback() || len(p.rt.chatDirectoryToken) < 32 || subtle.ConstantTimeCompare([]byte(r.Header.Get("Authorization")), []byte("Bearer "+p.rt.chatDirectoryToken)) != 1 {
			http.NotFound(w, r)
			return
		}
		var identity struct {
			Issuer  string `json:"issuer"`
			Subject string `json:"subject"`
		}
		decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096))
		decoder.DisallowUnknownFields()
		if decoder.Decode(&identity) != nil || identity.Issuer != p.rt.Config.Issuer || identity.Subject == "" || decoder.Decode(new(any)) != io.EOF {
			http.Error(w, "Invalid identity", 400)
			return
		}
		apps := []map[string]string{}
		for _, app := range p.rt.Config.Catalog.Apps {
			if app.ID == "monitor" && identity.Subject != p.rt.Config.Owner {
				continue
			}
			target := p.rt.Config.Origin + "/_workspace/open/" + app.ID
			if app.ID == "accounts" {
				target = p.rt.Config.Issuer
			}
			if app.ID == "chat" {
				target = chatOrigin(p.rt.Config) + "/"
			}
			apps = append(apps, map[string]string{"id": app.ID, "name": app.Name, "description": app.Description, "url": target, "icon": "img/apps/" + strings.TrimPrefix(branding.Icon(app.ID), "/branding/icons/")})
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"home": p.rt.Config.Origin, "apps": apps})
	})
}
