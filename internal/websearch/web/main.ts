import React, {createElement as h,useEffect,useRef,useState,createContext,useContext} from 'react';
import {createRoot} from 'react-dom/client';
import {Console} from './generated/Console';
import './style.css';

export const Fields=createContext<any>({fields:{},setField:()=>{},item:{current:null}});
const tag=(name:string)=>(p:any)=>h(name,{className:p.variant||undefined},p.children);
export const Page=tag('main'),Row=tag('div'),Card=tag('section'),Grid=tag('div'),Text=(p:any)=>h(p.variant==='pre'?'pre':'p',{},p.children),Badge=tag('strong'),Spacer=tag('span');
export const Heading=(p:any)=>h('h'+(p.level||2),{},p.children);
export const Alert=(p:any)=>p.children?h('p',{role:'status',className:'notice'},p.children):null;
export const Link=(p:any)=>h('a',{href:p.href,'data-workspace-apps-trigger':p.variant==='apps'?'true':undefined},p.children);
export const Image=(p:any)=>h('img',{src:p.src,alt:p.alt});
export const Button=(p:any)=>{const c=useContext(Fields);return h('button',{type:'button',disabled:p.disabled,onClick:()=>{c.item.current=p.item;p.onClick?.()}},p.children)};
export const Input=(p:any)=>{const c=useContext(Fields);return h('label',{},p.label,h(p.type==='textarea'?'textarea':'input',{'aria-label':p.label,placeholder:p.placeholder,value:c.fields[p.name]||'',onChange:(e:any)=>c.setField(p.name,e.target.value)}))};

const base=location.pathname.replace(/\/$/,'');
function App(){
 const[fields,setFields]=useState<any>({}),[status,setStatus]=useState<any>({state:'Connecting'}),[jobs,setJobs]=useState<any[]>([]),[grants,setGrants]=useState<any[]>([]),[notice,setNotice]=useState(''),[token,setToken]=useState(''),[preview,setPreview]=useState(''),[signedOut,setSignedOut]=useState(false);
 const session=useRef<any>(null),item=useRef<any>(null);
 function clearPrivate(){setSignedOut(true);setPreview('');setJobs([]);setGrants([]);setToken('');session.current=null}
 async function api(path:string,method='GET',body?:any){const r=await fetch(base+path,{method,headers:{'Content-Type':'application/json','X-CSRF-Token':session.current?.csrf||'','X-Workspace-Expected-Account':session.current?.user?.id||''},body:body===undefined?undefined:JSON.stringify(body)});if(r.status===401||r.status===403||(r.status===409&&r.headers.get('X-Workspace-Account')!==session.current?.user?.id)){clearPrivate();throw Error('Continue with the Workspace owner account to reconnect.')}const data=await r.json();if(!r.ok)throw Error(data.error||'Request could not finish');return data}
 async function refresh(){try{session.current=await api('/api/me');if(session.current?.user?.role!=='owner'){clearPrivate();throw Error('Web Search is managed by the Workspace owner.')}setSignedOut(false);const[s,j,g]=await Promise.all([api('/api/status'),api('/api/jobs'),api('/api/grants')]);setStatus(s);setJobs(j.map((x:any)=>({...x,answer:x.result?.answer||'',sources:(x.result?.sources||[]).map((s:any)=>({...s,excerpt:(s.text||'').slice(0,1000)}))})));setGrants(g)}catch(e:any){setNotice(e.message)}}
 useEffect(()=>{refresh();const timer=setInterval(refresh,3000);return()=>clearInterval(timer)},[]);
 async function action(name:string){try{setNotice('');switch(name){case'install':await api('/api/model/install','POST',{});break;case'enable':await api('/api/settings','POST',{enabled:true});break;case'disable':await api('/api/settings','POST',{enabled:false});setPreview('');break;case'search':await api('/api/search','POST',{query:fields.query||''});setFields((v:any)=>({...v,query:''}));break;case'create':setToken((await api('/api/grants','POST',{name:fields.name||''})).token);break;case'dismiss':setToken('');break;case'revoke':await api('/api/grants/'+item.current,'DELETE');break;case'cancel':await api('/api/jobs/'+item.current,'DELETE');break;case'preview':setPreview((await api('/api/browser/screenshot')).image);break}await refresh()}catch(e:any){setNotice(e.message)}}
 const data={workspace:'/',login:base+'/auth/login',signedOut,state:status.state,enabled:status.enabled,enableDisabled:status.enabled||!status.modelInstalled,disabled:!status.enabled,notReady:status.state!=='ready',notice,installDisabled:status.enabled||status.installingModel,setupLabel:status.installingModel?`Installing model: ${Math.round((status.downloadedBytes||0)/1000000)} MB downloaded`:status.setupError||(status.modelInstalled?'Local model installed.':'Model setup is a one-time verified download.'),queueLabel:`${status.queued||0} waiting${status.active?' · Research in progress':''}`,jobs,grants,token,preview,apiExample:`POST ${location.origin}${base}/api/client/search\nAuthorization: Bearer YOUR_KEY\n{ "query": "Your research question" }\n\nGET ${base}/api/client/jobs/REQUEST_ID`};
 return h(Fields.Provider,{value:{fields,setField:(n:string,v:string)=>setFields((old:any)=>({...old,[n]:v})),item}},h(Console,{data,onAction:action}));
}
createRoot(document.getElementById('root')!).render(h(App));
