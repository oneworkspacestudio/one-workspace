package suite

import (
	"net/http"
	"regexp"
	"strings"
)

var searchJobPath = regexp.MustCompile(`^/websearch/api/client/jobs/[a-f0-9]{48}$`)

// Server calls bypass only the hub's browser login. Native Web Search still
// authenticates each request, checks its job ownership and enforces revocation.
func remoteWebSearchRequest(r *http.Request) bool {
	if r.Header.Get("Origin") != "" || r.URL.RawQuery != "" || !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
		return false
	}
	path := r.URL.EscapedPath()
	return r.Method == http.MethodPost && path == "/websearch/api/client/search" || (r.Method == http.MethodGet || r.Method == http.MethodDelete) && searchJobPath.MatchString(path)
}
