package suite

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"path/filepath"
	"testing"
)

func TestNewsNativeIdentityAndStorage(t *testing.T) {
	a := catalog.App{ID: "news", Kind: "process", Executable: "MyNews.exe", Port: 9030, Route: "/news", OIDC: true}
	c := NewConfig(Request{Port: 29000, Mode: "local", Programs: "C:/Synthetic/Apps", Data: "C:/Synthetic/Data", Files: "C:/Synthetic/Files"}, catalog.Catalog{Apps: []catalog.App{a}})
	a, ok := c.Catalog.Find("news")
	if !ok || a.Port != 29040 || !foundationApp("news") || !knownUpdateStorageProfile("news") {
		t.Fatal("News managed runtime not registered", a)
	}
	if appWorkingDir(c, "news") != filepath.Join(c.Data, "news") || len(appUpdatePrivatePaths("news")) != 1 {
		t.Fatal("News private recovery missing")
	}
	found := false
	for _, s := range suiteSignIns {
		if s.id == "news" {
			found = s.login == "/news/auth/login" && s.logout == "/news/auth/logout" && s.cookies[0] == "news_session"
		}
	}
	if !found {
		t.Fatal("News shared identity missing")
	}
}
