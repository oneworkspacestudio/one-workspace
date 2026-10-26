package suite

import (
	"archive/zip"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

func TestAIConnectorNativeRouteTrailingSlash(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/drive/api/users" || r.URL.Query().Get("username") != "self" {
			t.Errorf("Wrong native route: %s", r.URL.String())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"username":"fixture"}`))
	}))
	defer server.Close()
	u, _ := url.Parse(server.URL)
	port, _ := strconv.Atoi(u.Port())
	p := &Portal{rt: &Runtime{Config: Config{Origin: "http://workspace.example.test"}}}
	for _, route := range []string{"/drive", "/drive/"} {
		for _, path := range []string{"/api/users?username=self", "/drive/api/users?username=self"} {
			if _, _, e := p.aiNativeRequest(context.Background(), catalog.App{ID: "drive", Route: route, Port: port}, aiAppConnection{}, "", "GET", path, nil, ""); e != nil {
				t.Fatal(route, path, e)
			}
		}
	}
}

func TestAIConnectorTubeNativeIdentityMapping(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "my-tube")
	if e := os.MkdirAll(dir, 0700); e != nil {
		t.Fatal(e)
	}
	db, e := sql.Open("sqlite", filepath.Join(dir, "mytube.db"))
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	if _, e = db.Exec(`CREATE TABLE users(id TEXT PRIMARY KEY,issuer TEXT,subject TEXT,disabled INTEGER); INSERT INTO users VALUES('native-id','https://accounts.example.test','one',0); INSERT INTO users VALUES('disabled-id','https://accounts.example.test','two',1)`); e != nil {
		t.Fatal(e)
	}
	p := &Portal{rt: &Runtime{Config: Config{Data: root}}}
	issuer, subject, e := p.aiTubeIdentity(context.Background(), "native-id")
	if e != nil || issuer != "https://accounts.example.test" || subject != "one" {
		t.Fatal("Native mapping unavailable", issuer, subject, e)
	}
	for _, id := range []string{"disabled-id", "unknown-id", "../outside"} {
		if _, _, e := p.aiTubeIdentity(context.Background(), id); e == nil {
			t.Fatal("Invalid native user accepted", id)
		}
	}
}

func TestAIConnectorGitArchiveUsesImmutableBlobs(t *testing.T) {
	root := t.TempDir()
	work := filepath.Join(root, "work")
	if e := os.MkdirAll(work, 0700); e != nil {
		t.Fatal(e)
	}
	rt := &Runtime{Config: Config{Data: filepath.Join(root, "private")}}
	p := myaiProject{Path: work, GitDir: filepath.Join(root, "private-git")}
	ctx := context.Background()
	git := func(args ...string) string {
		t.Helper()
		out, e := rt.myaiGitRun(ctx, p, "", args...)
		if e != nil {
			t.Fatal(e, args)
		}
		return strings.TrimSpace(out)
	}
	git("init")
	original := []byte("<!doctype html><title>Fixture</title>Saved revision")
	if e := os.WriteFile(filepath.Join(work, "index.html"), original, 0600); e != nil {
		t.Fatal(e)
	}
	git("add", "index.html")
	git("-c", "user.name=Fixture", "-c", "user.email=fixture@example.test", "commit", "-m", "Saved source")
	commit := git("rev-parse", "HEAD")
	_ = os.WriteFile(filepath.Join(work, "index.html"), []byte("Unfinished edits"), 0600)
	_ = os.WriteFile(filepath.Join(work, ".env"), []byte("PRIVATE_FIXTURE=do-not-upload"), 0600)
	archive, e := rt.aiGitArchive(ctx, p, commit)
	if e != nil {
		t.Fatal(e)
	}
	z, e := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	if e != nil || len(z.File) != 1 || z.File[0].Name != "index.html" {
		t.Fatal("Mutable or untracked files included", e)
	}
	f, e := z.File[0].Open()
	if e != nil {
		t.Fatal(e)
	}
	content, e := io.ReadAll(f)
	_ = f.Close()
	if e != nil || !bytes.Equal(content, original) {
		t.Fatal("Wrong source revision", e)
	}
	again, e := rt.aiGitArchive(ctx, p, commit)
	if e != nil || !bytes.Equal(again, archive) {
		t.Fatal("Source hash unstable across retries", e)
	}
	git("add", ".env")
	git("-c", "user.name=Fixture", "-c", "user.email=fixture@example.test", "commit", "-m", "Credential must block handoff")
	if _, e = rt.aiGitArchive(ctx, p, git("rev-parse", "HEAD")); e == nil {
		t.Fatal("Tracked credentials uploaded")
	}
}

func TestAIConnectorUsesOnlyCurrentBrowserNativeIdentity(t *testing.T) {
	var writes atomic.Int32
	var wrongIdentity atomic.Bool
	native := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" || strings.Contains(r.Header.Get("Cookie"), "workspace_session") || r.Header.Get("X-Workspace-Browser-Session") != "" {
			t.Error("another service's credentials forwarded")
		}
		c, e := r.Cookie("tasks_session")
		if e != nil || c.Value != "native-own-session" {
			http.Error(w, "native login required", 401)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/me" {
			subject := "one"
			if wrongIdentity.Load() {
				subject = "two"
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"user": map[string]string{"issuer": "https://accounts.example.test", "subject": subject, "id": "native-id"}, "csrf": "native-csrf"})
			return
		}
		if r.Method == "POST" {
			if r.Header.Get("X-CSRF-Token") != "native-csrf" || r.Header.Get("Origin") != "https://workspace.example.test" {
				http.Error(w, "native CSRF refused", 403)
				return
			}
			writes.Add(1)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"tasks": []any{}, "csrf": "must-not-reach-model", "nested": map[string]string{"access_token": "must-not-reach-model"}})
	}))
	defer native.Close()
	u, _ := url.Parse(native.URL)
	_, port, _ := net.SplitHostPort(u.Host)
	n, _ := strconv.Atoi(port)
	c := Config{Origin: "https://workspace.example.test", Issuer: "https://accounts.example.test", Owner: "owner", Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "tasks", Name: "My Tasks", Route: "/tasks", Port: n, StripPrefix: true, Kind: "process"}}}}
	rt := &Runtime{Config: c, Secrets: Secrets{Clients: map[string]Client{"myai": {Secret: "fixture-secret"}}}}
	p := &Portal{rt: rt, sessions: map[string]session{"current-browser": {Subject: "one", Expires: time.Now().Add(time.Hour), AIApps: map[string]aiAppConnection{"tasks": {Cookies: []*http.Cookie{{Name: "tasks_session", Value: "native-own-session"}}}}}}}
	mux := http.NewServeMux()
	p.myaiAppsRoutes(mux)
	call := func(method, body, root, subject string, remote bool) *httptest.ResponseRecorder {
		r := httptest.NewRequest(method, "http://localhost/_integration/myai/workspace", strings.NewReader(body))
		r.RemoteAddr = "127.0.0.1:4444"
		if remote {
			r.RemoteAddr = "192.0.2.3:4444"
		}
		r.Header.Set("Authorization", "Bearer "+myaiCapability("fixture-secret"))
		r.Header.Set("X-Workspace-Browser-Session", root)
		r.Header.Set("X-MyAI-Issuer", c.Issuer)
		r.Header.Set("X-MyAI-Subject", subject)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		return w
	}
	for _, test := range []struct {
		root, subject string
		remote        bool
	}{{"other-browser", "one", false}, {"current-browser", "two", false}, {"current-browser", "one", true}} {
		if w := call("GET", "", test.root, test.subject, test.remote); w.Code != 404 {
			t.Fatal("invalid authority accepted", w.Code)
		}
	}
	if w := call("GET", "", "current-browser", "one", false); w.Code != 200 || !strings.Contains(w.Body.String(), "My Tasks") {
		t.Fatal("installed catalogue missing", w.Code)
	}
	change := `{"app":"tasks","operation":"create_list","body":{"name":"Own list"}}`
	if w := call("POST", change, "current-browser", "one", false); w.Code != 409 || writes.Load() != 0 {
		t.Fatal("unreviewed mutation sent")
	}
	prepared := strings.TrimSuffix(change, "}") + `,"prepare":true}`
	if w := call("POST", prepared, "current-browser", "one", false); w.Code != 200 || writes.Load() != 0 {
		t.Fatal("review preview mutated", w.Code, w.Body.String())
	}
	approved := strings.TrimSuffix(change, "}") + `,"reviewed":true}`
	if w := call("POST", approved, "current-browser", "one", false); w.Code != 200 || writes.Load() != 1 || strings.Contains(w.Body.String(), "must-not-reach-model") {
		t.Fatal("native action/redaction", w.Code, w.Body.String())
	}
	wrongIdentity.Store(true)
	if w := call("POST", approved, "current-browser", "one", false); w.Code == 200 || writes.Load() != 1 {
		t.Fatal("cross-account native mutation accepted")
	}
	for _, body := range []string{`{"app":"tasks","operation":"invented_delete_everything","reviewed":true}`, `{"app":"tasks","operation":"task","parameters":{"id":"../admin"}}`} {
		if w := call("POST", body, "current-browser", "one", false); w.Code == 200 {
			t.Fatal("unregistered route accepted")
		}
	}
	p.mu.Lock()
	s := p.sessions["current-browser"]
	s.Expires = time.Now().Add(-time.Second)
	p.sessions["current-browser"] = s
	p.mu.Unlock()
	if w := call("GET", "", "current-browser", "one", false); w.Code != 404 {
		t.Fatal("expired portal authority survived")
	}
}

func TestAIConnectorRejectsUnregisteredArgumentsAndExecutableTrust(t *testing.T) {
	op := aiRead("test", "test", "/api/projects/{id}", "q")
	op.Parameters = []string{"id"}
	if _, e := validateAIAppCall(op, aiAppCall{Parameters: map[string]string{"id": "valid"}, Body: json.RawMessage(`{}`)}); e != nil {
		t.Fatal("Empty read body from structured contract rejected", e)
	}
	for _, call := range []aiAppCall{{Parameters: map[string]string{"id": "../private"}}, {Parameters: map[string]string{"id": "valid", "url": "https://example.test"}}, {Parameters: map[string]string{"id": "valid"}, Query: map[string]string{"target": "http://127.0.0.1"}}, {Parameters: map[string]string{"id": "valid"}, Body: json.RawMessage(`{"method":"DELETE"}`)}} {
		if _, e := validateAIAppCall(op, call); e == nil {
			t.Fatal("unsafe app route arguments accepted")
		}
	}
	if _, e := validateAIAppCall(aiChange("create_site", "create", "POST", "/api/projects", nil), aiAppCall{App: "hosting", Body: json.RawMessage(`{"trusted":true}`)}); e == nil {
		t.Fatal("AI enabled executable trust")
	}
	for _, path := range []string{".env", "a/.ssh/id_rsa", "../outside", "a\\b", "/root", "a/.git/config"} {
		if aiArchivePath(path) {
			t.Fatal("unsafe deployment source", path)
		}
	}
	for _, path := range []string{"index.html", "assets/main.js", "README.md", ".env.example"} {
		if !aiArchivePath(path) {
			t.Fatal("regular source rejected", path)
		}
	}
}

func TestAIConnectorNoNativeSessionMinting(t *testing.T) {
	p := &Portal{rt: &Runtime{Config: Config{Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "tasks", Route: "/tasks"}}}}}, sessions: map[string]session{"one-browser": {Subject: "one", Expires: time.Now().Add(time.Hour)}, "two-browser": {Subject: "two", Expires: time.Now().Add(time.Hour)}}}
	r := httptest.NewRequest("GET", "http://localhost/tasks/api/me", nil)
	r.AddCookie(&http.Cookie{Name: "workspace_session", Value: "one-browser"})
	r.AddCookie(&http.Cookie{Name: "tasks_session", Value: "one-native"})
	p.rememberAIAppSession(r)
	if len(p.sessions["one-browser"].AIApps["tasks"].Cookies) != 1 || len(p.sessions["two-browser"].AIApps) != 0 {
		t.Fatal("native session crossed browser boundary")
	}
	_, _, e := p.aiNativeIdentity(context.Background(), p.sessions["two-browser"], catalog.App{ID: "tasks", Name: "My Tasks"})
	if e == nil {
		t.Fatal("missing session minted")
	}
}

func TestAIConnectorDeploymentIntentCanonical(t *testing.T) {
	a := aiAppCall{App: "hosting", Operation: "deploy_project", Parameters: map[string]string{"site": "s", "project": "p"}, Body: json.RawMessage(`{"commit":"a","environment":"preview"}`)}
	b := a
	b.Reviewed = true
	b.Prepare = true
	b.Query = map[string]string{}
	b.Body = json.RawMessage(`{ "environment": "preview", "commit": "a" }`)
	if aiRequestDigest("one", a) != aiRequestDigest("one", b) || aiRequestDigest("one", a) == aiRequestDigest("two", b) {
		t.Fatal("deployment retry authority is unstable or shared")
	}
}

func TestAIConnectorValidatesNativeBodyBeforeReview(t *testing.T) {
	var create aiOperation
	for _, op := range aiOperations("tasks") {
		if op.ID == "create_task" {
			create = op
		}
	}
	if _, e := validateAIAppCall(create, aiAppCall{App: "tasks", Body: json.RawMessage(`{"title":"Fixture","priority":0,"status":"todo"}`)}); e != nil {
		t.Fatal(e)
	}
	for _, body := range []string{`{"title":"Fixture","priority":5}`, `{"title":"Fixture","priority":"high"}`, `{"title":"Fixture","status":"finished"}`, `{"priority":0}`, `{"title":"Fixture","headers":{"Authorization":"forged"}}`} {
		if _, e := validateAIAppCall(create, aiAppCall{App: "tasks", Body: json.RawMessage(body), Prepare: true}); e == nil {
			t.Fatal("Invalid native proposal reached review", body)
		}
	}
	for _, app := range []string{"tasks", "notes", "hosting", "gitea", "drive", "calendar", "contacts", "social", "clips", "music", "maps", "meet", "my3d", "create", "pdf"} {
		for _, op := range aiOperations(app) {
			if op.Review && op.BodySchema == nil {
				t.Fatal("Missing registered request schema", app, op.ID)
			}
		}
	}
}
