package suite

import (
	"context"
	"encoding/json"
	"github.com/fahimc/one-workspace/internal/catalog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSharedClientHasOnlyExactSuiteCallbacks(t *testing.T) {
	var got map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/get-application" {
			w.Write([]byte(`{"status":"ok","data":{"name":"myaccounts","owner":"admin"}}`))
			return
		}
		if r.URL.Path == "/api/add-application" {
			json.NewDecoder(r.Body).Decode(&got)
			w.Write([]byte(`{"status":"ok"}`))
			return
		}
		t.Errorf("unexpected path %s", r.URL.Path)
	}))
	defer server.Close()
	api := newIdentity(Config{Issuer: server.URL})
	api.base = server.URL
	client, e := registerSharedClient(context.Background(), api, Config{Origin: "https://workspace.example.test"})
	if e != nil {
		t.Fatal(e)
	}
	if client.ID != "workspace-shared" {
		t.Fatal(client.ID)
	}
	redirects, ok := got["redirectUris"].([]any)
	if !ok || len(redirects) != 3 {
		t.Fatalf("callbacks: %v", got)
	}
	expected := map[string]bool{"https://workspace.example.test/calendar/auth/callback": true, "https://workspace.example.test/mail/auth/callback": true, "https://workspace.example.test/contacts/auth/callback": true}
	for _, uri := range redirects {
		if !expected[uri.(string)] {
			t.Fatal("unexpected callback", uri)
		}
	}
}
func TestSharedPackagesAreExplicitAndAppPortsDoNotDependOnMail(t *testing.T) {
	c := Config{Catalog: catalog.Catalog{Packages: map[string]catalog.Package{}}}
	if SharedEnabled(c) {
		t.Fatal("legacy catalog opted in")
	}
	c.Catalog.Packages["stalwart-shared"] = catalog.Package{File: "services/stalwart.exe"}
	if SharedEnabled(c) {
		t.Fatal("partial catalog opted in")
	}
	c.Catalog.Packages["services-shared"] = catalog.Package{File: "services/WorkspaceServices.exe"}
	if !SharedEnabled(c) {
		t.Fatal("shared packages not discovered")
	}
	if !sharedConsumer("calendar") || !sharedConsumer("contacts") || sharedConsumer("drive") {
		t.Fatal("incorrect shared dependency")
	}
}
