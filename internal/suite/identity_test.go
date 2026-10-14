package suite

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisteredClientSupportsSharedSignin(t *testing.T) {
	var registered map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/get-application" {
			json.NewEncoder(w).Encode(map[string]any{"status": "ok", "data": map[string]any{"owner": "admin", "organization": "myaccounts", "enableSigninSession": false}})
			return
		}
		if r.URL.Path != "/api/add-application" {
			t.Error("unexpected identity operation", r.URL.Path)
		}
		json.NewDecoder(r.Body).Decode(&registered)
		json.NewEncoder(w).Encode(map[string]any{"status": "ok", "data": "affected"})
	}))
	defer server.Close()
	client, e := registerClient(context.Background(), &identityAPI{base: server.URL, client: server.Client()}, "test", "Test app", "https://workspace.example.com/test/callback")
	if e != nil {
		t.Fatal(e)
	}
	if registered["enableSigninSession"] != true || registered["enableAutoSignin"] != true {
		t.Fatal("native login must establish and reuse an Accounts session")
	}
	if registered["organization"] != "myaccounts" {
		t.Fatal("client escaped ordinary realm")
	}
	redirects := registered["redirectUris"].([]any)
	if len(redirects) != 1 || redirects[0] != "https://workspace.example.com/test/callback" || client.Secret == "" {
		t.Fatal("client registration was incomplete")
	}
}
