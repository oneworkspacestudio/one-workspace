//go:build windows

package apphost

import (
	"context"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Worker reuses the supervisor's suspended-start job boundary and scrubs secrets.
func Worker(ctx context.Context, exe, dir string, input io.Reader, output io.Writer) (result error) {
	if e := ctx.Err(); e != nil {
		return e
	}
	job, e := platform.NewJob()
	if e != nil {
		return e
	}
	defer func() {
		stop, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if e := job.Stop(stop); e != nil && result == nil {
			result = e
		}
		job.Close()
	}()
	if e = job.WorkerLimits(512 << 20); e != nil {
		return e
	}
	cmd := exec.Command(exe, "--worker")
	cmd.Dir = dir
	cmd.Env = []string{"SystemRoot=" + os.Getenv("SystemRoot"), "WINDIR=" + os.Getenv("WINDIR"), "TEMP=" + dir, "TMP=" + dir, "USERPROFILE=" + dir, "APPDATA=" + filepath.Join(dir, "roaming"), "GOMAXPROCS=2"}
	cmd.Stdin = input
	cmd.Stdout = output
	cmd.Stderr = io.Discard
	if e = job.Start(cmd); e != nil {
		return e
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case e = <-done:
		return e
	case <-ctx.Done():
		stop, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if e = job.Stop(stop); e != nil {
			return e
		}
		<-done
		return ctx.Err()
	}
}
