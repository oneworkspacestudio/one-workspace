package suite

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

func TestVaultIdentityUsesDedicatedCallbackAndPreservesOptOut(t *testing.T) {
	var added map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data any
		switch r.URL.Path {
		case "/api/get-application":
			if r.URL.Query().Get("id") == "admin/myaccounts" {
				data = map[string]any{"owner": "admin", "name": "myaccounts"}
			}
		case "/api/add-application":
			if e := json.NewDecoder(r.Body).Decode(&added); e != nil {
				t.Error(e)
			}
		}
		json.NewEncoder(w).Encode(map[string]any{"status": "ok", "data": data})
	}))
	defer server.Close()
	u, _ := url.Parse(server.URL)
	port, _ := strconv.Atoi(u.Port())
	root := t.TempDir()
	rt := &Runtime{Config: Config{Data: root, Origin: "https://workspace.example.test", Issuer: server.URL, Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "accounts", Port: port}, {ID: "vault", Port: 18934}}}}}
	if e := rt.prepareVaultIdentity(context.Background()); e != nil {
		t.Fatal(e)
	}
	if added["clientId"] != "workspace-vault" || added["tokenFormat"] != "JWT-Standard" || added["enableSignUp"] != false {
		t.Fatal("unexpected client policy")
	}
	redirects := added["redirectUris"].([]any)
	if len(redirects) != 1 || redirects[0] != "https://localhost:18935/identity/connect/oidc-signin" {
		t.Fatal("callback uses shared Workspace origin", redirects)
	}
	var s map[string]any
	path := filepath.Join(root, "my-vault", "myvault-sso.json")
	if e := platform.LoadJSON(path, &s); e != nil {
		t.Fatal(e)
	}
	if s["issuer"] != server.URL || s["enabled"] != true || s["clientSecret"] == "" {
		t.Fatal("optional identity not saved")
	}
	if e := platform.SaveJSON(path, map[string]bool{"enabled": false}); e != nil {
		t.Fatal(e)
	}
	added = nil
	if e := rt.prepareVaultIdentity(context.Background()); e != nil || added != nil {
		t.Fatal("owner opt-out replaced")
	}
	if _, e := os.Stat(path); e != nil {
		t.Fatal(e)
	}
}
