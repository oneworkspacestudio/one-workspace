package apphost

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCaptureIsExplicitCapability(t *testing.T) {
	h := testHost(t)
	for _, allowed := range []bool{false, true} {
		h.Config.MediaCapture = allowed
		w := httptest.NewRecorder()
		h.Handler().ServeHTTP(w, httptest.NewRequest("GET", h.Config.Origin+"/healthz", nil))
		policy := w.Header().Get("Permissions-Policy")
		if strings.Contains(policy, "camera=(self)") != allowed {
			t.Fatal(policy)
		}
		if !strings.Contains(w.Header().Get("Content-Security-Policy"), "media-src 'self' blob:") {
			t.Fatal("recorded preview blocked")
		}
	}
}
