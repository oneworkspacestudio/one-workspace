package suite

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"path/filepath"
	"testing"
)

func TestCommuteManagedContract(t *testing.T) {
	c := NewConfig(Request{Port: 29000, Mode: "local", Programs: "C:/Synthetic/Apps", Data: "C:/Synthetic/Data", Files: "C:/Synthetic/Files"}, catalog.Catalog{Apps: []catalog.App{{ID: "commute", Kind: "process", Executable: "MyCommute.exe", Port: 9040, Route: "/commute", OIDC: true}}})
	a, ok := c.Catalog.Find("commute")
	if !ok || a.Port != 29050 || !foundationApp("commute") || !knownUpdateStorageProfile("commute") {
		t.Fatal("Commute runtime/auth/storage not registered", a)
	}
	if appWorkingDir(c, "commute") != filepath.Join(c.Data, "commute") || len(appUpdatePrivatePaths("commute")) != 1 {
		t.Fatal("private recovery scope missing")
	}
	for _, s := range suiteSignIns {
		if s.id == "commute" {
			if s.login != "/commute/auth/login" || s.logout != "/commute/auth/logout" || s.cookies[0] != "commute_session" {
				t.Fatal("sign-in/out contract drift")
			}
			return
		}
	}
	t.Fatal("shared sign-in/out missing")
}
