package suite

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestWebSearchLifecycle(t *testing.T) {
	c := NewConfig(Request{Port: 25100}, catalog.Catalog{Apps: []catalog.App{{ID: "websearch", Port: 9043}}})
	a, _ := c.Catalog.Find("websearch")
	if a.Port != 25153 || !foundationApp("websearch") || !knownUpdateStorageProfile("websearch") {
		t.Fatal("missing runtime contract")
	}
	c.Data = t.TempDir()
	if appWorkingDir(c, "websearch") != filepath.Join(c.Data, "websearch") || len(appUpdatePrivatePaths("websearch")) != 1 {
		t.Fatal("missing data ownership")
	}
	if !appOwnsUpdatePackage(catalog.App{Package: "websearch"}, "websearch-runtime-browser") {
		t.Fatal("unowned browser")
	}
}
func TestWebSearchServerIngressIsNarrow(t *testing.T) {
	id := strings.Repeat("a", 48)
	for _, tc := range []struct {
		method, path, origin string
		want                 bool
	}{{"POST", "/websearch/api/client/search", "", true}, {"GET", "/websearch/api/client/jobs/" + id, "", true}, {"DELETE", "/websearch/api/client/jobs/" + id, "", true}, {"POST", "/websearch/api/settings", "", false}, {"GET", "/websearch/api/browser/screenshot", "", false}, {"POST", "/websearch/api/client/search?token=x", "", false}, {"POST", "/websearch/api/client/search", "https://evil.test", false}, {"POST", "/websearch/api/client%2fsearch", "", false}, {"GET", "/websearch/api/client/jobs/../settings", "", false}, {"POST", "/websearch/api/client/search/", "", false}} {
		r := httptest.NewRequest(tc.method, "https://workspace.example.test"+tc.path, nil)
		r.Header.Set("Authorization", "Bearer synthetic")
		r.Header.Set("Origin", tc.origin)
		if got := remoteWebSearchRequest(r); got != tc.want {
			t.Errorf("%s %s = %v", tc.method, tc.path, got)
		}
	}
}
