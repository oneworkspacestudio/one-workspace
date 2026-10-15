package suite

import (
	"os/exec"
	"syscall"
)

func execJunction(path, target string) error {
	// Static command, with paths as distinct arguments. No filesystem enumeration
	// crosses shells and no deletion is delegated to cmd.
	cmd := exec.Command("cmd.exe", "/d", "/c", "mklink", "/J", path, target)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}
