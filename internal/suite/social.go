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

const socialPath = "/My Photos/My Social/"

func socialCapability(secret string) string {
	if secret == "" {
		return ""
	}
	sum := sha256.Sum256([]byte("workspace-social-storage-v1\x00" + secret))
	return hex.EncodeToString(sum[:])
}
func socialRoot(c Config) string { return filepath.Join(c.Files, "My Photos", "My Social") }

type socialRule struct {
	DenyAll bool                             `json:"denyAll"`
	Allow   struct{ Users, Groups []string } `json:"allow"`
}

// Quantum specific-user/group allows override denyAll, including ancestor rules.
// Refuse those configurations rather than claim a folder name isolates originals.
func safeSocialRules(rules map[string]socialRule) bool {
	denied := false
	for p, rule := range rules {
		normalized := "/" + strings.Trim(strings.ReplaceAll(p, "\\", "/"), "/") + "/"
		if normalized == "//" {
			normalized = "/"
		}
		relevant := strings.HasPrefix(strings.ToLower(socialPath), strings.ToLower(normalized)) || strings.HasPrefix(strings.ToLower(normalized), strings.ToLower(socialPath))
		if relevant && (len(rule.Allow.Users) > 0 || len(rule.Allow.Groups) > 0) {
			return false
		}
		if strings.EqualFold(normalized, socialPath) && rule.DenyAll {
			denied = true
		}
	}
	return denied
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
	marker := filepath.Join(root, ".workspace-social-v1")
	if _, e := os.Stat(marker); os.IsNotExist(e) {
		if !repair {
			return errors.New("storage marker is missing")
		}
		if entries, e := os.ReadDir(root); e == nil && len(entries) > 0 {
			return errors.New("My Photos/My Social already contains files; explicit migration is required")
		}
		if e = os.MkdirAll(root, 0700); e != nil {
			return e
		}
		if e = platform.AtomicWrite(marker, []byte("One Workspace My Social protected originals v1\n")); e != nil {
			return e
		}
	}
	token, e := rt.socialToken(ctx)
	if e != nil {
		return e
	}
	endpoint := "access?" + url.Values{"source": {"My Drive"}, "path": {socialPath}}.Encode()
	if repair {
		if _, e = rt.socialDrive(ctx, "POST", endpoint, map[string]any{"allow": false, "ruleCategory": "all", "value": ""}, token); e != nil {
			return e
		}
	}
	b, e := rt.socialDrive(ctx, "GET", "access?source=My%20Drive", nil, token)
	if e != nil {
		return e
	}
	rules := map[string]socialRule{}
	if json.Unmarshal(b, &rules) != nil || !safeSocialRules(rules) {
		return errors.New("Quantum access rules could expose Social originals; remove conflicting allow rules")
	}
	return nil
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
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"isolated":true,"contract":"quantum-social-v1"}`)
	})
}
