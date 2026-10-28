package suite

import (
	"fmt"

	"github.com/oneworkspacestudio/one-workspace/internal/catalog"
)

// SelectInstallApps removes omitted executables as well as launcher entries.
// A nil selection installs General and Administration; optional apps require opt-in.
func SelectInstallApps(c catalog.Catalog, r Request) (catalog.Catalog, error) {
	selected := map[string]bool{}
	if r.SelectedApps == nil {
		for _, a := range c.Apps {
			selected[a.ID] = catalog.DefaultInstall(a.ID)
		}
	} else {
		for _, id := range r.SelectedApps {
			if _, ok := c.Find(id); !ok {
				return c, fmt.Errorf("app %q is unavailable in this bundle; reopen setup", id)
			}
			selected[id] = true
		}
	}
	if !selected["accounts"] {
		return c, fmt.Errorf("MyAccounts is required for Workspace sign-in")
	}
	apps := []catalog.App{}
	keep := catalog.InstallPackageIDs(c, selected)
	selection := &catalog.AppSelection{}
	for _, a := range c.Apps {
		if !selected[a.ID] {
			selection.Excluded = append(selection.Excluded, a.ID)
			continue
		}
		for _, dep := range a.Depends {
			if !selected[dep] {
				return c, fmt.Errorf("%s requires %s; select both apps or uncheck %s", a.Name, dep, a.Name)
			}
		}
		apps = append(apps, a)
	}
	packages := map[string]catalog.Package{}
	for id, p := range c.Packages {
		// License/provenance text remains available even for omitted components.
		if keep[id] {
			packages[id] = p
		}
	}
	c.Apps, c.Packages, c.Selection = apps, packages, selection
	return c, c.Validate()
}
