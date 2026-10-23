package maintenance

// Gateway failures are connection state, not operation failures. Never render a
// proxy's HTML response or retry a mutation whose acknowledgement was lost.
const requestScript = `
function connectionError(){return Object.assign(new Error('Workspace is temporarily unreachable.'),{reconnect:true})}
function responseMessage(text,type,status){
 let message=text.trim();
 if(type.includes('json')){try{const value=JSON.parse(message);message=typeof value.error==='string'?value.error:typeof value.msg==='string'?value.msg:''}catch{message=''}}
 if(!message||type.includes('html')||/<[!/?a-z]/i.test(message))return 'The update request could not be completed (HTTP '+status+').';
 return message.replace(/[\u0000-\u001f\u007f]+/g,' ').slice(0,400);
}
async function api(path,body){
 const controller=new AbortController(),timeout=setTimeout(()=>controller.abort(),body?45000:20000);
 try{
  let r;try{r=await fetch((remote?'/updates/api/':'/api/')+path,{method:body?'POST':'GET',credentials:'same-origin',cache:'no-store',redirect:'error',signal:controller.signal,headers:{'X-Workspace-Update':token||'','Content-Type':'application/json'},body:body?JSON.stringify(body):undefined})}catch{throw connectionError()}
  if(remote&&r.status===401){location.replace('/auth/login?returnTo=%2Fupdates');throw Error('Sign in as the Workspace owner to continue')}
  if(r.status===408||r.status===429||r.status>=500)throw connectionError();
  if(!r.ok)throw Error(responseMessage(await r.text(),r.headers.get('Content-Type')||'',r.status));
  if(r.status===202)return null;
  let value;try{value=await r.json()}catch{throw connectionError()}
  if(path==='state'&&(!value||typeof value.busy!=='boolean'||typeof value.pending!=='boolean'||typeof value.phase!=='string'||typeof value.version!=='string'))throw connectionError();
  return value;
 }finally{clearTimeout(timeout)}
}
function connectionLost(uncertain=false){
 disconnected=true;awaitingResult=awaitingResult||uncertain;
 $('connection').hidden=false;
 $('connection').textContent='Connection to '+(remote?'Workspace':'the updater')+' interrupted. Reconnecting automatically. '+(awaitingResult?'The request may have been received; waiting for its status. Do not start it again.':'Your last received progress is kept below.');
 buttons();renderProgress();
}
`
