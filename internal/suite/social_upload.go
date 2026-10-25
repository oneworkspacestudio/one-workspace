package suite

import (
	"encoding/hex"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Only the archive-scoped native transfer API can pass without portal cookies.
// Social independently validates the hashed, expiring, single-job credential.
// This never grants access to sign-in, import/publication or ordinary app APIs.
func nativeSocialRequest(r *http.Request) bool {
	if r.Header.Get("Origin") != "" || r.URL.RawPath != "" {
		return false
	}
	const prefix = "/social/api/imports/instagram/native/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		return false
	}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, prefix), "/")
	if len(parts) != 2 || !socialTransferSecret(parts[0]) {
		return false
	}
	token := r.Header.Get("Authorization")
	if !strings.HasPrefix(token, "Bearer ") || !socialTransferSecret(strings.TrimPrefix(token, "Bearer ")) {
		return false
	}
	switch parts[1] {
	case "claim", "complete", "problem":
		return r.Method == "POST" && r.URL.RawQuery == ""
	case "status":
		return r.Method == "GET" && r.URL.RawQuery == ""
	case "chunks":
		q, e := url.ParseQuery(r.URL.RawQuery)
		if e != nil {
			return false
		}
		n, e := strconv.ParseInt(q.Get("offset"), 10, 64)
		return r.Method == "PUT" && len(q) == 1 && len(q["offset"]) == 1 && e == nil && n >= 0
	}
	return false
}
func socialTransferSecret(s string) bool {
	if len(s) != 48 || strings.ToLower(s) != s {
		return false
	}
	_, e := hex.DecodeString(s)
	return e == nil
}
