package suite

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

func TestStartupContinuesAndRetriesDependencies(t *testing.T) {
	apps := []catalog.App{{ID: "accounts", Name: "MyAccounts", Kind: "process"}, {ID: "slicer", Name: "My Slicer", Kind: "process", Depends: []string{"accounts"}}, {ID: "dependent", Name: "Dependent", Kind: "process", Depends: []string{"slicer"}}, {ID: "tasks", Name: "My Tasks", Kind: "process", Depends: []string{"accounts"}}}
	rt := &Runtime{Config: Config{Data: t.TempDir(), Catalog: catalog.Catalog{Apps: apps}}, status: map[string]string{}, active: map[string]bool{}}
	rt.progress = NewProgress("starting", StartupSteps(rt.Config), "")
	fail := true
	counts := map[string]int{}
	start := func(ctx context.Context, app catalog.App) error {
		counts[app.ID]++
		if app.ID == "slicer" && fail {
			return fmt.Errorf("fork/exec private-path: %w", syscall.Errno(1455))
		}
		rt.set(app.ID, "running")
		rt.active[app.ID] = true
		rt.progress.Step(app.ID, "done", "App health check passed")
		return nil
	}
	if err := rt.startAvailable(context.Background(), apps, start); err != nil {
		t.Fatal(err)
	}
	if counts["tasks"] != 1 || counts["dependent"] != 0 || rt.Status()["slicer"] != "failed" || rt.Status()["dependent"] != "blocked" {
		t.Fatal("independent app or dependency handling failed", counts, rt.Status())
	}
	p := rt.progress.Snapshot()
	for _, s := range p.Steps {
		if s.ID == "slicer" && (!strings.Contains(s.Detail, "Close unused apps") || strings.Contains(s.Detail, "private-path")) {
			t.Fatal(s)
		}
	}
	fail = false
	if err := rt.startAvailable(context.Background(), apps, start); err != nil {
		t.Fatal(err)
	}
	if counts["accounts"] != 1 || counts["tasks"] != 1 || counts["slicer"] != 2 || counts["dependent"] != 1 {
		t.Fatal("retry restarted healthy apps or missed dependencies", counts)
	}
}

func TestAccountsFailureRemainsFatalAndCancellationStopsPass(t *testing.T) {
	apps := []catalog.App{{ID: "accounts", Name: "MyAccounts", Kind: "process"}, {ID: "tasks", Kind: "process"}}
	rt := &Runtime{Config: Config{Data: t.TempDir()}, status: map[string]string{}, progress: NewProgress("starting", nil, "")}
	count := 0
	start := func(context.Context, catalog.App) error { count++; return errors.New("private failure") }
	if err := rt.startAvailable(context.Background(), apps, start); err == nil || strings.Contains(err.Error(), "private failure") {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal(count)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if !errors.Is(rt.startAvailable(ctx, apps, start), context.Canceled) || count != 1 {
		t.Fatal("cancelled startup launched an app")
	}
}

func TestMemoryStartRetryBoundedAndCancellable(t *testing.T) {
	for _, tc := range []struct {
		name               string
		err                error
		successAfter, want int
	}{{"transient", syscall.Errno(1455), 1, 2}, {"persistent", syscall.Errno(1455), 99, 3}, {"not memory", os.ErrPermission, 99, 1}} {
		t.Run(tc.name, func(t *testing.T) {
			n := 0
			err := startWithMemoryRetry(context.Background(), func() error {
				n++
				if n > tc.successAfter {
					return nil
				}
				return fmt.Errorf("wrapped: %w", tc.err)
			}, func() {})
			if n != tc.want || (tc.name == "transient" && err != nil) {
				t.Fatal(n, err)
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	n := 0
	err := startWithMemoryRetry(ctx, func() error { n++; return syscall.Errno(1455) }, cancel)
	if !errors.Is(err, context.Canceled) || n != 1 {
		t.Fatal(n, err)
	}
}

func TestRetryRoutesRejectUnauthorizedAndBusy(t *testing.T) {
	rt := &Runtime{Config: Config{Origin: "https://workspace.example", Owner: "owner"}, status: map[string]string{}, progress: NewProgress("ready", nil, "")}
	p := &Portal{rt: rt, sessions: map[string]session{}}
	mux := http.NewServeMux()
	p.startupRoutes(mux, "private-control")
	for _, tc := range []struct {
		path, token, remote string
		want                int
	}{{"/_control/retry-apps", "", "127.0.0.1:42", 403}, {"/_control/retry-apps", "private-control", "203.0.113.1:42", 403}, {"/_control/retry-apps", "private-control", "127.0.0.1:42", 409}, {"/_workspace/startup/retry", "", "127.0.0.1:42", 401}} {
		r := httptest.NewRequest("POST", tc.path, nil)
		r.RemoteAddr = tc.remote
		r.Header.Set("Authorization", "Bearer "+tc.token)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != tc.want {
			t.Fatalf("%s: %d", tc.path, w.Code)
		}
	}
	for _, tc := range []struct {
		subject, origin string
		want            int
	}{{"other", "https://workspace.example", 403}, {"owner", "https://evil.test", 403}, {"owner", "https://workspace.example", 409}} {
		p.sessions["fixture"] = session{Subject: tc.subject, Expires: time.Now().Add(time.Minute)}
		r := httptest.NewRequest("POST", "/_workspace/startup/retry", nil)
		r.AddCookie(&http.Cookie{Name: "workspace_session", Value: "fixture"})
		r.Header.Set("Origin", tc.origin)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)
		if w.Code != tc.want {
			t.Fatalf("owner %s origin %s: %d", tc.subject, tc.origin, w.Code)
		}
	}
	rt.operations.Lock()
	if err := rt.retryApps(); err == nil {
		t.Error("retry ignored maintenance lock")
	}
	rt.operations.Unlock()
}

func TestProgressPartialAndRecoveryClearStaleError(t *testing.T) {
	p := NewProgress("starting", []ProgressStep{{ID: "app"}, {ID: "dependent"}, {ID: "portal"}}, "")
	p.Step("app", "failed", "Retry")
	p.Step("dependent", "blocked", "Waiting")
	p.Step("portal", "done", "Ready")
	p.Finish("failed", "Interrupted", errors.New("old"))
	p.Finish("ready", "Ready", nil)
	s := p.Snapshot()
	if s.Error != "" || s.Percent != 100 || s.Steps[0].State != "failed" {
		t.Fatal(s)
	}
}
