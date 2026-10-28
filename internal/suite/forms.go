package suite

import (
	"bytes"
	"context"
	"crypto/sha256"
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
	"regexp"
	"strings"
	"time"
)

const formsSource = "My Forms"

func formsRoot(c Config) string { return filepath.Join(c.Data, "forms-drive") }
func formsCapability(secret string) string {
	if secret == "" {
		return ""
	}
	s := sha256.Sum256([]byte("workspace-forms-drive-v1\x00" + secret))
	return hex.EncodeToString(s[:])
}
func ensureFormsDriveSource(c Config) error {
	if _, ok := c.Catalog.Find("forms"); !ok {
		return nil
	}
	if _, ok := c.Catalog.Find("drive"); !ok {
		return errors.New("My Forms requires My Drive")
	}
	root := formsRoot(c)
	if containedPath(c.Files, root) || containedPath(root, c.Files) {
		return errors.New("Forms source overlaps shared files")
	}
	if e := platform.RejectReparse(root); e != nil && !os.IsNotExist(e) {
		return e
	}
	if e := os.MkdirAll(root, 0700); e != nil {
		return e
	}
	file := filepath.Join(c.Data, "drive", "config.yaml")
	b, e := os.ReadFile(file)
	if e != nil {
		return e
	}
	var cfg map[string]any
	if json.Unmarshal(b, &cfg) != nil {
		return errors.New("Unsupported Drive configuration")
	}
	server, ok := cfg["server"].(map[string]any)
	if !ok {
		return errors.New("Missing Drive server")
	}
	sources, ok := server["sources"].([]any)
	if !ok {
		return errors.New("Missing Drive sources")
	}
	policy := map[string]any{"defaultEnabled": false, "defaultPermissions": map[string]bool{"view": false, "download": false, "create": false, "modify": false, "delete": false}}
	found := false
	for _, v := range sources {
		src, ok := v.(map[string]any)
		if !ok {
			return errors.New("Invalid source")
		}
		p, _ := src["path"].(string)
		if src["name"] == formsSource {
			if !strings.EqualFold(filepath.Clean(p), filepath.Clean(root)) {
				return errors.New("Forms source path conflict")
			}
			src["config"] = policy
			found = true
		} else if p == "" || containedPath(p, root) || containedPath(root, p) {
			return errors.New("Drive source overlaps private Forms")
		}
	}
	if !found {
		sources = append(sources, map[string]any{"path": root, "name": formsSource, "config": policy})
	}
	server["sources"] = sources
	return platform.SaveJSON(file, cfg)
}

var formsPath = regexp.MustCompile(`^/(?:[a-f0-9]{64}(?:/[a-f0-9]{48}(?:/(?:events|responses|uploads|exports)(?:/[a-zA-Z0-9._-]{1,120})?)?)?)?$`)

type formsDriveCall struct {
	Operation string `json:"operation"`
	Path      string `json:"path"`
	Data      []byte `json:"data,omitempty"`
	Subject   string `json:"subject,omitempty"`
}

