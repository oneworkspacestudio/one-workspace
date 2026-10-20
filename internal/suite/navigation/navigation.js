(() => {
 if (window.top !== window || document.getElementById('one-workspace-nav')) return;
 if(location.pathname === '/') {
  const apply = value => {document.documentElement.dataset.theme=value==='dark'?'dark':'light'};
  try{apply(localStorage.getItem('workspace-theme'))}catch{apply(null)}
  window.addEventListener('storage',e=>{if(e.key==='workspace-theme')apply(e.newValue)});
  const appearance=document.createElement('link');appearance.rel='stylesheet';appearance.href='/_workspace/appearance.css';document.head.append(appearance);
 }
 fetch('/_workspace/apps', {credentials:'same-origin',cache:'no-store'}).then(r => r.ok ? r.json() : null).then(catalog => {
  if (!catalog || !Array.isArray(catalog.apps)) return;
  if ('serviceWorker' in navigator && catalog.apps.some(app => app.ID === 'music')) {
   const scope = new URL('/music/', location.origin).href;
   navigator.serviceWorker.getRegistration(scope).then(r => r?.scope === scope ? r.update() : undefined).catch(() => {});
  }
  const host = document.createElement('div'); host.id = 'one-workspace-nav';
  const shadow = host.attachShadow({mode:'open'}), style = document.createElement('link');
  style.rel = 'stylesheet'; style.href = '/_workspace/navigation.css'; shadow.append(style);
  const make = (tag, text, className) => {const el = document.createElement(tag); if(text) el.textContent = text; if(className) el.className = className; return el};
  const icon = (url, size=24) => {const img = make('img'); img.src=url; img.alt=''; img.width=size; img.height=size; return img};
  const toggle = make('button', '', 'launcher'); toggle.type='button';
  toggle.append(icon('/branding/icons/one_workspace.svg'), make('span','Workspace'));
  toggle.setAttribute('aria-expanded','false'); toggle.setAttribute('aria-label','Open Workspace apps'); toggle.setAttribute('aria-haspopup','dialog');
  const panel = make('dialog'); panel.setAttribute('aria-label','Workspace apps');
  const heading = make('header'), home = make('a','One Workspace'); home.href='/'; heading.append(home);
  const close = make('button','×','close'); close.type='button'; close.setAttribute('aria-label','Close Workspace apps'); heading.append(close); panel.append(heading);
  panel.append(make('p',catalog.name || catalog.email || 'Your workspace','account'));
  const current = make('p','','current'); panel.append(current);
  let searchInput;
  if(catalog.apps.some(a => a.ID === 'search')) {
   const search = make('form','','search'); search.action='/search/'; search.method='get'; search.setAttribute('role','search');
   searchInput=make('input'); searchInput.type='search'; searchInput.name='q'; searchInput.placeholder='Search your Workspace'; searchInput.setAttribute('aria-label','Search Workspace'); searchInput.maxLength=200;
   const submit=make('button','Search'); submit.type='submit'; search.append(searchInput,submit); panel.append(search);
  }
  const links = make('nav','','apps'); links.setAttribute('aria-label','Installed apps'); const entries=[];
  for(const app of catalog.apps) {
   const a=make('a'); a.href=app.URL; a.title=app.Description || app.Name;
   if(app.Icon) a.append(icon(app.Icon,32)); a.append(make('span',app.Name));
   links.append(a); entries.push({app,a});
  }
  panel.append(links);
  const footer=make('footer'), drive=catalog.apps.find(a=>a.ID==='drive');
  if(drive){const files=make('a','Your files in Drive'); files.href=drive.URL; footer.append(files)}
  const form=make('form'); form.action='/auth/logout'; form.method='post'; const logout=make('button','Sign out of Workspace'); logout.type='submit'; form.append(logout); footer.append(form); panel.append(footer);
  const triggerSelector='[data-workspace-apps-trigger]';
  const triggers=()=>Array.from(document.querySelectorAll(triggerSelector));
  const visible=el=>{const r=el.getBoundingClientRect();return r.width>0&&r.height>0&&r.bottom>0&&r.right>0&&r.top<innerHeight&&r.left<innerWidth&&getComputedStyle(el).visibility==='visible'};
  const expanded=value=>{toggle.setAttribute('aria-expanded',value);for(const el of triggers())el.setAttribute('aria-expanded',value)};
  let opener;
  const open = (focusSearch=false,trigger) => {if(!panel.open){opener=trigger || shadow.activeElement || document.activeElement; panel.showModal()} expanded('true'); (focusSearch && searchInput ? searchInput : close).focus()};
  toggle.onclick=()=>open(); close.onclick=()=>panel.close();
  panel.addEventListener('close',()=>{expanded('false');const target=opener?.isConnected&&visible(opener)?opener:triggers().find(visible)||toggle;target.focus()});
  panel.addEventListener('click',e=>{if(e.target===panel){const r=panel.getBoundingClientRect();if(e.clientX<r.left||e.clientX>r.right||e.clientY<r.top||e.clientY>r.bottom)panel.close()}});
  window.addEventListener('workspace:open-apps',()=>open());
  // Apps may provide a normal navigation link as an accessible picker trigger.
  // Keep its href as a fallback, and hide only the floating button, never the dialog.
  const syncTriggers=()=>{
   const links=triggers();
   for(const el of links){el.setAttribute('role','button');el.setAttribute('aria-haspopup','dialog');el.setAttribute('aria-expanded',String(panel.open))}
   toggle.hidden=links.some(visible);
  };
  let triggerFrame;
  const scheduleTriggers=()=>{if(!triggerFrame)triggerFrame=requestAnimationFrame(()=>{triggerFrame=0;syncTriggers()})};
  new MutationObserver(scheduleTriggers).observe(document.body,{childList:true,subtree:true,attributes:true,attributeFilter:['class','style','hidden','data-workspace-apps-trigger']});
  window.addEventListener('resize',scheduleTriggers);
  document.addEventListener('scroll',scheduleTriggers,true);
  document.addEventListener('click',e=>{
   const trigger=e.target.closest?.(triggerSelector);
   if(!trigger||e.defaultPrevented||e.button!==0||e.ctrlKey||e.metaKey||e.altKey||e.shiftKey)return;
   e.preventDefault();open(false,trigger);
  });
  document.addEventListener('keydown',e=>{
   const trigger=e.target.closest?.(triggerSelector);
   if(trigger&&!e.defaultPrevented&&e.key===' '){e.preventDefault();open(false,trigger)}
  });
  document.addEventListener('keydown',e=>{if(!e.defaultPrevented && !e.isComposing && (e.ctrlKey||e.metaKey) && e.key.toLowerCase()==='k' && searchInput){e.preventDefault();open(true)}});
  const decode = path => {try{return decodeURIComponent(path).replace(/\/$/,'')}catch{return path}};
  let currentApp;
  const title = () => {
   let next=document.title;
   if(currentApp?.ID==='music') next=next.replace(/Swing Music$/, 'Music · One Workspace');
   if(currentApp?.ID==='photos' && next==='Photos') next='Photos · One Workspace';
   if(currentApp?.ID==='studio' && next==='Workspace Drive - Studio') next='Photo Editor · One Workspace';
   if(next!==document.title) document.title=next;
  };
  const active = () => {
   const path=decode(location.pathname), matched=entries.filter(({app})=>app.ID!=='accounts' && app.Route && decode(app.Route) && (path===decode(app.Route)||path.startsWith(decode(app.Route)+'/'))).sort((a,b)=>b.app.Route.length-a.app.Route.length)[0];
   for(const {a} of entries)a.removeAttribute('aria-current');
   if(matched)matched.a.setAttribute('aria-current','page');
   current.textContent=matched ? 'You’re in '+matched.app.Name : 'Your apps, together';
   toggle.title=matched ? 'Workspace apps · '+matched.app.Name : 'Workspace apps';
   currentApp=matched?.app; title();
  };
  new MutationObserver(title).observe(document.head,{childList:true,subtree:true,characterData:true});
  const theme=()=>{host.dataset.theme=document.documentElement.dataset.theme || (document.documentElement.classList.contains('dark-mode')||document.documentElement.classList.contains('dark')||document.body.classList.contains('dark')?'dark':'light')};
  new MutationObserver(theme).observe(document.documentElement,{attributes:true,attributeFilter:['data-theme','class']});
  new MutationObserver(theme).observe(document.body,{attributes:true,attributeFilter:['class']});
  window.addEventListener('popstate',active); window.addEventListener('hashchange',active);
  // Vue/React routers can change the path without popstate or a document load.
  let lastPath=location.pathname;
  setInterval(()=>{if(!document.hidden&&location.pathname!==lastPath){lastPath=location.pathname;active()}},1000);
  shadow.append(panel,toggle); document.body.append(host); active(); theme(); syncTriggers();
  document.documentElement.dataset.workspaceNavigation='ready';
 }).catch(()=>{});
})();
