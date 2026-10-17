package suite

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fahimc/one-workspace/internal/catalog"
)

func TestMusicWorkerKeepsOriginalCachingButUsesLiveNavigation(t *testing.T) {
	original := `self.addEventListener('fetch', originalWorkboxHandler);`
	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Fixture-Expect-Identity") == "1" && r.Header.Get("Accept-Encoding") != "identity" {
			t.Error("worker was requested compressed")
		}
		w.Header().Set("Content-Type", "application/javascript")
		w.Header().Set("ETag", "old-worker")
		io.WriteString(w, original)
	}))
	defer up.Close()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "https://workspace.example.test/music/sw.js", nil)
	req.Header.Set("X-Fixture-Expect-Identity", "1")
	proxy(up.URL, "https://workspace.example.test", "/music", true, false).ServeHTTP(w, req)
	if w.Code != 200 || !strings.HasPrefix(w.Body.String(), string(musicWorkerNavigation)) || !strings.HasSuffix(w.Body.String(), original) {
		t.Fatal("worker navigation adapter/original caching missing")
	}
	if w.Header().Get("Service-Worker-Allowed") != "/music/" || w.Header().Get("Cache-Control") != "no-cache" || w.Header().Get("ETag") != "" {
		t.Fatal("worker scope/cache headers incorrect")
	}
	for _, tc := range []struct{ prefix, path string }{{"/notes", "/notes/sw.js"}, {"/music", "/music/assets/app.js"}} {
		w := httptest.NewRecorder()
		proxy(up.URL, "https://workspace.example.test", tc.prefix, true, false).ServeHTTP(w, httptest.NewRequest("GET", "https://workspace.example.test"+tc.path, nil))
		if w.Body.String() != original {
			t.Error("modified an unrelated worker/asset", tc.path)
		}
	}
}

func TestMusicOpeningRetiresOnlyMusicWorkerAndPreservesDestination(t *testing.T) {
	p := &Portal{rt: &Runtime{Config: Config{Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "music", Route: "/music"}}}}}, sessions: map[string]session{"fixture": {Expires: time.Now().Add(time.Hour)}}}
	r := httptest.NewRequest("GET", "https://workspace.example.test/_workspace/open/music?returnTo=%2Fmusic%2F%3Ftab%3Dlibrary", nil)
	r.SetPathValue("id", "music")
	r.AddCookie(&http.Cookie{Name: "workspace_session", Value: "fixture"})
	w := httptest.NewRecorder()
	p.appOpen(w, r)
	if w.Code != 200 || !strings.Contains(w.Body.String(), "registration.scope===scope") || !strings.Contains(w.Body.String(), `new URL('/music/',location.origin)`) || !strings.Contains(w.Body.String(), `location.replace("/music/?tab=library")`) {
		t.Fatal("music handoff lost scope or destination", w.Body.String())
	}
	if strings.Contains(w.Body.String(), "caches.delete") || strings.Contains(w.Body.String(), "localStorage.clear") || !strings.Contains(w.Header().Get("Content-Security-Policy"), "nonce-") {
		t.Fatal("music handoff changed app data or omitted script protection")
	}
}
