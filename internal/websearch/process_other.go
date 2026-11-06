//go:build !windows

package websearch

import "os/exec"

func hideProcess(cmd *exec.Cmd) {}
