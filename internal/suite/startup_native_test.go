//go:build windows

package suite

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

// A real native process fixture; never reads an installed app or user database.
func TestStartupFixtureProcess(t *testing.T) {
	if os.Getenv("WORKSPACE_STARTUP_CHILD") != "yes" {
		return
	}
	f, _ := os.OpenFile(os.Getenv("STARTUP_PID_FILE"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	fmt.Fprintln(f, os.Getpid())
	f.Close()
	if os.Getenv("STARTUP_CHILD_MODE") == "tree" {
		exe, _ := os.Executable()
		cmd := exec.Command(exe, "-test.run=^TestStartupFixtureProcess$")
		cmd.Env = append(os.Environ(), "STARTUP_CHILD_MODE=sleeper", "STARTUP_PID_FILE="+os.Getenv("STARTUP_DESCENDANT_FILE"))
		if err := cmd.Start(); err != nil {
			os.Exit(9)
		}
		_ = cmd.Wait()
		os.Exit(0)
	}
	if os.Getenv("STARTUP_CHILD_MODE") == "sleeper" {
		time.Sleep(time.Hour)
		os.Exit(0)
	}
	if p := os.Getenv("STARTUP_FAIL_FILE"); p != "" {
		if _, err := os.Stat(p); err == nil {
			os.Exit(19)
		}
	}
	_ = http.ListenAndServe("127.0.0.1:"+os.Getenv("STARTUP_PORT"), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "ok") }))
	os.Exit(0)
}

func TestNativePartialStartupRetryAndTreeCleanup(t *testing.T) {
	if os.Getenv("WORKSPACE_NATIVE_STARTUP_TEST") != "1" {
		t.Skip("opt-in native processes in a fresh NTFS fixture")
	}
	root := t.TempDir()
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	bytes, err := os.ReadFile(exe)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "fixture.exe")
	if err = os.WriteFile(path, bytes, 0700); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(bytes)
	c := Config{Data: root, Programs: root, Catalog: catalog.Catalog{Packages: map[string]catalog.Package{"fixture": {File: "fixture.exe", Size: int64(len(bytes)), SHA256: hex.EncodeToString(sum[:])}}}}
	makeApp := func(id string, deps ...string) catalog.App {
		l, e := net.Listen("tcp4", "127.0.0.1:0")
		if e != nil {
			t.Fatal(e)
		}
		port := l.Addr().(*net.TCPAddr).Port
		l.Close()
		return catalog.App{ID: id, Name: id, Kind: "process", Package: "fixture", Executable: "fixture.exe", Port: port, Health: "/healthz", Args: []string{"-test.run=^TestStartupFixtureProcess$"}, Depends: deps, Env: map[string]string{"WORKSPACE_STARTUP_CHILD": "yes", "STARTUP_PORT": strconv.Itoa(port), "STARTUP_PID_FILE": filepath.Join(root, id+".pids")}}
	}
	accounts, slicer, later, dependent := makeApp("accounts"), makeApp("slicer", "accounts"), makeApp("later", "accounts"), makeApp("dependent", "slicer")
	failure := filepath.Join(root, "fail-slicer")
	os.WriteFile(failure, []byte("fixture"), 0600)
	slicer.Env["STARTUP_FAIL_FILE"] = failure
	c.Catalog.Apps = []catalog.App{accounts, slicer, later, dependent}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	job, err := platform.NewJob()
	if err != nil {
		t.Fatal(err)
	}
	rt := &Runtime{Config: c, status: map[string]string{}, rootContext: ctx, job: job, progress: NewProgress("starting", StartupSteps(c), "")}
	defer func() {
		cancel()
		if err := rt.stopApps(); err != nil {
			t.Error(err)
		}
	}()
	if err = rt.startApps(ctx, false); err != nil {
		t.Fatal(err)
	}
	if rt.Status()["slicer"] != "failed" || rt.Status()["later"] != "running" || rt.Status()["dependent"] != "blocked" {
		t.Fatal(rt.Status())
	}
	before, _ := os.ReadFile(filepath.Join(root, "later.pids"))
	accountBefore, _ := os.ReadFile(filepath.Join(root, "accounts.pids"))
	os.Remove(failure)
	if err = rt.retryApps(); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(15 * time.Second)
	for rt.progress.Snapshot().State == "retrying" && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	if rt.progress.Snapshot().State != "ready" || rt.Status()["dependent"] != "running" {
		t.Fatal(rt.Status(), rt.progress.Snapshot())
	}
	after, _ := os.ReadFile(filepath.Join(root, "later.pids"))
	accountAfter, _ := os.ReadFile(filepath.Join(root, "accounts.pids"))
	if string(before) != string(after) || string(accountBefore) != string(accountAfter) {
		t.Fatal("retry restarted healthy processes")
	}
	// A timed-out primary with a real descendant must be completely removed
	// before another attempt can use its files or port.
	tree := makeApp("tree")
	tree.Env["STARTUP_CHILD_MODE"] = "tree"
	tree.Env["STARTUP_DESCENDANT_FILE"] = filepath.Join(root, "descendant.pids")
	timeout, stop := context.WithTimeout(ctx, time.Second)
	defer stop()
	if err = rt.start(timeout, tree); err == nil {
		t.Fatal("unhealthy process accepted")
	}
	for _, name := range []string{"tree.pids", "descendant.pids"} {
		pids, e := os.ReadFile(filepath.Join(root, name))
		if e != nil {
			t.Fatal(e)
		}
		for _, line := range strings.Fields(string(pids)) {
			pid, _ := strconv.Atoi(line)
			for platform.ProcessRunning(pid, path) && time.Now().Before(deadline) {
				time.Sleep(20 * time.Millisecond)
			}
			if platform.ProcessRunning(pid, path) {
				t.Fatal("failed attempt left a child running", pid)
			}
		}
	}
	// A second retry is harmless and does not duplicate healthy processes.
	if err = rt.retryApps(); err != nil {
		t.Fatal(err)
	}
	for rt.progress.Snapshot().State == "retrying" && time.Now().Before(deadline) {
		time.Sleep(20 * time.Millisecond)
	}
	again, _ := os.ReadFile(filepath.Join(root, "later.pids"))
	if string(again) != string(before) {
		t.Fatal("duplicate healthy app")
	}
}
