package suite

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestSocialStorageDeniesAncestorAndDescendantAllows(t *testing.T) {
	base := map[string]socialRule{socialPath: {DenyAll: true}}
	if !safeSocialRules(base) {
		t.Fatal("reserved root denied")
	}
	for _, p := range []string{"/", "/My Photos/", "/My Photos/My Social/", "/My Photos/My Social/users/a/"} {
		r := map[string]socialRule{socialPath: {DenyAll: true}}
		allow := socialRule{}
		allow.Allow.Users = []string{"another-user"}
		r[p] = allow
		if safeSocialRules(r) {
			t.Fatal("allow rule could expose originals", p)
		}
	}
	if safeSocialRules(map[string]socialRule{}) {
		t.Fatal("missing deny accepted")
	}
}

func TestSocialGuardReusesOnlyAuthenticationAndRechecksACL(t *testing.T) {
	logins, reads := 0, 0
	conflict := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "auth/login"):
			logins++
			w.Write([]byte("fixture.jwt.token"))
		case r.Method == "POST":
			w.Write([]byte(`{}`))
		default:
			reads++
			if conflict {
				w.Write([]byte(`{"/My Photos/My Social/":{"denyAll":true},"/":{"allow":{"users":["another"]}}}`))
			} else {
				w.Write([]byte(`{"/My Photos/My Social/":{"denyAll":true}}`))
			}
		}
	}))
	defer server.Close()
	_, raw, _ := net.SplitHostPort(strings.TrimPrefix(server.URL, "http://"))
	port, _ := strconv.Atoi(raw)
	rt := &Runtime{Config: Config{Files: t.TempDir(), Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "drive", Port: port}}}}, Secrets: Secrets{RecoveryPassword: "synthetic"}}
	if e := rt.verifySocialStorage(t.Context(), true); e != nil {
		t.Fatal(e)
	}
	if e := rt.verifySocialStorage(t.Context(), false); e != nil {
		t.Fatal(e)
	}
	conflict = true
	if e := rt.verifySocialStorage(t.Context(), false); e == nil {
		t.Fatal("changed ACL accepted")
	}
	if logins != 1 || reads != 3 {
		t.Fatal("must reuse auth without caching authorization", logins, reads)
	}
}
