//go:build windows

package platform

import (
	"bytes"
	"context"
	"errors"
	"golang.org/x/sys/windows"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

// A real nested Windows job reproduces the updater failure. No installation,
// credentials or app data are involved; all children are fixture-owned.
func TestStartDetachedUnderRestrictedHostJob(t *testing.T) {
	const roleKey = "WORKSPACE_DETACHED_TEST_ROLE"
	switch os.Getenv(roleKey) {
	case "leaf":
		return
	case "worker":
		if os.Getenv("WORKSPACE_DETACHED_TEST_SILENT") != "1" {
			original := exec.Command(os.Args[0], "-test.run=^TestStartDetachedUnderRestrictedHostJob$")
			original.Env = append(os.Environ(), roleKey+"=leaf")
			original.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: windows.CREATE_NO_WINDOW | windows.CREATE_BREAKAWAY_FROM_JOB}
			if err := original.Run(); !errors.Is(err, windows.ERROR_ACCESS_DENIED) {
				t.Fatalf("original launch did not reproduce access denied: %v", err)
			}
		}
		flags, err := detachedCreationFlags()
		if err != nil {
			t.Fatal(err)
		}
		if flags&windows.CREATE_BREAKAWAY_FROM_JOB != 0 {
			t.Fatal("requested breakaway from a job that forbids it")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestStartDetachedUnderRestrictedHostJob$")
		cmd.Env = append(os.Environ(), roleKey+"=leaf")
		if err := StartDetached(cmd); err != nil {
			t.Fatal("restricted host prevented child launch", err)
		}
		if err := cmd.Wait(); err != nil {
			t.Fatal(err)
		}
		return
	}
	job, err := NewJob()
	if err != nil {
		t.Fatal(err)
	}
	defer job.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	defer job.Stop(context.Background())
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestStartDetachedUnderRestrictedHostJob$")
	cmd.Env = append(os.Environ(), roleKey+"=worker")
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	if err := job.Start(cmd); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		t.Fatalf("fixture worker failed: %v\n%s", err, output.String())
	}
	if !strings.Contains(output.String(), "PASS") {
		t.Fatalf("fixture did not complete: %s", output.String())
	}
}

func TestStartDetachedUnderSilentBreakawayHost(t *testing.T) {
	job, err := NewJob()
	if err != nil {
		t.Fatal(err)
	}
	defer job.Close()
	var info windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION
	info.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE | windows.JOB_OBJECT_LIMIT_SILENT_BREAKAWAY_OK
	if _, err := windows.SetInformationJobObject(job.handle, windows.JobObjectExtendedLimitInformation, uintptr(unsafe.Pointer(&info)), uint32(unsafe.Sizeof(info))); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	defer job.Stop(context.Background())
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestStartDetachedUnderRestrictedHostJob$")
	cmd.Env = append(os.Environ(), "WORKSPACE_DETACHED_TEST_ROLE=worker", "WORKSPACE_DETACHED_TEST_SILENT=1")
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	if err := job.Start(cmd); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		t.Fatalf("silent-breakaway fixture failed: %v\n%s", err, output.String())
	}
}
