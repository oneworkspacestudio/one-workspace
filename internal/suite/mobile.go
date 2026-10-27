package suite

import (
	"encoding/json"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
)

// Public connection metadata contains no accounts, secrets, local paths or ports.
// Android obtains the actual account's app list separately after normal sign-in.
func (p *Portal) mobileRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /.well-known/one-workspace.json", func(w http.ResponseWriter, r *http.Request) {
		origins := map[string]bool{}
		add := func(raw string) {
			u, e := url.Parse(raw)
			if e == nil && u.Scheme == "https" && u.Host != "" && u.User == nil {
				origins[u.Scheme+"://"+u.Host] = true
			}
		}
		add(p.rt.Config.Origin)
		add(p.rt.Config.Issuer)
		chat := chatOrigin(p.rt.Config)
		add(chat)
		if _, ok := p.rt.Config.Catalog.Find("vault"); ok {
			var vault struct {
				Origin string `json:"origin"`
			}
			if platform.LoadJSON(filepath.Join(p.rt.Config.Data, "my-vault", "myvault.json"), &vault) == nil {
				u, e := url.Parse(vault.Origin)
				if e == nil && !localVaultHostname(u.Hostname()) {
					add(vault.Origin)
				}
			}
		}
		// App links continue through Workspace's existing authenticated launch routes.
		for _, tile := range p.tiles(false) {
			add(tile.URL)
		}
		list := []string{}
		for origin := range origins {
			list = append(list, origin)
		}
		sort.Strings(list)
		if u, e := url.Parse(chat); e != nil || u.Scheme != "https" {
			chat = ""
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		json.NewEncoder(w).Encode(map[string]any{"schema": 1, "workspace": p.rt.Config.Origin, "issuer": p.rt.Config.Issuer, "chat": chat, "origins": list, "catalog": "/_workspace/apps", "androidPackage": "uk.co.m8e.workspace"})
	})
}
