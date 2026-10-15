package branding

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMyTubeIconIsEmbeddedAndServed(t *testing.T) {
	p := Icon("my-tube")
	if p != "/branding/icons/my_tube.svg" {
		t.Fatal(p)
	}
	w := httptest.NewRecorder()
	Handler().ServeHTTP(w, httptest.NewRequest("GET", p, nil))
	if w.Code != 200 || !strings.Contains(w.Body.String(), "my-tube-title") {
		t.Fatal("My Tube artwork unavailable")
	}
}
