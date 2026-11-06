package suite

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

func TestRemoteAIServiceOnlyExemptsExactServerEndpoint(t *testing.T) {
	for _, tc := range []struct {
		method, path, auth, origin string
		want                       bool
	}{
		{"POST", "/aiservice/api/client/chat", "Bearer fixture", "", true},
		{"POST", "/aiservice/api/client/chat", "", "", false},
		{"POST", "/aiservice/api/client/chat", "Bearer fixture", "https://evil.test", false},
		{"GET", "/aiservice/api/client/chat", "Bearer fixture", "", false},
		{"POST", "/aiservice/api/client/chat?token=fixture", "Bearer fixture", "", false},
		{"POST", "/aiservice/api/service/settings", "Bearer fixture", "", false},
		{"POST", "/aiservice/_internal/chat", "Bearer fixture", "", false},
		{"POST", "/aiservice/api/chat", "Bearer fixture", "", false},
		{"POST", "/aiservice/api/client%2fchat", "Bearer fixture", "", false},
		{"POST", "/aiservice/api/client/chat/", "Bearer fixture", "", false},
	} {
		r := httptest.NewRequest(tc.method, "https://workspace.example.test"+tc.path, nil)
		r.Header.Set("Authorization", tc.auth)
		r.Header.Set("Origin", tc.origin)
		if got := remoteAIServiceRequest(r); got != tc.want {
			t.Errorf("%s %s: %v", tc.method, tc.path, got)
		}
	}
}

func TestRemoteAIServiceReachesNativeAuthWithoutExposingConsole(t *testing.T) {
	calls := 0
	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.URL.Path != "/api/client/chat" || r.Header.Get("Authorization") != "Bearer synthetic" {
			t.Error("wrong upstream request")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		w.Write([]byte(`{"error":"remote credential required"}`))
	}))
	defer up.Close()
	u, _ := url.Parse(up.URL)
	port, _ := strconv.Atoi(u.Port())
	rt := &Runtime{Config: Config{Origin: "http://localhost:28990", Port: 28990, Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "aiservice", Kind: "process", Route: "/aiservice", StripPrefix: true, Port: port}}}}, status: map[string]string{}}
	p := &Portal{rt: rt, flows: map[string]loginFlow{}, sessions: map[string]session{}}
	h := p.Handler("synthetic-control")
	for _, path := range []string{"/aiservice/api/client/chat", "/aiservice/api/service/settings", "/aiservice/_internal/chat", "/aiservice/api/chat"} {
		r := httptest.NewRequest("POST", "http://localhost:28990"+path, strings.NewReader(`{}`))
		r.Header.Set("Authorization", "Bearer synthetic")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 401 {
			t.Fatalf("%s: %d", path, w.Code)
		}
		if path == "/aiservice/api/client/chat" && !strings.Contains(w.Body.String(), "remote credential required") {
			t.Fatal("native auth hidden by SSO")
		}
	}
	if calls != 1 {
		t.Fatal("private route reached native app", calls)
	}
}
