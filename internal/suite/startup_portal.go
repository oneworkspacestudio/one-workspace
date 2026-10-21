package suite

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/oneworkspacestudio/one-workspace/internal/branding"
)

func (p *Portal) startupRoutes(mux *http.ServeMux, control string) {
	mux.HandleFunc("POST /_control/retry-apps", func(w http.ResponseWriter, r *http.Request) {
		host, _, _ := net.SplitHostPort(r.RemoteAddr)
		ip := net.ParseIP(host)
		if ip == nil || !ip.IsLoopback() || control == "" || subtle.ConstantTimeCompare([]byte(r.Header.Get("Authorization")), []byte("Bearer "+control)) != 1 {
			http.Error(w, "Forbidden", 403)
			return
		}
		if err := p.rt.retryApps(); err != nil {
			http.Error(w, err.Error(), 409)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	})
	owner := func(w http.ResponseWriter, r *http.Request) bool {
		u, ok := p.user(r)
		if !ok {
			http.Error(w, "Sign in to Workspace", 401)
			return false
		}
		if u.Subject != p.rt.Config.Owner {
			http.Error(w, "Workspace owner access required", 403)
			return false
		}
		if r.Method != "GET" && r.Header.Get("Origin") != p.rt.Config.Origin {
			http.Error(w, "Origin rejected", 403)
			return false
		}
		w.Header().Set("Cache-Control", "no-store")
		return true
	}
	mux.HandleFunc("GET /_workspace/startup", func(w http.ResponseWriter, r *http.Request) {
		if !owner(w, r) {
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, branding.VersionHTML(startupPage))
	})
	mux.HandleFunc("GET /_workspace/startup/state", func(w http.ResponseWriter, r *http.Request) {
		if !owner(w, r) {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p.rt.progress.Snapshot())
	})
	mux.HandleFunc("POST /_workspace/startup/retry", func(w http.ResponseWriter, r *http.Request) {
		if !owner(w, r) {
			return
		}
		if err := p.rt.retryApps(); err != nil {
			http.Error(w, err.Error(), 409)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	})
}

const startupPage = `<!doctype html><html lang="en"><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>App status · One Workspace</title><style>body{font:16px/1.6 system-ui,sans-serif;background:#f4f6fc;color:#17233b;margin:0}main{max-width:760px;margin:36px auto;padding:28px;background:white;border-radius:20px}img{width:48px}a{color:#345de5}button{font:inherit;border:0;padding:12px 20px;background:#345de5;color:white;border-radius:10px;cursor:pointer}button:disabled{opacity:.55}ul{list-style:none;padding:0}li{padding:14px 0;border-bottom:1px solid #e7ecf5}small{display:block;color:#60708a;overflow-wrap:anywhere}.failed,.blocked{color:#a42636}#error{color:#a42636}@media(max-width:800px){main{margin:12px;padding:22px}}</style><main><img src="/branding/icons/one_workspace.svg" alt=""><h1>Your apps</h1><p>Keep using the apps that are ready. If an app needs attention, you can retry it here without restarting the others.</p><p id="phase" role="status" aria-live="polite"></p><ul id="steps"></ul><p id="error" role="alert"></p><button id="retry" onclick="retry()" disabled>Retry apps that need attention</button><p><a href="/">Back to Workspace →</a></p></main><script>const $=id=>document.getElementById(id);let timer;async function poll(){clearTimeout(timer);try{const r=await fetch('/_workspace/startup/state');if(!r.ok)throw Error('Sign in as the Workspace owner to check your apps.');const p=await r.json();$('phase').textContent=p.phase;$('steps').replaceChildren();for(const s of p.steps){const li=document.createElement('li'),b=document.createElement('b'),d=document.createElement('small');li.className=s.state;b.textContent=(({done:'✓ ',warning:'! ',failed:'× ',blocked:'· ',running:'… '})[s.state]||'· ')+s.label;d.textContent=s.detail||'Waiting';li.append(b,d);$('steps').append(li)}$('retry').disabled=p.state==='retrying'||!p.steps.some(s=>['failed','blocked'].includes(s.state))}catch(e){$('error').textContent=e.message}timer=setTimeout(poll,2000)}async function retry(){$('retry').disabled=true;$('error').textContent='';try{const r=await fetch('/_workspace/startup/retry',{method:'POST'});if(!r.ok)throw Error(await r.text())}catch(e){$('error').textContent=e.message}poll()}poll()</script></html>`
