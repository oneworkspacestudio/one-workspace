//go:build linux

package portablesecure

import (
	"runtime"
	"strings"
)

func Executable(name string) string {
	if runtime.GOOS == "windows" {
		return name
	}
	return strings.TrimSuffix(name, ".exe")
}
