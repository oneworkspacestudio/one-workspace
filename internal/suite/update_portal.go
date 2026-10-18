package suite

import (
	"fmt"
	"net/http"

	"github.com/oneworkspacestudio/one-workspace/internal/branding"
)

func (p *Portal) updateRoutes(mux *http.ServeMux) {
	owner := func(w http.ResponseWriter, r *http.Request) bool {
		u, ok := p.user(r)
		if !ok || u.Subject != p.rt.Config.Owner {
			http.Error(w, "Workspace owner access required", 403)
			return false
		}
		return true
	}
	mux.HandleFunc("GET /updates", func(w http.ResponseWriter, r *http.Request) {
		if !owner(w, r) {
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, branding.VersionHTML(updatePortalPage))
	})
	mux.HandleFunc("POST /updates/open", func(w http.ResponseWriter, r *http.Request) {
		if !owner(w, r) {
			return
		}
		if r.Header.Get("Origin") != p.rt.Config.Origin {
			http.Error(w, "Origin rejected", 403)
			return
		}
		if e := LaunchUpdates(p.rt.Config); e != nil {
			http.Error(w, e.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, `<!doctype html><meta charset="utf-8"><title>Workspace updates opened</title><h1>Update window opened on the server PC</h1><p>Continue in that Windows browser window. Update credentials stay on the PC that stores your installation.</p><p><a href="/">Back to Workspace</a></p>`)
	})
}

const updatePortalPage = `<!doctype html><html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Workspace updates</title><style>body{font:17px system-ui;background:#f5f7fc;color:#18233b;max-width:720px;margin:70px auto;padding:24px}main{background:white;padding:32px;border-radius:20px}p{line-height:1.7;color:#5e6c84}button{border:0;background:#345de5;color:white;border-radius:12px;padding:14px 20px;font:inherit;cursor:pointer}a{color:#345de5}</style><main><img src="/branding/icons/one_workspace.svg" width="64" height="64" alt=""><h1>Keep your workspace up to date</h1><p>Check One Workspace and your apps for compatible releases, choose updates and keep verified backups for recovery. The update window runs locally on the Windows PC that holds your installation.</p><form method="post" action="/updates/open"><button>Open update window on this PC</button></form><p>If you are connected remotely, use this control when you are at the server PC.</p><p><a href="/">Back to Workspace</a></p></main></html>`
