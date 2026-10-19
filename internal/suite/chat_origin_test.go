package suite

import (
	"encoding/json"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChatCommandUsesExactDedicatedCallbackAndManagerCapability(t *testing.T) {
	c := chatTestConfig(t)
	c.Programs = t.TempDir()
	c.Port = 8990
	app, _ := c.Catalog.Find("chat")
	app.OIDC = true
	app.Executable = "Chat.exe"
	app.Args = []string{"--origin", "${ORIGIN}/chat"}
	rt := &Runtime{Config: c, chatDirectoryToken: strings.Repeat("x", 40), Secrets: Secrets{Clients: map[string]Client{"chat": {ID: "synthetic-chat", Secret: "synthetic-secret"}}}}
	command := rt.command(app)
	if command.Args[2] != chatOrigin(c) {
		t.Fatal("command retained shared origin")
	}
	env := map[string]string{}
	for _, pair := range command.Env {
		key, value, ok := strings.Cut(pair, "=")
		if ok {
			env[key] = value
		}
	}
	if env["WORKSPACE_OIDC_REDIRECT_URI"] != chatOrigin(c)+"/auth/callback" || env["WORKSPACE_CHAT_MANAGER_CAPABILITIES"] != "isolated-chat-origin-v1" {
		t.Fatal("manager integration does not match exact origin contract")
	}
}

func chatTestConfig(t *testing.T) Config {
	return Config{Data: t.TempDir(), Origin: "http://localhost:8990", Issuer: "http://localhost:8991", Owner: "alice", Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "chat", Name: "Chat", Port: 9006, Route: "/chat"}, {ID: "monitor", Name: "My Monitor", Route: "/monitor"}, {ID: "drive", Name: "Drive", Route: "/drive"}}}}
}
func TestChatDedicatedOriginAndLegacyLinks(t *testing.T) {
	c := chatTestConfig(t)
	if chatOrigin(c) != "http://127.0.0.1:9006/chat" {
		t.Fatal("local hostname is not isolated")
	}
	for _, raw := range []string{"https://workspace.example.test:123/chat", "https://chat.example.test/chat?evil=1", "https://chat.example.test/other"} {
		if validChatOrigin(raw, "https://workspace.example.test") {
			t.Fatal("invalid origin accepted", raw)
		}
	}
	p := &Portal{rt: &Runtime{Config: c}}
	mux := http.NewServeMux()
	p.chatRoutes(mux)
	for _, path := range []string{"/chat/", "/chat/?topic=usrAAAAAAAAAAA", "/chat/auth/callback?code=must-not-forward&state=old"} {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, httptest.NewRequest("GET", "http://localhost:8990"+path, nil))
		if w.Code != 307 || !strings.HasPrefix(w.Header().Get("Location"), chatOrigin(c)) || strings.Contains(w.Header().Get("Location"), "must-not-forward") {
			t.Fatal("legacy entry migration failed")
		}
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest("GET", "http://localhost:8990/chat/v0/file/s/synthetic?auth=must-not-forward&SeCrEt=must-not-forward&token=must-not-forward&sid=must-not-forward&apikey=must-not-forward&download=1", nil))
	if w.Code != 307 || w.Header().Get("Location") != chatOrigin(c)+"/v0/file/s/synthetic?download=1" {
		t.Fatal("legacy media redirect transplanted a transport credential")
	}
}
func TestChatCatalogPrivateRoleScope(t *testing.T) {
	c := chatTestConfig(t)
	p := &Portal{rt: &Runtime{Config: c, chatDirectoryToken: strings.Repeat("x", 40)}}
	mux := http.NewServeMux()
	p.chatCatalogRoute(mux)
	for _, test := range []struct {
		subject, auth, remote string
		code                  int
		monitor               bool
	}{{"bob", "", "127.0.0.1:1000", 404, false}, {"bob", "Bearer " + p.rt.chatDirectoryToken, "192.0.2.4:1000", 404, false}, {"bob", "Bearer " + p.rt.chatDirectoryToken, "127.0.0.1:1000", 200, false}, {"alice", "Bearer " + p.rt.chatDirectoryToken, "127.0.0.1:1000", 200, true}} {
		r := httptest.NewRequest("POST", "http://localhost:8990/_integration/chat/catalog", strings.NewReader(`{"issuer":"`+c.Issuer+`","subject":"`+test.subject+`"}`))
		r.RemoteAddr = test.remote
		r.Header.Set("Authorization", test.auth)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != test.code {
			t.Fatalf("status %d", w.Code)
		}
		if w.Code != 200 {
			continue
		}
		var body struct {
			Home string              `json:"home"`
			Apps []map[string]string `json:"apps"`
		}
		if json.Unmarshal(w.Body.Bytes(), &body) != nil || body.Home != c.Origin {
			t.Fatal("invalid catalog")
		}
		found := false
		for _, app := range body.Apps {
			if app["id"] == "monitor" {
				found = true
			}
			if !strings.HasPrefix(app["url"], "http://") {
				t.Fatal("relative destination")
			}
		}
		if found != test.monitor {
			t.Fatal("catalog exposed owner-only app")
		}
	}
}
