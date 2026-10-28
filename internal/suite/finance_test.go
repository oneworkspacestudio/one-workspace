package suite

import (
	"github.com/oneworkspacestudio/one-workspace/internal/branding"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"path/filepath"
	"strings"
	"testing"
)

func TestFinanceManagedIdentityAndProtectedLedger(t *testing.T) {
	a := catalog.App{ID: "finance", Name: "My Finance", Kind: "process", Package: "finance", Executable: "MyFinance.exe", Port: 9024, Route: "/finance", OIDC: true, Callback: "/finance/auth/callback", Args: []string{"--data", "${DATA}/finance", "--listen", "127.0.0.1:${PORT}", "--origin", "${ORIGIN}/finance", "--no-browser"}}
	c := NewConfig(Request{Programs: filepath.Join(t.TempDir(), "Programs"), Data: filepath.Join(t.TempDir(), "Private"), Files: filepath.Join(t.TempDir(), "Files"), Port: 23900, Mode: "local"}, catalog.Catalog{Apps: []catalog.App{a}})
	c.Owner = "synthetic-owner"
	a = c.Catalog.Apps[0]
	if a.Port != 23934 || !foundationApp("finance") || appWorkingDir(c, "finance") != filepath.Join(c.Data, "finance") || !knownUpdateStorageProfile("finance") {
		t.Fatal("Finance must have an isolated managed storage/identity contract")
	}
	if p := appUpdatePrivatePaths("finance"); len(p) != 1 || p[0] != "finance" || preservedUpdateData("finance/actual/server-files/account.sqlite") {
		t.Fatal("The actual ledger and credential store must be in protected recovery")
	}
	rt := &Runtime{Config: c, Secrets: Secrets{Clients: map[string]Client{"finance": {ID: "synthetic-client", Secret: "synthetic-secret"}}}}
	cmd := rt.command(a)
	for _, want := range []string{"WORKSPACE_OWNER_SUBJECT=" + c.Owner, "WORKSPACE_HOME_URL=" + c.Origin, "WORKSPACE_OIDC_REDIRECT_URI=" + c.Origin + "/finance/auth/callback", "MYFINANCE_BACKUP_DIR=" + filepath.Join(c.Files, "One Workspace Backups", "My Finance")} {
		found := false
		for _, got := range cmd.Env {
			found = found || got == want
		}
		if !found {
			t.Fatal("missing Finance runtime contract")
		}
	}
	if !strings.Contains(strings.Join(cmd.Args, " "), c.Origin+"/finance") || branding.Icon("finance") == branding.Icon("unknown") {
		t.Fatal("Finance origin/icon missing")
	}
}
