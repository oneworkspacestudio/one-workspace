package maintenance

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"github.com/oneworkspacestudio/one-workspace/internal/suite"
)

// Exact packaged page and native capability gate; all ownership and API state
// used by the browser are synthetic. No owner installation or tunnel is used.
func TestNativeReconnectPage(t *testing.T) {
	exe := os.Getenv("WORKSPACE_RECONNECT_EXE")
	if exe == "" {
		t.Skip("opt-in exact executable acceptance")
	}
	root := `C:\Projects\.cache\update-reconnect-20261023\native-verified`
	data := filepath.Join(root, "Private")
	if _, err := os.Stat(root); !os.IsNotExist(err) {
		t.Fatal("use a new disposable native root")
	}
	for _, folder := range []string{data, filepath.Join(root, "Apps"), filepath.Join(root, "Files")} {
		if err := os.MkdirAll(folder, 0700); err != nil {
			t.Fatal(err)
		}
	}
	if err := platform.ProtectDir(data); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(os.Getenv("WORKSPACE_RECONNECT_CATALOG"))
	if err != nil {
		t.Fatal(err)
	}
	cat, err := catalog.Read(raw)
	if err != nil {
		t.Fatal(err)
	}
	c := suite.Config{Schema: 1, Version: cat.Version, Programs: filepath.Join(root, "Apps"), Data: data, Files: filepath.Join(root, "Files"), Origin: "http://localhost:29970", Issuer: "http://localhost:29971", Port: 29970, Owner: "synthetic-owner", Catalog: cat}
	if err = suite.Save(c, suite.Secrets{AdminEmail: "fixture@example.invalid", AdminPassword: platform.Token(), Clients: map[string]suite.Client{}}); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, exe, "-data", data, "updates-background")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err = cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
		}
	})
	var window struct {
		URL      string `json:"url"`
		Protocol int    `json:"remoteProtocol"`
	}
	for i := 0; i < 100; i++ {
		if platform.LoadJSON(filepath.Join(data, "updates", "window.json"), &window) == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	u, err := url.Parse(window.URL)
	if err != nil || u.Host == "" || window.Protocol != suite.RemoteUpdateProtocol {
		t.Fatal("native controller metadata missing/current revision mismatch")
	}
	origin := "http://" + u.Host
	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Get(origin + "/api/state")
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != 403 {
		t.Fatal("anonymous state must stay protected")
	}
	req, _ := http.NewRequest("GET", origin+"/api/state", nil)
	req.Header.Set("X-Workspace-Update", u.Fragment)
	response, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != 200 {
		t.Fatal("synthetic authenticated state unavailable")
	}
	script, err := filepath.Abs(filepath.Join("..", "..", "scripts", "update-reconnect-acceptance.mjs"))
	if err != nil {
		t.Fatal(err)
	}
	node := exec.CommandContext(ctx, "node", script)
	node.Env = append(os.Environ(), "WORKSPACE_QA_PAGE_URL="+origin+"/")
	output, err := node.CombinedOutput()
	if err != nil {
		t.Fatalf("native page browser acceptance: %v\n%s", err, output)
	}
	raw, err = os.ReadFile(exe)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(raw)
	receipt := map[string]any{"exactExecutableSHA256": hex.EncodeToString(sum[:]), "nativePageBrowserPassed": true, "anonymousStateDenied": true, "authenticatedSyntheticStatePassed": true, "remoteProtocol": window.Protocol, "syntheticOnly": true, "testRoot": root}
	out := os.Getenv("WORKSPACE_QA_OUTPUT")
	if !strings.HasPrefix(filepath.ToSlash(out), "I:/Projects/tooling/update-reconnect-20261023/") {
		t.Fatal("use task-owned evidence output")
	}
	b, _ := json.MarshalIndent(receipt, "", "  ")
	if err = os.WriteFile(filepath.Join(out, "native-acceptance.json"), b, 0600); err != nil {
		t.Fatal(err)
	}
	req, _ = http.NewRequest("POST", origin+"/api/close", nil)
	req.Header.Set("Origin", origin)
	req.Header.Set("X-Workspace-Update", u.Fragment)
	response, err = client.Do(req)
	if err == nil {
		response.Body.Close()
	}
	t.Log("Exact Windows page, protected state and reconnection browser checks passed")
}
