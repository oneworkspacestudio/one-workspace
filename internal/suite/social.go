package suite

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const socialSource = "My Social Originals"

func socialCapability(secret string) string {
	if secret == "" {
		return ""
	}
	sum := sha256.Sum256([]byte("workspace-social-storage-v2\x00" + secret))
	return hex.EncodeToString(sum[:])
}
func socialRoot(c Config) string {
	return filepath.Join(c.Data, "photos-social", "My Photos", "My Social")
}
func socialSourceConfig(c Config) map[string]any {
	return map[string]any{"path": socialRoot(c), "name": socialSource, "config": map[string]any{"defaultEnabled": false, "defaultPermissions": map[string]bool{"view": false, "download": false, "create": false, "modify": false, "delete": false}}}
}
func containedPath(root, candidate string) bool {
	rel, e := filepath.Rel(root, candidate)
	return e == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// The dedicated source is outside every broad source. Quantum's physical scope
// containment, rather than case-sensitive folder ACLs, isolates each account.
func ensureSocialDriveSource(c Config) error {
	if _, ok := c.Catalog.Find("social"); !ok {
		return nil
	}
	if e := os.MkdirAll(socialRoot(c), 0700); e != nil {
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
		if name == socialSource {
			if !strings.EqualFold(filepath.Clean(p), filepath.Clean(socialRoot(c))) {
				return errors.New("Social source path conflict")
			}
			src["config"] = socialSourceConfig(c)["config"]
			found = true
			continue
		}
		if p == "" || containedPath(p, socialRoot(c)) || containedPath(socialRoot(c), p) {
			return errors.New("broad Drive source overlaps private Social storage")
		}
	}
	if !found {
		sources = append(sources, socialSourceConfig(c))
	}
	server["sources"] = sources
	return platform.SaveJSON(file, cfg)
}
func (rt *Runtime) socialDrive(ctx context.Context, method, endpoint string, body any, token string) ([]byte, error) {
	drive, ok := rt.Config.Catalog.Find("drive")
	if !ok {
		return nil, errors.New("My Social needs Drive/My Photos")
	}
	var reader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reader = bytes.NewReader(b)
	}
	req, e := http.NewRequestWithContext(ctx, method, fmt.Sprintf("http://127.0.0.1:%d/drive/api/%s", drive.Port, endpoint), reader)
	if e != nil {
		return nil, e
	}
	req.Header.Set("Origin", rt.Config.Origin)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Password", url.QueryEscape(rt.Secrets.RecoveryPassword))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	} else {
		req.Header.Set("X-Password", url.QueryEscape(rt.Secrets.RecoveryPassword))
	}
	client := http.Client{Timeout: 8 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	resp, e := client.Do(req)
	if e != nil {
		return nil, errors.New("Quantum storage policy unavailable")
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("Quantum denied the storage isolation check")
	}
	return io.ReadAll(io.LimitReader(resp.Body, 1<<20))
}
func (rt *Runtime) socialToken(ctx context.Context) (string, error) {
	rt.socialGuardMu.Lock()
	defer rt.socialGuardMu.Unlock()
	if rt.socialGuardToken != "" && time.Since(rt.socialGuardTokenAt) < 2*time.Minute {
		return rt.socialGuardToken, nil
	}
	b, e := rt.socialDrive(ctx, "POST", "auth/login?username=workspace-recovery", nil, "")
	if e != nil {
		return "", e
	}
	token := strings.Trim(strings.TrimSpace(string(b)), `"`)
	if strings.Count(token, ".") != 2 {
		return "", errors.New("Quantum did not authenticate storage policy maintenance")
	}
	rt.socialGuardToken = token
	rt.socialGuardTokenAt = time.Now()
	return token, nil
}
func (rt *Runtime) verifySocialStorage(ctx context.Context, repair bool) error {
	root := socialRoot(rt.Config)
	if containedPath(rt.Config.Files, root) || containedPath(root, rt.Config.Files) {
		return errors.New("Social storage overlaps shared Files")
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
	// Never silently relocate existing private originals from the old layout.
	legacy := filepath.Join(rt.Config.Files, "My Photos", "My Social")
	if entries, e := os.ReadDir(legacy); e == nil {
		for _, entry := range entries {
			if entry.Name() != ".workspace-social-v1" && entry.Name() != ".workspace-social-v1-paused" {
				return errors.New("legacy Social media requires an explicit private migration")
			}
		}
	}
	marker := filepath.Join(root, ".workspace-social-v2")
	if _, e := os.Stat(marker); os.IsNotExist(e) {
		if !repair {
			return errors.New("storage marker is missing")
		}
		if entries, e := os.ReadDir(root); e == nil && len(entries) > 0 {
			return errors.New("private Social destination is not empty")
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
func socialScope(subject, issuer string) string {
	sum := sha256.Sum256([]byte(issuer + "\x00" + subject))
	return "/users/" + hex.EncodeToString(sum[:])
}
func (rt *Runtime) bindSocialAccount(ctx context.Context, subject string) error {
	if subject == "" || len(subject) > 256 || strings.ContainsAny(subject, "\r\n\x00") {
		return errors.New("invalid Social identity")
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
		return errors.New("native identity does not match Social OIDC identity")
	}
	scopes, _ := user["scopes"].([]any)
	want := socialScope(subject, rt.Config.Issuer)
	for _, v := range scopes {
		scope, ok := v.(map[string]any)
		if !ok {
			continue
		}
		if scope["name"] == socialSource {
			perms, _ := scope["permissions"].(map[string]any)
			if scope["scope"] != want || perms["view"] != true || perms["download"] != true || perms["modify"] == true || perms["create"] == true || perms["delete"] == true {
				return errors.New("unsafe Social native scope")
			}
			return nil
		}
	}
	if e = os.MkdirAll(filepath.Join(socialRoot(rt.Config), filepath.FromSlash(strings.TrimPrefix(want, "/"))), 0700); e != nil {
		return e
	}
	scopes = append(scopes, map[string]any{"name": socialSource, "scope": want, "permissions": map[string]bool{"view": true, "download": true, "create": false, "modify": false, "delete": false}})
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
func (p *Portal) socialStorageRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /_integration/social/storage", func(w http.ResponseWriter, r *http.Request) {
		host, _, e := net.SplitHostPort(r.RemoteAddr)
		ip := net.ParseIP(host)
		want := socialCapability(p.rt.Secrets.Clients["social"].Secret)
		got := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if e != nil || ip == nil || !ip.IsLoopback() || want == "" || subtle.ConstantTimeCompare([]byte(want), []byte(got)) != 1 {
			http.NotFound(w, r)
			return
		}
		if e = p.rt.verifySocialStorage(r.Context(), false); e != nil {
			http.Error(w, "My Photos isolation needs attention", 503)
			return
		}
		subject := r.Header.Get("X-Social-Subject")
		if subject != "" {
			if r.Header.Get("X-Social-Issuer") != p.rt.Config.Issuer {
				http.Error(w, "identity mismatch", 403)
				return
			}
			if e = p.rt.bindSocialAccount(r.Context(), subject); e != nil {
				http.Error(w, "Private My Photos scope unavailable", 503)
				return
			}
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"isolated":true,"contract":"quantum-social-v2"}`)
	})
}
