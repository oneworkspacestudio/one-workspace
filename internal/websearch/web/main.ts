import {createElement as h,useEffect,useRef,useState,createContext,useContext} from 'react';
import {createRoot} from 'react-dom/client';
import {Console} from './generated/Console';
import {RemoteScreen,RemoteInput} from './RemoteScreen';
import './style.css';

export const Fields=createContext<any>({fields:{},setField:()=>{},item:{current:null},action:()=>{}});
const tag=(name:string)=>(p:any)=>h(name,{className:p.variant||undefined,id:p.id},p.children);
export const Page=tag('main'),Row=tag('div'),Card=tag('section'),Grid=tag('div'),Text=(p:any)=>h(p.variant==='pre'?'pre':'p',{},p.children),Badge=tag('strong'),Spacer=tag('span');
export const Heading=(p:any)=>h('h'+(p.level||2),{},p.children);
export const Alert=(p:any)=>p.children?h('p',{role:'status',className:'notice'},p.children):null;
export const Link=(p:any)=>h('a',{href:p.href,'data-workspace-apps-trigger':p.variant==='apps'?'true':undefined},p.children);
export const Image=(p:any)=>p.variant==='remote'?h(RemoteScreen,p):h('img',{src:p.src,alt:p.alt});
export const Button=(p:any)=>{const c=useContext(Fields);return h('button',{type:'button',disabled:p.disabled,onClick:()=>{c.item.current=p.item;p.onClick?.()}},p.children)};
export const Input=(p:any)=>{const c=useContext(Fields);return h('label',{},p.label,h(p.type==='textarea'?'textarea':'input',{'aria-label':p.label,type:p.type==='password'?'password':undefined,autoComplete:p.type==='password'?'new-password':undefined,placeholder:p.placeholder,value:c.fields[p.name]||'',onChange:(e:any)=>c.setField(p.name,e.target.value),onKeyDown:(e:any)=>{if(e.key==='Enter'&&p.name==='address'){e.preventDefault();c.action('navigate')}if(e.key==='Enter'&&p.name==='remoteText'){e.preventDefault();c.action('sendText')}}}))};

