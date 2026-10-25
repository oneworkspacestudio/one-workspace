//go:build windows

package suite

import "syscall"

func hiddenProcessAttrs() *syscall.SysProcAttr { return &syscall.SysProcAttr{HideWindow: true} }
