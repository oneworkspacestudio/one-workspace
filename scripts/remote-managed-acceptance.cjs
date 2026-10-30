// Full supported installer/update paths. Synthetic identities; never enrolls a Windows host.
const fs=require('node:fs'),path=require('node:path'),crypto=require('node:crypto'),assert=require('node:assert/strict'),{spawn,spawnSync}=require('node:child_process');
const {chromium}=require(process.env.PLAYWRIGHT_MODULE||'playwright');
const root=process.env.WORKSPACE_REMOTE_QA_ROOT,installer=process.env.WORKSPACE_REMOTE_QA_INSTALLER,previous=process.env.WORKSPACE_REMOTE_QA_PREVIOUS,bundle=process.env.WORKSPACE_REMOTE_QA_BUNDLE;
assert(root&&installer&&previous&&bundle);assert(!fs.existsSync(root),'Fresh named root required');fs.mkdirSync(root,{recursive:true});
const password=crypto.randomBytes(24).toString('hex'),checks=[],running=[];let browser;
const env={...process.env,LOCALAPPDATA:path.join(root,'local-app-data')};delete env.GH_TOKEN;delete env.GITHUB_TOKEN;
function command(exe,args,input){const p=spawnSync(exe,args,{input:input&&JSON.stringify(input),env,windowsHide:true,encoding:'utf8',timeout:240000,maxBuffer:1024*1024});if(p.status!==0)throw Error(path.basename(exe)+' failed: '+p.stdout+' '+p.stderr);return p.stdout;}
async function waitHealth(port){for(let n=0;n<180;n++){try{const r=await fetch('http://localhost:'+port+'/healthz');if(r.ok){const v=await r.json();if(Object.values(v.apps).every(x=>x==='running'))return v;}}catch{}await new Promise(r=>setTimeout(r,500));}throw Error('Startup timeout '+port);}
async function verifySetupInventory(){
 const child=spawn(installer,['--offline','--no-browser'],{env,windowsHide:true,stdio:['ignore','pipe','pipe']});let output='';child.stdout.on('data',b=>output+=b.toString());
 try{let match;for(let n=0;n<120;n++){match=output.match(/http:\/\/localhost:[0-9]+\/#[A-Za-z0-9_-]+/);if(match)break;if(child.exitCode!==null)throw Error('Setup closed');await new Promise(r=>setTimeout(r,250));}assert(match,'Setup URL absent');const u=new URL(match[0]),r=await fetch(u.origin+'/api/state',{headers:{Origin:u.origin,'X-Setup-Token':u.hash.slice(1)}});assert(r.ok);const state=await r.json();assert.equal(state.availableApps.length,33);const remote=state.availableApps.find(a=>a.id==='remote');assert(remote&&remote.defaultInstall&&remote.category==='administration');}
 finally{child.kill();}
}
function install(exe,name,port,apps){const r=path.join(root,name);fs.mkdirSync(r);command(exe,['--headless','--offline','--no-browser'],{programs:path.join(r,'Apps'),data:path.join(r,'Private'),files:path.join(r,'Files'),name:'Remote QA owner',email:'remote-owner@example.test',password,confirm:password,port,mode:'local',noShortcuts:true,autostart:false,selectedApps:apps});return r;}
function config(r){return JSON.parse(fs.readFileSync(path.join(r,'Private/workspace.json'),'utf8').replace(/^\uFEFF/,''));}
function start(r){const log=fs.openSync(path.join(r,'launcher.log'),'a');const child=spawn(path.join(r,'Private/Workspace.exe'),['--data',path.join(r,'Private'),'run'],{env,windowsHide:true,stdio:['ignore',log,log]});running.push(r);return child;}
async function stop(r){command(path.join(r,'Private/Workspace.exe'),['--data',path.join(r,'Private'),'stop']);for(let n=0;n<40;n++){try{await fetch(config(r).origin+'/healthz')}catch{return;}await new Promise(r=>setTimeout(r,250));}throw Error('Fixture did not stop');}
async function verify(r,label){const c=config(r),origin=c.origin;await waitHealth(c.port);
 assert(c.catalog.apps.some(a=>a.id==='remote'&&a.version==='0.1.0-preview.5'));assert(!fs.existsSync(path.join(r,'Private/remote/host-secret.json')));
 assert.equal((await fetch(origin+'/remote/')).status,200);assert.equal((await fetch(origin+'/remote/desktop',{headers:{Origin:origin}})).status,401);assert.equal((await fetch(origin+'/remote/desktop',{headers:{Origin:'https://invalid.example'}})).status,403);assert([401,302,303].includes((await fetch(origin+'/remote/api/remote/devices',{redirect:'manual'})).status));
 const context=await browser.newContext({viewport:{width:1366,height:900}}),page=await context.newPage(),errors=[];page.on('pageerror',e=>errors.push(e.message));
 try{
  await page.goto(origin+'/_workspace/open/remote');await page.locator('#username').fill('remote-owner@example.test');await page.getByRole('button',{name:'Next',exact:true}).click();await page.locator('input[type=password]').fill(password);await page.getByRole('button',{name:/Sign in/i}).click();
  await page.getByRole('button',{name:'Set up remote access',exact:true}).waitFor({timeout:30000});assert(page.url().startsWith(origin+'/remote'));
  const me=await(await context.request.get(origin+'/remote/api/me')).json();assert.equal(me.user.role,'owner');const remoteToken=(await context.cookies()).find(c=>c.name==='remote_session')?.value;assert(remoteToken);
  const devices=await(await context.request.get(origin+'/remote/api/remote/devices')).json();assert.deepEqual(devices.devices,[]);
  const apps=await(await context.request.get(origin+'/_workspace/apps')).json();assert(apps.apps.some(a=>a.ID==='remote'&&a.URL==='/_workspace/open/remote'));
  await page.getByRole('button',{name:'Apps',exact:true}).click();await page.getByRole('dialog',{name:'Workspace apps',exact:true}).getByRole('link',{name:'My Remote',exact:true}).waitFor();await page.keyboard.press('Escape');
  assert.equal((await context.request.post(origin+'/remote/_integration/revoke-owner-sessions',{headers:{Origin:origin}})).status(),403);
  const headers={Origin:origin,'X-CSRF-Token':me.csrf,'X-Workspace-Expected-Account':me.user.id};
  assert.equal((await context.request.post(origin+'/remote/api/remote/enroll',{headers,data:{pin:'1'}})).status(),400);
  assert.equal((await context.request.post(origin+'/remote/api/remote/enroll',{headers:{Origin:origin},data:{pin:'654321'}})).status(),403);
  await page.screenshot({path:path.join(root,label+'-desktop.png'),fullPage:true});await page.setViewportSize({width:390,height:844});await page.screenshot({path:path.join(root,label+'-phone.png'),fullPage:true});assert(await page.evaluate(()=>document.documentElement.scrollWidth<=innerWidth));assert.deepEqual(errors,[]);
  // A browser with no Workspace/Remote session may request support, never owner access.
  const guest=await browser.newContext();try{assert.equal((await guest.request.get(origin+'/remote/')).status(),200);assert.equal((await guest.request.post(origin+'/remote/support-api/redeem',{headers:{Origin:origin},data:{code:'invalid'}})).status(),400);}finally{await guest.close()}
  await context.request.post(origin+'/auth/logout',{headers:{Origin:origin}});assert.equal((await context.request.get('http://localhost:'+(c.port+35)+'/remote/api/me',{headers:{Cookie:'remote_session='+remoteToken}})).status(),401);assert(!fs.existsSync(path.join(r,'Private/remote/host-secret.json')));checks.push({path:label,actualInstaller:true,MyAccountsOIDC:true,sharedPicker:true,ownerAPI:true,guestShell:true,unauthorizedEnrollmentDenied:true,idleWithoutHostCredentials:true,desktopAndPhoneNoErrorsOrOverflow:true});
 }catch(e){await page.screenshot({path:path.join(root,label+'-failure.png'),fullPage:true}).catch(()=>{});throw e;}finally{await context.close()}
}
(async()=>{try{
 browser=await chromium.launch({channel:'msedge',headless:true});await verifySetupInventory();
 const fresh=install(installer,'fresh',30200,['accounts','remote']);start(fresh);await verify(fresh,'fresh');await stop(fresh);
 const old=install(previous,'upgrade',30400,['accounts']);const before=config(old);start(old);await waitHealth(before.port);
 const incompatible=spawnSync(previous,['--update',path.join(old,'Private'),'--apply','--apps','workspace','--bundle',bundle,'--no-browser'],{env,windowsHide:true,encoding:'utf8',timeout:60000});assert.notEqual(incompatible.status,0);assert((incompatible.stdout+incompatible.stderr).includes('newer Workspace installer'));assert.deepEqual(config(old),before);
 command(installer,['--update',path.join(old,'Private'),'--apply','--apps','workspace','--bundle',bundle,'--no-browser']);running.push(old);
 const after=config(old);assert.equal(after.version,'0.3.0-preview.104');assert.equal(after.catalog.installerIntegration,4);assert.deepEqual(after.catalog.selection,before.catalog.selection);assert.deepEqual(after.catalog.apps.find(a=>a.id==='accounts'),before.catalog.apps.find(a=>a.id==='accounts'));assert.equal(after.owner,before.owner);
 const tx=JSON.parse(fs.readFileSync(path.join(old,'Private/updates/transaction.json'),'utf8'));assert.equal(tx.phase,'complete');assert(tx.backup);await verify(old,'upgrade');await stop(old);
 fs.writeFileSync(path.join(root,'acceptance.json'),JSON.stringify({schema:1,checks,remoteVersion:'0.1.0-preview.5',freshSetupOffersRemoteDefault:true,priorVersion:'0.3.0-preview.103',backupAndCompletedTransaction:true,oldControllerRejectsNewRuntimeContractBeforeMutation:true,oldAccountAndSelectionRetained:true,ownerDataUsed:false,actualHostEnrollment:false,physicalPhone:false},null,2));console.log('PASS fresh104 installer and protected103-to104 update, MyAccounts/owner API, guest boundary, picker, desktop/mobile');
 }finally{if(browser)await browser.close();for(const r of new Set(running)){try{await stop(r)}catch{}}}})().catch(e=>{console.error(e.stack.replaceAll(password,'[redacted]'));process.exitCode=1});
