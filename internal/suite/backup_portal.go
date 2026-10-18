package suite

import (
	"encoding/json"
	"fmt"
	"github.com/oneworkspacestudio/one-workspace/internal/branding"
	"github.com/oneworkspacestudio/one-workspace/internal/platform"
	"net/http"
	"strings"
)

func backupAppPath(c Config, path string) bool {
	for _, prefix := range []string{"/drive", "/api", "/photos-app", "/photos-static", "/sdkjs", "/fonts"} {
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return true
		}
	}
	for _, a := range c.Catalog.Apps {
		prefix := strings.TrimRight(a.Route, "/")
		if prefix != "" && (path == prefix || strings.HasPrefix(path, prefix+"/")) {
			return true
		}
	}
	return false
}
func (p *Portal) backupRoutes(mux *http.ServeMux) {
	owner := func(w http.ResponseWriter, r *http.Request) bool {
		u, ok := p.user(r)
		if !ok {
			if r.Method == "GET" && r.URL.Path == "/backups" {
				http.Redirect(w, r, "/auth/login?returnTo=%2Fbackups", 303)
			} else {
				http.Error(w, "Sign in to Workspace", 401)
			}
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
	mux.HandleFunc("GET /backups", func(w http.ResponseWriter, r *http.Request) {
		if !owner(w, r) {
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, branding.VersionHTML(backupPage))
	})
	mux.HandleFunc("GET /backups/api/state", func(w http.ResponseWriter, r *http.Request) {
		if !owner(w, r) {
			return
		}
		settings, state := p.rt.backups.snapshot()
		receipts, e := UpdateBackups(p.rt.Config)
		if e != nil {
			http.Error(w, "Could not read backup history. Check the protected updates folder.", 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"settings": settings, "state": state, "backups": receipts, "storage": map[string]string{"programs": p.rt.Config.Programs, "private": p.rt.Config.Data, "files": p.rt.Config.Files}})
	})
	mux.HandleFunc("POST /backups/api/settings", func(w http.ResponseWriter, r *http.Request) {
		if !owner(w, r) {
			return
		}
		var in BackupSettings
		d := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8192))
		d.DisallowUnknownFields()
		if d.Decode(&in) != nil {
			http.Error(w, "Invalid backup setup", 400)
			return
		}
		if e := p.rt.backups.save(in); e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		w.WriteHeader(204)
	})
	mux.HandleFunc("POST /backups/api/check", func(w http.ResponseWriter, r *http.Request) {
		if !owner(w, r) {
			return
		}
		var in struct {
			Destination string `json:"destination"`
		}
		if json.NewDecoder(http.MaxBytesReader(w, r.Body, 8192)).Decode(&in) != nil {
			http.Error(w, "Choose a backup folder", 400)
			return
		}
		destination, e := CheckBackupDestination(p.rt.Config, in.Destination)
		if e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		volume, e := platform.CheckDataVolume(destination)
		if e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		source, e := platform.CheckDataVolume(p.rt.Config.Data)
		if e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"destination": destination, "sameDrive": strings.EqualFold(volume.Root, source.Root), "volume": volume.Root})
	})
	mux.HandleFunc("POST /backups/api/run", func(w http.ResponseWriter, r *http.Request) {
		if !owner(w, r) {
			return
		}
		if e := p.rt.beginBackup(p.rt.rootContext); e != nil {
			http.Error(w, e.Error(), 409)
			return
		}
		w.WriteHeader(202)
	})
	mux.HandleFunc("POST /backups/api/browse", func(w http.ResponseWriter, r *http.Request) {
		if !owner(w, r) {
			return
		}
		if p.rt.backups.busy() {
			http.Error(w, "Wait for the current backup to finish", 409)
			return
		}
		if !p.rt.backups.picker.TryLock() {
			http.Error(w, "The folder picker is already open on the server PC", 409)
			return
		}
		defer p.rt.backups.picker.Unlock()
		var in struct {
			Destination string `json:"destination"`
		}
		if json.NewDecoder(http.MaxBytesReader(w, r.Body, 8192)).Decode(&in) != nil {
			http.Error(w, "Choose a backup folder", 400)
			return
		}
		folder, e := platform.ChooseDataDir(in.Destination)
		if e != nil {
			http.Error(w, "Folder selection cancelled", 400)
			return
		}
		folder, e = CheckBackupDestination(p.rt.Config, folder)
		if e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"destination": folder})
	})
}

