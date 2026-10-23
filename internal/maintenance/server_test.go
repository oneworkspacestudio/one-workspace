package maintenance

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMaintenanceCapabilityOriginHostAndBusyGates(t *testing.T) {
	s := &Server{Token: "fixture-token", Origin: "http://localhost:38100", Context: context.Background(), Close: func() {}}
	for _, tc := range []struct {
		path, host, token, origin string
		busy                      bool
		want                      int
	}{{"/api/check", "localhost:38100", "", "http://localhost:38100", false, 403}, {"/api/check", "localhost:38100", "fixture-token", "https://attacker.example", false, 403}, {"/api/check", "attacker.example", "fixture-token", "http://localhost:38100", false, 403}, {"/api/check", "localhost:38100", "fixture-token", "http://localhost:38100", true, 409}, {"/api/close", "localhost:38100", "fixture-token", "http://localhost:38100", false, 202}} {
		s.Busy = tc.busy
		r := httptest.NewRequest("POST", tc.path, strings.NewReader("{}"))
		r.Host = tc.host
		r.Header.Set("X-Workspace-Update", tc.token)
		r.Header.Set("Origin", tc.origin)
		w := httptest.NewRecorder()
		s.Handler().ServeHTTP(w, r)
		if w.Code != tc.want {
			t.Fatal(tc, w.Code)
		}
	}
}

func TestReinstallRequiresCapabilityOriginAndFreshPlan(t *testing.T) {
	s := &Server{Token: "fixture-token", Origin: "http://localhost:38100", Context: context.Background()}
	for _, tc := range []struct {
		token, origin string
		busy          bool
		want          int
	}{{"", "http://localhost:38100", false, 403}, {"fixture-token", "https://attacker.example", false, 403}, {"fixture-token", "http://localhost:38100", true, 409}, {"fixture-token", "http://localhost:38100", false, 409}} {
		s.Busy = tc.busy
		r := httptest.NewRequest("POST", "/api/reinstall", strings.NewReader("{}"))
		r.Host = "localhost:38100"
		r.Header.Set("X-Workspace-Update", tc.token)
		r.Header.Set("Origin", tc.origin)
		w := httptest.NewRecorder()
		s.Handler().ServeHTTP(w, r)
		if w.Code != tc.want {
			t.Fatal(tc, w.Code)
		}
	}
}
