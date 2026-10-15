//go:build windows

package platform

import (
	"errors"
	"golang.org/x/sys/windows"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
)

// Protect only the directory boundary; existing maintenance lock handles must
// remain open. New children inherit this DACL. Retain the installation owner's
// access even when a different local administrator performs maintenance.
func ProtectMaintenanceDir(dir, data string) error {
	if _, e := CheckDataVolume(dir); e != nil {
		return e
	}
	if e := RejectReparse(dir); e != nil {
		return e
	}
	u, e := user.Current()
	if e != nil {
		return e
	}
	owner, e := os.ReadFile(filepath.Join(data, ".owner-sid"))
	if e != nil {
		return e
	}
	sid := strings.TrimSpace(string(owner))
	if _, e = windows.StringToSid(sid); e != nil {
		return e
	}
	sd, e := windows.SecurityDescriptorFromString("D:P(A;OICI;FA;;;" + u.Uid + ")(A;OICI;FA;;;" + sid + ")(A;OICI;FA;;;SY)(A;OICI;FA;;;BA)")
	if e != nil {
		return e
	}
	acl, _, e := sd.DACL()
	if e != nil {
		return e
	}
	return windows.SetNamedSecurityInfo(dir, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION, nil, nil, acl, nil)
}

func RejectReparse(path string) error {
	p, e := windows.UTF16PtrFromString(path)
	if e != nil {
		return e
	}
	flags, e := windows.GetFileAttributes(p)
	if e != nil {
		return e
	}
	if flags&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0 {
		return errors.New("updates do not follow symbolic links or junctions: " + path)
	}
	return nil
}
func FreeBytes(path string) (uint64, error) {
	p, e := windows.UTF16PtrFromString(path)
	if e != nil {
		return 0, e
	}
	var available, total, free uint64
	e = windows.GetDiskFreeSpaceEx(p, &available, &total, &free)
	return available, e
}

// Maintenance must outlive the manager's kill-on-close job. Only a verified,
// protected copy of the current manager is launched through this function.
func StartDetached(cmd *exec.Cmd) error {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: windows.CREATE_NO_WINDOW | windows.CREATE_BREAKAWAY_FROM_JOB}
	return cmd.Start()
}
