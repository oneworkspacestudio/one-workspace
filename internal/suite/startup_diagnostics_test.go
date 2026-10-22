package suite

import (
	"context"
	"errors"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestStartupDiagnosisDoesNotExposeRawProviderErrors(t *testing.T) {
	for _, tc := range []struct {
		err  error
		code string
	}{
		{startupStage("chat-sign-in", "Chat could not connect to MyAccounts.", errors.New("secret=private-provider-token")), "chat-sign-in"},
		{startupStage("shared-backend", "The shared background service needs attention.", errors.New("private configuration")), "shared-backend"},
		{syscall.Errno(1455), "windows-memory"}, {os.ErrPermission, "access-denied"}, {os.ErrNotExist, "file-missing"},
	} {
		code, message := startupDiagnosis(tc.err)
		if code != tc.code || strings.Contains(message, "private") || strings.Contains(message, "secret=") {
			t.Fatal(code, message)
		}
	}
}

func TestSharedReadinessDetectsEarlyExitAndExistingOwner(t *testing.T) {
	done := make(chan error, 1)
	done <- errors.New("child failed")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	start := time.Now()
	if err := waitSharedReady(ctx, func() bool { return false }, done); err == nil || !strings.Contains(err.Error(), "startup.log") || time.Since(start) > time.Second/2 {
		t.Fatal(err)
	}
	done <- errors.New("another manager already owns it")
	if err := waitSharedReady(ctx, func() bool { return true }, done); err != nil {
		t.Fatal(err)
	}
	cancelled, stop := context.WithCancel(context.Background())
	stop()
	if err := waitSharedReady(cancelled, func() bool { return false }, nil); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
