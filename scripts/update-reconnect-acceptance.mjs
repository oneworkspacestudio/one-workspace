// NODE_PATH must provide @playwright/test. All API responses and data are synthetic.
import {createRequire} from 'node:module';
import {readFile,mkdir,writeFile} from 'node:fs/promises';
import {fileURLToPath} from 'node:url';
import path from 'node:path';
const {chromium,expect}=createRequire(import.meta.url)('@playwright/test');
const root=path.resolve(path.dirname(fileURLToPath(import.meta.url)),'..'),out=process.env.WORKSPACE_QA_OUTPUT;
if(!out)throw Error('Set WORKSPACE_QA_OUTPUT to a task-owned artifact folder');
const raw=async file=>(await readFile(path.join(root,'internal/maintenance',file),'utf8')).match(/const \w+ = `([\s\S]*)`/)[1];
let html=process.env.WORKSPACE_QA_PAGE_URL?await(await fetch(process.env.WORKSPACE_QA_PAGE_URL)).text():await raw('page.go');
if(!process.env.WORKSPACE_QA_PAGE_URL)for(const [key,file] of [['requestScript','request_page.go'],['progressScript','progress_page.go']])html=html.replace('` + '+key+' + `',await raw(file));
await mkdir(out,{recursive:true});
const browser=await chromium.launch({channel:'msedge',headless:true}),checks=[],errors=[];
try{
 const context=await browser.newContext({viewport:{width:1424,height:940}}),page=await context.newPage();page.on('pageerror',e=>errors.push(e.message));
 let mode='idle',posts=0,reject=false,stalled=[];
 const now=new Date().toISOString();
 const state=()=>({busy:mode==='running',pending:false,version:'0.3.0-preview.test',origin:'http://127.0.0.1:39130',phase:mode==='complete'?'Update complete. All apps are healthy.':'Ready for a synthetic update.',error:'',plan:{choices:[{id:'workspace',name:'Workspace',current:'test.1',version:'test.2',available:true,source:'https://example.invalid/release'}]},backups:[],progress:['running','complete'].includes(mode)?{kind:'apply',status:mode,detail:mode==='complete'?'Update complete. All apps are healthy.':'Copying synthetic files',lastActivityAt:now,elapsedSeconds:5,remainingSeconds:10,steps:[{state:mode==='complete'?'done':'running',label:'Backup'}],history:[],done:50,total:100,unit:'files'}:null});
 await page.route('http://127.0.0.1:39130/**',async route=>{
  const url=new URL(route.request().url());
  if(url.pathname==='/updates')return route.fulfill({contentType:'text/html',body:html});
  if(url.pathname==='/updates/api/state'){
   if(mode==='stall'){stalled.push(route);return}
   if(mode==='html'||mode==='html200')return route.fulfill({status:mode==='html'?502:200,contentType:'text/html',body:'<!doctype html><html><title>Cloudflare Tunnel error</title><body>NEVER DISPLAY GATEWAY HTML</body></html>'});
   if(mode==='invalid')return route.fulfill({contentType:'application/json',body:'{"unexpected":true}'});
   return route.fulfill({contentType:'application/json',body:JSON.stringify(state())});
  }
  if(route.request().method()==='POST'){
   posts++;
   if(reject)return route.fulfill({status:409,contentType:'text/plain',body:'The selected release changed. Check for updates again.'});
   mode='html';return route.abort('connectionreset');
  }
  return route.fulfill({status:404,body:''});
 });
 await page.goto('http://127.0.0.1:39130/updates');await expect(page.locator('#apply')).toBeEnabled();
 await page.locator('#apply').click();await expect(page.locator('#connection')).toContainText('may have been received');
 await expect(page.locator('#apply')).toBeDisabled();await expect(page.locator('#check')).toBeDisabled();await expect(page.locator('#backup')).toBeDisabled();
 expect(posts).toBe(1);expect(await page.locator('body').innerText()).not.toContain('NEVER DISPLAY');
 mode='running';await expect(page.locator('#current-detail')).toHaveText('Copying synthetic files',{timeout:6000});await expect(page.locator('#connection')).toBeHidden();
 mode='html';await expect(page.locator('#connection')).toBeVisible();await expect(page.locator('#current-detail')).toHaveText('Copying synthetic files');expect(await page.locator('#error').innerText()).toBe('');
 await page.screenshot({path:path.join(out,'desktop-reconnecting.png')});await page.setViewportSize({width:390,height:844});await page.screenshot({path:path.join(out,'phone-reconnecting.png')});expect(await page.evaluate(()=>document.documentElement.scrollWidth<=innerWidth)).toBeTruthy();
 checks.push('502 gateway HTML becomes a short reconnect notice; last progress remains visible on desktop and phone');
 mode='complete';await expect(page.locator('#operation-title')).toHaveText('Update complete',{timeout:6000});await expect(page.locator('#connection')).toBeHidden();expect(posts).toBe(1);checks.push('Lost install acknowledgement is not replayed; controls stay disabled until fresh authoritative state');
 mode='idle';await expect(page.locator('#apply')).toBeEnabled();reject=true;await page.locator('#apply').click();await expect(page.locator('#error')).toContainText('selected release changed');await page.waitForTimeout(1700);await expect(page.locator('#error')).toContainText('selected release changed');checks.push('Real HTTP409 action errors survive successful background polls');
 mode='html200';await expect(page.locator('#connection')).toBeVisible();expect(await page.locator('body').innerText()).not.toContain('NEVER DISPLAY');mode='idle';await expect(page.locator('#connection')).toBeHidden();
 mode='invalid';await expect(page.locator('#connection')).toBeVisible();await expect(page.locator('#installed')).toHaveText('Installed Workspace 0.3.0-preview.test');mode='idle';await expect(page.locator('#connection')).toBeHidden();checks.push('HTTP200 HTML and malformed state never erase the last valid UI');
 mode='stall';await expect.poll(()=>stalled.length).toBeGreaterThan(0);await expect(page.locator('#connection')).toBeVisible({timeout:25000});mode='idle';for(const route of stalled)await route.abort().catch(()=>{});await expect(page.locator('#connection')).toBeHidden({timeout:6000});checks.push('Stalled polls time out and polling resumes');
 expect(posts).toBe(2);expect(errors).toEqual([]);await writeFile(path.join(out,'browser-acceptance.json'),JSON.stringify({checks,errors,syntheticOnly:true,nativePage:!!process.env.WORKSPACE_QA_PAGE_URL},null,2)+'\n');console.log('Update reconnection browser checks passed');
}finally{await browser.close()}
