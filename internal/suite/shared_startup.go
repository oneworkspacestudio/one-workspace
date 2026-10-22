package suite

import (
	"context"
	"fmt"
	"github.com/oneworkspacestudio/calendar/shared"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"os/exec"
	"path/filepath"
	"time"
)

func ensureSharedReady(ctx context.Context, root, manager string) error {
	if shared.Healthy(ctx, root) {
		return nil
	}
	// Capture the service's bootstrap errors, which shared.Ensure previously
	// discarded. Observe an early exit instead of waiting out two silent minutes.
	log := &platform.Log{Path: filepath.Join(root, "startup.log")}
	var cmd *exec.Cmd
	err := startWithMemoryRetry(ctx, func() error {
		cmd = exec.Command(manager, "run", "--root", root)
		cmd.Dir = filepath.Dir(manager)
		cmd.Stdout = log
		cmd.Stderr = log
		shared.Hide(cmd)
		return cmd.Start()
	}, func() {})
	if err != nil {
		return err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	return waitSharedReady(ctx, func() bool { return shared.Healthy(ctx, root) }, done)
}

func waitSharedReady(ctx context.Context, healthy func() bool, done <-chan error) error {
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()
	for {
		if healthy() {
			return nil
		}
		select {
		case err := <-done:
			// Another registered manager may have won the service mutex in the meantime.
			if healthy() {
				return nil
			}
			if err == nil {
				return fmt.Errorf("shared background service exited before readiness; inspect shared/startup.log")
			}
			return fmt.Errorf("shared background service exited; inspect shared/startup.log: %w", err)
		case <-ctx.Done():
			return fmt.Errorf("shared background service did not become ready; inspect shared/startup.log and shared/engine.log: %w", ctx.Err())
		case <-ticker.C:
		}
	}
}
