package suite

import (
	"errors"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"os"
	"path/filepath"
)

// RemoveApp opts out of future automatic installation. Only the four independent
// executable contracts are supported here; shared engines need their own policy.
// It deliberately preserves the app's private data, OIDC client and backups so a
// later explicit reinstall reconnects the same encrypted/user-owned state.
func RemoveApp(data, id string) error {
	if !foundationApp(id) {
		return errors.New("safe removal supports Workspace foundation apps, including My Social")
	}
	lock, e := platform.Lock(filepath.Join(data, "runtime.lock"))
	if e != nil {
		return errors.New("stop Workspace before removing an app")
	}
	defer lock.Close()
	c, _, e := Load(data)
	if e != nil {
		return e
	}
	a, ok := c.Catalog.Find(id)
	if !ok {
		return errors.New("app is not installed")
	}
	for _, other := range c.Catalog.Apps {
		for _, dep := range other.Depends {
			if dep == id {
				return fmt.Errorf("%s depends on this app", other.Name)
			}
		}
	}
	pkg, ok := c.Catalog.Packages[a.Package]
	if !ok || pkg.File != a.Executable {
		return errors.New("app executable contract does not match")
	}
	path, e := platform.Inside(c.Programs, pkg.File)
	if e != nil {
		return e
	}
	if e = platform.RejectReparse(path); e != nil && !os.IsNotExist(e) {
		return e
	}
	if _, e = os.Stat(path); e == nil {
		held, e := platform.LockedBinary(path, pkg.SHA256)
		if e != nil {
			return errors.New("app executable is changed or still running; stop it before removal")
		}
		held.Close()
	} else if !os.IsNotExist(e) {
		return e
	}
	old, e := os.ReadFile(filepath.Join(data, "workspace.json"))
	if e != nil {
		return e
	}
	selection := c.Catalog.EffectiveSelection()
	if !selection.Excludes(id) {
		selection.Excluded = append(selection.Excluded, id)
	}
	c.Catalog.Selection = selection
	for i, v := range c.Catalog.Apps {
		if v.ID == id {
			c.Catalog.Apps = append(c.Catalog.Apps[:i], c.Catalog.Apps[i+1:]...)
			break
		}
	}
	delete(c.Catalog.Packages, a.Package)
	if e = c.Catalog.Validate(); e != nil {
		return e
	}
	if e = platform.SaveJSON(filepath.Join(data, "workspace.json"), c); e != nil {
		return e
	}
	if e = os.Remove(path); e != nil && !os.IsNotExist(e) {
		if restore := platform.AtomicWrite(filepath.Join(data, "workspace.json"), old); restore != nil {
			return fmt.Errorf("remove failed and config needs recovery: %w", restore)
		}
		return e
	}
	return nil
}
