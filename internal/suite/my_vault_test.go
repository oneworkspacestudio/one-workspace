package suite

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fahimc/one-workspace/internal/catalog"
	"github.com/fahimc/one-workspace/internal/platform"
)

func TestVaultUsesIndependentOriginWithoutPortalSession(t *testing.T) {
	root := t.TempDir()
	if e := os.MkdirAll(filepath.Join(root, "my-vault"), 0700); e != nil {
		t.Fatal(e)
	}
	app := catalog.App{ID: "vault", Kind: "process", Route: "/vault", Port: 19014}
	p := &Portal{rt: &Runtime{Config: Config{Data: root, Origin: "https://workspace.example.test", Catalog: catalog.Catalog{Apps: []catalog.App{app}}}}}
	mux := http.NewServeMux()
	p.vaultRoutes(mux)
	r := httptest.NewRequest("GET", "https://workspace.example.test/vault/", nil)
	if p.isAppRequest(r) {
		t.Fatal("vault launcher gated behind MyAccounts")
	}
	for _, origin := range []string{"", "http://insecure.example.test", "https://workspace.example.test", "https://WORKSPACE.example.test:443", "https://localhost:19015"} {
		platform.SaveJSON(filepath.Join(root, "my-vault", "myvault.json"), map[string]string{"origin": origin})
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != 200 || !strings.Contains(w.Body.String(), "https://localhost:19015/") {
			t.Fatal("unsafe/missing public origin did not show independent local setup")
		}
		if strings.Contains(w.Body.String(), "/_workspace/navigation.js") {
			t.Fatal("Workspace script inserted into vault launch page")
		}
	}
	platform.SaveJSON(filepath.Join(root, "my-vault", "myvault.json"), map[string]string{"origin": "https://vault.example.test"})
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	if w.Code != 303 || w.Header().Get("Location") != "https://vault.example.test/" {
		t.Fatal("dedicated vault origin not linked directly")
	}
}

func TestVaultAssignedIndependentPortWithoutOIDC(t *testing.T) {
	c := NewConfig(Request{Port: 19000}, catalog.Catalog{Apps: []catalog.App{{ID: "vault", Kind: "process", Route: "/vault"}}})
	vault, _ := c.Catalog.Find("vault")
	if vault.Port != 19014 || vault.OIDC || len(vault.Depends) != 0 {
		t.Fatal("vault lost independent launch contract")
	}
}
