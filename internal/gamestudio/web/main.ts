import {createElement as h, createContext, useContext, useEffect, useRef, useState} from 'react';
import {createRoot} from 'react-dom/client';
import {Studio} from './generated/Studio';
import './style.css';

const Fields = createContext<any>(null);
const tag = (name: string) => (p: any) => h(name, {className:p.variant,id:p.id}, p.children);
export const Page = tag('main'), Row = tag('div'), Stack = tag('div'), Card = tag('section'),
  Badge = tag('strong'), Text = tag('p'), Divider = tag('hr');
export const Image = (p: any) => h('img', {src:p.src,alt:p.alt,width:38,height:38});
export const Heading = (p: any) => h('h'+(p.level || 2), {}, p.children);
export const Alert = (p: any) => p.children && h('p',{role:'status',className:'notice'},p.children);
export const Link = (p: any) => h('a',{
  href:p.href, className:p.variant, target:p.target, rel:p.target?'noopener noreferrer':undefined,
  'data-workspace-apps-trigger':p.variant==='apps'?'true':undefined,
},p.children);
export const Button = (p: any) => {
  const c = useContext(Fields);
  return h('button',{type:'button',className:p.variant,disabled:p.disabled,onClick:()=>{
    c.selected.current=p.item;
    p.onClick?.();
  }},p.children);
};
export const Input = (p: any) => {
  const c=useContext(Fields);
  return h('label',{},p.label,h(p.type==='textarea'?'textarea':'input',{
    type:p.type==='password'?'password':undefined, maxLength:p.limit,
    autoComplete:p.type==='password'?'new-password':undefined,
    placeholder:p.placeholder, value:c.fields[p.name]||'',
    onChange:(e:any)=>c.change(p.name,e.target.value),
    onKeyDown:(e:any)=>{
      if(p.name==='prompt' && e.key==='Enter' && (e.ctrlKey||e.metaKey)) {
        e.preventDefault();c.action('send');
      }
    },
  }));
};
export const Checkbox = (p: any) => {
  const c=useContext(Fields);
  return h('label',{className:'check'},h('input',{type:'checkbox',checked:!!c.fields[p.name],
    onChange:(e:any)=>c.change(p.name,e.target.checked)}),p.label);
};
export const Select = (p: any) => {
  const c=useContext(Fields);
  return h('label',{},p.label,h('select',{value:c.fields[p.name]||'',
    onChange:(e:any)=>c.change(p.name,e.target.value)},[
    h('option',{key:'off',value:''},'Off'), h('option',{key:'local',value:'local'},'Local DeepSeek via Ollama'),
    h('option',{key:'cloud',value:'cloud'},'DeepSeek Cloud'),
  ]));
};

type Game = {id:string,build:string,editable:boolean,spec:{title:string,template:string}};
type Message = {author:string,text:string,style:string};
const starters: Record<string,string> = {
  catch:'Make a game where a rocket catches stars in space.',
  maze:'Make a forest maze where a frog looks for a gem.',
  runner:'Make a candy runner where a dragon jumps over blocks.',
  memory:'Make an ocean matching game with shells and treasure.',
};
const kinds: Record<string,string> = {catch:'Treasure catcher',maze:'Maze adventure',runner:'Jump & run',memory:'Matching cards'};
const base = location.pathname.replace(/\/$/,'');
const preview = (build:string) => base+'/play/'+build+'/';
const hostingOpen = '/_workspace/open/hosting';
const buildID = /^[a-f0-9]{32}$/;

function reserveTab(message:string) {
  const tab=window.open('about:blank','_blank');
  if(tab) {
    tab.opener=null;
    tab.document.title='Game Studio';
    const heading=tab.document.createElement('h1');
    heading.textContent=message;
    tab.document.body.append(heading);
  }
  return tab;
}

async function readJSON(r:Response) {
  const body=await r.json().catch(()=>null);
  if(!r.ok||body===null) throw Error(body?.error || 'The request could not finish.');
  return body;
}

