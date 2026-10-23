package suite

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/distribution"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

type installerFileTransport struct {
	base      http.RoundTripper
	url, path string
}

func (tr installerFileTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.URL.String() != tr.url {
		return tr.base.RoundTrip(r)
	}
	f, e := os.Open(tr.path)
	if e != nil {
		return nil, e
	}
	return &http.Response{StatusCode: 200, Body: f, Header: http.Header{}, Request: r}, nil
}

// Opt-in only: exact release setup execution against an explicitly named fixture.
// The installer asset is served through the normal transport/digest path without
// publishing an untested release or giving tests access to a live installation.
func TestNativeVerifiedInstallerRepair(t *testing.T) {
	installer := os.Getenv("WORKSPACE_NATIVE_REPAIR_INSTALLER")
	if installer == "" {
		t.Skip("requires named disposable native fixture")
	}
	data := filepath.Clean(os.Getenv("WORKSPACE_NATIVE_REPAIR_DATA"))
	want := filepath.Clean(`C:\Projects\.cache\release-check-403-20261023\final-fixture\Private`)
	if !strings.EqualFold(data, want) {
		t.Fatal("not the named disposable fixture")
	}
	if !strings.Contains(strings.ToLower(os.Getenv("APPDATA")), "release-check-403-20261023") {
		t.Fatal("fixture AppData not isolated")
	}
	c, s, e := Load(data)
	if e != nil {
		t.Fatal(e)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Minute)
	defer cancel()
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		if e := stopUpdateRuntime(cleanup, c); e != nil {
			t.Errorf("fixture cleanup: %v", e)
		}
	})
	if e = stopUpdateRuntime(ctx, c); e != nil {
		t.Fatal(e)
	}
	tasks, ok := c.Catalog.Find("tasks")
	if !ok {
		t.Fatal("fixture needs Tasks")
	}
	if e = os.WriteFile(filepath.Join(c.Programs, tasks.Executable), []byte("synthetic damaged program"), 0600); e != nil {
		t.Fatal(e)
	}
	marker := filepath.Join(c.Files, "Documents", "reinstall-preserve.txt")
	if e = os.WriteFile(marker, []byte("keep this user document"), 0600); e != nil {
		t.Fatal(e)
	}
	f, e := os.Open(installer)
	if e != nil {
		t.Fatal(e)
	}
	hash := sha256.New()
	st, _ := f.Stat()
	_, e = io.Copy(hash, f)
	f.Close()
	if e != nil {
		t.Fatal(e)
	}
	digest := hex.EncodeToString(hash.Sum(nil))
	version := os.Getenv("WORKSPACE_NATIVE_REPAIR_VERSION")
	tag := "v" + version
	repo := distribution.Repository
	asset := distribution.Asset{Name: "WorkspaceSetup-windows-amd64.exe", State: "uploaded", Size: st.Size(), Digest: "sha256:" + digest, APIURL: "https://api.github.com/repos/" + repo + "/releases/assets/999999", URL: "https://github.com/" + repo + "/releases/download/" + tag + "/WorkspaceSetup-windows-amd64.exe"}
	b, _ := json.Marshal(c.Catalog)
	identity := sha256.Sum256(b)
	plan := distribution.UpdatePlan{CheckedAt: time.Now(), BaseCatalogSHA256: hex.EncodeToString(identity[:]), Installer: &distribution.InstallerRelease{Version: version, Source: "https://github.com/" + repo + "/releases/tag/" + tag, Download: distribution.UpdateDownload{Asset: asset, Repository: repo, Tag: tag}}}
	original := http.DefaultTransport
	http.DefaultTransport = installerFileTransport{original, asset.URL, installer}
	defer func() { http.DefaultTransport = original }()
	result, e := RunInstallerRepair(ctx, data, plan, func(s string) { t.Log(s) })
	if e != nil {
		t.Fatal(e)
	}
	after, secrets, e := Load(data)
	if e != nil {
		t.Fatal(e)
	}
	if result.Version != version || len(result.NeedsAttention) != 0 || after.Owner != c.Owner || after.Autostart != c.Autostart || !reflect.DeepEqual(s, secrets) || !reflect.DeepEqual(c.Catalog.Apps, after.Catalog.Apps) || !reflect.DeepEqual(c.Catalog.Selection, after.Catalog.Selection) {
		t.Fatal("identity, app selection, contracts or health changed")
	}
	for id, p := range c.Catalog.Packages {
		if id != "manager" && id != "notice-Workspace-release-status.md" && id != "notice-organization-components.json" && after.Catalog.Packages[id] != p {
			t.Fatalf("installed app package changed: %s", id)
		}
	}
	content, _ := os.ReadFile(marker)
	if string(content) != "keep this user document" {
		t.Fatal("user document changed")
	}
	if e = verifyInstalledPrograms(after); e != nil {
		t.Fatal(e)
	}
	if e = platform.SaveJSON(os.Getenv("WORKSPACE_NATIVE_REPAIR_RECEIPT"), map[string]any{"sourceVersion": c.Version, "version": after.Version, "result": result, "identityPreserved": true, "appContractsPreserved": true, "damagedTaskRepaired": true, "documentPreserved": true}); e != nil {
		t.Fatal(e)
	}
}
