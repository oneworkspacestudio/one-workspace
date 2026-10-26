package suite

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

func TestMailboxSyncBeforeFirstLogin(t *testing.T) {
	users := []map[string]any{{"id": "owner-id", "name": "owner", "owner": "myaccounts", "displayName": "Owner"}, {"id": "disabled", "name": "disabled", "owner": "myaccounts", "isForbidden": true}}
	loggedIn := false
	idp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/login" {
			loggedIn = true
			_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok", "data": true})
			return
		}
		if !loggedIn {
			http.Error(w, "unauthorized", 401)
			return
		}
		if r.URL.Path != "/api/get-users" || r.URL.Query().Get("owner") != "myaccounts" {
			t.Error("unexpected identity request")
			http.Error(w, "bad request", 400)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok", "data": users})
	}))
	defer idp.Close()
	var provisioned []mailboxUser
	mail := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/status" {
			if r.Header.Get("Authorization") != "" {
				t.Error("service role probe forwarded provision credentials")
			}
			w.Write([]byte(`{"mainMailURL":""}`))
			return
		}
		if r.URL.Path != "/api/mailboxes/provision" || r.Header.Get("Authorization") != "Bearer fixture-only-token" || r.Header.Get("X-Forwarded-For") != "" {
			t.Error("wrong provision authorization")
			http.Error(w, "forbidden", 403)
			return
		}
		var in struct {
			Users []mailboxUser `json:"users"`
		}
		if json.NewDecoder(r.Body).Decode(&in) != nil {
			t.Fatal("bad provision payload")
		}
		provisioned = in.Users
		addresses := map[string]string{}
		for _, u := range in.Users {
			addresses[u.Subject] = u.Username + "@example.test"
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"mailboxes": addresses})
	}))
	defer mail.Close()
	u, _ := url.Parse(mail.URL)
	port, _ := strconv.Atoi(u.Port())
	rt := &Runtime{Config: Config{Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "mail", Port: port}}}}, mailProvisionToken: "fixture-only-token"}
	api := &identityAPI{base: idp.URL, client: idp.Client()}
	count, err := rt.syncMailboxes(context.Background(), api)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 || len(provisioned) != 1 || provisioned[0].Subject != "owner-id" {
		t.Fatal("owner mailbox was not provisioned or disabled user was included")
	}
	users = append(users, map[string]any{"id": "new-id", "name": "jane", "owner": "myaccounts", "displayName": "Jane"})
	count, err = rt.syncMailboxes(context.Background(), api)
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 || len(provisioned) != 2 || provisioned[1].Subject != "new-id" {
		t.Fatal("new account required a Mail login before provisioning")
	}
}
