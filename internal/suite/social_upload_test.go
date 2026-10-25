package suite

import (
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestNativeSocialTransferBoundary(t *testing.T) {
	id := strings.Repeat("a", 48)
	token := "Bearer " + strings.Repeat("b", 48)
	base := "/social/api/imports/instagram/native/" + id
	for _, tc := range []struct {
		method, path, token, origin string
		allow                       bool
	}{
		{"POST", base + "/claim", token, "", true}, {"GET", base + "/status", token, "", true}, {"PUT", base + "/chunks?offset=0", token, "", true}, {"POST", base + "/complete", token, "", true},
		{"POST", base + "/import", token, "", false}, {"POST", "/social/api/posts", token, "", false}, {"GET", "/social/", token, "", false}, {"POST", "/social/api/imports/instagram/native", token, "", false},
		{"GET", base + "/claim", token, "", false}, {"POST", base + "/claim", "", "", false}, {"POST", base + "/claim", token, "https://evil.example", false}, {"PUT", base + "/chunks?offset=0&extra=1", token, "", false},
		{"POST", base + "/claim/extra", token, "", false}, {"POST", strings.Replace(base, "/native/", "/native/%61", 1) + "/claim", token, "", false},
	} {
		r := httptest.NewRequest(tc.method, tc.path, nil)
		r.Header.Set("Authorization", tc.token)
		r.Header.Set("Origin", tc.origin)
		if got := nativeSocialRequest(r); got != tc.allow {
			t.Errorf("%s %s: %v", tc.method, tc.path, got)
		}
	}
}

func TestNativeSocialProxyRetainsOrdinarySignInGate(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/imports/instagram/native/") {
			t.Error("ordinary app route reached backend anonymously")
		}
		w.Header().Set("X-Scoped-Backend", "social")
		w.WriteHeader(401)
		fmt.Fprint(w, `{"error":"invalid scoped permission"}`)
	}))
	defer backend.Close()
	u, _ := url.Parse(backend.URL)
	port, _ := strconv.Atoi(u.Port())
	rt := &Runtime{Config: Config{Data: t.TempDir(), Origin: "https://workspace.example.test", Port: 8990, Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "social", Kind: "process", Route: "/social", Port: port, StripPrefix: true}}}}, status: map[string]string{}}
	p := &Portal{rt: rt, flows: map[string]loginFlow{}, sessions: map[string]session{}}
	h := p.Handler("synthetic")
	token := "Bearer " + strings.Repeat("b", 48)
	for _, tc := range []struct {
		method, path string
		forward      bool
	}{{"GET", "/social/", false}, {"POST", "/social/api/posts", false}, {"POST", "/social/api/imports/instagram/native", false}, {"POST", "/social/api/imports/instagram/native/" + strings.Repeat("a", 48) + "/claim", true}} {
		r := httptest.NewRequest(tc.method, "https://workspace.example.test"+tc.path, nil)
		r.Header.Set("Authorization", token)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 401 || (w.Header().Get("X-Scoped-Backend") == "social") != tc.forward {
			t.Fatalf("%s: status%d forward%v", tc.path, w.Code, w.Header().Get("X-Scoped-Backend"))
		}
	}
}
