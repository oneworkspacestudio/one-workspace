package suite

import (
	"github.com/fahimc/one-workspace/internal/platform"
	"os"
	"path/filepath"
)

func runtimeTemp(c Config, app string) string { return filepath.Join(c.Data, "runtime", app, "temp") }

func appWorkingDir(c Config, app string) string {
	switch app {
	case "accounts", "drive", "photos", "gitea", "music", "mail", "calendar", "contacts", "notes", "my-tube", "chat":
		return filepath.Join(c.Data, app)
	case "vault":
		return filepath.Join(c.Data, "my-vault")
	default:
		return c.Data // Preserve existing third-party extension contracts.
	}
}
func prepareRuntimeTemp(c Config, app string) error {
	working := appWorkingDir(c, app)
	if working != c.Data {
		if _, err := os.Lstat(working); err == nil {
			if err = platform.RejectReparse(working); err != nil {
				return err
			}
		} else if !os.IsNotExist(err) {
			return err
		}
		if err := os.MkdirAll(working, 0700); err != nil {
			return err
		}
	}

	root := runtimeTemp(c, app)
	for current := root; current != c.Data; current = filepath.Dir(current) {
		if _, err := os.Lstat(current); err == nil {
			if err = platform.RejectReparse(current); err != nil {
				return err
			}
		} else if !os.IsNotExist(err) {
			return err
		}
	}
	if err := os.MkdirAll(root, 0700); err != nil {
		return err
	}
	for current := root; current != c.Data; current = filepath.Dir(current) {
		if err := platform.RejectReparse(current); err != nil {
			return err
		}
	}
	return nil
}

// The shared service and manager helpers inherit this protected scratch area.
// Process apps receive their own override in command().
func setWorkspaceTemp(c Config) (func(), error) {
	if err := prepareRuntimeTemp(c, "workspace"); err != nil {
		return nil, err
	}
	oldTemp, tempSet := os.LookupEnv("TEMP")
	oldTmp, tmpSet := os.LookupEnv("TMP")
	restore := func() {
		if tempSet {
			_ = os.Setenv("TEMP", oldTemp)
		} else {
			_ = os.Unsetenv("TEMP")
		}
		if tmpSet {
			_ = os.Setenv("TMP", oldTmp)
		} else {
			_ = os.Unsetenv("TMP")
		}
	}
	if err := os.Setenv("TEMP", runtimeTemp(c, "workspace")); err != nil {
		return nil, err
	}
	if err := os.Setenv("TMP", runtimeTemp(c, "workspace")); err != nil {
		restore()
		return nil, err
	}
	return restore, nil
}