const base=location.pathname.replace(/\/$/,'');
function App(){
 const[fields,setFields]=useState<any>({}),[status,setStatus]=useState<any>({state:'Connecting'}),[jobs,setJobs]=useState<any[]>([]),[grants,setGrants]=useState<any[]>([]),[notice,setNotice]=useState(''),[token,setToken]=useState(''),[signedOut,setSignedOut]=useState(false);
 const[control,setControl]=useState(''),[remoteState,setRemoteState]=useState('closed'),[remoteImage,setRemoteImage]=useState(''),[tabs,setTabs]=useState<any[]>([]),[remotePage,setRemotePage]=useState('');
 const session=useRef<any>(null),item=useRef<any>(null),controlRef=useRef(''),frameBusy=useRef(false),inputChain=useRef(Promise.resolve()),pending=useRef(0),remoteStateRef=useRef('closed'),addressDirty=useRef(false);
 function state(value:string){remoteStateRef.current=value;setRemoteState(value)}
 function clearRemote(){controlRef.current='';setControl('');setRemoteImage('');setRemotePage('');setTabs([]);state('closed');setFields((v:any)=>({...v,remoteText:'',address:''}));addressDirty.current=false}
 function clearPrivate(){clearRemote();setSignedOut(true);setJobs([]);setGrants([]);setToken('');setFields({});session.current=null}
 async function api(path:string,method='GET',body?:any,handle=''){
  const r=await fetch(base+path,{method,cache:'no-store',headers:{'Content-Type':'application/json','X-CSRF-Token':session.current?.csrf||'','X-Workspace-Expected-Account':session.current?.user?.id||'',...(handle?{'X-WebSearch-Control':handle}:{})},body:body===undefined?undefined:JSON.stringify(body)});
  if(r.status===401||r.status===403||(r.status===409&&r.headers.get('X-Workspace-Account')!==session.current?.user?.id))clearPrivate();
  const data=await r.json();if(!r.ok){const e:any=Error(data.error||'Request could not finish');e.status=r.status;throw e}return data;
 }
 async function desktop(path:string,method='GET',body?:any,handle=controlRef.current){
  try{return await api('/api/desktop/'+path,method,body,handle)}catch(e:any){if(controlRef.current===handle){if(e.status===409)clearRemote();else{setRemoteImage('');state('disconnected')}}throw e}
 }
 async function refresh(){try{
  session.current=await api('/api/me');if(session.current?.user?.role!=='owner'){clearPrivate();throw Error('Web Search is managed by the Workspace owner.')}setSignedOut(false);
  const[s,j,g]=await Promise.all([api('/api/status'),api('/api/jobs'),api('/api/grants')]);setStatus(s);setJobs(j.map((x:any)=>({...x,answer:x.result?.answer||'',sources:(x.result?.sources||[]).map((s:any)=>({...s,excerpt:(s.text||'').slice(0,1000)}))})));setGrants(g);
 }catch(e:any){setNotice(e.message)}}
 useEffect(()=>{refresh();const timer=setInterval(refresh,3000);return()=>clearInterval(timer)},[]);
 async function frame(){
  const id=controlRef.current;if(!id||!session.current||frameBusy.current)return;frameBusy.current=true;
  try{const data=await desktop('frame','GET',undefined,id);if(controlRef.current!==id)return;state(data.state);if(data.image){setRemoteImage(data.image);setTabs((data.tabs||[]).map((x:any)=>({...x,label:(x.active?'Current: ':'')+(x.title||x.url||'New tab').slice(0,65)})));setRemotePage(data.title+' · '+data.url);if(!addressDirty.current)setFields((v:any)=>({...v,address:data.url.startsWith('https://')?data.url:''}))}else setRemoteImage('');if(data.error)setNotice(data.error)}catch(e:any){if(controlRef.current===id)setNotice(e.message)}finally{frameBusy.current=false}
 }
 useEffect(()=>{
  if(!control)return;void frame();
  const frames=setInterval(frame,900),heartbeat=setInterval(async()=>{const id=controlRef.current;if(!id||!session.current)return;try{const d=await desktop('heartbeat','POST',{},id);if(controlRef.current===id){state(d.state);if(d.error){setRemoteImage('');setNotice(d.error)}}}catch(e:any){setNotice(e.message)}},10000);
  return()=>{clearInterval(frames);clearInterval(heartbeat)};
 },[control]);
 useEffect(()=>{const close=()=>{const id=controlRef.current,s=session.current;if(id&&s)void fetch(base+'/api/desktop/stop',{method:'POST',keepalive:true,headers:{'Content-Type':'application/json','X-CSRF-Token':s.csrf,'X-Workspace-Expected-Account':s.user.id,'X-WebSearch-Control':id},body:'{}'}).catch(()=>{})};addEventListener('pagehide',close);return()=>removeEventListener('pagehide',close)},[]);
 async function send(input:any){
  const id=controlRef.current;if(!id||remoteStateRef.current!=='ready')return;
  if(pending.current>4&&(input.type==='pointerMove'||input.type==='wheel'))return;
  if(pending.current>=64){setNotice('Input is catching up. Wait a moment before typing more.');return}
  pending.current++;
  const request=inputChain.current.catch(()=>{}).then(async()=>{if(controlRef.current!==id||remoteStateRef.current!=='ready')return;const d=await desktop('input','POST',input,id);if(d.state){state(d.state);setRemoteImage('')}});
  inputChain.current=request.catch((e:any)=>{if(controlRef.current===id)setNotice(e.message)}).finally(()=>{pending.current--});
  await inputChain.current;
 }
 async function action(name:string){const selected=item.current;try{
  setNotice('');
  switch(name){
   case'install':await api('/api/model/install','POST',{});break;
   case'enable':await api('/api/settings','POST',{enabled:true});break;
   case'disable':await api('/api/settings','POST',{enabled:false});clearRemote();break;
   case'search':await api('/api/search','POST',{query:fields.query||''});setFields((v:any)=>({...v,query:''}));break;
   case'create':setToken((await api('/api/grants','POST',{name:fields.name||''})).token);break;
   case'dismiss':setToken('');break;
   case'revoke':await api('/api/grants/'+selected,'DELETE');break;
   case'cancel':await api('/api/jobs/'+selected,'DELETE');break;
   case'openRemote':{const d=await api('/api/desktop/start','POST',{});controlRef.current=d.id;setControl(d.id);state(d.state);break}
   case'closeRemote':{const id=controlRef.current;clearRemote();if(id)await desktop('stop','POST',{},id);break}
   case'reconnectRemote':setRemoteImage('');state('starting');await desktop('reconnect','POST',{});break;
   case'navigate':await send({type:'navigate',url:fields.address||''});addressDirty.current=false;break;
   case'back':case'forward':case'reload':case'newTab':await send({type:name});addressDirty.current=false;break;
   case'selectTab':await send({type:'selectTab',tab:selected});addressDirty.current=false;break;
   case'closeTab':await send({type:'closeTab'});addressDirty.current=false;break;
   case'sendText':{const text=fields.remoteText||'';setFields((v:any)=>({...v,remoteText:''}));if(text)await send({type:'text',text});break}
   case'enter':await send({type:'key',keys:['Enter']});break;
   case'tab':await send({type:'key',keys:['Tab']});break;
   case'scrollUp':case'scrollDown':await send({type:'wheel',x:640,y:400,deltaY:name==='scrollUp'?-550:550,deltaX:0});break;
   case'fullscreen':await document.getElementById('remote-browser-screen')?.requestFullscreen();return;
  }
  await refresh();
 }catch(e:any){setNotice(e.message)}}
 const data={workspace:'/',login:base+'/auth/login',signedOut,state:status.state,enabled:status.enabled,enableDisabled:status.enabled||!status.modelInstalled,disabled:!status.enabled&&!control,notReady:status.state!=='ready'||status.queuePaused,notice,installDisabled:status.enabled||status.installingModel,setupLabel:status.installingModel?`Installing model: ${Math.round((status.downloadedBytes||0)/1000000)} MB downloaded`:status.setupError||(status.modelInstalled?'Local model installed.':'Model setup is a one-time verified download.'),queueLabel:`${status.queued||0} waiting${status.queuePaused?' · Paused for remote control':status.active?' · Research in progress':''}`,jobs,grants,token,remoteOpen:!!control,remoteClosed:!control,remoteOpenDisabled:!!control||signedOut,remoteNotReady:remoteState!=='ready',remoteImage,tabs,remotePage,remoteLabel:control?`Remote browser: ${remoteState}. Research is paused.`:status.controlActive?'Another owner session is controlling the browser.':'Remote control is closed.',apiExample:`POST ${location.origin}${base}/api/client/search\nAuthorization: Bearer YOUR_KEY\n{ "query": "Your research question" }\n\nGET ${base}/api/client/jobs/REQUEST_ID`};
 return h(Fields.Provider,{value:{fields,setField:(n:string,v:string)=>{if(n==='address')addressDirty.current=true;setFields((old:any)=>({...old,[n]:v}))},item,action}},h(RemoteInput.Provider,{value:{send,warn:setNotice}},h(Console,{data,onAction:action})));
}
createRoot(document.getElementById('root')!).render(h(App));
