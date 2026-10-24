package catalog

import (
	"path/filepath"
	"strings"
)

// EffectiveSelection migrates installations which predate recorded opt-outs.
// Preserve missing apps from the pre-Hosting catalog; Hosting and future app IDs
// are new defaults. Old managers discarded this distinction on every update.
func (c Catalog) EffectiveSelection() *AppSelection {
	s := &AppSelection{}
	if c.Selection != nil {
		s.Excluded = append([]string(nil), c.Selection.Excluded...)
		return s
	}
	for _, id := range []string{"drive", "docs", "photos", "gitea", "music", "mail", "studio", "notes", "my-tube", "vault", "calendar", "contacts", "chat"} {
		if _, installed := c.Find(id); !installed {
			s.Excluded = append(s.Excluded, id)
		}
	}
	return s
}

func (s *AppSelection) Excludes(id string) bool {
	if s != nil {
		for _, excluded := range s.Excluded {
			if excluded == id {
				return true
			}
		}
	}
	return false
}

// DefaultApps closes exclusions over dependencies rather than re-enabling an
// app the user unchecked. The caller preserves those exclusions for later runs.
func (c Catalog) DefaultApps(s *AppSelection) []App {
	allowed := map[string]bool{}
	for _, a := range c.Apps {
		allowed[a.ID] = !s.Excludes(a.ID)
	}
	for changed := true; changed; {
		changed = false
		for _, a := range c.Apps {
			if !allowed[a.ID] {
				continue
			}
			for _, dep := range a.Depends {
				if !allowed[dep] {
					allowed[a.ID] = false
					changed = true
					break
				}
			}
		}
	}
	var apps []App
	for _, a := range c.Apps {
		if allowed[a.ID] {
			apps = append(apps, a)
		}
	}
	return apps
}

// InstallPackageIDs is shared by fresh setup and new-app updates. In particular,
// a coupled release must never extract an explicitly omitted executable.
func InstallPackageIDs(c Catalog, selected map[string]bool) map[string]bool {
	keep := map[string]bool{"manager": true, "cloudflared": true, "license": true, "notice": true}
	for _, a := range c.Apps {
		if !selected[a.ID] {
			continue
		}
		keep[a.Package] = true
		if a.ID == "mail" {
			for id := range c.Packages {
				if strings.HasPrefix(id, "msvc-") {
					keep[id] = true
				}
			}
		}
		if a.ID == "studio" {
			keep["studio"] = true
		}
		if a.ID == "calendar" || a.ID == "contacts" || (a.ID == "mail" && a.Env["WORKSPACE_MAIL_BACKEND"] != "private") {
			keep["stalwart-shared"], keep["services-shared"] = true, true
			for id := range c.Packages {
				if strings.HasPrefix(id, "msvc-") {
					keep[id] = true
				}
			}
		}
	}
	for id, p := range c.Packages {
		if strings.HasPrefix(filepath.ToSlash(p.File), "licenses/") || (selected["vault"] && strings.HasPrefix(id, "vault-")) || (selected["cinema"] && strings.HasPrefix(id, "cinema-")) {
			keep[id] = true
		}
	}
	return keep
}
