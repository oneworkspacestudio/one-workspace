package suite

import (
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRemoteOwnerLogoutRequiresConfirmedStreamRevocation(t *testing.T) {
	for _, success := range []bool{true, false} {
		t.Run(fmt.Sprint(success), func(t *testing.T) {
			calls := 0
			peer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.Method != "POST" || r.URL.Path != "/_integration/revoke-owner-sessions" || r.Header.Get("Origin") != "https://workspace.example.test" || r.Header.Get("Authorization") != "Bearer fixture-remote-client" || r.Header.Get("Cookie") != "" {
					t.Error("incorrect private revocation request")
				}
				if success {
					fmt.Fprint(w, `{"ok":true}`)
				} else {
					http.Error(w, "unavailable", 503)
				}
			}))
			defer peer.Close()
			_, port, _ := net.SplitHostPort(strings.TrimPrefix(peer.URL, "http://"))
			n, _ := strconv.Atoi(port)
			p := &Portal{rt: &Runtime{Config: Config{Owner: "owner", Origin: "https://workspace.example.test", Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "remote", Kind: "process", Port: n}}}}, Secrets: Secrets{Clients: map[string]Client{"remote": {Secret: "fixture-remote-client"}}}, status: map[string]string{"remote": "running"}}, sessions: map[string]session{"fixture-session": {Subject: "owner", Expires: time.Now().Add(time.Hour)}}}
			r := httptest.NewRequest("POST", "https://workspace.example.test/auth/logout", nil)
			r.AddCookie(&http.Cookie{Name: "workspace_session", Value: "fixture-session"})
			w := httptest.NewRecorder()
			if p.revokeRemoteDesktop(w, r) != success {
				t.Fatal("revocation acknowledgement ignored")
			}
			if calls != 1 {
				t.Fatal("owner revocation omitted")
			}
			if !success && w.Code != 503 {
				t.Fatal("unconfirmed sign-out did not fail closed")
			}
			p.sessions["fixture-session"] = session{Subject: "ordinary-member", Expires: time.Now().Add(time.Hour)}
			if !p.revokeRemoteDesktop(httptest.NewRecorder(), r) || calls != 1 {
				t.Fatal("ordinary member could revoke owner streams")
			}
		})
	}
}
