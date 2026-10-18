package suite

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestHostingWebhookHasExactAnonymousBoundary(t *testing.T) {
	p := &Portal{rt: &Runtime{Config: Config{Catalog: catalog.Catalog{Apps: []catalog.App{{ID: "hosting", Route: "/hosting"}}}}}}
	id := strings.Repeat("a", 32)
	for _, test := range []struct {
		method, path string
		allowed      bool
	}{{"POST", "/hosting/hooks/github/" + id, true}, {"GET", "/hosting/hooks/github/" + id, false}, {"POST", "/hosting/hooks/github/invalid", false}, {"POST", "/hosting/hooks/github/" + id + "/extra", false}, {"POST", "/hosting/api/projects", false}, {"POST", "/hosting/auth/logout", false}, {"POST", "/mail/hooks/github/" + id, false}} {
		r := httptest.NewRequest(test.method, test.path, nil)
		if p.hostingWebhook(r) != test.allowed {
			t.Errorf("unexpected anonymous route: %s %s", test.method, test.path)
		}
	}
	p.rt.Config.Catalog.Apps = nil
	if p.hostingWebhook(httptest.NewRequest("POST", "/hosting/hooks/github/"+id, nil)) {
		t.Fatal("uninstalled Hosting route was anonymous")
	}
}
func TestHostingOwnerAndStorageLaunchContract(t *testing.T) {
	a := catalog.App{ID: "hosting", OIDC: true, Executable: "MyHosting.exe", Port: 19007, Callback: "/hosting/auth/callback", Args: []string{"--data", "${DATA}/my-hosting", "--listen", "127.0.0.1:${PORT}", "--origin", "${ORIGIN}/hosting", "--no-browser"}}
	rt := &Runtime{Config: Config{Data: "C:/fixture/private", Programs: "C:/fixture/apps", Owner: "configured-owner", HostingDomain: "example.test", Origin: "https://workspace.example.test", Issuer: "https://accounts.example.test"}, Secrets: Secrets{Clients: map[string]Client{"hosting": {ID: "fixture-client", Secret: "fixture-only"}}}}
	cmd := rt.command(a)
	env := strings.Join(cmd.Env, "\n")
	for _, expected := range []string{"WORKSPACE_OWNER_SUBJECT=configured-owner", "WORKSPACE_OIDC_CLIENT_ID=fixture-client", "WORKSPACE_SITES_DOMAIN=sites.example.test", "WORKSPACE_SITES_TRANSPORT=direct"} {
		if !strings.Contains(env, expected) {
			t.Fatal("Hosting identity contract missing", expected)
		}
	}
	if !strings.Contains(strings.Join(cmd.Args, " "), "my-hosting") {
		t.Fatal("Hosting private path missing")
	}
	var signIn *appSignIn
	for i := range suiteSignIns {
		if suiteSignIns[i].id == "hosting" {
			signIn = &suiteSignIns[i]
		}
	}
	if signIn == nil || signIn.probe != "/api/me" || signIn.login != "/hosting/auth/login" {
		t.Fatal("Hosting lost native seamless sign-in")
	}
}

func TestHostingTunnelUsesConnectedRootZone(t *testing.T) {
	c := Config{Data: t.TempDir(), Tunnel: true, HostingDomain: "inbox.example.test", Origin: "https://workspace-123.example.test"}
	for _, zoneName := range []string{"", "example.test"} {
		if e := platform.SaveJSON(filepath.Join(c.Data, "cloudflare-resources.json"), map[string]string{"accountId": strings.Repeat("a", 32), "zoneId": strings.Repeat("b", 32), "zoneName": zoneName, "tunnelId": "00000000-0000-0000-0000-000000000001"}); e != nil {
			t.Fatal(e)
		}
		env := strings.Join(hostingSiteEnvironment(c, Secrets{HostingToken: "synthetic-provider-only"}), "\n")
		for _, expected := range []string{"WORKSPACE_SITES_DOMAIN=example.test", "WORKSPACE_SITES_STYLE=flat", "WORKSPACE_SITES_TRANSPORT=tunnel", "WORKSPACE_SITES_CF_TOKEN=synthetic-provider-only"} {
			if !strings.Contains(env, expected) {
				t.Fatal("Managed site publication contract", expected)
			}
		}
	}
}

func TestHostingLocalInstallHasNoInvalidPublicDomain(t *testing.T) {
	env := strings.Join(hostingSiteEnvironment(Config{}, Secrets{}), "\n")
	if !strings.HasPrefix(env, "WORKSPACE_SITES_DOMAIN=\n") || strings.Contains(env, "sites.") {
		t.Fatal("Local-only install produced a public domain", env)
	}
	c := NewConfig(Request{Mode: "cloudflare", Domain: "inbox.example.test", ZoneID: strings.Repeat("a", 32), Port: 19300}, catalog.Catalog{Apps: []catalog.App{{ID: "accounts"}, {ID: "hosting"}}})
	if c.HostingDomain != "inbox.example.test" || c.DNSZoneID == "" {
		t.Fatal("Hosting-only installation lost its domain connection")
	}
}
