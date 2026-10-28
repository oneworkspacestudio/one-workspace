package catalog

// Categories are presentation/setup policy, not serialized launch contracts.
// Older managers strictly decode catalogs, so no new fields enter that format.
type Category struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	DefaultInstall bool   `json:"defaultInstall"`
}

func Categories() []Category {
	return []Category{
		{"general", "General", "Everyday work, communication and creativity", true},
		{"administration", "Administration", "Accounts, protection and running your Workspace", true},
		{"development", "Development", "Source code, AI development and website hosting", false},
	}
}

func CategoryFor(id string) Category {
	index := 0
	switch id {
	case "accounts", "backups", "monitor", "shield", "push", "updates", "app-status":
		index = 1
	case "gitea", "hosting", "myai":
		index = 2
	}
	return Categories()[index]
}

func DefaultInstall(id string) bool { return CategoryFor(id).DefaultInstall }

// WantedApps preserves existing opt-ins/exclusions while adding only new defaults.
// Used both when an older installer fetches a newer catalog and during updates.
func (c Catalog) WantedApps(previous Catalog) []App {
	var out []App
	for _, app := range c.DefaultApps(previous.EffectiveSelection()) {
		_, alreadyWanted := previous.Find(app.ID)
		if alreadyWanted || DefaultInstall(app.ID) {
			out = append(out, app)
		}
	}
	return out
}
