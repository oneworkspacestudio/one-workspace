package suite

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var hostingDomainPattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?)+$`)

func NormalizeDomain(domain string) (string, error) {
	domain = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(domain), "."))
	if len(domain) > 220 || !hostingDomainPattern.MatchString(domain) {
		return "", errors.New("enter your domain, for example example.com, without https:// or a path")
	}
	return domain, nil
}

// CloudflareLoginPath is the official connector's login destination. Existing
// credentials are read only; connecting another account never overwrites them.
func CloudflareLoginPath() (string, error) {
	home, e := os.UserHomeDir()
	if e != nil {
		return "", e
	}
	return filepath.Join(home, ".cloudflared", "cert.pem"), nil
}

// ResolveHosting discovers provider details after an explicit owner connection.
// Returned credentials stay in the installer process, never in browser state.
func ResolveHosting(ctx context.Context, domain, certPath string) (Request, error) {
	var r Request
	domain, e := NormalizeDomain(domain)
	if e != nil {
		return r, e
	}
	raw, e := os.ReadFile(certPath)
	if e != nil {
		return r, errors.New("connect your Cloudflare account first")
	}
	var cert struct {
		ZoneID    string `json:"zoneID"`
		AccountID string `json:"accountID"`
		APIToken  string `json:"apiToken"`
		Endpoint  string `json:"endpoint"`
	}
	found := false
	for len(raw) > 0 {
		block, rest := pem.Decode(raw)
		if block == nil {
			break
		}
		raw = rest
		if block.Type == "ARGO TUNNEL TOKEN" {
			if found || json.Unmarshal(block.Bytes, &cert) != nil {
				return r, errors.New("invalid Cloudflare login credential")
			}
			found = true
		}
	}
	id := regexp.MustCompile(`^[a-fA-F0-9]{32}$`)
	if !found || !id.MatchString(cert.ZoneID) || !id.MatchString(cert.AccountID) || cert.APIToken == "" || (cert.Endpoint != "" && cert.Endpoint != "us") {
		return r, errors.New("Cloudflare login is incomplete or uses an unsupported API region; use Advanced connection")
	}
	a := cfAPI{cert.APIToken, &http.Client{Timeout: 25 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}
	var zone struct {
		Name    string `json:"name"`
		Status  string `json:"status"`
		Account struct {
			ID string `json:"id"`
		} `json:"account"`
	}
	if e = a.call(ctx, "GET", "/zones/"+cert.ZoneID, nil, &zone); e != nil {
		return r, e
	}
	if zone.Status != "active" {
		return r, errors.New("finish activating this domain in Cloudflare first; its nameservers are not active yet")
	}
	if zone.Account.ID != cert.AccountID || (domain != zone.Name && !strings.HasSuffix(domain, "."+zone.Name)) {
		return r, errors.New("the saved Cloudflare connection belongs to another domain; use Advanced connection for this account")
	}
	// First-level hostnames work with Cloudflare's standard wildcard certificate,
	// including when the selected email domain is itself a subdomain.
	suffix := ""
	if domain != zone.Name {
		sum := sha256.Sum256([]byte(domain))
		suffix = "-" + hex.EncodeToString(sum[:4])
	}
	hosts := []string{}
	for _, prefix := range []string{"workspace", "accounts"} {
		chosen := ""
		for n := 0; n < 20; n++ {
			label := prefix + suffix
			if n > 0 {
				label += fmt.Sprintf("-%d", n+1)
			}
			host := label + "." + zone.Name
			var records []json.RawMessage
			if e = a.call(ctx, "GET", "/zones/"+cert.ZoneID+"/dns_records?name="+url.QueryEscape(host), nil, &records); e != nil {
				return r, e
			}
			if len(records) == 0 {
				chosen = "https://" + host
				break
			}
		}
		if chosen == "" {
			return r, errors.New("could not find an unused Workspace address; use Advanced connection")
		}
		hosts = append(hosts, chosen)
	}
	return Request{Mode: "cloudflare", Domain: domain, Origin: hosts[0], Issuer: hosts[1], CFToken: cert.APIToken, AccountID: cert.AccountID, ZoneID: cert.ZoneID, Authorize: true}, nil
}
