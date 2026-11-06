package suite

import (
	"net/http"
	"strings"
)

// Only the dedicated, capability-authenticated server API bypasses browser SSO.
// AI Service validates the explicit remote grant, queue bounds and revocation.
// Its owner console, browser chat and loopback integration routes remain gated.
func remoteAIServiceRequest(r *http.Request) bool {
	return r.Method == http.MethodPost && r.URL.EscapedPath() == "/aiservice/api/client/chat" && r.URL.RawQuery == "" && r.Header.Get("Origin") == "" && strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ")
}
