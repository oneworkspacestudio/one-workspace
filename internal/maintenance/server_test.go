package maintenance

import (
	"context"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpgradedControllerRejectsMutationsUntilReopened(t *testing.T) {
	data := t.TempDir()
	if err := os.WriteFile(filepath.Join(data, "workspace.json"), []byte(`{"version":"future-version","catalog":{"schema":999,"packages":{"manager":{"sha256":"new-binary"}}}}`), 0600); err != nil {
		t.Fatal(err)
	}
	s := &Server{Data: data, ControllerSHA256: "old-binary", Token: "fixture-token", Origin: "http://localhost:38100", Context: context.Background(), Close: func() {}}
	request := httptest.NewRequest("GET", "/api/state", nil)
	request.Host = "localhost:38100"
	request.Header.Set("X-Workspace-Update", s.Token)
	result := httptest.NewRecorder()
	s.Handler().ServeHTTP(result, request)
	if result.Code != 200 || !strings.Contains(result.Body.String(), `"controllerChanged":true`) || !strings.Contains(result.Body.String(), `"version":"future-version"`) {
		t.Fatalf("future catalog handoff: %d %s", result.Code, result.Body.String())
	}
	for _, path := range []string{"check", "apply", "reinstall", "rollback", "recover"} {
		request := httptest.NewRequest("POST", "/api/"+path, strings.NewReader("{}"))
		request.Host = "localhost:38100"
		request.Header.Set("Origin", s.Origin)
		request.Header.Set("X-Workspace-Update", s.Token)
		result := httptest.NewRecorder()
		s.Handler().ServeHTTP(result, request)
		if result.Code != 409 || !strings.Contains(result.Body.String(), "Workspace was upgraded") {
			t.Fatalf("%s: %d %s", path, result.Code, result.Body.String())
		}
	}
	// Reopening must also respect an operation still running in the old process.
	s.Busy = true
	request = httptest.NewRequest("POST", "/api/reopen", strings.NewReader("{}"))
	request.Host = "localhost:38100"
	request.Header.Set("Origin", s.Origin)
	request.Header.Set("X-Workspace-Update", s.Token)
	result = httptest.NewRecorder()
	s.Handler().ServeHTTP(result, request)
	if result.Code != 409 {
		t.Fatalf("busy reopen: %d", result.Code)
	}
}

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
