//go:build windows

package platform

import (
	"golang.org/x/sys/windows"
	"path/filepath"
	"strings"
)

// Verify the recorded launcher is still running at this installation's path;
// stale progress files must never leave the installer saying it is ready.
func ProcessRunning(pid int, executable string) bool {
	if pid <= 0 {
		return false
	}
	handle, e := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION|windows.SYNCHRONIZE, false, uint32(pid))
	if e != nil {
		return false
	}
	defer windows.CloseHandle(handle)
	state, e := windows.WaitForSingleObject(handle, 0)
	if e != nil || state != uint32(windows.WAIT_TIMEOUT) {
		return false
	}
	path := make([]uint16, 32768)
	size := uint32(len(path))
	if windows.QueryFullProcessImageName(handle, 0, &path[0], &size) != nil {
		return false
	}
	return strings.EqualFold(filepath.Clean(windows.UTF16ToString(path[:size])), filepath.Clean(executable))
}
