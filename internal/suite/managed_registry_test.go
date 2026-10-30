package suite

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"net/http/httptest"
	"testing"
)

func TestManagedFoundationRegistrations(t *testing.T) {
	for _, app := range catalog.ManagedApps() {
		if !app.Foundation {
			continue
		}
		if !foundationApp(app.ID) || !knownUpdateStorageProfile(app.ID) || (app.Recovery == "app" && len(appUpdatePrivatePaths(app.ID)) == 0) {
			t.Errorf("%s lacks managed runtime/recovery registration", app.ID)
		}
		found := false
		for _, s := range suiteSignIns {
			if s.id == app.ID {
				found = true
			}
		}
		if !found {
			t.Errorf("%s lacks SSO navigation registration", app.ID)
		}
	}
}

func TestEveryManagedFoundationCookieIsClearedOnSignOut(t *testing.T) {
	w := httptest.NewRecorder()
	clearAppCookies(w, httptest.NewRequest("GET", "http://localhost/", nil), false)
	cleared := map[string]bool{}
	for _, c := range w.Result().Cookies() {
		if c.MaxAge < 0 {
			cleared[c.Name+" "+c.Path] = true
		}
	}
	for _, app := range catalog.ManagedApps() {
		if !app.Foundation {
			continue
		}
		for _, suffix := range []string{"_session", "_login"} {
			if !cleared[app.ID+suffix+" /"+app.ID+"/"] {
				t.Errorf("%s missing cookie cleanup", app.ID)
			}
		}
	}
}