async function deploy(g:Game, zip:Blob, account:{issuer:string,subject:string}, update:(text:string)=>void) {
  const root='/hosting/api';
  const me=await fetch(root+'/me',{redirect:'error'}).then(readJSON).catch(()=>{
    throw Error('Open My Hosting to connect your owner session, then try again.');
  });
  if(me.user?.role!=='owner' || me.user?.issuer!==account?.issuer || me.user?.subject!==account?.subject || !me.csrf) {
    throw Error('My Hosting must be signed in with the same Workspace owner account.');
  }
  const options={redirect:'error' as RequestRedirect,headers:{'X-CSRF-Token':me.csrf,'Content-Type':'application/json'}};
  const slug='game-studio-'+g.id.slice(0,30);
  const projects=()=>fetch(root+'/projects',options).then(readJSON);
  const find=(items:any[])=>items.find(p=>p.slug===slug);
  let project=find(await projects());
  if(!project) {
    const result=await fetch(root+'/projects',{...options,method:'POST',
      body:JSON.stringify({name:'Game Studio '+g.id.slice(0,12),slug,adapter:'static',output:'.',baseDirectory:'.'})});
    project=result.status===409?find(await projects()):await readJSON(result);
  }
  if(!project || !buildID.test(project.id) || project.name!=='Game Studio '+g.id.slice(0,12) ||
    project.adapter!=='static' || project.trusted || project.repository || project.output!=='.' ||
    !['','.'].includes(project.baseDirectory) || ['install','build','start'].some(key=>project[key]?.length)) {
    throw Error('The reserved Hosting project was changed. Restore its static Game Studio settings before deploying.');
  }
  const endpoint=root+'/projects/'+project.id+'/deployments';
  const checksum=Array.from(new Uint8Array(await crypto.subtle.digest('SHA-256',await zip.arrayBuffer())),
    b=>b.toString(16).padStart(2,'0')).join('');
  const history=await fetch(endpoint,options).then(readJSON);
  let release=history.deployments?.find((d:any)=>d.sourceHash===checksum && d.environment==='production' &&
    (d.id===project.live || ['queued','preparing','uploading','building','publishing'].includes(d.state)));
  if(!release) release=await fetch(endpoint+'?environment=production',{...options,method:'POST',
    headers:{...options.headers,'Content-Type':'application/zip'},body:zip}).then(readJSON);
  if(!buildID.test(release.id)) throw Error('My Hosting returned an invalid deployment receipt.');
  update('My Hosting is deploying your game…');
  for(let i=0;i<50;i++) {
    const next=await fetch(endpoint+'/'+release.id,options).then(readJSON);
    if(next.state==='published') {
      const link=new URL(next.url);
      if(link.username || link.password || !(link.protocol==='https:' ||
        link.protocol==='http:' && (link.hostname==='localhost'||link.hostname.endsWith('.localhost')))) {
        throw Error('My Hosting returned an unexpected game address.');
      }
      return link.href;
    }
    if(['failed','cancelled','interrupted'].includes(next.state)) throw Error('My Hosting could not publish the game. Open its deployment log for details.');
    await new Promise(resolve=>setTimeout(resolve,1250));
  }
  throw Error('The build was queued in My Hosting. Open My Hosting to check its progress.');
}

