package branding

import (
	"crypto/sha256"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMyTubeIconIsEmbeddedAndServed(t *testing.T) {
	p := Icon("my-tube")
	if !strings.HasPrefix(p, "/branding/icons/my_tube-") || !strings.HasSuffix(p, ".svg") {
		t.Fatal(p)
	}
	w := httptest.NewRecorder()
	Handler().ServeHTTP(w, httptest.NewRequest("GET", p, nil))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "my-tube-title") {
		t.Fatal("My Tube artwork unavailable")
	}
}

func TestEveryIconUsesItsContentAndRevalidatesLegacyURLs(t *testing.T) {
	files, err := icons.ReadDir("icons")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		body, _ := icons.ReadFile("icons/" + file.Name())
		sum := sha256.Sum256(body)
		legacy := "/branding/icons/" + file.Name()
		versioned := fmt.Sprintf("/branding/icons/%s-%x.svg", strings.TrimSuffix(file.Name(), ".svg"), sum[:8])
		page := VersionHTML(`<link rel="icon" href="` + legacy + `"><img src="` + legacy + `">`)
		if strings.Count(page, versioned) != 2 {
			t.Fatal("stale page icon", file.Name(), page)
		}
		for _, path := range []string{legacy, versioned} {
			w := httptest.NewRecorder()
			Handler().ServeHTTP(w, httptest.NewRequest("GET", path, nil))
			if w.Code != 200 || w.Body.String() != string(body) {
				t.Fatal("incorrect artwork", path)
			}
			cache := w.Header().Get("Cache-Control")
			if path == legacy && cache != "no-cache" || path == versioned && !strings.Contains(cache, "immutable") {
				t.Fatal("wrong cache policy", path, cache)
			}
			r := httptest.NewRequest("GET", path, nil)
			r.Header.Set("If-None-Match", w.Header().Get("ETag"))
			cached := httptest.NewRecorder()
			Handler().ServeHTTP(cached, r)
			if cached.Code != 304 {
				t.Fatal("cache validator ignored", path)
			}
		}
	}
	if Icon("unknown-extension") != iconURLs["one_workspace"] {
		t.Fatal("Workspace fallback is stale")
	}
}

func TestUnknownRevisionCannotServeAnImmutableIcon(t *testing.T) {
	w := httptest.NewRecorder()
	Handler().ServeHTTP(w, httptest.NewRequest("GET", "/branding/icons/one_workspace-0000000000000000.svg", nil))
	if w.Code != 404 || strings.Contains(w.Header().Get("Cache-Control"), "immutable") {
		t.Fatal("unverified icon revision served")
	}
}
