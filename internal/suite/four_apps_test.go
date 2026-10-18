package suite

import (
	"encoding/json"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestScopedMonitorRegistryAndNoSecretLeak(t *testing.T) {
	rt := &Runtime{Config: Config{Data: t.TempDir(), Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "tasks", Name: "My Tasks", Kind: "process", Port: 29008, Health: "/healthz"}}}}, Secrets: Secrets{AdminPassword: "must-never-leak", Clients: map[string]Client{"monitor": {Secret: "scoped-fixture-secret"}}}, status: map[string]string{"tasks": "running"}}
	p := &Portal{rt: rt}
	mux := http.NewServeMux()
	p.monitorRegistry(mux)
	for _, tc := range []struct {
		remote, token string
		code          int
	}{{"127.0.0.1:1234", monitorCapability("scoped-fixture-secret"), 200}, {"203.0.113.1:1234", monitorCapability("scoped-fixture-secret"), 404}, {"127.0.0.1:1234", "", 404}, {"127.0.0.1:1234", "wrong", 404}} {
		r := httptest.NewRequest("GET", "http://localhost/_integration/monitor", nil)
		r.RemoteAddr = tc.remote
		r.Header.Set("Authorization", "Bearer "+tc.token)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != tc.code {
			t.Fatal(w.Code, tc)
		}
		if tc.code == 200 {
			if strings.Contains(w.Body.String(), "secret") || strings.Contains(w.Body.String(), rt.Config.Data) {
				t.Fatal("registry leaked a secret or private path")
			}
			var info map[string]any
			if e := json.Unmarshal(w.Body.Bytes(), &info); e != nil {
				t.Fatal(e)
			}
		}
	}
}
func TestFoundationSignOutRevokesEachCapturedBrowserSession(t *testing.T) {
	count := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/auth/logout" || r.Header.Get("Origin") != "http://localhost:29000" {
			t.Error("wrong logout contract")
		}
		if _, e := r.Cookie("workspace_session"); e == nil {
			t.Error("forwarded portal cookie")
		}
		count++
		w.WriteHeader(200)
	}))
	defer server.Close()
	_, port, _ := net.SplitHostPort(strings.TrimPrefix(server.URL, "http://"))
	n, _ := strconv.Atoi(port)
	rt := &Runtime{Config: Config{Origin: "http://localhost:29000", Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "tasks", Port: n}, {ID: "pdf", Port: n}}}}}
	p := &Portal{rt: rt, sessions: map[string]session{"fixture-portal": {Subject: "owner", Expires: time.Now().Add(time.Hour), AppSessions: map[string]string{"tasks": strings.Repeat("a", 96), "pdf": strings.Repeat("b", 96)}}}}
	r := httptest.NewRequest("POST", "http://localhost:29000/auth/logout", nil)
	r.AddCookie(&http.Cookie{Name: "workspace_session", Value: "fixture-portal"})
	p.closeFoundationSessions(r)
	if count != 2 {
		t.Fatal("sessions not revoked", count)
	}
	w := httptest.NewRecorder()
	clearAppCookies(w, r, false)
	for _, id := range []string{"tasks", "pdf", "monitor", "search"} {
		found := false
		for _, cookie := range w.Result().Cookies() {
			if cookie.Name == id+"_session" && cookie.Path == "/"+id+"/" && cookie.MaxAge < 0 {
				found = true
			}
		}
		if !found {
			t.Fatal("missing scoped cookie cleanup", id)
		}
	}
}
func TestRemoveAppPreservesFilesIdentityAndFutureOptOut(t *testing.T) {
	c, s := updateFixture(t)
	pkg := updatePackage("MyPDF.exe", "pdf-original")
	c.Catalog.Packages["pdf"] = pkg
	c.Catalog.Apps = append(c.Catalog.Apps, catalog.App{ID: "pdf", Name: "My PDF", Kind: "process", Package: "pdf", Executable: "MyPDF.exe", Port: 29119, Route: "/pdf", Version: "0.1.0-preview.1"})
	s.Clients["pdf"] = Client{ID: "existing-pdf", Secret: "preserve-client"}
	os.WriteFile(filepath.Join(c.Programs, "MyPDF.exe"), []byte("pdf-original"), 0600)
	os.MkdirAll(filepath.Join(c.Data, "pdf"), 0700)
	dataFile := filepath.Join(c.Data, "pdf", "state.sqlite")
	os.WriteFile(dataFile, []byte("preserve-private-data"), 0600)
	if e := Save(c, s); e != nil {
		t.Fatal(e)
	}
	lock, e := platform.Lock(filepath.Join(c.Data, "runtime.lock"))
	if e != nil {
		t.Fatal(e)
	}
	if RemoveApp(c.Data, "pdf") == nil {
		t.Fatal("removed while Workspace runs")
	}
	lock.Close()
	if e = RemoveApp(c.Data, "pdf"); e != nil {
		t.Fatal(e)
	}
	after, secrets, e := Load(c.Data)
	if e != nil {
		t.Fatal(e)
	}
	if _, found := after.Catalog.Find("pdf"); found || !after.Catalog.Selection.Excludes("pdf") {
		t.Fatal("removal did not persist opt-out")
	}
	if secrets.Clients["pdf"] != s.Clients["pdf"] {
		t.Fatal("identity was deleted")
	}
	if b, e := os.ReadFile(dataFile); e != nil || string(b) != "preserve-private-data" {
		t.Fatal("data deleted")
	}
	if _, e := os.Stat(filepath.Join(c.Programs, "MyPDF.exe")); !os.IsNotExist(e) {
		t.Fatal("executable retained")
	}
}
