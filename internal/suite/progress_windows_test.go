package suite

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/platform"
)

func TestReadyProgressSurvivesWindowsReader(t *testing.T) {
	path := filepath.Join(t.TempDir(), "startup-progress.json")
	p := NewProgress("starting", []ProgressStep{{ID: "portal", Label: "Portal"}}, path)
	p.Step("portal", "running", "Connecting")
	reader, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	// os.Open on Windows does not share delete access. Replacing this file must
	// wait for the updater's ordinary read to finish, rather than lose readiness.
	if err := platform.SaveJSON(path, Progress{State: "probe"}); err == nil {
		reader.Close()
		t.Fatal("expected the held Windows reader to block replacement")
	}
	done := make(chan struct{})
	go func() {
		time.Sleep(70 * time.Millisecond)
		reader.Close()
		close(done)
	}()
	p.Finish("ready", "Your Workspace is running", nil)
	<-done
	var got Progress
	if err := platform.LoadJSON(path, &got); err != nil {
		t.Fatal(err)
	}
	if got.State != "ready" || got.PID != os.Getpid() {
		t.Fatalf("final startup acknowledgement was lost: %+v", got)
	}
}
