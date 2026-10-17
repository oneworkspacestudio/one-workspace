//go:build windows

package platform

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
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

// The manager is outside its app kill-on-close job. Maintenance and restarted
// managers must remain outside that app job too. A host (terminal, installer or
// process supervisor) may put the manager itself in a separate job which forbids
// breakaway. Respect that host job rather than requesting a forbidden flag: on
// Windows that flag makes CreateProcess fail with ERROR_ACCESS_DENIED even when
// the executable's ACL permits execution.
var isProcessInJob = windows.NewLazySystemDLL("kernel32.dll").NewProc("IsProcessInJob")

func detachedCreationFlags() (uint32, error) {
	flags := uint32(windows.CREATE_NO_WINDOW)
	var inJob uint32
	if ok, _, e := isProcessInJob.Call(uintptr(windows.CurrentProcess()), 0, uintptr(unsafe.Pointer(&inJob))); ok == 0 {
		return 0, fmt.Errorf("inspect Windows host process group: %w", e)
	}
	if inJob == 0 {
		return flags, nil
	}
	var info windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION
	if e := windows.QueryInformationJobObject(0, windows.JobObjectExtendedLimitInformation, uintptr(unsafe.Pointer(&info)), uint32(unsafe.Sizeof(info)), nil); e != nil {
		return 0, fmt.Errorf("inspect Windows host process group limits: %w", e)
	}
	// A silent-breakaway job already lets ordinary children leave; it does not
	// require the explicit flag. Request it only when explicitly permitted.
	if info.BasicLimitInformation.LimitFlags&windows.JOB_OBJECT_LIMIT_BREAKAWAY_OK != 0 {
		flags |= windows.CREATE_BREAKAWAY_FROM_JOB
	}
	return flags, nil
}

// Only a verified, protected maintenance/runtime executable is launched here.
func StartDetached(cmd *exec.Cmd) error {
	flags, e := detachedCreationFlags()
	if e != nil {
		return e
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: flags}
	return cmd.Start()
}
