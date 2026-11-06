package suite

import (
	"context"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"os"
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

// Descriptor updates must retain the shared metadata recovery boundary.
func TestPodcastsDescriptorUpdateCreatesWorkspaceRecovery(t *testing.T) {
	c, _ := updateFixture(t)
	c.Version = "0.3.0-preview.122"
	c.Catalog.Version = c.Version
	c.Catalog.Selection = &catalog.AppSelection{}
	c.Catalog.Apps = append(c.Catalog.Apps, catalog.App{ID: "podcasts", Name: "My Podcasts", Description: "Original description", Kind: "process", Package: "podcasts", Executable: "MyPodcasts.exe", Port: 29154, Route: "/podcasts", Version: "0.1.0-preview.2", OIDC: true, Callback: "/podcasts/auth/callback", Args: []string{"--data", "${DATA}/podcasts"}, Source: "https://github.com/oneworkspacestudio/my-podcasts/releases/tag/v0.1.0-preview.2"})
	c.Catalog.Packages["podcasts"] = updatePackage("MyPodcasts.exe", "podcasts-original")
	os.WriteFile(filepath.Join(c.Programs, "MyPodcasts.exe"), []byte("podcasts-original"), 0600)
	target := c.Catalog
	target.Apps = append([]catalog.App(nil), c.Catalog.Apps...)
	target.Packages = map[string]catalog.Package{}
	for id, p := range c.Catalog.Packages {
		target.Packages[id] = p
	}
	target.Apps[1].Version = "0.1.0-preview.3"
	target.Packages["podcasts"] = updatePackage("MyPodcasts.exe", "podcasts-updated")
	if planAppUpdateScope(c, target) == nil {
		t.Fatal("compatible version-only update lost independent recovery")
	}
	target.Apps[1].Description = "DeepSeek and local AI writing with background production"
	if planAppUpdateScope(c, target) != nil {
		t.Fatal("descriptor change incorrectly limited recovery to the app")
	}
	policy, err := planUpdateMetadata(c, target, false)
	if err != nil {
		t.Fatal(err)
	}
	receipt, err := makeScopedBackupDetailed(context.Background(), c, "", func(string) {}, nil, "update", nil, policy)
	if err != nil {
		t.Fatal(err)
	}
	backup, err := verifyUpdateBackup(c, receipt)
	if err != nil {
		t.Fatal(err)
	}
	if backup.Metadata == nil || backup.Scope != nil {
		t.Fatal("shared metadata recovery missing")
	}
	hasOtherMetadata := false
	for _, f := range backup.Files {
		if f.Area == "private" && f.Path == "notes.db" {
			hasOtherMetadata = true
		}
	}
	if !hasOtherMetadata {
		t.Fatal("other application metadata was omitted from workspace recovery")
	}
}
