package suite

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fahimc/one-workspace/internal/platform"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

// CFEndpoint can be replaced by a local test server; installers use the official API.
var CFEndpoint = "https://api.cloudflare.com/client/v4"

type cfAPI struct {
	token  string
	client *http.Client
}

func (a cfAPI) call(ctx context.Context, method, path string, body any, out any) error {
	var b bytes.Buffer
	if body != nil {
		if e := json.NewEncoder(&b).Encode(body); e != nil {
			return e
		}
	}
	req, e := http.NewRequestWithContext(ctx, method, CFEndpoint+path, &b)
	if e != nil {
		return e
	}
	req.Header.Set("Authorization", "Bearer "+a.token)
	req.Header.Set("Content-Type", "application/json")
	resp, e := a.client.Do(req)
	if e != nil {
		return errors.New("Cloudflare connection failed")
	}
	defer resp.Body.Close()
	var result struct {
		Success bool            `json:"success"`
		Result  json.RawMessage `json:"result"`
	}
	if json.NewDecoder(io.LimitReader(resp.Body, 2<<20)).Decode(&result) != nil || !result.Success || resp.StatusCode >= 300 {
		return fmt.Errorf("Cloudflare rejected %s %s (HTTP %d); check token permissions and account/zone", method, strings.Split(path, "?")[0], resp.StatusCode)
	}
	if out != nil {
		return json.Unmarshal(result.Result, out)
	}
	return nil
}
func ProvisionCloudflare(ctx context.Context, r Request, c Config) (token string, err error) {
	a := cfAPI{r.CFToken, &http.Client{Timeout: 25 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}
	var zone struct {
		Name string `json:"name"`
	}
	if err = a.call(ctx, "GET", "/zones/"+r.ZoneID, nil, &zone); err != nil {
		return
	}
	hosts := []string{}
	for _, origin := range []string{c.Origin, c.Issuer} {
		u, _ := url.Parse(origin)
		h := u.Hostname()
		if h == zone.Name || !strings.HasSuffix(h, "."+zone.Name) {
			return "", errors.New("both hostnames must be subdomains in the selected Cloudflare zone")
		}
		var records []json.RawMessage
		if err = a.call(ctx, "GET", "/zones/"+r.ZoneID+"/dns_records?name="+url.QueryEscape(h), nil, &records); err != nil {
			return
		}
		if len(records) > 0 {
			return "", fmt.Errorf("DNS record %s already exists; choose an unused hostname", h)
		}
		hosts = append(hosts, h)
	}
	var tunnel struct {
		ID string `json:"id"`
	}
	if err = a.call(ctx, "POST", "/accounts/"+r.AccountID+"/cfd_tunnel", map[string]string{"name": "workspace-" + strings.ToLower(strings.Split(hosts[0], ".")[0]) + "-" + time.Now().UTC().Format("20060102150405"), "config_src": "cloudflare"}, &tunnel); err != nil {
		return
	}
	base := "/accounts/" + r.AccountID + "/cfd_tunnel/" + tunnel.ID
	created := []string{}
	defer func() {
		if err == nil {
			return
		}
		cleanup, cancel := context.WithTimeout(context.Background(), 45*time.Second)
		defer cancel()
		failed := false
		for _, id := range created {
			if a.call(cleanup, "DELETE", "/zones/"+r.ZoneID+"/dns_records/"+id, nil, nil) != nil {
				failed = true
			}
		}
		if a.call(cleanup, "DELETE", base, nil, nil) != nil {
			failed = true
		}
		if failed {
			err = fmt.Errorf("%w; automatic cleanup was incomplete: remove tunnel %s and its Workspace DNS records in Cloudflare", err, tunnel.ID)
		}
	}()
	ingress := []any{map[string]string{"hostname": hosts[0], "service": fmt.Sprintf("http://localhost:%d", c.Port)}, map[string]string{"hostname": hosts[1], "service": fmt.Sprintf("http://localhost:%d", c.Port+1)}, map[string]string{"service": "http_status:404"}}
	if err = a.call(ctx, "PUT", base+"/configurations", map[string]any{"config": map[string]any{"ingress": ingress}}, nil); err != nil {
		return
	}
	for _, h := range hosts {
		var record struct {
			ID string `json:"id"`
		}
		if err = a.call(ctx, "POST", "/zones/"+r.ZoneID+"/dns_records", map[string]any{"type": "CNAME", "name": h, "content": tunnel.ID + ".cfargotunnel.com", "proxied": true, "ttl": 1}, &record); err != nil {
			return
		}
		created = append(created, record.ID)
	}
	err = a.call(ctx, "GET", base+"/token", nil, &token)
	if err == nil && c.Data != "" {
		err = platform.SaveJSON(filepath.Join(c.Data, "cloudflare-resources.json"), map[string]any{"accountId": r.AccountID, "zoneId": r.ZoneID, "zoneName": zone.Name, "tunnelId": tunnel.ID, "dnsRecordIds": created, "hostnames": hosts})
	}
	return
}
