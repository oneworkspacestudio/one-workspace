package suite

import (
	"context"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHostingConnectionStillAllowsReviewWhenSuggestionsOccupied(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatal("connection must be read only")
		}
		result := `[{"id":"occupied"}]`
		if strings.HasPrefix(r.URL.Path, "/zones/") && !strings.HasSuffix(r.URL.Path, "dns_records") {
			result = `{"name":"example.com","status":"active","account":{"id":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}}`
		}
		fmt.Fprintf(w, `{"success":true,"result":%s}`, result)
	}))
	defer server.Close()
	previous := CFEndpoint
	CFEndpoint = server.URL
	defer func() { CFEndpoint = previous }()
	path := filepath.Join(t.TempDir(), "cert.pem")
	cert := pem.EncodeToMemory(&pem.Block{Type: "ARGO TUNNEL TOKEN", Bytes: []byte(`{"zoneID":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb","accountID":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","apiToken":"fixture"}`)})
	if err := os.WriteFile(path, cert, 0600); err != nil {
		t.Fatal(err)
	}
	request, err := ResolveHosting(context.Background(), "example.com", path)
	if err != nil || request.Origin != "https://workspace.example.com" || request.Issuer != "https://accounts.example.com" {
		t.Fatalf("owner cannot reclaim originals: %+v %v", request, err)
	}
}

func TestHostingReviewReadsEveryPageAndRejectsSameAddress(t *testing.T) {
	pages := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var result any = map[string]string{"name": "example.com"}
		if strings.HasSuffix(r.URL.Path, "dns_records") {
			records := []map[string]string{}
			if r.URL.Query().Get("name") == "work.example.com" {
				pages++
				count, offset := 100, 0
				if r.URL.Query().Get("page") == "2" {
					count, offset = 1, 100
				}
				for i := 0; i < count; i++ {
					records = append(records, map[string]string{"id": fmt.Sprint(i + offset), "name": "work.example.com", "type": "A", "content": "192.0.2.1"})
				}
			}
			result = records
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "result": result})
	}))
	defer server.Close()
	previous := CFEndpoint
	CFEndpoint = server.URL
	defer func() { CFEndpoint = previous }()
	request := Request{ZoneID: "zone", Origin: "work.example.com", Issuer: "login.example.com"}
	review, err := ReviewHosting(context.Background(), request)
	if err != nil || len(review.Conflicts) != 101 || pages != 2 {
		t.Fatalf("incomplete review: %d pages=%d %v", len(review.Conflicts), pages, err)
	}
	request.Issuer = "https://WORK.example.com/"
	if _, err = ReviewHosting(context.Background(), request); err == nil {
		t.Fatal("same canonical addresses accepted")
	}
}
