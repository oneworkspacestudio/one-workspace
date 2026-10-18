//go:build windows

package apphost

import (
	"errors"
	"golang.org/x/sys/windows"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func validateOwner(root string) error {
	raw, e := os.ReadFile(filepath.Join(root, ".owner-sid"))
	if os.IsNotExist(e) {
		return nil
	}
	if e != nil {
		return e
	}
	current, e := user.Current()
	if e != nil {
		return e
	}
	owner := strings.TrimSpace(string(raw))
	if owner == current.Uid || owner == "S-1-5-18" {
		return nil
	}
	if current.Uid != "S-1-5-18" {
		return errors.New("this app data belongs to another Windows account; restore an encrypted backup into a new folder")
	}
	sid, e := windows.StringToSid(owner)
	if e != nil {
		return errors.New("invalid data owner marker")
	}
	_, _, kind, e := sid.LookupAccount("")
	if e != nil || kind != windows.SidTypeUser {
		return errors.New("data owner marker must identify a Windows user")
	}
	return nil
}
