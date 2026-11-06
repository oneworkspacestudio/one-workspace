package suite

import (
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"path/filepath"
	"testing"
)

func TestPodcastsManagedContract(t *testing.T) {
	c := NewConfig(Request{Port: 29000, Mode: "local", Programs: "C:/Synthetic/Apps", Data: "C:/Synthetic/Data", Files: "C:/Synthetic/Files"}, catalog.Catalog{Apps: []catalog.App{{ID: "podcasts", Kind: "process", Executable: "MyPodcasts.exe", Port: 9044, Route: "/podcasts", OIDC: true}}})
	a, ok := c.Catalog.Find("podcasts")
	if !ok || a.Port != 29054 || !foundationApp("podcasts") || !knownUpdateStorageProfile("podcasts") {
		t.Fatal("Podcasts runtime/auth/storage not registered", a)
	}
	if appWorkingDir(c, "podcasts") != filepath.Join(c.Data, "podcasts") || len(appUpdatePrivatePaths("podcasts")) != 1 {
		t.Fatal("private recovery scope missing")
	}
	for _, s := range suiteSignIns {
		if s.id == "podcasts" {
			if s.login != "/podcasts/auth/login" || s.logout != "/podcasts/auth/logout" || s.cookies[0] != "podcasts_session" {
				t.Fatal("sign-in/out contract drift")
			}
			return
		}
	}
	t.Fatal("shared sign-in/out missing")
}
