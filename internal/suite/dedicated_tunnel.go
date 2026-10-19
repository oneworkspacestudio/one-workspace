package suite

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

// Allocate a first-level hostname covered by standard Cloudflare edge TLS.
// Preserve unrelated DNS/ingress and serialize with the other managed apps.
func (rt *Runtime) ensureDedicatedTunnelHost(ctx context.Context, id, saved string, pending func(string) error) (string, string, error) {
	fail := func(message string) (string, string, error) { return "", "", errors.New(message) }
	app, ok := rt.Config.Catalog.Find(id)
	if !ok {
		return fail("app is not installed")
	}
	var resources struct {
		AccountID string `json:"accountId"`
		ZoneID    string `json:"zoneId"`
		TunnelID  string `json:"tunnelId"`
	}
	if platform.LoadJSON(filepath.Join(rt.Config.Data, "cloudflare-resources.json"), &resources) != nil {
		return fail("Reconnect the managed Cloudflare hosting connection.")
	}
	hexID := regexp.MustCompile(`^[a-fA-F0-9]{32}$`)
	uuid := regexp.MustCompile(`^[a-fA-F0-9]{8}(-[a-fA-F0-9]{4}){3}-[a-fA-F0-9]{12}$`)
	if !hexID.MatchString(resources.AccountID) || !hexID.MatchString(resources.ZoneID) || !uuid.MatchString(resources.TunnelID) {
		return fail("Managed Cloudflare identifiers need repair.")
	}
	unlock, err := platform.LockCloudflareTunnel(ctx, resources.AccountID, resources.TunnelID)
	if err != nil {
		return "", "", err
	}
	defer unlock()
	token := rt.Secrets.HostingToken
	if token == "" {
		token = rt.Secrets.MailDNSToken
	}
	if token == "" {
		return fail("Reconnect Cloudflare to authorize Chat's dedicated address.")
	}
	api := cfAPI{token, &http.Client{Timeout: 25 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}
	var zone struct {
		Name    string `json:"name"`
		Account struct {
			ID string `json:"id"`
		} `json:"account"`
	}
	if err = api.call(ctx, "GET", "/zones/"+resources.ZoneID, nil, &zone); err != nil {
		return "", "", err
	}
	if zone.Account.ID != resources.AccountID || !hostingDomainPattern.MatchString(zone.Name) {
		return fail("The Cloudflare zone does not match this installation.")
	}
	target := resources.TunnelID + ".cfargotunnel.com"
	type dnsRecord struct {
		ID      string `json:"id"`
		Type    string `json:"type"`
		Content string `json:"content"`
		Proxied bool   `json:"proxied"`
	}
	lookup := func(host string) ([]dnsRecord, error) {
		var records []dnsRecord
		err := api.call(ctx, "GET", "/zones/"+resources.ZoneID+"/dns_records?name="+url.QueryEscape(host), nil, &records)
		return records, err
	}
	ownedDNS := func(records []dnsRecord) bool {
		return len(records) == 1 && records[0].Type == "CNAME" && records[0].Proxied && strings.EqualFold(strings.TrimSuffix(records[0].Content, "."), target)
	}
	base := "/accounts/" + resources.AccountID + "/cfd_tunnel/" + resources.TunnelID + "/configurations"
	var remote struct {
		Config map[string]any `json:"config"`
	}
	if err = api.call(ctx, "GET", base, nil, &remote); err != nil {
		return "", "", err
	}
	ingress, ok := remote.Config["ingress"].([]any)
	if !ok || len(ingress) == 0 {
		return fail("The managed tunnel has no valid routes.")
	}
	service := fmt.Sprintf("http://127.0.0.1:%d", app.Port)
	host := ""
	if saved != "" {
		u, e := url.Parse(saved)
		if e != nil || !validChatOrigin(saved, rt.Config.Origin) || u.Scheme != "https" || u.Host != u.Hostname() || !strings.HasSuffix(u.Hostname(), "."+zone.Name) || strings.Contains(strings.TrimSuffix(u.Hostname(), "."+zone.Name), ".") {
			return fail("The saved dedicated Chat address needs owner repair.")
		}
		host = u.Hostname()
	} else {
		home, _ := url.Parse(rt.Config.Origin)
		for n := 1; n <= 20; n++ {
			label := id
			if n > 1 {
				label += fmt.Sprintf("-%d", n)
			}
			candidate := label + "." + zone.Name
			if strings.EqualFold(candidate, home.Hostname()) {
				continue
			}
			records, e := lookup(candidate)
			if e != nil {
				return "", "", e
			}
			owned := false
			if ownedDNS(records) {
				for _, entry := range ingress {
					if rule, ok := entry.(map[string]any); ok && rule["hostname"] == candidate && rule["service"] == service && (rule["path"] == nil || rule["path"] == "") {
						owned = true
					}
				}
			}
			if len(records) == 0 || owned {
				host = candidate
				break
			}
		}
		if host == "" {
			return fail("No unused Chat hostname is available; existing DNS was preserved.")
		}
		if err = pending("https://" + host + "/chat"); err != nil {
			return "", "", err
		}
	}
	found := false
	for _, entry := range ingress {
		rule, ok := entry.(map[string]any)
		if !ok {
			return fail("The managed tunnel contains an invalid route.")
		}
		if rule["hostname"] == host {
			if rule["service"] != service || (rule["path"] != nil && rule["path"] != "") {
				return fail("Chat's hostname routes to another service; existing routes were preserved.")
			}
			found = true
		}
	}
	records, err := lookup(host)
	if err != nil {
		return "", "", err
	}
	if len(records) > 0 && !ownedDNS(records) {
		return fail("Chat's hostname is now used by another DNS record; it was preserved.")
	}
	if !found {
		last, ok := ingress[len(ingress)-1].(map[string]any)
		if !ok || (last["hostname"] != nil && last["hostname"] != "") {
			return fail("The tunnel needs a terminal catch-all route.")
		}
		next := append([]any{}, ingress[:len(ingress)-1]...)
		next = append(next, map[string]any{"hostname": host, "service": service, "originRequest": map[string]any{"httpHostHeader": host}}, last)
		remote.Config["ingress"] = next
		if err = api.call(ctx, "PUT", base, map[string]any{"config": remote.Config}, nil); err != nil {
			return "", "", err
		}
	}
	recordID := ""
	if len(records) == 0 {
		var record dnsRecord
		if err = api.call(ctx, "POST", "/zones/"+resources.ZoneID+"/dns_records", map[string]any{"type": "CNAME", "name": host, "content": target, "proxied": true, "ttl": 1}, &record); err != nil {
			return "", "", err
		}
		recordID = record.ID
	} else {
		recordID = records[0].ID
	}
	return "https://" + host + "/chat", recordID, nil
}
