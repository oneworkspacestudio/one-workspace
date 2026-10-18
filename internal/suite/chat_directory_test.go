package suite

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

func TestChatDirectorySnapshot(t *testing.T) {
	var received []chatDirectoryUser
	chat := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-directory-token" {
			t.Error("missing private token")
		}
		var request struct {
			Users []chatDirectoryUser `json:"users"`
		}
		json.NewDecoder(r.Body).Decode(&request)
		received = request.Users
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"members":2}`))
	}))
	defer chat.Close()
	_, port, _ := net.SplitHostPort(chat.Listener.Addr().String())
	p, _ := strconv.Atoi(port)
	accounts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","data":[{"owner":"myaccounts","id":"alice-id","name":"alice","displayName":"Alice","email":"alice@example.test","emailVerified":true},{"owner":"myaccounts","id":"bob-id","name":"bob","email":"unverified@example.test"},{"owner":"myaccounts","id":"disabled-id","name":"disabled","isForbidden":true}]}`))
	}))
	defer accounts.Close()
	rt := &Runtime{Config: Config{Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "chat", Port: p}}}}, chatDirectoryToken: "test-directory-token"}
	count, err := rt.syncChatDirectory(context.Background(), &identityAPI{base: accounts.URL, client: accounts.Client()})
	if err != nil || count != 2 || len(received) != 3 {
		t.Fatalf("sync: %d %v", count, err)
	}
	if received[0].Subject != "alice-id" || received[0].Email != "alice@example.test" || received[1].Email != "" || !received[2].Disabled {
		t.Fatal("identity, verified email or disabled state not retained")
	}
}

func TestChatDirectoryRejectsInvalidIdentity(t *testing.T) {
	accounts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","data":[{"owner":"another-organization","id":"id","name":"alice"}]}`))
	}))
	defer accounts.Close()
	rt := &Runtime{Config: Config{Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "chat", Port: 1}}}}}
	if _, err := rt.syncChatDirectory(context.Background(), &identityAPI{base: accounts.URL, client: accounts.Client()}); err == nil {
		t.Fatal("foreign organization was accepted")
	}
}
