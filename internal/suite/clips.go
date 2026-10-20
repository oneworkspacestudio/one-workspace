package suite

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const clipsSource = "My Clips Originals"

func clipsCapability(secret string) string {
	if secret == "" {
		return ""
	}
	sum := sha256.Sum256([]byte("workspace-clips-storage-v1\x00" + secret))
	return hex.EncodeToString(sum[:])
}
func clipsRoot(c Config) string {
	return filepath.Join(c.Data, "videos-clips", "My Photos", "My Clips")
}
func clipsSourceConfig(c Config) map[string]any {
	return map[string]any{"path": clipsRoot(c), "name": clipsSource, "config": map[string]any{"defaultEnabled": false, "defaultPermissions": map[string]bool{"view": false, "download": false, "create": false, "modify": false, "delete": false}}}
}

// The dedicated source is outside every broad source. Quantum's physical scope
// containment, rather than case-sensitive folder ACLs, isolates each account.
func ensureClipsDriveSource(c Config) error {
	if _, ok := c.Catalog.Find("clips"); !ok {
		return nil
	}
	if e := os.MkdirAll(clipsRoot(c), 0700); e != nil {
		return e
	}
	file := filepath.Join(c.Data, "drive", "config.yaml")
	b, e := os.ReadFile(file)
	if e != nil {
		return e
	}
	var cfg map[string]any
	if json.Unmarshal(b, &cfg) != nil {
		return errors.New("Drive config is not supported JSON")
	}
	server, ok := cfg["server"].(map[string]any)
	if !ok {
		return errors.New("Drive server settings missing")
	}
	sources, ok := server["sources"].([]any)
	if !ok {
		return errors.New("Drive sources missing")
	}
	found := false
	for _, raw := range sources {
		src, ok := raw.(map[string]any)
		if !ok {
			return errors.New("invalid Drive source")
		}
		p, _ := src["path"].(string)
		name, _ := src["name"].(string)
		if name == clipsSource {
			if !strings.EqualFold(filepath.Clean(p), filepath.Clean(clipsRoot(c))) {
				return errors.New("Clips source path conflict")
			}
			src["config"] = clipsSourceConfig(c)["config"]
			found = true
			continue
		}
		if p == "" || containedPath(p, clipsRoot(c)) || containedPath(clipsRoot(c), p) {
			return errors.New("broad Drive source overlaps private Clips storage")
		}
	}
	if !found {
		sources = append(sources, clipsSourceConfig(c))
	}
	server["sources"] = sources
	return platform.SaveJSON(file, cfg)
}
func (rt *Runtime) verifyClipsStorage(ctx context.Context, repair bool) error {
	root := clipsRoot(rt.Config)
	if containedPath(rt.Config.Files, root) || containedPath(root, rt.Config.Files) {
		return errors.New("Clips storage overlaps shared Files")
	}
	for p := root; ; p = filepath.Dir(p) {
		if _, e := os.Lstat(p); e == nil {
			if e = platform.RejectReparse(p); e != nil {
				return e
			}
		}
		if filepath.Dir(p) == p {
			break
		}
	}
	marker := filepath.Join(root, ".workspace-clips-v1")
	if _, e := os.Stat(marker); os.IsNotExist(e) {
		if !repair {
			return errors.New("storage marker is missing")
		}
		if entries, e := os.ReadDir(root); e == nil && len(entries) > 0 {
			return errors.New("private Clips destination is not empty")
		}
		if e = os.MkdirAll(root, 0700); e != nil {
			return e
		}
		if e = platform.AtomicWrite(marker, []byte("One Workspace scoped My Photos originals v2\n")); e != nil {
			return e
		}
	}
	return nil
}
func clipsScope(subject, issuer string) string {
	sum := sha256.Sum256([]byte(issuer + "\x00" + subject))
	return "/users/" + hex.EncodeToString(sum[:])
}
func (rt *Runtime) bindClipsAccount(ctx context.Context, subject string) error {
	if subject == "" || len(subject) > 256 || strings.ContainsAny(subject, "\r\n\x00") {
		return errors.New("invalid Clips identity")
	}
	token, e := rt.socialToken(ctx)
	if e != nil {
		return e
	}
	rt.socialGuardMu.Lock()
	defer rt.socialGuardMu.Unlock()
	endpoint := "users?" + url.Values{"username": {subject}}.Encode()
	b, e := rt.socialDrive(ctx, "GET", endpoint, nil, token)
	var user map[string]any
	creating := e != nil
	if creating {
		user = map[string]any{"username": subject, "loginMethod": "oidc", "permissions": map[string]bool{"admin": false, "api": false, "share": false}, "scopes": []any{map[string]any{"name": "My Drive", "scope": "/", "permissions": map[string]bool{"view": true, "download": true, "create": true, "modify": true, "delete": true}}}}
	} else if json.Unmarshal(b, &user) != nil {
		return errors.New("invalid native user")
	}
	if user["username"] != subject || user["loginMethod"] != "oidc" {
		return errors.New("native identity does not match Clips OIDC identity")
	}
	scopes, _ := user["scopes"].([]any)
	want := clipsScope(subject, rt.Config.Issuer)
	for _, v := range scopes {
		scope, ok := v.(map[string]any)
		if !ok {
			continue
		}
		if scope["name"] == clipsSource {
			perms, _ := scope["permissions"].(map[string]any)
			if scope["scope"] != want || perms["view"] != true || perms["download"] != true || perms["modify"] == true || perms["create"] == true || perms["delete"] == true {
				return errors.New("unsafe Clips native scope")
			}
			return nil
		}
	}
	if e = os.MkdirAll(filepath.Join(clipsRoot(rt.Config), filepath.FromSlash(strings.TrimPrefix(want, "/"))), 0700); e != nil {
		return e
	}
	scopes = append(scopes, map[string]any{"name": clipsSource, "scope": want, "permissions": map[string]bool{"view": true, "download": true, "create": false, "modify": false, "delete": false}})
	user["scopes"] = scopes
	body := map[string]any{"data": user, "which": []string{"scopes"}}
	if creating {
		// Native CreateUser applies global defaults to incoming scope permissions.
		// Create without the private scope, then apply its explicit read-only grant.
		initial := map[string]any{}
		for key, value := range user {
			initial[key] = value
		}
		initial["scopes"] = scopes[:len(scopes)-1]
		if _, e = rt.socialDrive(ctx, "POST", "users", map[string]any{"data": initial}, token); e != nil {
			return e
		}
	}
	_, e = rt.socialDrive(ctx, "PATCH", endpoint, body, token)
	return e
}
func (p *Portal) clipsStorageRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /_integration/clips/storage", func(w http.ResponseWriter, r *http.Request) {
		host, _, e := net.SplitHostPort(r.RemoteAddr)
		ip := net.ParseIP(host)
		want := clipsCapability(p.rt.Secrets.Clients["clips"].Secret)
		got := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if e != nil || ip == nil || !ip.IsLoopback() || want == "" || subtle.ConstantTimeCompare([]byte(want), []byte(got)) != 1 {
			http.NotFound(w, r)
			return
		}
		if e = p.rt.verifyClipsStorage(r.Context(), false); e != nil {
			http.Error(w, "My Photos isolation needs attention", 503)
			return
		}
		subject := r.Header.Get("X-Clips-Subject")
		if subject != "" {
			if r.Header.Get("X-Clips-Issuer") != p.rt.Config.Issuer {
				http.Error(w, "identity mismatch", 403)
				return
			}
			if e = p.rt.bindClipsAccount(r.Context(), subject); e != nil {
				http.Error(w, "Private My Photos scope unavailable", 503)
				return
			}
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"isolated":true,"contract":"quantum-clips-v1"}`)
	})
}
