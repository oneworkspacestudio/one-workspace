package branding

import (
	"embed"
	"net/http"
)

//go:embed icons/*.svg
var icons embed.FS

func Handler() http.Handler {
	return http.StripPrefix("/branding/", http.FileServer(http.FS(icons)))
}

func Icon(app string) string {
	name := map[string]string{
		"calendar": "my_calendar", "contacts": "my_contacts", "drive": "my_drive", "docs": "my_office", "office": "my_office",
		"photos": "my_photos", "gitea": "my_git", "git": "my_git",
		"music": "my_music", "studio": "photo_editor", "photo-editor": "photo_editor",
		"mail": "my_mail", "accounts": "my_accounts", "myaccounts": "my_accounts", "notes": "my_notes",
	}[app]
	if name == "" {
		name = "one_workspace"
	}
	return "/branding/icons/" + name + ".svg"
}
