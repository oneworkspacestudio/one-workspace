package suite

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestRemoteRejectsMissingOrModifiedRuntime(t *testing.T) {
	root := t.TempDir()
	good := []byte("verified sidecar")
	c := Config{Programs: root, Catalog: catalog.Catalog{Packages: map[string]catalog.Package{
		"remote-runtime-test": {File: "sidecar.cjs", SHA256: fmt.Sprintf("%x", sha256.Sum256(good)), Size: int64(len(good))},
	}}}
	for _, data := range [][]byte{nil, []byte("modified")} {
		if data != nil {
			if err := os.WriteFile(filepath.Join(root, "sidecar.cjs"), data, 0600); err != nil {
				t.Fatal(err)
			}
		}
		if close, err := lockRemotePayloads(c, catalog.App{ID: "remote"}); err == nil {
			close()
			t.Fatal("untrusted sidecar accepted")
		}
	}
	if err := os.WriteFile(filepath.Join(root, "sidecar.cjs"), good, 0600); err != nil {
		t.Fatal(err)
	}
	close, err := lockRemotePayloads(c, catalog.App{ID: "remote"})
	if err != nil {
		t.Fatal(err)
	}
	close()
}

func TestRemoteRuntimeContract(t *testing.T) {
	c := Config{Data: t.TempDir(), Programs: t.TempDir(), Port: 29000}
	if !foundationApp("remote") || !knownUpdateStorageProfile("remote") || appWorkingDir(c, "remote") != filepath.Join(c.Data, "remote") {
		t.Fatal("Remote lacks lifecycle/storage registration")
	}
	if p := appUpdatePrivatePaths("remote"); len(p) != 1 || p[0] != "remote" {
		t.Fatal("Remote recovery scope missing")
	}
	if !appOwnsUpdatePackage(catalog.App{Package: "remote"}, "remote-runtime-node") {
		t.Fatal("Remote does not own its runtime")
	}
	if catalog.CategoryFor("remote").ID != "administration" {
		t.Fatal("Remote category missing")
	}
}
func TestRemotePublicAdmission(t *testing.T) {
	for _, tt := range []struct {
		method, path string
		allowed      bool
	}{
		{"GET", "/remote/", true}, {"GET", "/remote/assets/index-abc.js", true}, {"POST", "/remote/support-api/redeem", true}, {"GET", "/remote/desktop", true},
		{"GET", "/remote/api/remote/devices", false}, {"POST", "/remote/api/remote/enroll", false}, {"POST", "/remote/support-api/status", false}, {"GET", "/remote/assets/../private", false}, {"GET", "/remote/assets/a/b.js", false},
	} {
		r := httptest.NewRequest(tt.method, tt.path, nil)
		if publicRemoteRequest(r) != tt.allowed {
			t.Errorf("%s %s", tt.method, tt.path)
		}
	}
}
func TestRemoteRuntimeUpdateAndFailedHealthRecovery(t *testing.T) {
	for _, fail := range []bool{false, true} {
		t.Run(fmt.Sprint("failed-health=", fail), func(t *testing.T) {
			c, s := updateFixture(t)
			c.Catalog.Selection = &catalog.AppSelection{}
			c.Catalog.Apps = append(c.Catalog.Apps, catalog.App{ID: "remote", Name: "My Remote", Kind: "process", Package: "remote", Executable: "MyRemote.exe", Port: 29135, Version: "0.1.0-preview.3", Route: "/remote", OIDC: true, Callback: "/remote/auth/callback", Args: []string{"--data", "${DATA}/remote"}, Source: "https://github.com/oneworkspacestudio/my-remote/releases/tag/v0.1.0-preview.3"})
			c.Catalog.Packages["remote"] = updatePackage("MyRemote.exe", "remote-original")
			c.Catalog.Packages["remote-runtime-test"] = updatePackage("remote-runtime/sidecar.cjs", "runtime-original")
			s.Clients["remote"] = Client{ID: "remote-fixture", Secret: "fixture-only"}
			if err := Save(c, s); err != nil {
				t.Fatal(err)
			}
			for _, dir := range []string{filepath.Join(c.Programs, "remote-runtime"), filepath.Join(c.Data, "remote")} {
				if err := os.MkdirAll(dir, 0700); err != nil {
					t.Fatal(err)
				}
			}
			files := map[string]string{filepath.Join(c.Programs, "MyRemote.exe"): "remote-original", filepath.Join(c.Programs, "remote-runtime/sidecar.cjs"): "runtime-original", filepath.Join(c.Data, "remote/settings.json"): "private-original"}
			for p, b := range files {
				if err := os.WriteFile(p, []byte(b), 0600); err != nil {
					t.Fatal(err)
				}
			}
			target := c.Catalog
			target.Apps = append([]catalog.App(nil), c.Catalog.Apps...)
			target.Apps[1].Version = "0.1.0-preview.4"
			target.Packages = map[string]catalog.Package{}
			for id, p := range c.Catalog.Packages {
				target.Packages[id] = p
			}
			target.Packages["remote-runtime-test"] = updatePackage("remote-runtime/sidecar.cjs", "runtime-updated")
			stage := t.TempDir()
			os.MkdirAll(filepath.Join(stage, "remote-runtime"), 0700)
			os.WriteFile(filepath.Join(stage, "remote-runtime/sidecar.cjs"), []byte("runtime-updated"), 0600)
			starts := 0
			ops := updateOperations{stop: func(context.Context, Config) error { return nil }, start: func(_ context.Context, next Config, nonce string) error {
				starts++
				if fail && starts == 1 {
					os.WriteFile(filepath.Join(c.Data, "remote/settings.json"), []byte("migrated-private"), 0600)
					os.WriteFile(filepath.Join(c.Data, "notes.db"), []byte("unrelated-newer"), 0600)
					return errors.New("injected remote health failure")
				}
				return nil
			}}
			err := applyStagedWith(context.Background(), c, s, target, stage, "", quietUpdate, ops)
			if (err != nil) != fail {
				t.Fatal(err)
			}
			for p, want := range map[string]string{filepath.Join(c.Data, "Workspace.exe"): "manager-original", filepath.Join(c.Programs, "MyNotes.exe"): "notes-original", filepath.Join(c.Data, "remote/settings.json"): "private-original"} {
				got, _ := os.ReadFile(p)
				if string(got) != want {
					t.Fatal(p, string(got))
				}
			}
			wantRuntime := "runtime-updated"
			if fail {
				wantRuntime = "runtime-original"
				got, _ := os.ReadFile(filepath.Join(c.Data, "notes.db"))
				if string(got) != "unrelated-newer" {
					t.Fatal("rollback replaced unrelated app state")
				}
			}
			got, _ := os.ReadFile(filepath.Join(c.Programs, "remote-runtime/sidecar.cjs"))
			if string(got) != wantRuntime {
				t.Fatal("runtime not committed/restored")
			}
			backups, e := UpdateBackups(c)
			if e != nil || len(backups) < 1 {
				t.Fatal(e)
			}
			b, e := verifyUpdateBackup(c, backups[0])
			if e != nil || b.Scope == nil || len(b.Scope.Apps) != 1 || b.Scope.Apps[0] != "remote" {
				t.Fatal("independent recovery scope missing", e)
			}
		})
	}
}
