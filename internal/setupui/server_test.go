package setupui

import (
	"net/http/httptest"
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
