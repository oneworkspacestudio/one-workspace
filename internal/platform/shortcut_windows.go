package platform

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// A Start Menu shortcut starts the manager in the ACL-protected data directory.
// Values travel as environment variables, never interpolated into shell code.
func Shortcut(exe, data string) error {
	path := filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "One Workspace.lnk")
	code := `$ErrorActionPreference='Stop'; $shell=New-Object -ComObject WScript.Shell; $link=$shell.CreateShortcut($env:WORKSPACE_SHORTCUT); $link.TargetPath=$env:WORKSPACE_MANAGER; $link.WorkingDirectory=$env:WORKSPACE_DATA; $link.Description='Open One Workspace'; $link.Save()`
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-Command", code)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Env = append(os.Environ(), "WORKSPACE_SHORTCUT="+path, "WORKSPACE_MANAGER="+exe, "WORKSPACE_DATA="+data)
	return cmd.Run()
}
