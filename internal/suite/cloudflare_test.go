package suite

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCloudflareRollbackDoesNotDeleteExistingRecords(t *testing.T) {
	for _, failure := range []string{"existing", "second DNS", "none"} {
		t.Run(failure, func(t *testing.T) {
			calls := []string{}
			dns := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls = append(calls, r.Method+" "+r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				if r.Header.Get("Authorization") != "Bearer scoped" {
					t.Error("missing auth")
				}
				result := "{}"
				success := true
				switch {
				case r.URL.Path == "/zones/zone":
					result = `{"name":"example.com"}`
				case r.Method == "GET" && strings.Contains(r.URL.Path, "dns_records"):
					result = "[]"
					if failure == "existing" {
						result = `[{"id":"existing"}]`
					}
				case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "cfd_tunnel"):
					result = `{"id":"new-tunnel"}`
				case r.Method == "POST" && strings.Contains(r.URL.Path, "dns_records"):
					dns++
					result = fmt.Sprintf(`{"id":"new-%d"}`, dns)
					if failure == "second DNS" && dns == 2 {
						success = false
					}
				case strings.HasSuffix(r.URL.Path, "/token"):
					result = `"connector-token"`
				}
				fmt.Fprintf(w, `{"success":%t,"result":%s}`, success, result)
			}))
			defer server.Close()
			before := CFEndpoint
			CFEndpoint = server.URL
			defer func() { CFEndpoint = before }()
			token, e := ProvisionCloudflare(context.Background(), Request{CFToken: "scoped", ZoneID: "zone", AccountID: "account"}, Config{Origin: "https://work.example.com", Issuer: "https://login.example.com", Port: 8990})
			joined := strings.Join(calls, "\n")
			if failure == "none" {
				if e != nil || token != "connector-token" || strings.Contains(joined, "DELETE") {
					t.Fatalf("success failed: %v %s", e, joined)
				}
			} else if failure == "existing" {
				if e == nil || strings.Contains(joined, "POST") || strings.Contains(joined, "DELETE") {
					t.Fatalf("changed existing zone: %v %s", e, joined)
				}
			} else {
				if e == nil || !strings.Contains(joined, "DELETE /zones/zone/dns_records/new-1") || !strings.Contains(joined, "DELETE /accounts/account/cfd_tunnel/new-tunnel") {
					t.Fatalf("rollback failed: %v %s", e, joined)
				}
			}
		})
	}
}
