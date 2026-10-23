package suite

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"testing"
)

func TestMy3DNativeContract(t *testing.T) {
	c := NewConfig(Request{Port: 29000, Mode: "local", Programs: "C:/Synthetic/Apps", Data: "C:/Synthetic/Data", Files: "C:/Synthetic/Files"}, catalog.Catalog{Apps: []catalog.App{{ID: "my3d", Kind: "process", Port: 9018, Route: "/my3d"}}})
	a, ok := c.Catalog.Find("my3d")
	if !ok || a.Port != 29028 || !foundationApp(a.ID) {
		t.Fatal("native contract", a)
	}
	if appWorkingDir(c, "my3d") == c.Data {
		t.Fatal("app working folder escaped")
	}
	found := false
	for _, s := range suiteSignIns {
		if s.id == "my3d" {
			found = s.route == "/my3d" && s.probe == "/api/me"
		}
	}
	if !found {
		t.Fatal("missing normal OIDC session contract")
	}
}
