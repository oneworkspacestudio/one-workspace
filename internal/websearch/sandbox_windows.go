//go:build windows

package websearch

import (
	"context"
	"golang.org/x/sys/windows"
	"os/exec"
	"path/filepath"
	"time"
)

func allowSandboxCode(path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	root, e := windows.GetWindowsDirectory()
	if e != nil {
		return e
	}
	cmd := exec.CommandContext(ctx, filepath.Join(root, "System32", "icacls.exe"), path, "/grant", "*S-1-15-2-1:(OI)(CI)(RX)", "/T", "/Q")
	hideProcess(cmd)
	return cmd.Run()
}
