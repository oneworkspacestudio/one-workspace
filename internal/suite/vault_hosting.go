package suite

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fahimc/one-workspace/internal/catalog"
	"github.com/fahimc/one-workspace/internal/platform"
)

// Vault uses its own browser origin. This record is installation state, never an
// authentication credential. A pending intent makes interrupted provisioning
// retry the same owned hostname instead of allocating a new address each boot.
type vaultHostingRecord struct {
	Origin      string `json:"origin"`
	DNSRecordID string `json:"dnsRecordId,omitempty"`
	State       string `json:"state"`
	Message     string `json:"message,omitempty"`
}

func (rt *Runtime) vaultHostingProgress(ctx context.Context, app catalog.App) {
	path := filepath.Join(rt.Config.Data, "my-vault", "workspace-hosting.json")
	var record vaultHostingRecord
	if platform.LoadJSON(path, &record) != nil || record.State != "ready" {
		return
	}
	var settings struct {
		Origin string `json:"origin"`
	}
	if platform.LoadJSON(filepath.Join(rt.Config.Data, "my-vault", "myvault.json"), &settings) != nil || settings.Origin != record.Origin {
		return
	}
	u, err := url.Parse(record.Origin)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return
	}
	rt.progress.Step(app.ID, "running", "Connecting My Vault at its secure remote address")
	err = WaitHTTPDetailed(ctx, record.Origin+"/healthz", 90*time.Second, nil, func(detail string) {
		rt.progress.Step(app.ID, "running", "Connecting remote My Vault: "+detail)
	})
	if err != nil {
		record.State = "attention"
		record.Message = "Remote My Vault could not be reached. Local independent access is available; restart Workspace to retry its remote connection."
		_ = platform.SaveJSON(path, record)
		rt.progress.Step(app.ID, "done", record.Message)
		return
	}
	rt.progress.Step(app.ID, "done", "My Vault is available locally and at its permanent HTTPS address")
}

