package suite

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type HostingConflict struct {
	Hostname string `json:"hostname"`
	Type     string `json:"type"`
	Content  string `json:"content"`
}

type HostingReview struct {
	Origin    string            `json:"origin"`
	Issuer    string            `json:"issuer"`
	Review    string            `json:"review"`
	Conflicts []HostingConflict `json:"conflicts"`
}

func (review HostingReview) validateConsent(r Request) error {
	if r.DNSReview != "" && r.DNSReview != review.Review {
		return errors.New("the addresses or DNS records changed since your review; check the addresses again")
	}
	if len(review.Conflicts) > 0 && (!r.ReplaceDNS || r.DNSReview == "") {
		return errors.New("these addresses are already in use; check them and explicitly confirm replacing their web address records, or choose unused addresses")
	}
	return nil
}

type hostingDNSRecord map[string]json.RawMessage

func (r hostingDNSRecord) text(key string) string {
	var value string
	_ = json.Unmarshal(r[key], &value)
	return value
}

// Only writable fields are replayed. Provider IDs and read-only metadata are
// retained in the recovery snapshot but never submitted as record properties.
func (r hostingDNSRecord) writable() map[string]json.RawMessage {
	out := map[string]json.RawMessage{}
	for _, key := range []string{"type", "name", "content", "ttl", "proxied", "comment", "tags", "settings"} {
		if value, ok := r[key]; ok {
			out[key] = value
		}
	}
	return out
}

func (a cfAPI) hostingRecords(ctx context.Context, zone, host string) ([]hostingDNSRecord, error) {
	all := []hostingDNSRecord{}
	for page := 1; page <= 100; page++ {
		var records []hostingDNSRecord
		err := a.call(ctx, "GET", fmt.Sprintf("/zones/%s/dns_records?name=%s&per_page=100&page=%d", zone, url.QueryEscape(host), page), nil, &records)
		if err != nil {
			return nil, err
		}
		for _, record := range records {
			if !strings.EqualFold(strings.TrimSuffix(record.text("name"), "."), host) || record.text("id") == "" {
				return nil, errors.New("Cloudflare returned an unexpected DNS record; check the addresses again")
			}
			all = append(all, record)
		}
		if len(records) < 100 {
			sort.Slice(all, func(i, j int) bool { return all[i].text("id") < all[j].text("id") })
			return all, nil
		}
	}
	return nil, errors.New("too many DNS records for this address; manage it in Cloudflare first")
}

func hostingAddress(raw, zone string, standardCertificate bool) (string, error) {
	raw = strings.TrimSpace(raw)
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.User != nil || u.Port() != "" || u.RawQuery != "" || u.Fragment != "" || (u.Path != "" && u.Path != "/") {
		return "", errors.New("enter a web address such as workspace.example.com, without a port or path")
	}
	host, err := NormalizeDomain(u.Hostname())
	if err != nil || host == zone || !strings.HasSuffix(host, "."+zone) {
		return "", errors.New("both addresses must be subdomains of the connected Cloudflare domain")
	}
	// The automatic flow uses the standard Cloudflare wildcard certificate.
	if standardCertificate && strings.Contains(strings.TrimSuffix(host, "."+zone), ".") {
		return "", errors.New("use a direct subdomain such as workspace." + zone + " so the standard secure certificate covers it")
	}
	return "https://" + host, nil
}

// ReviewHosting is read-only. Its digest binds consent to these exact addresses
// and DNS records, including provider modification metadata.
func ReviewHosting(ctx context.Context, r Request) (HostingReview, error) {
	a := cfAPI{r.CFToken, &http.Client{Timeout: 25 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}
	review, _, err := a.reviewHosting(ctx, r)
	return review, err
}

func (a cfAPI) reviewHosting(ctx context.Context, r Request) (HostingReview, [][]hostingDNSRecord, error) {
	var review HostingReview
	var zone struct {
		Name string `json:"name"`
	}
	if err := a.call(ctx, "GET", "/zones/"+r.ZoneID, nil, &zone); err != nil {
		return review, nil, err
	}
	origin, err := hostingAddress(r.Origin, zone.Name, r.Domain != "")
	if err != nil {
		return review, nil, err
	}
	issuer, err := hostingAddress(r.Issuer, zone.Name, r.Domain != "")
	if err != nil {
		return review, nil, err
	}
	if origin == issuer {
		return review, nil, errors.New("Workspace and Accounts need different addresses")
	}
	review = HostingReview{Origin: origin, Issuer: issuer, Conflicts: []HostingConflict{}}
	sets := [][]hostingDNSRecord{}
	for _, address := range []string{origin, issuer} {
		host := strings.TrimPrefix(address, "https://")
		records, err := a.hostingRecords(ctx, r.ZoneID, host)
		if err != nil {
			return review, nil, err
		}
		for _, record := range records {
			typ := record.text("type")
			if typ != "A" && typ != "AAAA" && typ != "CNAME" {
				return review, nil, fmt.Errorf("%s has a %s record. Choose another address; setup only replaces web address records (A, AAAA or CNAME)", host, typ)
			}
			review.Conflicts = append(review.Conflicts, HostingConflict{host, typ, record.text("content")})
		}
		sets = append(sets, records)
	}
	raw, _ := json.Marshal([]any{r.AccountID, r.ZoneID, origin, issuer, sets})
	sum := sha256.Sum256(raw)
	review.Review = hex.EncodeToString(sum[:])
	return review, sets, nil
}
