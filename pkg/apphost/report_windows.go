//go:build windows

package apphost

import (
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"strings"
	"unsafe"
)

// GUI builds still provide actionable startup errors. Managed/no-browser and
// command-line maintenance invocations report to their existing stderr logs.
func Report(e error) {
	fmt.Fprintln(os.Stderr, e)
	for _, arg := range os.Args[1:] {
		if arg == "--no-browser" || arg == "--worker" || strings.HasPrefix(arg, "--backup") || strings.HasPrefix(arg, "--restore") || arg == "--version" {
			return
		}
	}
	text, _ := windows.UTF16PtrFromString("The app could not start.\n\n" + e.Error() + "\n\nApp data must be in a dedicated NTFS folder. Managed apps are started by One Workspace. See the included README for startup and recovery options.")
	title, _ := windows.UTF16PtrFromString("One Workspace")
	windows.NewLazySystemDLL("user32.dll").NewProc("MessageBoxW").Call(0, uintptr(unsafe.Pointer(text)), uintptr(unsafe.Pointer(title)), 0x10|0x10000)
}
