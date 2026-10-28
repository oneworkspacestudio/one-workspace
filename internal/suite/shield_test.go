package suite

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/oneworkspacestudio/one-workspace/internal/branding"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

func TestShieldManagedOwnerAndStorageContract(t *testing.T) {
	a := catalog.App{ID: "shield", Name: "My Shield", Kind: "process", Package: "shield", Executable: "MyShield.exe", Port: 9022, Route: "/shield", OIDC: true, Callback: "/shield/auth/callback", Args: []string{"--data", "${DATA}/shield", "--listen", "127.0.0.1:${PORT}", "--origin", "${ORIGIN}/shield", "--no-browser"}}
	c := NewConfig(Request{Programs: filepath.Join(t.TempDir(), "Programs"), Data: filepath.Join(t.TempDir(), "Private"), Files: filepath.Join(t.TempDir(), "Files"), Port: 23900, Mode: "local"}, catalog.Catalog{Apps: []catalog.App{a}})
	c.Owner = "synthetic-owner-subject"
	a = c.Catalog.Apps[0]
	if a.Port != 23932 || !foundationApp("shield") || appWorkingDir(c, "shield") != filepath.Join(c.Data, "shield") || !knownUpdateStorageProfile("shield") {
		t.Fatal("Shield is not a managed, separately stored application")
	}
	paths := appUpdatePrivatePaths("shield")
	if len(paths) != 1 || paths[0] != "shield" || preservedUpdateData("shield/state.sqlite") || preservedUpdateData("shield/state.sqlite-wal") {
		t.Fatal("Shield state must remain inside verified update recovery")
	}
	rt := &Runtime{Config: c, Secrets: Secrets{Clients: map[string]Client{"shield": {ID: "synthetic-shield", Secret: "synthetic-secret"}}}}
	cmd := rt.command(a)
	for _, want := range []string{"WORKSPACE_OWNER_SUBJECT=" + c.Owner, "WORKSPACE_HOME_URL=" + c.Origin, "WORKSPACE_OIDC_ISSUER=" + c.Issuer, "WORKSPACE_OIDC_REDIRECT_URI=" + c.Origin + "/shield/auth/callback", "MYSHIELD_LAB=0"} {
		found := false
		for _, got := range cmd.Env {
			found = found || got == want
		}
		if !found {
			t.Fatal("missing managed identity/environment contract")
		}
	}
	if !strings.Contains(strings.Join(cmd.Args, " "), c.Origin+"/shield") || cmd.Dir != filepath.Join(c.Data, "shield") {
		t.Fatal("managed origin or working directory is wrong")
	}
	p := &Portal{rt: rt}
	if len(p.tiles(false)) != 0 || len(p.tiles(true)) != 2 {
		t.Fatal("host security controls must be visible only to the owner, alongside Backups")
	}
	if branding.Icon("shield") == branding.Icon("unknown") {
		t.Fatal("Shield icon not registered")
	}
}
