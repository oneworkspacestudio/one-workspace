package suite

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fahimc/one-workspace/internal/catalog"
)

// SelectInstallApps removes omitted executables as well as launcher entries.
// A nil selection retains the legacy Photos switch for older setup clients.
func SelectInstallApps(c catalog.Catalog, r Request) (catalog.Catalog, error) {
	selected := map[string]bool{}
	if r.SelectedApps == nil {
		for _, a := range c.Apps {
			selected[a.ID] = a.ID != "photos" || r.Photos
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
	keep := map[string]bool{"manager": true, "cloudflared": true, "license": true, "notice": true}
	shared := false
	for _, a := range c.Apps {
		if !selected[a.ID] {
			continue
		}
		for _, dep := range a.Depends {
			if !selected[dep] {
				return c, fmt.Errorf("%s requires %s; select both apps or uncheck %s", a.Name, dep, a.Name)
			}
		}
		apps = append(apps, a)
		keep[a.Package] = true
		if a.ID == "studio" {
			keep["studio"] = true
		}
		if sharedConsumer(a.ID) && !(a.ID == "mail" && a.Env["WORKSPACE_MAIL_BACKEND"] == "private") {
			shared = true
		}
	}
	keep["stalwart-shared"], keep["services-shared"] = shared, shared
	packages := map[string]catalog.Package{}
	for id, p := range c.Packages {
		// License/provenance text remains available even for omitted components.
		if keep[id] || strings.HasPrefix(filepath.ToSlash(p.File), "licenses/") || (selected["vault"] && strings.HasPrefix(id, "vault-")) {
			packages[id] = p
		}
	}
	c.Apps, c.Packages = apps, packages
	return c, c.Validate()
}
