package suite

import (
	"net/http"
	"strings"
)

// Only the support shell, static assets and grant endpoints are public.
// Apphost and the private broker still authorize every desktop connection.
func publicRemoteRequest(r *http.Request) bool {
	p := r.URL.Path
	if strings.Contains(p, "..") || strings.ContainsAny(p, "\\%") {
		return false
	}
	if r.Method == http.MethodGet {
		switch p {
		case "/remote/", "/remote/index.html", "/remote/api/bootstrap", "/remote/api/me", "/remote/support-api/status", "/remote/auth/login", "/remote/auth/callback", "/remote/desktop":
			return true
		}
		return strings.HasPrefix(p, "/remote/assets/") && !strings.Contains(strings.TrimPrefix(p, "/remote/assets/"), "/")
	}
	return r.Method == http.MethodPost && (p == "/remote/support-api/redeem" || p == "/remote/support-api/stop")
}
