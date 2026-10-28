package suite

import (
	"net/http"
	"regexp"
	"strings"
)

var formsPublicID = regexp.MustCompile(`^[a-f0-9]{48}$`)

// Public admission exposes only the embedded shell and native respondent APIs.
// Authoring and response administration still require both Workspace and Forms SSO.
func publicFormsRequest(r *http.Request) bool {
	path := r.URL.Path
	if r.Method == http.MethodGet {
		switch path {
		case "/forms/", "/forms/index.html", "/forms/app.js", "/forms/style.css", "/forms/logo.svg", "/forms/api/bootstrap", "/forms/api/me", "/forms/auth/login", "/forms/auth/callback":
			return true
		}
	}
	if !strings.HasPrefix(path, "/forms/api/v1/public/") {
		return false
	}
	rest := strings.TrimPrefix(path, "/forms/api/v1/public/")
	parts := strings.Split(rest, "/")
	if len(parts) == 0 || !formsPublicID.MatchString(parts[0]) {
		return false
	}
	if r.Method == http.MethodGet {
		return len(parts) == 1 || (len(parts) == 3 && (parts[1] == "receipt" || parts[1] == "media") && formsPublicID.MatchString(parts[2]))
	}
	return r.Method == http.MethodPost && len(parts) == 2 && (parts[1] == "responses" || parts[1] == "uploads")
}
