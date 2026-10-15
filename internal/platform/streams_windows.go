//go:build windows

package platform

import (
	"errors"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

type NamedStream struct {
	Name string
	Size int64
}

func EncryptedFileAttributes(attrs uint32) bool { return attrs&windows.FILE_ATTRIBUTE_ENCRYPTED != 0 }

var findFirstStream = windows.NewLazySystemDLL("kernel32.dll").NewProc("FindFirstStreamW")
var findNextStream = windows.NewLazySystemDLL("kernel32.dll").NewProc("FindNextStreamW")

func ValidStreamName(name string) bool { return name != "" && !strings.ContainsAny(name, "/:\\\x00") }

// FindFirstStreamW documents ERROR_INVALID_PARAMETER for filesystems without
// streams, and ERROR_HANDLE_EOF for a directory without named streams.
func NamedStreams(path string) ([]NamedStream, error) {
	p, e := windows.UTF16PtrFromString(path)
	if e != nil {
		return nil, e
	}
	var data struct {
		Size int64
		Name [296]uint16
	}
	handle, _, err := findFirstStream.Call(uintptr(unsafe.Pointer(p)), 0, uintptr(unsafe.Pointer(&data)), 0)
	if handle == uintptr(windows.InvalidHandle) {
		if err == windows.ERROR_HANDLE_EOF || err == windows.ERROR_INVALID_PARAMETER {
			return nil, nil
		}
		return nil, err
	}
	defer windows.FindClose(windows.Handle(handle))
	var streams []NamedStream
	for {
		name := windows.UTF16ToString(data.Name[:])
		if name != "::$DATA" {
			if !strings.HasPrefix(name, ":") || !strings.HasSuffix(name, ":$DATA") {
				return nil, errors.New("unsupported Windows data stream")
			}
			name = strings.TrimSuffix(strings.TrimPrefix(name, ":"), ":$DATA")
			if !ValidStreamName(name) {
				return nil, errors.New("invalid Windows data stream name")
			}
			streams = append(streams, NamedStream{Name: name, Size: data.Size})
		}
		ok, _, err := findNextStream.Call(handle, uintptr(unsafe.Pointer(&data)))
		if ok == 0 {
			if err == windows.ERROR_HANDLE_EOF {
				return streams, nil
			}
			return nil, err
		}
	}
}
func FileAttributes(path string) (uint32, error) {
	p, e := windows.UTF16PtrFromString(path)
	if e != nil {
		return 0, e
	}
	return windows.GetFileAttributes(p)
}
func RestoreFileAttributes(path string, attrs uint32) error {
	p, e := windows.UTF16PtrFromString(path)
	if e != nil {
		return e
	}
	attrs &= windows.FILE_ATTRIBUTE_READONLY | windows.FILE_ATTRIBUTE_HIDDEN | windows.FILE_ATTRIBUTE_SYSTEM | windows.FILE_ATTRIBUTE_ARCHIVE | windows.FILE_ATTRIBUTE_NOT_CONTENT_INDEXED
	if attrs == 0 {
		attrs = windows.FILE_ATTRIBUTE_NORMAL
	}
	return windows.SetFileAttributes(p, attrs)
}
func ClearFileReadOnly(path string) error {
	p, e := windows.UTF16PtrFromString(path)
	if e != nil {
		return e
	}
	attrs, e := windows.GetFileAttributes(p)
	if e != nil {
		return e
	}
	return windows.SetFileAttributes(p, attrs&^windows.FILE_ATTRIBUTE_READONLY)
}
func FileDACL(path string) (string, error) {
	sd, e := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if e == windows.ERROR_NOT_SUPPORTED || e == windows.ERROR_INVALID_FUNCTION {
		return "", nil
	}
	if e != nil {
		return "", e
	}
	return sd.String(), nil
}
func RestoreFileDACL(path, sddl string) error {
	if sddl == "" {
		return nil
	}
	sd, e := windows.SecurityDescriptorFromString(sddl)
	if e != nil {
		return e
	}
	acl, _, e := sd.DACL()
	if e != nil {
		return e
	}
	control, _, e := sd.Control()
	if e != nil {
		return e
	}
	flags := windows.DACL_SECURITY_INFORMATION | windows.UNPROTECTED_DACL_SECURITY_INFORMATION
	if control&windows.SE_DACL_PROTECTED != 0 {
		flags = windows.DACL_SECURITY_INFORMATION | windows.PROTECTED_DACL_SECURITY_INFORMATION
	}
	e = windows.SetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.SECURITY_INFORMATION(flags), nil, nil, acl, nil)
	if e == windows.ERROR_NOT_SUPPORTED || e == windows.ERROR_INVALID_FUNCTION {
		return nil
	}
	return e
}