func (rt *Runtime) ensureVaultHosting(ctx context.Context) (result error) {
	c := rt.Config
	app, exists := c.Catalog.Find("vault")
	if !exists || !c.Tunnel {
		return nil
	}
	root := filepath.Join(c.Data, "my-vault")
	if err := os.MkdirAll(root, 0700); err != nil {
		return err
	}
	if err := platform.ProtectDir(root); err != nil {
		return err
	}
	settings := map[string]any{}
	settingsPath := filepath.Join(root, "myvault.json")
	if err := platform.LoadJSON(settingsPath, &settings); err != nil && !os.IsNotExist(err) {
		return err
	}
	if settings == nil {
		return errors.New("the Vault settings need local owner repair")
	}
	// A pending owner-origin change takes precedence over automatic hosting.
	if _, err := os.Lstat(filepath.Join(root, "myvault-pending.json")); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	recordPath := filepath.Join(root, "workspace-hosting.json")
	var record vaultHostingRecord
	if err := platform.LoadJSON(recordPath, &record); err != nil && !os.IsNotExist(err) {
		return err
	}
	configured, _ := settings["origin"].(string)
	if value, exists := settings["origin"]; exists && value != configured {
		return errors.New("the Vault HTTPS address is invalid")
	}
	if configured != "" {
		u, err := url.Parse(configured)
		if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
			return errors.New("the Vault HTTPS address needs local owner repair")
		}
	}
	if record.State == "ready" {
		// Persist the address across upgrades. Respect later owner changes too.
		return nil
	}
	if u, err := url.Parse(configured); err == nil && u.Scheme == "https" && u.Host != "" && !localVaultHostname(u.Hostname()) && configured != record.Origin {
		return nil // A manually configured dedicated origin belongs to the owner.
	}
	defer func() {
		if result != nil {
			record.State = "attention"
			record.Message = "Remote My Vault setup could not finish. Restart Workspace to retry; local encrypted-vault access remains available. " + result.Error()
			_ = platform.SaveJSON(recordPath, record)
		}
	}()
	var resources struct {
		AccountID string `json:"accountId"`
		ZoneID    string `json:"zoneId"`
		TunnelID  string `json:"tunnelId"`
	}
	if err := platform.LoadJSON(filepath.Join(c.Data, "cloudflare-resources.json"), &resources); err != nil {
		return errors.New("the managed Cloudflare connection is missing; reconnect hosting in Workspace")
	}
	id := regexp.MustCompile(`^[a-fA-F0-9]{32}$`)
	tunnelID := regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$`)
	if !id.MatchString(resources.AccountID) || !id.MatchString(resources.ZoneID) || !tunnelID.MatchString(resources.TunnelID) {
		return errors.New("managed Cloudflare resource identifiers need repair")
	}
	unlock,lockErr:=platform.LockCloudflareTunnel(ctx,resources.AccountID,resources.TunnelID)
	if lockErr!=nil{return lockErr}
	defer unlock()
	token := rt.Secrets.HostingToken
	if token == "" {
		token = rt.Secrets.MailDNSToken
	} // Existing hosted installations.
	if token == "" {
		return errors.New("reconnect Cloudflare hosting to authorize the additional Vault address")
	}
	api := cfAPI{token, &http.Client{Timeout: 25 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}
	var zone struct {
		Name    string `json:"name"`
		Account struct {
			ID string `json:"id"`
		} `json:"account"`
	}
	if err := api.call(ctx, "GET", "/zones/"+resources.ZoneID, nil, &zone); err != nil {
		return err
	}
	if zone.Account.ID != resources.AccountID || !hostingDomainPattern.MatchString(zone.Name) {
		return errors.New("Cloudflare zone does not match this installation")
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
	base := "/accounts/" + resources.AccountID + "/cfd_tunnel/" + resources.TunnelID + "/configurations"
	var remote struct {
		Config map[string]any `json:"config"`
	}
	if err := api.call(ctx, "GET", base, nil, &remote); err != nil {
		return err
	}
	ingress, ok := remote.Config["ingress"].([]any)
	if !ok || len(ingress) == 0 {
		return errors.New("the managed tunnel has no valid ingress configuration")
	}
	service := fmt.Sprintf("http://127.0.0.1:%d", app.Port)
	host := ""
	if record.Origin != "" {
		u, err := url.Parse(record.Origin)
		if err != nil || u.Scheme != "https" || u.Host != u.Hostname() || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" || !strings.HasSuffix(u.Hostname(), "."+zone.Name) || strings.Contains(strings.TrimSuffix(u.Hostname(), "."+zone.Name), ".") {
			return errors.New("saved Vault hosting address needs owner repair")
		}
		host = u.Hostname()
	} else {
		// Find an unused first-level hostname covered by Cloudflare edge TLS.
		for n := 1; n <= 20; n++ {
			label := "vault"
			if n > 1 {
				label += fmt.Sprintf("-%d", n)
			}
			candidate := label + "." + zone.Name
			records, err := lookup(candidate)
			if err != nil {
				return err
			}
			owned := false
			if len(records) == 1 && records[0].Type == "CNAME" && records[0].Proxied && strings.EqualFold(strings.TrimSuffix(records[0].Content, "."), target) {
				// Recover an earlier successful route after local update rollback,
				// but only when both the tunnel and exact Vault gateway match.
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
			return errors.New("no unused Vault hostname is available in this zone")
		}
		record = vaultHostingRecord{Origin: "https://" + host, State: "pending"}
		if err := platform.SaveJSON(recordPath, record); err != nil {
			return err
		}
	}
	found := false
	for _, entry := range ingress {
		rule, ok := entry.(map[string]any)
		if !ok {
			return errors.New("the managed tunnel has an invalid route")
		}
		if rule["hostname"] == host {
			if rule["service"] != service || (rule["path"] != nil && rule["path"] != "") {
				return errors.New("the selected Vault hostname already routes to another service; existing routes were preserved")
			}
			found = true
		}
	}
	records, err := lookup(host)
	if err != nil {
		return err
	}
	if len(records) > 0 && (len(records) != 1 || records[0].Type != "CNAME" || !strings.EqualFold(strings.TrimSuffix(records[0].Content, "."), target) || !records[0].Proxied) {
		return errors.New("the selected Vault hostname is now used by another DNS record; existing records were preserved")
	}
	if !found {
		last, ok := ingress[len(ingress)-1].(map[string]any)
		if !ok || (last["hostname"] != nil && last["hostname"] != "") {
			return errors.New("the managed tunnel needs a terminal catch-all route")
		}
		route := map[string]any{"hostname": host, "service": service, "originRequest": map[string]any{"httpHostHeader": host}}
		next := append([]any{}, ingress[:len(ingress)-1]...)
		next = append(next, route, last)
		remote.Config["ingress"] = next
		if err := api.call(ctx, "PUT", base, map[string]any{"config": remote.Config}, nil); err != nil {
			return err
		}
	}
	if len(records) == 0 {
		var created dnsRecord
		if err := api.call(ctx, "POST", "/zones/"+resources.ZoneID+"/dns_records", map[string]any{"type": "CNAME", "name": host, "content": target, "proxied": true, "ttl": 1}, &created); err != nil {
			return err
		}
		record.DNSRecordID = created.ID
	} else {
		record.DNSRecordID = records[0].ID
	}
	settings["origin"] = record.Origin
	if err := platform.SaveJSON(settingsPath, settings); err != nil {
		return err
	}
	record.State, record.Message = "ready", "My Vault is available at its permanent HTTPS address. Updates retain this address."
	return platform.SaveJSON(recordPath, record)
}
