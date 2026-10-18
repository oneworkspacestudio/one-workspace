package branding

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"fmt"
	"net/http"
	"strings"
	"time"
)

//go:embed icons/*.svg
var icons embed.FS

type asset struct {
	body      []byte
	etag      string
	immutable bool
}

var assets, iconURLs, htmlURLs = loadIcons()

func loadIcons() (map[string]asset, map[string]string, *strings.Replacer) {
	files, err := icons.ReadDir("icons")
	if err != nil {
		panic(err)
	}
	served := map[string]asset{}
	urls := map[string]string{}
	var replacements []string
	for _, file := range files {
		body, err := icons.ReadFile("icons/" + file.Name())
		if err != nil {
			panic(err)
		}
		sum := sha256.Sum256(body)
		name := strings.TrimSuffix(file.Name(), ".svg")
		legacy := "/branding/icons/" + file.Name()
		versioned := fmt.Sprintf("/branding/icons/%s-%x.svg", name, sum[:8])
		etag := fmt.Sprintf(`"%x"`, sum)
		served[legacy] = asset{body, etag, false}
		served[versioned] = asset{body, etag, true}
		urls[name] = versioned
		replacements = append(replacements, legacy, versioned)
	}
	return served, urls, strings.NewReplacer(replacements...)
}

// VersionHTML gives static page logos and favicons the same content-addressed
// URLs as app tiles. A changed SVG must never reuse a browser/CDN cache entry.
func VersionHTML(page string) string { return htmlURLs.Replace(page) }

func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a, ok := assets[r.URL.Path]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("ETag", a.etag)
		if a.immutable {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache")
		}
		http.ServeContent(w, r, r.URL.Path, time.Time{}, bytes.NewReader(a.body))
	})
}

func Icon(app string) string {
	name := map[string]string{
		"backups": "backups", "calendar": "my_calendar", "contacts": "my_contacts", "drive": "my_drive", "docs": "my_office", "office": "my_office",
		"photos": "my_photos", "gitea": "my_git", "git": "my_git",
		"music": "my_music", "studio": "photo_editor", "photo-editor": "photo_editor",
		"hosting": "my_hosting", "chat": "my_chat", "mail": "my_mail", "accounts": "my_accounts", "myaccounts": "my_accounts", "notes": "my_notes", "my-tube": "my_tube", "vault": "my_vault",
	}[app]
	if name == "" {
		name = "one_workspace"
	}
	return iconURLs[name]
}