function App() {
  const [fields,setFields]=useState<any>({});
  const [session,setSession]=useState<any>(null), sessionRef=useRef<any>(null);
  const [bootstrap,setBootstrap]=useState<any>(null), [config,setConfig]=useState<any>({});
  const [games,setGames]=useState<Game[]>([]), [messages,setMessages]=useState<Message[]>([]);
  const [editing,setEditing]=useState<Game|null>(null), [publishing,setPublishing]=useState<Game|null>(null);
  const [hosted,setHosted]=useState(''), [lastPlay,setLastPlay]=useState('');
  const [notice,setNotice]=useState(''), [loading,setLoading]=useState(true), [busy,setBusy]=useState(false);
  const selected=useRef<string>(''), running=useRef(false);
  function change(name:string,value:any) {setFields((f:any)=>({...f,[name]:value}));}
  function clear() {
    sessionRef.current=null;
    setSession(null);setGames([]);setMessages([]);setConfig({});setFields({});
    setEditing(null);setPublishing(null);setHosted('');setLastPlay('');
  }
  async function fetchAPI(path:string,method='GET',body?:any) {
    const account=sessionRef.current;
    const r=await fetch(base+path,{method,redirect:'error',cache:'no-store',headers:{
      'Content-Type':'application/json','X-CSRF-Token':account?.csrf||'',
      'X-Workspace-Expected-Account':account?.user.id||'',
      ...(path==='/api/setup'?{'X-Setup-Token':fields.setupCode||''}:{}),
    },body:body===undefined?undefined:JSON.stringify(body)});
    if(r.status===401 || r.status===409 && account && r.headers.get('X-Workspace-Account')!==account.user.id) clear();
    return r;
  }
  async function api(path:string,method?:string,body?:any) {
    return readJSON(await fetchAPI(path,method,body));
  }
  async function refresh() {
    const me=await api('/api/me');
    sessionRef.current=me;setSession(me);
    const next=await api('/api/studio');
    setConfig(next);
    if(me.user.role==='owner') {
      setFields((f:any)=>({...f,child:next.child,localModel:next.localModel,provider:next.provider,
        gitRepo:next.gitRepo,cloudKey:'',gitKey:'',cloudConsent:next.provider==='cloud',clearCloud:false,clearGit:false}));
    }
    if(next.allowed) setGames(await api('/api/games'));
  }
  useEffect(()=>{
    void (async()=>{
      try {
        const b=await api('/api/bootstrap');setBootstrap(b);
        const r=await fetchAPI('/api/me');
        if(r.ok) await refresh();
      } catch(e:any) {setNotice(e.message);}
      finally {setLoading(false);}
    })();
  },[]);
  useEffect(()=>{
    if(!session) return;
    const check=async()=>{
      try {
        const current=await api('/api/studio');
        setConfig(current);
        if(!current.allowed) {
          setGames([]);setMessages([]);setEditing(null);setLastPlay('');
          setFields((f:any)=>({...f,prompt:''}));
        }
      } catch { /* The API clears the session when the cookie expires. */ }
    };
    const timer=setInterval(check,15000);
    window.addEventListener('focus',check);
    return ()=>{clearInterval(timer);window.removeEventListener('focus',check);};
  },[session]);
  async function action(name:string) {
    if(running.current) return;
    if(name in starters) {
      change('prompt',starters[name]);document.querySelector('textarea')?.focus();return;
    }
    if(name==='newGame') {
      setEditing(null);setMessages([]);setLastPlay('');change('prompt','');return;
    }
    const game=games.find(g=>g.id===selected.current);
    if(name==='edit' && game?.editable) {
      setEditing(game);setMessages([]);setLastPlay(preview(game.build));change('prompt','');
      document.querySelector('textarea')?.focus();return;
    }
    if(name==='publish' && game) {
      setPublishing(game);setHosted('');setNotice('');requestAnimationFrame(()=>document.querySelector('#publishPanel')?.scrollIntoView({block:'nearest'}));return;
    }
    if(name==='cancelPublish') {setPublishing(null);return;}
    const prompt=fields.prompt?.trim();
    if(name==='send' && (!prompt||!config.allowed||!config.provider)) return;
    const tab=name==='send'?reserveTab('Building your game…'):
      name==='confirmPublish'?reserveTab('Publishing your game…'):null;
    running.current=true;setBusy(true);setNotice('');
    try {
      switch(name) {
      case 'localLogin': {
        const body={email:fields.loginName||'',password:fields.password||''};
        await api(bootstrap.needsSetup?'/api/setup':'/api/login','POST',body);
        setFields({});setBootstrap({...bootstrap,needsSetup:false});await refresh();break;
      }
      case 'signOut':
        await api('/auth/logout','POST');clear();break;
      case 'saveSettings':
        await api('/api/settings','PUT',Object.fromEntries(['provider','localModel','child','gitRepo',
          'cloudKey','gitKey','cloudConsent','clearGit','clearCloud'].map(k=>[k,fields[k]])));
        await refresh();setNotice('Settings saved.');break;
      case 'send': {
        const previous=messages.slice(-4).map(m=>m.author+': '+m.text).join('\n');
        const out=await api('/api/turn','POST',{game:editing?.id||'',build:editing?.build||'',
          prompt,context:previous.slice(-1500)});
        setMessages([...messages,{author:'You',text:prompt,style:'message mine'},
          {author:'Game Studio',text:out.message,style:'message studio'}].slice(-20));
        change('prompt','');
        if(out.game) {
          setEditing(out.game);setLastPlay(preview(out.game.build));
          if(tab) tab.location.replace(preview(out.game.build));
          setGames([out.game,...games.filter(g=>g.id!==out.game.id)]);
        } else tab?.close();
        break;
      }
      case 'confirmPublish': {
        if(!publishing || !crypto.subtle) throw Error('Open Game Studio over HTTPS to publish.');
        const selected=publishing;
        setNotice('Saving this build to private My Git…');
        await api('/api/builds/'+selected.build+'/git','POST',{});
        const file=await fetchAPI('/api/builds/'+selected.build+'/archive');
        if(!file.ok) throw Error('The build archive could not be downloaded.');
        const link=await deploy(selected,await file.blob(),sessionRef.current?.user,setNotice);
        setHosted(link);setPublishing(null);setNotice('Game published.');
        if(tab) tab.location.replace(link);
        break;
      }
      }
    } catch(e:any) {tab?.close();setNotice(e.message);}
    finally {running.current=false;setBusy(false);}
  }
  const data={
    icon:base+'/icon.svg', workspace:bootstrap?.workspace||'/', loading,
    signedOut:!loading&&!session, oidc:bootstrap?.oidc, setup:bootstrap?.needsSetup, session,
    account:session?.user.id, waiting:!!session&&config.allowed===false, allowed:config.allowed,
    owner:session?.user.role==='owner', notice, busy, paused:!config.provider, cloud:config.provider==='cloud',
    sendDisabled:busy||!fields.prompt?.trim()||!config.provider, sendLabel:busy?'Working…':'Send idea',
    chatTitle:editing?'Changing '+editing.spec.title:'What shall we build?', messages, lastPlay,
    empty:games.length===0, games:games.map(g=>({...g,type:kinds[g.spec.template],play:preview(g.build)})),
    publishing, publishTitle:publishing?'Publish '+publishing.spec.title:'', hostingOpen, hosted,
  };
  return h(Fields.Provider,{value:{fields,change,selected,action}},h(Studio,{data,onAction:action}));
}
createRoot(document.querySelector('#root')!).render(h(App));
