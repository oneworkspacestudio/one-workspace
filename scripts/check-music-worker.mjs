import fs from 'node:fs';
import vm from 'node:vm';
import assert from 'node:assert/strict';

const source = fs.readFileSync(new URL('../internal/suite/assets/music-worker-navigation.js', import.meta.url), 'utf8');
const handlers = new Map();
let offline = false;
let cacheReads = 0;
const context = {
  URL, Response,
  self: {
    location: new URL('https://workspace.example.test/music/sw.js'),
    registration: {scope: 'https://workspace.example.test/music/'},
    addEventListener(type, handler) { handlers.set(type, handler); },
    skipWaiting: async () => {}, clients: {claim: async () => {}}
  },
  fetch: async () => {if (offline) throw new Error('fixture offline'); return new Response('network');},
  caches: {match: async (url, options) => {cacheReads++; assert.equal(url, 'https://workspace.example.test/music/index.html'); assert.equal(options.ignoreSearch, true); return new Response('cached music shell');}}
};
vm.runInNewContext(source, context);
async function navigate(path, mode='navigate') {
  let stopped=false; let response;
  handlers.get('fetch')({request:{url:'https://workspace.example.test'+path,mode},stopImmediatePropagation(){stopped=true;},respondWith(value){assert.equal(response,undefined);response=value;}});
  return {stopped,response:response && await response};
}
for (const path of ['/music/','/music/auth/myaccounts/login','/music/auth/myaccounts/callback']) {
  const result=await navigate(path);assert.equal(result.stopped,true);assert.equal(await result.response.text(),'network');
}
for (const [path,mode] of [['/notes/','navigate'],['/music/assets/app.js','cors'],['/drive/','navigate']]) {
  const result=await navigate(path,mode);assert.equal(result.stopped,false);assert.equal(result.response,undefined);
}
offline=true;
let result=await navigate('/music/');assert.equal(await result.response.text(),'cached music shell');assert.equal(cacheReads,1);
result=await navigate('/music/auth/myaccounts/callback');assert.equal(result.response.status,503);assert.equal(cacheReads,1);
console.log('Music worker checks passed: live authentication/navigation, scoped asset handling and safe offline shell fallback.');
