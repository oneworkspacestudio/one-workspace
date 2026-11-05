package suite

import (
	"net/http"
	"path"
	"regexp"
	"strings"
)

var survivalAsset = regexp.MustCompile(`^/survival/assets/[A-Za-z0-9_-]+\.(js|css)$`)

// Only immutable embedded handbook/shell assets are public. Personal backup APIs
// retain both the Workspace gate and native per-user apphost authorization.
func publicSurvivalRequest(r *http.Request) bool {
	if (r.Method != "GET" && r.Method != "HEAD") || r.URL.RawQuery != "" || path.Clean(r.URL.Path) != strings.TrimSuffix(r.URL.Path, "/") {
		return false
	}
	switch r.URL.Path {
	case "/survival", "/survival/", "/survival/index.html", "/survival/portable.html", "/survival/sw.js", "/survival/manifest.webmanifest", "/survival/offline-manifest.json", "/survival/icon.svg", "/survival/icon-192.png", "/survival/icon-512.png":
		return true
	}
	return survivalAsset.MatchString(r.URL.Path)
}
