//go:build windows

package platform

import (
	"errors"
	"golang.org/x/sys/windows"
	"os"
	"path/filepath"
	"runtime"
	"unsafe"
)

var pickerUser = windows.NewLazySystemDLL("user32.dll")
var pickerInitialized = windows.NewCallback(func(hwnd, message, selection, initial uintptr) uintptr {
	if message != 1 {
		return 0
	} // BFFM_INITIALIZED
	if initial != 0 {
		pickerUser.NewProc("SendMessageW").Call(hwnd, 0x467, 1, initial)
	} // BFFM_SETSELECTIONW
	foreground, _, _ := pickerUser.NewProc("GetForegroundWindow").Call()
	thread, _, _ := pickerUser.NewProc("GetWindowThreadProcessId").Call(foreground, 0)
	current := uintptr(windows.GetCurrentThreadId())
	if thread != 0 && thread != current {
		attached, _, _ := pickerUser.NewProc("AttachThreadInput").Call(current, thread, 1)
		if attached != 0 {
			defer pickerUser.NewProc("AttachThreadInput").Call(current, thread, 0)
		}
	}
	pickerUser.NewProc("ShowWindow").Call(hwnd, 5)
	pickerUser.NewProc("BringWindowToTop").Call(hwnd)
	pickerUser.NewProc("SetForegroundWindow").Call(hwnd)
	pickerUser.NewProc("SetFocus").Call(hwnd)
	return 0
})

// Native Windows folder selection keeps first-run data placement independent of developer runtimes.
func ChooseDataDir(initialPath string) (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	ole := windows.NewLazySystemDLL("ole32.dll")
	coInit := ole.NewProc("CoInitializeEx")
	coUninit := ole.NewProc("CoUninitialize")
	free := ole.NewProc("CoTaskMemFree")
	hr, _, _ := coInit.Call(0, 2)
	if hr == 0 || hr == 1 {
		defer coUninit.Call()
	}
	shell := windows.NewLazySystemDLL("shell32.dll")
	browse := shell.NewProc("SHBrowseForFolderW")
	getPath := shell.NewProc("SHGetPathFromIDListW")
	title, _ := windows.UTF16PtrFromString("Choose your Workspace parent folder on a local NTFS drive (usually C:). All Workspace apps, settings and files will be kept here.")
	owner, _, _ := pickerUser.NewProc("GetForegroundWindow").Call()
	// Start at the current choice, or its closest existing parent.
	for initialPath != "" {
		if info, e := os.Stat(initialPath); e == nil && info.IsDir() {
			break
		}
		next := filepath.Dir(initialPath)
		if next == initialPath {
			initialPath = ""
			break
		}
		initialPath = next
	}
	var initial *uint16
	if initialPath != "" {
		initial, _ = windows.UTF16PtrFromString(initialPath)
	}
	display := make([]uint16, 260)
	info := struct {
		Owner    uintptr
		Root     uintptr
		Display  *uint16
		Title    *uint16
		Flags    uint32
		Callback uintptr
		Param    uintptr
		Image    int32
	}{Owner: owner, Display: &display[0], Title: title, Flags: 0x51, Callback: pickerInitialized, Param: uintptr(unsafe.Pointer(initial))}
	pidl, _, _ := browse.Call(uintptr(unsafe.Pointer(&info)))
	runtime.KeepAlive(initial)
	if pidl == 0 {
		return "", errors.New("no data folder was selected")
	}
	defer free.Call(pidl)
	buffer := make([]uint16, 32768)
	ok, _, _ := getPath.Call(pidl, uintptr(unsafe.Pointer(&buffer[0])))
	if ok == 0 {
		return "", errors.New("selected folder has no filesystem path")
	}
	return windows.UTF16ToString(buffer), nil
}
