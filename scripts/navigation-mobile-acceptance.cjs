// Actual isolated Workspace/MyAccounts/Chat; synthetic users only.
const fs=require('node:fs'),assert=require('node:assert/strict');
const {chromium}=require(process.env.PLAYWRIGHT_MODULE||'playwright');
const root=process.env.CHAT_TEST_PRIVATE_ROOT,origin=process.env.CHAT_TEST_ORIGIN;
assert(root&&/chat-parity-[0-9]+[\\/]fixture$/.test(root)&&/^http:\/\/(localhost|127\.0\.0\.1):[0-9]+$/.test(origin));
const users=JSON.parse(fs.readFileSync(root+'/test-users.json','utf8'));assert(users.every(u=>u.email.endsWith('@example.test')));
const checks=[],errors=[];
(async()=>{
 const browser=await chromium.launch({executablePath:'C:/Program Files/Google/Chrome/Application/chrome.exe',headless:true});let page;
 try{
  const context=await browser.newContext({viewport:{width:390,height:700},isMobile:true,hasTouch:true});page=await context.newPage();page.on('pageerror',e=>errors.push(e.message));
  await page.goto(origin+'/chat/');await page.locator('#username').fill(users[0].username);await page.getByRole('button',{name:'Next',exact:true}).click();await page.locator('input[type=password]').fill(users[0].password);await page.getByRole('button',{name:/Sign in/i}).click();await page.locator('.workspace-mobile-header').waitFor();
  await page.goto(origin+'/auth/login');const opener=page.getByRole('button',{name:'Open Workspace apps',exact:true});await opener.waitFor();
  for(const theme of ['light','dark'])for(const [width,height]of[[320,480],[390,700],[640,360]]){
   await page.setViewportSize({width,height});await page.evaluate(theme=>document.documentElement.dataset.theme=theme,theme);
   await opener.click();const panel=page.getByRole('dialog',{name:'Workspace apps',exact:true}),grid=panel.getByRole('navigation',{name:'Installed apps',exact:true});
   await grid.getByRole('link',{name:'Chat',exact:true}).waitFor();
   await grid.evaluate(grid=>{const source=grid.querySelector('a');for(let i=grid.querySelectorAll('a').length;i<20;i++){const link=source.cloneNode(true);link.querySelector('span').textContent='Synthetic layout '+(i+1);grid.append(link)}});
   assert(await grid.evaluate(n=>n.clientHeight>=44),'Catalog retains usable scroll space');
   const close=panel.getByRole('button',{name:'Close Workspace apps',exact:true}),before=await close.boundingBox();
   assert(before.y>=0&&before.y+before.height<=height&&before.height>=44);
   await grid.evaluate(n=>n.scrollTop=n.scrollHeight);assert.deepEqual(await close.boundingBox(),before);
   await close.click();assert(await opener.evaluate(n=>n===n.getRootNode().activeElement));
   await opener.click();await page.keyboard.press('Escape');await panel.waitFor({state:'hidden'});checks.push({theme,width,height,closeVisible:true,scrollingCatalogOnly:true});
  }
  assert.deepEqual(errors,[]);fs.writeFileSync(root+'/workspace-mobile-navigation-results.json',JSON.stringify({checks,errors,ownerDataUsed:false,realServices:true,twentyEntryStressIsDOMOnly:true,physicalPhone:false},null,2));
  console.log('PASS Shared Workspace picker: six narrow/short/theme cases, usable20-entry scrolling layout, fixed Close and Escape/focus restoration');
 }catch(error){if(page)await page.screenshot({path:root+'/workspace-mobile-navigation-failure.png'}).catch(()=>{});throw error}finally{await browser.close()}
})().catch(error=>{console.error(users.reduce((text,u)=>text.replaceAll(u.password,'[redacted]'),error.message));process.exitCode=1});
