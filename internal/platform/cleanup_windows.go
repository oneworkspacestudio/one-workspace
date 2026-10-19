//go:build windows

package platform

import (
	"errors"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Refuse cleanup while an executable in the candidate tree is still running.
// Protected system processes cannot execute from owner-only maintenance folders.
func DirectoryHasProcesses(root string) (bool, error) {
	snap, e := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if e != nil {
		return false, e
	}
	defer windows.CloseHandle(snap)
	var p windows.ProcessEntry32
	p.Size = uint32(unsafe.Sizeof(p))
	e = windows.Process32First(snap, &p)
	for e == nil {
		h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, p.ProcessID)
		if err == nil {
			buf := make([]uint16, 32768)
			n := uint32(len(buf))
			err = windows.QueryFullProcessImageName(h, 0, &buf[0], &n)
			windows.CloseHandle(h)
			if err == nil {
				path := filepath.Clean(windows.UTF16ToString(buf[:n]))
				rel, err := filepath.Rel(root, path)
				if err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && !filepath.IsAbs(rel) {
					return true, nil
				}
			}
		}
		e = windows.Process32Next(snap, &p)
	}
	if !errors.Is(e, windows.ERROR_NO_MORE_FILES) {
		return false, e
	}
	return false, nil
}
