package suite

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"syscall"
)

func execJunction(path, target string) error {
	// Static command, with paths as distinct arguments. No filesystem enumeration
	// crosses shells and no deletion is delegated to cmd.
	cmd := exec.Command("cmd.exe", "/d", "/c", "mklink", "/J", filepath.FromSlash(path), filepath.FromSlash(target))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, e := cmd.CombinedOutput()
	if e != nil {
		return fmt.Errorf("%w: %s", e, out)
	}
	return nil
}
