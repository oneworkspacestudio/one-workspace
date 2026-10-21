package suite

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCloudflareAddressReviewAndRecovery(t *testing.T) {
	for _, scenario := range []string{"unused", "refuse", "missing review", "replace", "second DNS fails", "token fails", "stale review", "TXT protected", "MX protected", "DNS changes during tunnel", "recovery fails", "lost POST response"} {
		t.Run(scenario, func(t *testing.T) {
			ctx := context.Background()
			old := map[string]any{"id": "old-work", "name": "work.example.com", "type": "A", "content": "192.0.2.1", "ttl": 300, "proxied": false, "comment": "keep comment", "tags": []any{"team:demo"}, "settings": map[string]any{"ipv4_only": false}}
			records := []map[string]any{old, {"id": "old-login", "name": "login.example.com", "type": "CNAME", "content": "old.cfargotunnel.com", "ttl": 1, "proxied": true}, {"id": "old-v6", "name": "work.example.com", "type": "AAAA", "content": "2001:db8::1", "ttl": 300, "proxied": false}, {"id": "untouched", "name": "other.example.com", "type": "TXT", "content": "keep"}}
			if scenario == "unused" {
				records = records[3:]
			}
			if scenario == "TXT protected" {
				old["type"] = "TXT"
			}
			if scenario == "MX protected" {
				old["type"] = "MX"
			}
			calls := []string{}
			newDNS, restored := 0, 0
			snapshotWritten := false
			data := t.TempDir()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				calls = append(calls, req.Method+" "+req.URL.Path)
				if req.Header.Get("Authorization") != "Bearer scoped" {
					t.Error("missing auth")
				}
				w.Header().Set("Content-Type", "application/json")
				var result any = map[string]any{}
				success := true
				switch {
				case req.URL.Path == "/zones/zone":
					result = map[string]string{"name": "example.com"}
				case req.Method == "GET" && strings.HasSuffix(req.URL.Path, "/dns_records"):
					found := []map[string]any{}
					for _, record := range records {
						if record["name"] == req.URL.Query().Get("name") {
							found = append(found, record)
						}
					}
					result = found
				case req.Method == "POST" && strings.HasSuffix(req.URL.Path, "cfd_tunnel"):
					result = map[string]string{"id": "new-tunnel"}
					if scenario == "DNS changes during tunnel" {
						old["content"] = "192.0.2.99"
					}
				case req.Method == "DELETE" && strings.Contains(req.URL.Path, "/dns_records/"):
					if _, err := os.Stat(filepath.Join(data, "cloudflare-dns-before.json")); err == nil {
						snapshotWritten = true
					}
					id := req.URL.Path[strings.LastIndex(req.URL.Path, "/")+1:]
					kept := []map[string]any{}
					for _, record := range records {
						if record["id"] != id {
							kept = append(kept, record)
						}
					}
					records = kept
				case req.Method == "POST" && strings.HasSuffix(req.URL.Path, "/dns_records"):
					var record map[string]any
					_ = json.NewDecoder(req.Body).Decode(&record)
					if record["content"] == "new-tunnel.cfargotunnel.com" {
						newDNS++
						if (scenario == "second DNS fails" || scenario == "recovery fails") && newDNS == 2 {
							success = false
							break
						}
						record["id"] = fmt.Sprintf("new-%d", newDNS)
						if scenario == "lost POST response" {
							records = append(records, record)
							success = false
							break
						}
					} else {
						restored++
						if scenario == "recovery fails" {
							success = false
							break
						}
						record["id"] = fmt.Sprintf("restored-%d", restored)
						if record["name"] == "work.example.com" && record["type"] == "A" {
							if record["comment"] != "keep comment" || record["ttl"] != float64(300) || record["proxied"] != false || record["tags"] == nil || record["settings"] == nil {
								t.Error("lost record metadata")
							}
						}
					}
					if success {
						records = append(records, record)
						result = record
					}
				case strings.HasSuffix(req.URL.Path, "/token"):
					result = "connector-token"
					success = scenario != "token fails"
				}
				_ = json.NewEncoder(w).Encode(map[string]any{"success": success, "result": result})
			}))
			defer server.Close()
			before := CFEndpoint
			CFEndpoint = server.URL
			defer func() { CFEndpoint = before }()
			r := Request{CFToken: "scoped", ZoneID: "zone", AccountID: "account", Domain: "example.com", Origin: "work.example.com", Issuer: "https://login.example.com"}
			review, err := ReviewHosting(ctx, r)
			if strings.HasSuffix(scenario, "protected") {
				if err == nil || strings.Contains(strings.Join(calls, "\n"), "POST") {
					t.Fatal("protected record was allowed")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if review.Origin != "https://work.example.com" || review.Issuer != "https://login.example.com" {
				t.Fatal(review)
			}
			r.Origin, r.Issuer = review.Origin, review.Issuer
			r.DNSReview, r.ReplaceDNS = review.Review, scenario != "refuse"
			if scenario == "missing review" {
				r.DNSReview = ""
			}
			if scenario == "stale review" {
				old["content"] = "192.0.2.99"
			}
			token, err := ProvisionCloudflare(ctx, r, Config{Origin: r.Origin, Issuer: r.Issuer, Port: 8990, Data: data})
			joined := strings.Join(calls, "\n")
			switch scenario {
			case "unused", "replace":
				if err != nil || token != "connector-token" {
					t.Fatalf("%v: %s", err, joined)
				}
				if scenario == "replace" && !snapshotWritten {
					t.Error("no durable snapshot before deletion")
				}
				if strings.Contains(joined, "DELETE /accounts/") {
					t.Error("deleted successful tunnel")
				}
			case "refuse", "missing review", "stale review":
				if err == nil || strings.Contains(joined, "POST") || strings.Contains(joined, "DELETE") {
					t.Fatalf("mutated unapproved records: %v %s", err, joined)
				}
			case "DNS changes during tunnel":
				if err == nil || strings.Contains(joined, "DELETE /zones/") || old["content"] != "192.0.2.99" {
					t.Fatalf("clobbered concurrent change: %v", err)
				}
			case "recovery fails":
				if err == nil || !strings.Contains(err.Error(), "recovery was incomplete") {
					t.Fatalf("missing recovery warning: %v", err)
				}
			default:
				wantRestored := 3
				if scenario == "lost POST response" {
					wantRestored = 2
				}
				if err == nil || restored != wantRestored || strings.Contains(err.Error(), "recovery was incomplete") {
					t.Fatalf("failed recovery: %v restored=%d calls=%s", err, restored, joined)
				}
				if !strings.Contains(joined, "DELETE /accounts/account/cfd_tunnel/new-tunnel") {
					t.Error("did not remove new tunnel")
				}
			}
			for _, record := range records {
				if record["id"] == "untouched" && record["content"] == "keep" {
					return
				}
			}
			t.Fatal("unrelated record changed")
		})
	}
}

func TestHostingAddressValidation(t *testing.T) {
	for _, raw := range []string{"example.com", "https://work.evil.com", "http://work.example.com", "work.example.com:443", "https://user@work.example.com", "https://work.example.com/path", "https://work.example.com/?q=1", "https://work.example.com/#x", "a.b.example.com"} {
		if _, err := hostingAddress(raw, "example.com", true); err == nil {
			t.Errorf("accepted %q", raw)
		}
	}
	got, err := hostingAddress("https://Work.Example.com.", "example.com", true)
	if err != nil || got != "https://work.example.com" {
		t.Fatalf("%s %v", got, err)
	}
}
