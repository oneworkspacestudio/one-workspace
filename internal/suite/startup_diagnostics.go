package suite

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type startupFailure struct {
	code, message string
	cause         error
}

func (e *startupFailure) Error() string { return e.message + ": " + e.cause.Error() }
func (e *startupFailure) Unwrap() error { return e.cause }
func startupStage(code, message string, err error) error {
	if err == nil {
		return nil
	}
	return &startupFailure{code, message, err}
}

// Only fixed explanations and numeric exit codes enter progress returned to a
// browser. Raw provider errors, command lines and logs remain in protected data.
func startupDiagnosis(err error) (string, string) {
	if memoryStartError(err) {
		return "windows-memory", "Windows is low on available memory. Close unused apps, then retry."
	}
	if errors.Is(err, os.ErrPermission) || errors.Is(err, syscall.Errno(5)) {
		return "access-denied", "Windows denied access to an app file or folder. Check Windows Security protection history and folder permissions, then retry."
	}
	if errors.Is(err, os.ErrNotExist) {
		return "file-missing", "A required app file is missing. Check Windows Security protection history and repair Workspace using the latest installer."
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		code := uint32(exit.ExitCode())
		if code == 0xc0000135 || code == 0xc000007b || code == 0xc0000142 {
			return "windows-runtime", fmt.Sprintf("Windows could not load an app runtime (0x%08X). Repair Workspace with the latest installer.", code)
		}
		return "process-exited", fmt.Sprintf("The app closed before it became ready (exit code %d / 0x%08X). Retry; if it happens again, inspect its diagnostic log.", exit.ExitCode(), code)
	}
	var failure *startupFailure
	if errors.As(err, &failure) {
		return failure.code, failure.message
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return "startup-timeout", "The app did not become ready in time. Retry after closing unused apps."
	}
	return "startup-failed", "The app could not finish starting. Retry; if it happens again, inspect its diagnostic log."
}
