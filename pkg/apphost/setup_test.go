package apphost

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFirstSetupCreatesSessionAndRejectsUntrustedOwner(t *testing.T) {
	h := testHost(t)
	r := httptest.NewRequest("POST", h.Config.Origin+"/api/setup", strings.NewReader(`{"email":"first@example.test","password":"fixture-password"}`))
	r.RemoteAddr = "127.0.0.1:1234"
	r.Header.Set("Origin", h.Config.Origin)
	r.Header.Set("X-Setup-Token", h.setup)
	w := httptest.NewRecorder()
	h.Handler().ServeHTTP(w, r)
	if w.Code != 201 || len(w.Result().Cookies()) == 0 {
		t.Fatal(w.Code, w.Body)
	}
	me := call(h, "GET", "/api/me", nil, w.Result().Cookies()[0], "")
	var s Session
	json.Unmarshal(me.Body.Bytes(), &s)
	if me.Code != 200 || s.User.Role != "owner" {
		t.Fatal("first setup required another login")
	}
	root := filepath.Join(t.TempDir(), "bad-owner")
	os.MkdirAll(root, 0700)
	os.WriteFile(filepath.Join(root, "app.json"), []byte(`{"schema":1,"id":"tasks"}`), 0600)
	os.WriteFile(filepath.Join(root, ".owner-sid"), []byte("S-1-1-0"), 0600)
	if PrepareData(root, "tasks") == nil {
		t.Fatal("Everyone owner marker accepted")
	}
}
