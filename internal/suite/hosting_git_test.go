package suite

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHostingGitBrokerBoundary(t *testing.T) {
	rt := Runtime{Config: Config{Issuer: "https://accounts.example.test"}, Secrets: Secrets{Clients: map[string]Client{"hosting": {Secret: "synthetic-hosting"}}}}
	p := Portal{rt: &rt}
	for _, tc := range []struct{ addr, token string }{{"203.0.113.1:1234", hostingGitCapability("synthetic-hosting")}, {"127.0.0.1:1234", ""}, {"127.0.0.1:1234", "wrong"}} {
		r := httptest.NewRequest("POST", "http://localhost/_integration/hosting/git", strings.NewReader(`{"operation":"resolve","repository":"private"}`))
		r.RemoteAddr = tc.addr
		r.Header.Set("Authorization", "Bearer "+tc.token)
		w := httptest.NewRecorder()
		p.hostingGitRequest(w, r)
		if w.Code != 404 {
			t.Fatal("broker exposed", w.Code)
		}
	}
	r := httptest.NewRequest("POST", "http://localhost/_integration/hosting/git", strings.NewReader(`{"operation":"resolve","issuer":"https://other.test","subject":"person","repository":"private"}`))
	r.RemoteAddr = "127.0.0.1:1234"
	r.Header.Set("Authorization", "Bearer "+hostingGitCapability("synthetic-hosting"))
	w := httptest.NewRecorder()
	p.hostingGitRequest(w, r)
	if w.Code != 400 {
		t.Fatal("issuer accepted", w.Code)
	}
}
func TestHostingGitIdentityIsLinkedNotGuessed(t *testing.T) {
	root := t.TempDir()
	if e := os.MkdirAll(filepath.Join(root, "gitea"), 0700); e != nil {
		t.Fatal(e)
	}
	rt := Runtime{Config: Config{Data: root, Issuer: "https://accounts.example.test"}, Secrets: Secrets{Clients: map[string]Client{"gitea": {ID: "git-fixture", Secret: "synthetic-client"}}}}
	db, e := sql.Open("sqlite", filepath.Join(root, "gitea", "gitea.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	cfg, _ := json.Marshal(map[string]string{"ClientID": "git-fixture", "ClientSecret": "synthetic-client", "OpenIDConnectAutoDiscoveryURL": rt.Config.Issuer + "/.well-known/openid-configuration"})
	for _, q := range []string{`CREATE TABLE login_source(id INTEGER,name TEXT,type INTEGER,is_active INTEGER,cfg TEXT)`, `CREATE TABLE user(id INTEGER,name TEXT,is_active INTEGER,prohibit_login INTEGER,login_source INTEGER,login_type INTEGER,login_name TEXT)`, `CREATE TABLE external_login_user(user_id INTEGER,login_source_id INTEGER,external_id TEXT)`} {
		if _, e = db.Exec(q); e != nil {
			t.Fatal(e)
		}
	}
	if _, e = db.Exec(`INSERT INTO login_source VALUES(1,'myaccounts',6,1,?)`, string(cfg)); e != nil {
		t.Fatal(e)
	}
	if _, e = db.Exec(`INSERT INTO user VALUES(2,'synthetic-alice',1,0,1,6,'subject-alice')`); e != nil {
		t.Fatal(e)
	}
	name, e := rt.hostingGitUsername(context.Background(), "subject-alice")
	if e != nil || name != "synthetic-alice" {
		t.Fatal(name, e)
	}
	if _, e = rt.hostingGitUsername(context.Background(), "other-subject"); e == nil {
		t.Fatal("unlinked subject accepted")
	}
	_, _ = db.Exec(`UPDATE user SET prohibit_login=1`)
	if _, e = rt.hostingGitUsername(context.Background(), "subject-alice"); e == nil {
		t.Fatal("blocked account accepted")
	}
}
func TestHostingGitResolvesOnlyActingUserRepositories(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != "workspace-recovery" || pass != "synthetic-admin" || r.Header.Get("Sudo") != "alice" {
			t.Error("missing server-side impersonation")
			http.Error(w, "denied", 403)
			return
		}
		switch r.URL.Path {
		case "/api/v1/user/repos":
			w.Write([]byte(`[{"id":7,"name":"site","full_name":"alice/site","owner":{"login":"alice"},"private":true},{"id":8,"name":"shared","full_name":"team/shared","owner":{"login":"team"}},{"id":9,"name":"shared","full_name":"other/shared","owner":{"login":"other"}}]`))
		case "/api/v1/repos/alice/site":
			w.Write([]byte(`{"id":7,"full_name":"alice/site"}`))
		default:
			http.Error(w, "not found", 404)
		}
	}))
	defer server.Close()
	g := &myaiGitClient{base: server.URL, password: "synthetic-admin", client: server.Client()}
	repo, status, e := hostingResolveRepo(context.Background(), g, "alice", "site")
	if e != nil || status != 200 || repo.ID != 7 {
		t.Fatal(repo, status, e)
	}
	if _, status, e = hostingResolveRepo(context.Background(), g, "alice", "shared"); e == nil || status != 409 {
		t.Fatal("ambiguous repository chosen")
	}
	if _, status, e = hostingResolveRepo(context.Background(), g, "alice", "bob/private"); e == nil || status != 404 {
		t.Fatal("unauthorized repository accepted")
	}
	for _, in := range []hostingGitCall{{Operation: "delete"}, {Operation: "archive", Ref: "main"}, {Operation: "contents", File: "../secrets"}, {Operation: "contents", File: "/absolute"}, {Operation: "commit", Ref: "main/../../other"}} {
		if _, e := hostingGitEndpoint(in, "alice/site"); e == nil {
			t.Fatal("unsafe operation accepted", in)
		}
	}
	for _, raw := range []string{"site", "alice/site", "https://workspace.example.test/git/alice/site.git/"} {
		if _, e := hostingRepoName(raw, "https://workspace.example.test/git"); e != nil {
			t.Fatal(raw, e)
		}
	}
}
