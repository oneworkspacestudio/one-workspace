package suite

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"io"
	"net/http"
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
	r.Origin, r.Issuer = c.Origin, c.Issuer
	review, originals, e := a.reviewHosting(ctx, r)
	if e != nil {
		return "", e
	}
	if err = review.validateConsent(r); err != nil {
		return
	}
	hosts := []string{strings.TrimPrefix(review.Origin, "https://"), strings.TrimPrefix(review.Issuer, "https://")}
	// Save the complete original records before any external mutation. This file
	// contains no API credential and remains in the protected installation folder.
	if len(review.Conflicts) > 0 && c.Data != "" {
		if err = platform.SaveJSON(filepath.Join(c.Data, "cloudflare-dns-before.json"), map[string]any{"zoneId": r.ZoneID, "hostnames": hosts, "records": originals}); err != nil {
			return
		}
	}
	var tunnel struct {
		ID string `json:"id"`
	}
	if err = a.call(ctx, "POST", "/accounts/"+r.AccountID+"/cfd_tunnel", map[string]string{"name": "workspace-" + strings.ToLower(strings.Split(hosts[0], ".")[0]) + "-" + time.Now().UTC().Format("20060102150405"), "config_src": "cloudflare"}, &tunnel); err != nil {
		return
	}
	base := "/accounts/" + r.AccountID + "/cfd_tunnel/" + tunnel.ID
	created := []string{}
	removed := map[string]bool{}
	touched := map[string]bool{}
	defer func() {
		if err == nil {
			return
		}
		cleanup, cancel := context.WithTimeout(context.Background(), 45*time.Second)
		defer cancel()
		failed := false
		for i, h := range hosts {
			if !touched[h] {
				continue
			}
			records, e := a.hostingRecords(cleanup, r.ZoneID, h)
			if e != nil {
				failed = true
				continue
			}
			// Also catches a successful POST whose response was lost. Never remove
			// a record pointing elsewhere, even during recovery.
			for _, record := range records {
				if record.text("type") == "CNAME" && record.text("content") == tunnel.ID+".cfargotunnel.com" {
					if a.call(cleanup, "DELETE", "/zones/"+r.ZoneID+"/dns_records/"+record.text("id"), nil, nil) != nil {
						failed = true
					}
				}
			}
			records, e = a.hostingRecords(cleanup, r.ZoneID, h)
			if e != nil {
				failed = true
				continue
			}
			expected := []hostingDNSRecord{}
			for _, record := range originals[i] {
				if !removed[record.text("id")] {
					expected = append(expected, record)
				}
			}
			if !sameHostingRecords(records, expected) {
				failed = true
				continue
			}
			for _, record := range originals[i] {
				if removed[record.text("id")] && a.call(cleanup, "POST", "/zones/"+r.ZoneID+"/dns_records", record.writable(), nil) != nil {
					failed = true
				}
			}
		}
		if a.call(cleanup, "DELETE", base, nil, nil) != nil {
			failed = true
		}
		if failed {
			err = fmt.Errorf("%w; automatic recovery was incomplete: inspect tunnel %s in Cloudflare and the original records in Private/cloudflare-dns-before.json before retrying", err, tunnel.ID)
		}
	}()
	ingress := []any{map[string]string{"hostname": hosts[0], "service": fmt.Sprintf("http://localhost:%d", c.Port)}, map[string]string{"hostname": hosts[1], "service": fmt.Sprintf("http://localhost:%d", c.Port+1)}, map[string]string{"service": "http_status:404"}}
	if err = a.call(ctx, "PUT", base+"/configurations", map[string]any{"config": map[string]any{"ingress": ingress}}, nil); err != nil {
		return
	}
	for i, h := range hosts {
		// Refuse a change made after the initial review, including during tunnel
		// creation. Cloudflare has no conditional DNS transaction across hosts.
		current, e := a.hostingRecords(ctx, r.ZoneID, h)
		if e != nil {
			err = e
			return
		}
		if !sameHostingRecords(current, originals[i]) {
			err = errors.New("DNS records changed during setup; check the addresses again")
			return
		}
		touched[h] = true
		for _, old := range originals[i] {
			if err = a.call(ctx, "DELETE", "/zones/"+r.ZoneID+"/dns_records/"+old.text("id"), nil, nil); err != nil {
				return
			}
			removed[old.text("id")] = true
		}
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
		var zone struct {
			Name string `json:"name"`
		}
		err = a.call(ctx, "GET", "/zones/"+r.ZoneID, nil, &zone)
		if err == nil {
			err = platform.SaveJSON(filepath.Join(c.Data, "cloudflare-resources.json"), map[string]any{"accountId": r.AccountID, "zoneId": r.ZoneID, "zoneName": zone.Name, "tunnelId": tunnel.ID, "dnsRecordIds": created, "hostnames": hosts})
		}
	}
	return
}

func sameHostingRecords(a, b []hostingDNSRecord) bool {
	x, _ := json.Marshal(a)
	y, _ := json.Marshal(b)
	return bytes.Equal(x, y)
}