const backupPage = `<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Backups · One Workspace</title><link rel="icon" href="/branding/icons/backups.svg"><style>
*{box-sizing:border-box}body{margin:0;background:#f6f8fc;color:#1f2937;font:16px system-ui,sans-serif}header,main{max-width:1000px;margin:auto;padding:24px}header{display:flex;align-items:center;gap:12px}header img{width:42px}a{color:#1967d2;text-decoration:none}header a:last-child{margin-left:auto}h1{font-size:36px;letter-spacing:-1px}h2{font-size:21px;margin-top:0}.muted,p,small{color:#596579;line-height:1.6}.card{background:white;border:1px solid #e1e6ee;border-radius:20px;padding:28px;margin:22px 0}.row{display:flex;gap:12px;flex-wrap:wrap;align-items:center}label{display:block;margin:18px 0 8px}input[type=text],input[type=number]{width:100%;padding:13px;border:1px solid #bdc7d5;border-radius:8px;font:inherit}input[type=checkbox]{width:20px;height:20px;vertical-align:middle;margin-right:9px}button,.button{background:#1967d2;color:white;border:0;border-radius:24px;padding:12px 22px;font:inherit;cursor:pointer}button.secondary{background:#e8f0fe;color:#1967d2}button:disabled{opacity:.5;cursor:wait}.note{background:#eaf1ff;padding:16px;border-radius:12px}.error{background:#fce8e6;color:#a12720;padding:16px;border-radius:12px;white-space:pre-wrap;overflow-wrap:anywhere}.ok{color:#137333}progress{width:100%;height:16px;accent-color:#1967d2}dl{display:grid;grid-template-columns:auto 1fr;gap:10px 22px}dd{margin:0;overflow-wrap:anywhere}li{padding:14px 0;border-bottom:1px solid #edf0f4;overflow-wrap:anywhere}code{overflow-wrap:anywhere;white-space:normal}small{display:block} [hidden]{display:none!important}@media(max-width:600px){.card{padding:20px}dl{display:block}dt{margin-top:12px}}
</style></head><body><header><img src="/branding/icons/backups.svg" alt=""><b>One Workspace / Backups</b><a href="/">All apps</a></header><main><h1>Keep a second copy.</h1><p>Back up accounts, app data and your files together. Scheduling is optional and starts switched off.</p><div id="error" class="error" role="alert" hidden></div>
<section class="card"><h2>Backup setup</h2><p>Choose a folder on another local NTFS drive. Backups contain private account data and use Windows permissions to protect it. FAT, exFAT and network shares are not supported.</p><label for="destination">Backup folder on this Windows PC</label><input id="destination" type="text" placeholder="E:\Workspace backups" autocomplete="off"><small>This is a path on the server, including when you use Workspace remotely.</small><p><button class="secondary" id="browse">Browse on server PC</button> <button class="secondary" id="check">Check folder</button></p><p id="destinationStatus" role="status"></p><label for="interval">Back up every (hours)</label><input id="interval" type="number" min="1" max="168" value="24"><label><input id="enabled" type="checkbox">Enable scheduled backups</label><p class="note">Each backup briefly pauses the apps while copying, then verifies the copy and restarts them. This page stays available. Save your work and close external editors before a manual backup. Scheduled backups run while Workspace and the PC are on; a missed run is picked up after startup.</p><div class="row"><button id="save">Save setup</button><button id="run" class="secondary">Back up now</button></div><p id="saved" role="status"></p><small>Copies are kept until you remove them. Nothing is automatically deleted. If a drive is unavailable or full, the run fails visibly and scheduled retries wait 30 minutes.</small></section>
<section class="card" aria-live="polite"><h2 id="phase">Backup status</h2><progress id="progress" max="100" value="0" hidden></progress><p id="detail"></p><p id="failure" class="error" hidden></p><dl><dt>Schedule</dt><dd id="schedule"></dd><dt>Last verified backup</dt><dd id="last"></dd><dt>Next scheduled run</dt><dd id="next"></dd></dl></section>
<section class="card"><h2>What is protected</h2><dl id="storage"></dl><p>Private settings, databases and runtime storage have their own app folders under Private. Ordinary documents, pictures and music stay in Files. Apps with private content, such as Mail, Notes, My Tube and My Vault, keep it in their protected app data. All these locations are included; temporary files, logs and previous backups are excluded.</p><p>This is a file backup, not a Windows system image. Restoring account secrets requires the original Windows installation and security context. Keep My Vault’s independent encrypted recovery export too.</p></section>
<section class="card"><h2>Verified backup history</h2><p>Only successfully verified snapshots appear here. Incomplete folders are never offered as complete backups.</p><ul id="history"></ul><a class="button" href="/updates">Open recovery tools</a><p>Recovery opens on the Windows server PC. It stops the suite and preserves current data in another archive before restoring. Keep the entire snapshot folder, including backup.json and receipt.json, together.</p></section></main><script>
const $=id=>document.getElementById(id);let loaded=false,busy=false,dirty=false,checking=false;
async function api(path,body){const r=await fetch('/backups/api/'+path,{method:body===undefined?'GET':'POST',headers:{'Content-Type':'application/json'},body:body===undefined?undefined:JSON.stringify(body)});if(r.status===401){location.href='/auth/login?returnTo=%2Fbackups';throw Error('Sign in to continue')}if(!r.ok)throw Error(await r.text());return r.status===202||r.status===204?{}:r.json()}
function error(e){$('error').hidden=false;$('error').textContent=e.message||e}
const date=s=>s&&Date.parse(s)>0?new Date(s).toLocaleString():'—';const mb=n=>(n/1048576).toFixed(1)+' MB';
function setup(){return {enabled:$('enabled').checked,destination:$('destination').value.trim(),intervalHours:Number($('interval').value)}}
for(const id of ['destination','interval','enabled'])$(id).oninput=()=>{dirty=true;$('saved').textContent='Unsaved changes';$('destinationStatus').textContent=''};
$('browse').onclick=async()=>{try{$('browse').disabled=true;$('destinationStatus').textContent='Choose a folder in the picker on the Windows server PC…';const r=await api('browse',{destination:$('destination').value.trim()});$('destination').value=r.destination;dirty=true;$('check').click()}catch(e){error(e)}finally{$('browse').disabled=false}};
$('check').onclick=async()=>{if(checking)return;checking=true;$('check').disabled=true;$('error').hidden=true;try{const s=await api('check',{destination:$('destination').value.trim()});$('destinationStatus').textContent=s.sameDrive?'This is on the same drive as your private data. Another drive is recommended to protect against drive failure.':'Ready: '+s.volume+' supports protected backups.'}catch(e){error(e)}finally{checking=false;$('check').disabled=false}};
async function save(){await api('settings',setup());dirty=false;$('saved').textContent='Setup saved. '+($('enabled').checked?'Scheduling is enabled.':'Scheduling is off.');await poll(false)}
$('save').onclick=async()=>{try{$('error').hidden=true;await save()}catch(e){error(e)}};
$('run').onclick=async()=>{if(busy)return;try{$('error').hidden=true;await save();await api('run',{});await poll(false)}catch(e){error(e)}};
async function poll(repeat=true){try{const r=await api('state'),s=r.state;busy=s.busy;if(!loaded){$('destination').value=r.settings.destination||'';$('interval').value=r.settings.intervalHours;$('enabled').checked=r.settings.enabled;loaded=true}for(const id of ['save','run','destination','interval','enabled','browse'])$(id).disabled=busy;$('phase').textContent=s.phase||'No backup yet';$('detail').textContent=busy?(s.bytesTotal?mb(s.bytesDone)+' of '+mb(s.bytesTotal)+' copied. ':'')+(s.appsPaused?'Apps resume after verification and restart checks.':'Preparing the snapshot. Apps will pause for copying.'):'Backups include all installed apps and ordinary files.';$('progress').hidden=!busy;if(busy&&s.bytesTotal&&s.phase.startsWith('Copying'))$('progress').value=Math.min(100,100*s.bytesDone/s.bytesTotal);else $('progress').removeAttribute('value');$('failure').hidden=!s.error;$('failure').textContent=s.error||'';$('schedule').textContent=r.settings.enabled?'Every '+r.settings.intervalHours+' hours':'Off';$('last').textContent=date(s.lastSuccess);$('next').textContent=r.settings.enabled?date(s.nextRun):'Scheduling is off';$('storage').replaceChildren();for(const [key,value]of Object.entries(r.storage)){const dt=document.createElement('dt'),dd=document.createElement('dd');dt.textContent=key==='private'?'Private app data':key==='files'?'User files':'Programs';dd.textContent=value;$('storage').append(dt,dd)}$('history').replaceChildren();for(const b of r.backups||[]){const li=document.createElement('li');li.textContent=date(b.created)+' · '+b.version+(b.recoveryOnly?' · Recovery archive':'');const path=document.createElement('small');path.textContent=b.path;li.append(path);$('history').append(li)}if(!r.backups?.length)$('history').textContent='No verified snapshots yet.'}catch(e){error(e)}finally{if(repeat)setTimeout(()=>poll(),2000)}}poll();
</script><script defer src="/_workspace/navigation.js"></script></body></html>`
