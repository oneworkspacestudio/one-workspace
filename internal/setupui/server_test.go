package setupui

import (
	"encoding/json"
	"github.com/oneworkspacestudio/one-workspace/internal/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWizardCapabilityOriginAndHost(t *testing.T) {
	s := &Server{token: "secret", origin: "http://localhost:12345"}
	for _, tc := range []struct {
		name, method, path, host, token, origin string
		want                                    int
	}{
		{"page", "GET", "/", "localhost:12345", "", "", 200},
		{"no capability", "GET", "/api/state", "localhost:12345", "", "", 403},
		{"wrong capability", "GET", "/api/state", "localhost:12345", "wrong", "", 403},
		{"state", "GET", "/api/state", "localhost:12345", "secret", "", 200},
		{"DNS rebinding", "GET", "/api/state", "evil.test", "secret", "", 403},
		{"CSRF", "POST", "/api/install", "localhost:12345", "secret", "http://evil.test", 403},
		{"retry capability", "POST", "/api/retry-apps", "localhost:12345", "wrong", "http://localhost:12345", 403},
		{"retry CSRF", "POST", "/api/retry-apps", "localhost:12345", "secret", "http://evil.test", 403},
		{"retry before install", "POST", "/api/retry-apps", "localhost:12345", "secret", "http://localhost:12345", 409},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, s.origin+tc.path, nil)
			r.Host = tc.host
			r.Header.Set("X-Setup-Token", tc.token)
			r.Header.Set("Origin", tc.origin)
			w := httptest.NewRecorder()
			s.Handler().ServeHTTP(w, r)
			if w.Code != tc.want {
				t.Fatalf("got %d want %d", w.Code, tc.want)
			}
		})
	}
}

func TestReviewCustomAddressesKeepsCredentialsServerSide(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.Header.Get("Authorization") != "Bearer private-provider-token" {
			t.Error("review must be authenticated and read only")
		}
		result := `[]`
		if r.URL.Path == "/zones/zone" {
			result = `{"name":"example.com"}`
		}
		_, _ = w.Write([]byte(`{"success":true,"result":` + result + `}`))
	}))
	defer api.Close()
	before := suite.CFEndpoint
	suite.CFEndpoint = api.URL
	defer func() { suite.CFEndpoint = before }()
	s := &Server{token: "secret", origin: "http://localhost:12345", hosting: &suite.Request{Mode: "cloudflare", Domain: "example.com", Origin: "https://workspace-2.example.com", Issuer: "https://accounts-2.example.com", CFToken: "private-provider-token", AccountID: "account", ZoneID: "zone"}}
	r := httptest.NewRequest("POST", s.origin+"/api/hosting/review", strings.NewReader(`{"domain":"example.com","origin":"my-space.example.com","issuer":"my-login.example.com"}`))
	r.Header.Set("Origin", s.origin)
	r.Header.Set("X-Setup-Token", s.token)
	w := httptest.NewRecorder()
	s.Handler().ServeHTTP(w, r)
	if w.Code != 200 || strings.Contains(w.Body.String(), "private-provider-token") {
		t.Fatalf("%d %s", w.Code, w.Body.String())
	}
	var review suite.HostingReview
	if json.Unmarshal(w.Body.Bytes(), &review) != nil || review.Review == "" || review.Origin != "https://my-space.example.com" {
		t.Fatal(w.Body.String())
	}
	request := suite.Request{Origin: review.Origin, Issuer: review.Issuer, DNSReview: review.Review, ReplaceDNS: true, CFToken: "browser-injected"}
	if err := applyHostingSelection(&request, *s.hosting); err != nil {
		t.Fatal(err)
	}
	if request.Origin != review.Origin || request.Issuer != review.Issuer || !request.ReplaceDNS || request.DNSReview != review.Review || request.CFToken != s.hosting.CFToken {
		t.Fatal("custom selection or credential boundary lost")
	}
	if err := applyHostingSelection(&suite.Request{}, *s.hosting); err == nil {
		t.Fatal("unreviewed selection accepted")
	}
}
