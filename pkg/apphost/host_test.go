package apphost

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func testHost(t *testing.T) *Host {
	t.Helper()
	h, e := New(Config{ID: "tasks", Name: "My Tasks", Version: "test", Data: filepath.Join(t.TempDir(), "data"), Listen: "127.0.0.1:28908", Origin: "http://localhost:28908"}, nil)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { h.Close() })
	return h
}
func call(h *Host, method, path string, body any, cookie *http.Cookie, csrf string) *httptest.ResponseRecorder {
	var data io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		data = bytes.NewReader(b)
	}
	r := httptest.NewRequest(method, h.Config.Origin+path, data)
	r.RemoteAddr = "127.0.0.1:12345"
	r.Header.Set("Origin", h.Config.Origin)
	r.Header.Set("X-CSRF-Token", csrf)
	if cookie != nil {
		r.AddCookie(cookie)
	}
	w := httptest.NewRecorder()
	h.Handler().ServeHTTP(w, r)
	return w
}
func TestLocalIdentityCSRFIsolationAndManagedRestrictions(t *testing.T) {
	h := testHost(t)
	if e := h.Bootstrap("owner@example.test", "fixture-password"); e != nil {
		t.Fatal(e)
	}
	if e := h.Bootstrap("intruder@example.test", "fixture-password"); e == nil {
		t.Fatal("bootstrap repeated")
	}
	login := call(h, "POST", "/api/login", map[string]string{"email": "owner@example.test", "password": "fixture-password"}, nil, "")
	if login.Code != 200 {
		t.Fatal(login.Code, login.Body.String())
	}
	cookie := login.Result().Cookies()[0]
	me := call(h, "GET", "/api/me", nil, cookie, "")
	var session Session
	if json.Unmarshal(me.Body.Bytes(), &session) != nil || session.User.Role != "owner" || session.CSRF == "" {
		t.Fatal(me.Body.String())
	}
	if me.Header().Get("X-Workspace-Account") != session.User.ID {
		t.Fatal("missing authenticated response account")
	}
	r := httptest.NewRequest("POST", h.Config.Origin+"/api/local-users", strings.NewReader(`{"email":"wrong@example.test","password":"fixture-password"}`))
	r.AddCookie(cookie)
	r.Header.Set("Origin", h.Config.Origin)
	r.Header.Set("X-CSRF-Token", session.CSRF)
	r.Header.Set("X-Workspace-Expected-Account", "another-account")
	mismatch := httptest.NewRecorder()
	h.Handler().ServeHTTP(mismatch, r)
	if mismatch.Code != 409 || mismatch.Header().Get("X-Workspace-Account") != session.User.ID {
		t.Fatal("stale account request executed", mismatch.Code)
	}
	var wrong int
	h.DB.QueryRow("SELECT COUNT(*) FROM host_users WHERE email='wrong@example.test'").Scan(&wrong)
	if wrong != 0 {
		t.Fatal("account mismatch wrote a user")
	}
	input := map[string]string{"email": "member@example.test", "password": "fixture-password"}
	if w := call(h, "POST", "/api/local-users", input, cookie, ""); w.Code != 403 {
		t.Fatal("CSRF accepted", w.Code)
	}
	if w := call(h, "POST", "/api/local-users", input, cookie, session.CSRF); w.Code != 201 {
		t.Fatal(w.Code, w.Body.String())
	}
	member := call(h, "POST", "/api/login", input, nil, "")
	mCookie := member.Result().Cookies()[0]
	var mSession Session
	json.Unmarshal(call(h, "GET", "/api/me", nil, mCookie, "").Body.Bytes(), &mSession)
	if mSession.User.ID == session.User.ID {
		t.Fatal("users merged")
	}
	if w := call(h, "POST", "/api/local-users", map[string]string{"email": "third@example.test", "password": "fixture-password"}, mCookie, mSession.CSRF); w.Code != 403 {
		t.Fatal("member changed accounts", w.Code)
	}
	if w := call(h, "POST", "/auth/logout", map[string]string{}, cookie, ""); w.Code != 200 {
		t.Fatal(w.Code)
	}
	if w := call(h, "GET", "/api/me", nil, cookie, ""); w.Code != 401 {
		t.Fatal("session survived logout")
	}
	h.Config.Issuer = "https://accounts.example.test"
	if w := call(h, "POST", "/api/login", input, nil, ""); w.Code != 403 {
		t.Fatal("managed auth bypass")
	}
}
func TestHostAndOriginRebindingAndExpiry(t *testing.T) {
	h := testHost(t)
	r := httptest.NewRequest("GET", "http://evil.test/healthz", nil)
	w := httptest.NewRecorder()
	h.Handler().ServeHTTP(w, r)
	if w.Code != 403 {
		t.Fatal("Host rebinding accepted")
	}
	r = httptest.NewRequest("POST", h.Config.Origin+"/api/login", strings.NewReader(`{"email":"x","password":"x"}`))
	r.RemoteAddr = "127.0.0.1:12"
	r.Header.Set("Origin", "https://evil.test")
	w = httptest.NewRecorder()
	h.Handler().ServeHTTP(w, r)
	if w.Code != 403 {
		t.Fatal("cross-origin login accepted")
	}
	h.Bootstrap("one@example.test", "fixture-password")
	u, _ := h.addLocal("two@example.test", "fixture-password", "member")
	w = httptest.NewRecorder()
	h.createSession(w, u, time.Now().Add(-time.Second))
	if call(h, "GET", "/api/me", nil, w.Result().Cookies()[0], "").Code != 401 {
		t.Fatal("expired cookie accepted")
	}
	if Identity("issuer-a", "same") == Identity("issuer-b", "same") {
		t.Fatal("issuer not in identity key")
	}
}
func TestBackupFreshRootIntegrityAndSessionRevocation(t *testing.T) {
	h := testHost(t)
	h.Bootstrap("owner@example.test", "fixture-password")
	if _, e := h.DB.Exec("CREATE TABLE documents(id TEXT PRIMARY KEY,title TEXT); INSERT INTO documents VALUES('a','private fixture')"); e != nil {
		t.Fatal(e)
	}
	blob, n, e := h.PutBlob(strings.NewReader("ordinary content"), 1024)
	if e != nil || n != 16 {
		t.Fatal(n, e)
	}
	login := call(h, "POST", "/api/login", map[string]string{"email": "owner@example.test", "password": "fixture-password"}, nil, "")
	cookie := login.Result().Cookies()[0]
	var encrypted bytes.Buffer
	if e = h.Snapshot(&encrypted, "backup-fixture-password"); e != nil {
		t.Fatal(e)
	}
	if bytes.Contains(encrypted.Bytes(), []byte("private fixture")) {
		t.Fatal("plaintext backup")
	}
	archive := filepath.Join(t.TempDir(), "snapshot.owbackup")
	os.WriteFile(archive, encrypted.Bytes(), 0600)
	destination := filepath.Join(t.TempDir(), "restored")
	if e = Restore(destination, "tasks", archive, "wrong-password"); e == nil {
		t.Fatal("wrong password restored")
	}
	if _, e = os.Stat(destination); !os.IsNotExist(e) {
		t.Fatal("failed restore changed destination")
	}
	if e = Restore(destination, "tasks", archive, "backup-fixture-password"); e != nil {
		t.Fatal(e)
	}
	conf := h.Config
	conf.Data = destination
	restored, e := New(conf, nil)
	if e != nil {
		t.Fatal(e)
	}
	defer restored.Close()
	var title string
	restored.DB.QueryRow("SELECT title FROM documents WHERE id='a'").Scan(&title)
	if title != "private fixture" {
		t.Fatal(title)
	}
	path, _ := restored.BlobPath(blob)
	b, e := os.ReadFile(path)
	if e != nil || string(b) != "ordinary content" {
		t.Fatal(e, string(b))
	}
	if call(restored, "GET", "/api/me", nil, cookie, "").Code != 401 {
		t.Fatal("restored old session")
	}
	if e = Restore(destination, "tasks", archive, "backup-fixture-password"); e == nil {
		t.Fatal("restore overwrote existing data")
	}
	tampered := append([]byte(nil), encrypted.Bytes()...)
	tampered[len(tampered)/2] ^= 1
	os.WriteFile(archive, tampered, 0600)
	if e = Restore(filepath.Join(t.TempDir(), "bad"), "tasks", archive, "backup-fixture-password"); e == nil {
		t.Fatal("tampered archive accepted")
	}
}
func TestDataPathAndMigrationGuards(t *testing.T) {
	h := testHost(t)
	if e := h.Migrate(1, "CREATE TABLE app_items(id TEXT)"); e != nil {
		t.Fatal(e)
	}
	if e := h.Migrate(1, "invalid SQL"); e != nil {
		t.Fatal("migration not idempotent", e)
	}
	if e := h.Migrate(3, "SELECT 1"); e == nil {
		t.Fatal("skipped migration")
	}
	if _, _, e := h.PutBlob(strings.NewReader("oversized"), 2); e == nil {
		t.Fatal("size limit bypass")
	}
	if _, e := h.BlobPath("../private"); e == nil {
		t.Fatal("traversal accepted")
	}
	other := h.Config
	other.ID = "pdf"
	if _, e := New(other, nil); e == nil {
		t.Fatal("adopted another app's data")
	}
	if _, e := New(h.Config, nil); e == nil {
		t.Fatal("duplicate server admitted")
	}
}