func (rt *Runtime) formsRequest(ctx context.Context, method, endpoint string, body []byte, content bool) ([]byte, int, error) {
	a, ok := rt.Config.Catalog.Find("drive")
	if !ok {
		return nil, 503, errors.New("Drive missing")
	}
	token, e := rt.socialToken(ctx)
	if e != nil {
		return nil, 503, e
	}
	req, e := http.NewRequestWithContext(ctx, method, fmt.Sprintf("http://127.0.0.1:%d/drive/api/%s", a.Port, endpoint), bytes.NewReader(body))
	if e != nil {
		return nil, 503, e
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Origin", rt.Config.Origin)
	req.Header.Set("Content-Type", "application/json")
	if content {
		req.Header.Set("Content-Type", "application/octet-stream")
		req.Header.Set("X-File-Total-Size", fmt.Sprint(len(body)))
		sum := sha256.Sum256(body)
		req.Header.Set("X-File-Upload-Session", hex.EncodeToString(sum[:])[:32])
	}
	client := http.Client{Timeout: 45 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	resp, e := client.Do(req)
	if e != nil {
		return nil, 503, e
	}
	defer resp.Body.Close()
	b, e := io.ReadAll(io.LimitReader(resp.Body, 24<<20))
	return b, resp.StatusCode, e
}

// Existing Drive installations do not automatically add new disabled-by-default
// sources to the recovery service account. Grant only this managed source, through
// Drive's user API, while preserving every other account and source permission.
func formsRecoveryScope(user map[string]any) (bool, error) {
	permissions, _ := user["permissions"].(map[string]any)
	if user["username"] != "workspace-recovery" || user["loginMethod"] != "password" || permissions["admin"] != true {
		return false, errors.New("Drive recovery identity does not match the managed service account")
	}
	want := map[string]any{"name": formsSource, "scope": "/", "permissions": map[string]any{"view": true, "download": true, "create": true, "modify": false, "delete": false}}
	scopes, _ := user["scopes"].([]any)
	for i, v := range scopes {
		scope, _ := v.(map[string]any)
		if scope["name"] != formsSource {
			continue
		}
		if scope["scope"] != "/" {
			return false, errors.New("Unexpected Forms recovery scope; review its Drive policy")
		}
		old, _ := json.Marshal(scope)
		expected, _ := json.Marshal(want)
		if bytes.Equal(old, expected) {
			return false, nil
		}
		scopes[i] = want
		user["scopes"] = scopes
		return true, nil
	}
	user["scopes"] = append(scopes, want)
	return true, nil
}

func (rt *Runtime) ensureFormsRecoveryAccess(ctx context.Context) error {
	token, e := rt.socialToken(ctx)
	if e != nil {
		return e
	}
	rt.socialGuardMu.Lock()
	defer rt.socialGuardMu.Unlock()
	endpoint := "users?username=workspace-recovery"
	b, e := rt.socialDrive(ctx, "GET", endpoint, nil, token)
	if e != nil {
		return e
	}
	var user map[string]any
	if json.Unmarshal(b, &user) != nil {
		return errors.New("Invalid Drive recovery profile")
	}
	changed, e := formsRecoveryScope(user)
	if e != nil {
		return e
	}
	if changed {
		if _, e = rt.socialDrive(ctx, "PATCH", endpoint, map[string]any{"data": user, "which": []string{"scopes"}}, token); e != nil {
			return e
		}
		rt.socialGuardToken = ""
	}
	return nil
}
func (rt *Runtime) bindFormsAccount(ctx context.Context, subject string) error {
	if subject == "" || len(subject) > 256 || strings.ContainsAny(subject, "\x00\r\n") {
		return errors.New("Invalid identity")
	}
	// This contract accepts subjects only from the authenticated native Forms app.
	token, e := rt.socialToken(ctx)
	if e != nil {
		return e
	}
	rt.socialGuardMu.Lock()
	defer rt.socialGuardMu.Unlock()
	endpoint := "users?" + url.Values{"username": {subject}}.Encode()
	b, e := rt.socialDrive(ctx, "GET", endpoint, nil, token)
	creating := e != nil
	var user map[string]any
	if creating {
		user = map[string]any{"username": subject, "loginMethod": "oidc", "permissions": map[string]bool{"admin": false, "api": false, "share": false}, "scopes": []any{map[string]any{"name": "My Drive", "scope": "/", "permissions": map[string]bool{"view": true, "download": true, "create": true, "modify": true, "delete": true}}}}
	} else if json.Unmarshal(b, &user) != nil {
		return errors.New("Invalid Drive identity")
	}
	if user["username"] != subject || user["loginMethod"] != "oidc" {
		return errors.New("Drive identity mismatch")
	}
	sum := sha256.Sum256([]byte(rt.Config.Issuer + "\x00" + subject))
	scope := "/" + hex.EncodeToString(sum[:])
	scopes, _ := user["scopes"].([]any)
	for _, v := range scopes {
		s, ok := v.(map[string]any)
		if !ok {
			continue
		}
		if s["name"] == formsSource {
			p, _ := s["permissions"].(map[string]any)
			if s["scope"] != scope || p["view"] != true || p["download"] != true || p["create"] == true || p["modify"] == true || p["delete"] == true {
				return errors.New("Unsafe Forms scope")
			}
			return nil
		}
	}
	if creating {
		if _, e = rt.socialDrive(ctx, "POST", "users", map[string]any{"data": user}, token); e != nil {
			return e
		}
	}
	scopes = append(scopes, map[string]any{"name": formsSource, "scope": scope, "permissions": map[string]bool{"view": true, "download": true, "create": false, "modify": false, "delete": false}})
	user["scopes"] = scopes
	_, e = rt.socialDrive(ctx, "PATCH", endpoint, map[string]any{"data": user, "which": []string{"scopes"}}, token)
	return e
}
func (p *Portal) formsStorageRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /_integration/forms/drive", func(w http.ResponseWriter, r *http.Request) {
		host, _, e := net.SplitHostPort(r.RemoteAddr)
		ip := net.ParseIP(host)
		want := formsCapability(p.rt.Secrets.Clients["forms"].Secret)
		if e != nil || ip == nil || !ip.IsLoopback() || want == "" || !platformSecretEqual(want, strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")) {
			http.NotFound(w, r)
			return
		}
		var in formsDriveCall
		d := json.NewDecoder(http.MaxBytesReader(w, r.Body, 32<<20))
		d.DisallowUnknownFields()
		if d.Decode(&in) != nil {
			http.Error(w, "Invalid Forms storage request", 400)
			return
		}
		if in.Operation == "bind" {
			if e = p.rt.bindFormsAccount(r.Context(), in.Subject); e != nil {
				http.Error(w, "Drive permission setup failed", 503)
				return
			}
			w.WriteHeader(204)
			return
		}
		if !formsPath.MatchString(in.Path) || strings.Contains(in.Path, "..") || len(in.Data) > 20<<20 {
			http.Error(w, "Invalid Forms object path", 400)
			return
		}
		q := url.Values{"source": {formsSource}, "path": {in.Path}}
		method, endpoint, content := "GET", "", false
		switch in.Operation {
		case "list":
			endpoint = "resources/items?" + q.Encode()
		case "read":
			q.Del("path")
			q.Set("file", in.Path)
			endpoint = "raw?" + q.Encode()
		case "write":
			if strings.Count(in.Path, "/") != 4 {
				http.Error(w, "Object path required", 400)
				return
			}
			method = "POST"
			q.Set("override", "false")
			endpoint = "resources?" + q.Encode()
			content = true
		case "mkdir":
			method = "POST"
			q.Set("isDir", "true")
			q.Set("override", "false")
			endpoint = "resources?" + q.Encode()
		default:
			http.Error(w, "Unsupported storage operation", 400)
			return
		}
		b, status, e := p.rt.formsRequest(r.Context(), method, endpoint, in.Data, content)
		if e != nil {
			http.Error(w, "My Drive is unavailable", 503)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write(b)
	})
}
