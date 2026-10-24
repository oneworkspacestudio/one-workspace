package suite

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"testing"
)

func TestCinemaNativeContract(t *testing.T) {
	c := NewConfig(Request{Port: 29000, Mode: "local", Programs: "C:/Synthetic/Apps", Data: "C:/Synthetic/Private", Files: "C:/Synthetic/Files"}, catalog.Catalog{Apps: []catalog.App{{ID: "cinema", Kind: "process", Port: 9020, Route: "/cinema"}}})
	a, ok := c.Catalog.Find("cinema")
	if !ok || a.Port != 29030 || !foundationApp(a.ID) || appWorkingDir(c, a.ID) == c.Data {
		t.Fatal("Cinema contract missing")
	}
	found := false
	for _, s := range suiteSignIns {
		if s.id == "cinema" {
			found = s.probe == "/api/me" && s.route == "/cinema" && s.cookies[0] == "cinema_session"
		}
	}
	if !found {
		t.Fatal("native sign-in missing")
	}
}
