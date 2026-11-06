//go:build windows

package websearch

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

var processJob windows.Handle

// All portable browser/model descendants die when the service exits, including
// forced manager stops and crashes. Nothing is left as a detached host service.
func OwnChildren() error {
	job, e := windows.CreateJobObject(nil, nil)
	if e != nil {
		return e
	}
	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{}
	info.BasicLimitInformation.LimitFlags = windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE
	if _, e = windows.SetInformationJobObject(job, windows.JobObjectExtendedLimitInformation, uintptr(unsafe.Pointer(&info)), uint32(unsafe.Sizeof(info))); e != nil {
		windows.CloseHandle(job)
		return e
	}
	if e = windows.AssignProcessToJobObject(job, windows.CurrentProcess()); e != nil {
		windows.CloseHandle(job)
		return e
	}
	processJob = job
	return nil
}
