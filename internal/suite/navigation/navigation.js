(() => {
 if (window.top !== window || document.getElementById('one-workspace-nav')) return;
 fetch('/_workspace/apps', {credentials:'same-origin'}).then(r => r.ok ? r.json() : null).then(catalog => {
  if (!catalog) return;
  // Refresh an already-installed Music worker after suite upgrades. Old offline
  // shells can otherwise intercept sign-in before the new client ever loads.
  if ('serviceWorker' in navigator) {
   const music=catalog.apps.find(app=>app.Name==='Music');
   if (music) {
    const scope=new URL('/music/',location.origin).href;
    navigator.serviceWorker.getRegistration(scope)
     .then(registration=>registration?.scope===scope ? registration.update() : undefined).catch(()=>{});
   }
  }
  const host = document.createElement('div'); host.id='one-workspace-nav';
  const shadow=host.attachShadow({mode:'open'}), style=document.createElement('link');
  style.rel='stylesheet'; style.href='/_workspace/navigation.css'; shadow.append(style);
  const toggle=document.createElement('button'); toggle.textContent='Workspace apps'; toggle.setAttribute('aria-expanded','false'); toggle.setAttribute('aria-label','Open Workspace apps');
  const panel=document.createElement('section'); panel.hidden=true; panel.setAttribute('aria-label','Workspace apps');
  const name=document.createElement('strong'); name.textContent=catalog.name || catalog.email || 'Your workspace'; panel.append(name);
  const searchApp=catalog.apps.find(a=>a.Name==='My Search');if(searchApp){const search=document.createElement('form');search.action='/search/';search.method='get';search.setAttribute('role','search');const input=document.createElement('input');input.type='search';input.name='q';input.placeholder='Search Workspace';input.setAttribute('aria-label','Search Workspace');input.maxLength=200;search.append(input);panel.append(search);document.addEventListener('keydown',e=>{if((e.ctrlKey||e.metaKey)&&e.key.toLowerCase()==='k'){e.preventDefault();panel.hidden=false;toggle.setAttribute('aria-expanded','true');input.focus()}})}
  const home=document.createElement('a'); home.href='/'; home.textContent='One Workspace home'; panel.append(home);
  for(const app of catalog.apps){const a=document.createElement('a');a.href=app.URL;if(app.Icon){const img=document.createElement('img');img.src=app.Icon;img.alt='';img.width=24;img.height=24;a.append(img)}a.append(document.createTextNode(app.Name));panel.append(a)}
  const form=document.createElement('form');form.action='/auth/logout';form.method='post';const logout=document.createElement('button');logout.textContent='Sign out of Workspace apps';form.append(logout);panel.append(form);
  toggle.onclick=()=>{panel.hidden=!panel.hidden;toggle.setAttribute('aria-expanded',String(!panel.hidden))};
  document.addEventListener('keydown',e=>{if(e.key==='Escape'){panel.hidden=true;toggle.setAttribute('aria-expanded','false')}});
  shadow.append(panel,toggle);document.body.append(host);
 }).catch(()=>{});
})();
