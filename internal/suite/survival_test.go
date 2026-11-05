package suite

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestSurvivalPublicBoundary(t *testing.T) {
	for _, v := range []struct {
		method, path string
		want         bool
	}{{"GET", "/survival/", true}, {"HEAD", "/survival/sw.js", true}, {"GET", "/survival/assets/index-abc.js", true}, {"GET", "/survival/portable.html", true}, {"GET", "/survival/api/vault", false}, {"PUT", "/survival/", false}, {"GET", "/survival/assets/../api/vault", false}, {"GET", "/survival/sw.js?secret=x", false}, {"GET", "/survival/auth/login", false}, {"GET", "/survival/data/host.db", false}} {
		if got := publicSurvivalRequest(httptest.NewRequest(v.method, v.path, nil)); got != v.want {
			t.Fatalf("%s %s: %v", v.method, v.path, got)
		}
	}
}
func TestSurvivalOwnership(t *testing.T) {
	cfg := NewConfig(Request{Port: 25100}, catalog.Catalog{Apps: []catalog.App{{ID: "survival", Port: 9041}}})
	if cfg.Catalog.Apps[0].Port != 25151 {
		t.Fatal("fresh setup disagrees with generic upgrade port rebasing")
	}
	if !foundationApp("survival") || !knownUpdateStorageProfile("survival") || len(appUpdatePrivatePaths("survival")) != 1 {
		t.Fatal("missing identity/recovery registration")
	}
	c := Config{Data: t.TempDir()}
	if appWorkingDir(c, "survival") != filepath.Join(c.Data, "survival") {
		t.Fatal("wrong private root")
	}
	a := catalog.App{ID: "survival", Package: "survival"}
	if !appOwnsUpdatePackage(a, "survival-contract") || appOwnsUpdatePackage(a, "books") {
		t.Fatal("bad package ownership")
	}
}
