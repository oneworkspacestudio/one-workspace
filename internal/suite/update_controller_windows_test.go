package suite

import (
	"context"
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
)

// A real old helper must be retired even when its protocol matches the new
// manager. Only synthetic NTFS data and public release binaries are used.
func TestNativeControllerUpgrade(t *testing.T) {
	current, prior := os.Getenv("WORKSPACE_CONTROLLER_CURRENT"), os.Getenv("WORKSPACE_CONTROLLER_PRIOR")
	if current == "" || prior == "" {
		t.Skip("opt-in native controller upgrade")
	}
	root := t.TempDir()
	data := filepath.Join(root, "Private")
	for _, p := range []string{data, filepath.Join(root, "Apps"), filepath.Join(root, "Files")} {
		if err := os.MkdirAll(p, 0700); err != nil {
			t.Fatal(err)
		}
	}
	if err := platform.ProtectDir(data); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(os.Getenv("WORKSPACE_CONTROLLER_CATALOG"))
	if err != nil {
		t.Fatal(err)
	}
	cat, err := catalog.Read(raw)
	if err != nil {
		t.Fatal(err)
	}
	c := Config{Schema: 1, Version: "0.3.0-preview.106", Programs: filepath.Join(root, "Apps"), Data: data, Files: filepath.Join(root, "Files"), Origin: "http://localhost:29973", Issuer: "http://localhost:29974", Port: 29973, Owner: "synthetic-owner", Catalog: cat}
	c.Catalog.Version = c.Version
	priorHash, err := UpdateControllerDigest(prior)
	if err != nil {
		t.Fatal(err)
	}
	manager := c.Catalog.Packages["manager"]
	manager.SHA256 = priorHash
	priorInfo, err := os.Stat(prior)
	if err != nil {
		t.Fatal(err)
	}
	manager.Size = priorInfo.Size()
	c.Catalog.Packages["manager"] = manager
	secret := Secrets{AdminEmail: "fixture@example.invalid", AdminPassword: platform.Token(), Clients: map[string]Client{}}
	if err := Save(c, secret); err != nil {
		t.Fatal(err)
	}
	installed := filepath.Join(data, "Workspace.exe")
	if err := platform.CopyFile(prior, installed); err != nil {
		t.Fatal(err)
	}
	var controllers []updateWindow
	defer func() {
		for _, w := range controllers {
			if strings.HasPrefix(filepath.Clean(w.Executable), data+string(os.PathSeparator)) && platform.ProcessRunning(w.PID, w.Executable) {
				p, _ := os.FindProcess(w.PID)
				_ = p.Kill()
				for i := 0; i < 50 && platform.ProcessRunning(w.PID, w.Executable); i++ {
					time.Sleep(100 * time.Millisecond)
				}
				if platform.ProcessRunning(w.PID, w.Executable) {
					t.Error("fixture controller did not stop")
				}
			}
		}
	}()
	launch := func() {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, installed, "--data", data, "updates-background")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("launch: %v %s", err, output)
		}
	}
	waitWindow := func(priorPID int) updateWindow {
		t.Helper()
		for i := 0; i < 150; i++ {
			var w updateWindow
			if platform.LoadJSON(filepath.Join(data, "updates", "window.json"), &w) == nil && w.PID != priorPID && platform.ProcessRunning(w.PID, w.Executable) {
				controllers = append(controllers, w)
				return w
			}
			time.Sleep(100 * time.Millisecond)
		}
		t.Fatal("new controller marker missing")
		return updateWindow{}
	}
	launch()
	old := waitWindow(0)
	if old.RemoteProtocol != RemoteUpdateProtocol {
		t.Fatal("fixture must reproduce same-protocol legacy reuse")
	}
	raw, err = os.ReadFile(os.Getenv("WORKSPACE_CONTROLLER_CATALOG"))
	if err != nil {
		t.Fatal(err)
	}
	cat, err = catalog.Read(raw)
	if err != nil {
		t.Fatal(err)
	}
	c.Catalog = cat
	c.Version = cat.Version
	if err := Save(c, secret); err != nil {
		t.Fatal(err)
	}
	if err := platform.CopyFile(current, installed); err != nil {
		t.Fatal(err)
	}
	launch()
	fresh := waitWindow(old.PID)
	if platform.ProcessRunning(old.PID, old.Executable) {
		t.Fatal("legacy helper survived upgrade")
	}
	if !controllerMatches(c, fresh) {
		t.Fatal("replacement helper identity mismatch")
	}
	u, err := url.Parse(fresh.URL)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	request, _ := http.NewRequest("GET", "http://"+u.Host+"/api/state", nil)
	request.Header.Set("X-Workspace-Update", u.Fragment)
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var state struct {
		Version string `json:"version"`
		Changed bool   `json:"controllerChanged"`
	}
	if err := json.NewDecoder(response.Body).Decode(&state); err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 200 || state.Version != c.Version || state.Changed {
		t.Fatalf("replacement state: %+v HTTP%d", state, response.StatusCode)
	}
	// Repeated opens reuse this helper rather than replacing it again.
	launch()
	var reused updateWindow
	if err := platform.LoadJSON(filepath.Join(data, "updates", "window.json"), &reused); err != nil || reused.PID != fresh.PID {
		t.Fatal("current helper was unnecessarily replaced", err)
	}
	t.Log("Native same-protocol legacy helper retired; replacement identity/state verified; repeat open reuses current helper.")
}
