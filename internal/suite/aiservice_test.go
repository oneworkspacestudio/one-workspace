package suite

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"path/filepath"
	"testing"
)

func TestAIServiceManagedIdentityStorageAndPorts(t *testing.T) {
	c := NewConfig(Request{Port: 25100}, catalog.Catalog{Apps: []catalog.App{{ID: "aiservice", Port: 9042}}})
	app, _ := c.Catalog.Find("aiservice")
	if app.Port != 25152 {
		t.Fatal(app.Port)
	}
	if !foundationApp("aiservice") || !knownUpdateStorageProfile("aiservice") {
		t.Fatal("unregistered profile")
	}
	c.Data = t.TempDir()
	if appWorkingDir(c, "aiservice") != filepath.Join(c.Data, "aiservice") || len(appUpdatePrivatePaths("aiservice")) != 1 {
		t.Fatal("private state ownership")
	}
	if !appOwnsUpdatePackage(catalog.App{Package: "aiservice"}, "aiservice-contract") {
		t.Fatal("unowned payload")
	}
	found := false
	for _, v := range suiteSignIns {
		if v.id == "aiservice" {
			found = v.route == "/aiservice" && v.login == "/aiservice/auth/login" && v.cookies[0] == "aiservice_session"
		}
	}
	if !found {
		t.Fatal("missing SSO/logout")
	}
}
