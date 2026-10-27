package suite

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

func pushCapability(secret string) string {
	if secret == "" {
		return ""
	}
	h := sha256.Sum256([]byte("workspace-my-push-chat-v1\x00" + secret))
	return hex.EncodeToString(h[:])
}
func pushEnvironment(c Config, s Secrets) []string {
	a, ok := c.Catalog.Find("push")
	if !ok {
		return nil
	}
	key := pushCapability(s.Clients["chat"].Secret)
	if key == "" {
		return nil
	}
	return []string{fmt.Sprintf("WORKSPACE_PUSH_URL=http://127.0.0.1:%d", a.Port), "WORKSPACE_PUSH_PUBLIC_URL=" + strings.TrimRight(c.Origin, "/") + a.Route, "WORKSPACE_PUSH_CHAT_KEY=" + key}
}

// Phone transport bypasses only the portal cookie. My Push validates its own
// expiring device bearer/connect credential; console and service APIs stay gated.
func (p *Portal) nativePushRequest(r *http.Request) bool {
	a, ok := p.rt.Config.Catalog.Find("push")
	if !ok {
		return false
	}
	path := r.URL.Path
	return r.Method == "GET" && (path == a.Route+"/connection/websocket" || path == a.Route+"/api/device/pending") || r.Method == "POST" && path == a.Route+"/api/device/ack"
}
