package suite

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"testing"
)

func TestVoiceManagedContract(t *testing.T) {
	c := NewConfig(Request{Port: 29000, Mode: "local", Programs: "C:/Synthetic/Apps", Data: "C:/Synthetic/Data", Files: "C:/Synthetic/Files"}, catalog.Catalog{Apps: []catalog.App{{ID: "voice", Kind: "process", Port: 9028, Route: "/voice"}}})
	v, ok := c.Catalog.Find("voice")
	if !ok || v.Port != 29038 || !foundationApp("voice") {
		t.Fatal("Voice native port/identity missing", v)
	}
	if appWorkingDir(c, "voice") != c.Data+"\\voice" {
		t.Fatal("Voice needs dedicated private working data")
	}
	if len(appUpdatePrivatePaths("voice")) != 1 || appUpdatePrivatePaths("voice")[0] != "voice" {
		t.Fatal("Voice recovery scope missing")
	}
	found := false
	for _, s := range suiteSignIns {
		if s.id == "voice" {
			found = s.route == "/voice" && s.probe == "/api/me" && s.cookies[0] == "voice_session"
		}
	}
	if !found {
		t.Fatal("Voice session/logout missing")
	}
}
