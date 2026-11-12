import {createElement as h, createContext, useContext, useEffect, useRef, useState, type ReactNode} from 'react';
import {renderHabitat, paint} from './world';

export type Point = {x: number; z: number};
export type Item = Point & {id: number; kind: string; style: number; level?: number};
export type World = {
  name: string; seed: string; tick: number; saved: number; position: Point; target: Point; action: string; until: number;
  focusStyle: number; mood: string; energy: number; fullness: number; happiness: number; curiosity: number;
  minute: number; day: number; weather: string; source: string; treatAfter: number;
  food: Item[]; creations: Item[]; journal: {tick: number; kind: string; text: string}[];
};
export type SettingsState = {mode: string; modelStatus: string; localModel: string; paired: boolean; hasService: boolean};
export type PetState = Partial<SettingsState> & {world: World | null};
type Props = {variant?: string; children?: ReactNode};
const tag = (name: string) => (p: Props) => h(name, {className: p.variant}, p.children);
export const Page = tag('main'), Row = tag('div'), Stack = tag('div'), Card = tag('section'), Text = tag('p'), Badge = tag('span');
export const Heading = (p: Props & {level?: number}) => h('h'+(p.level || 2), {}, p.children);
export const Alert = (p: Props) => p.children ? h('p', {className:'notice', role:'status'}, p.children) : null;
export const Image = (p: {src: string; alt: string}) => h('img', {src:p.src, alt:p.alt, width:37, height:37});
export const Link = (p: Props & {href: string}) => h('a', {
  href:p.href, className:p.variant, 'data-workspace-apps-trigger':p.variant==='apps'?'true':undefined,
}, p.children);
export const Button = (p: Props & {onClick: () => void; disabled?: boolean}) =>
  h('button', {type:'button', className:p.variant, disabled:p.disabled, onClick:p.onClick}, p.children);
type FormData = {mode:string; model:string; key:string; clearKey:boolean};
export const Fields = createContext({name:'', setName:(_value:string) => {}, save:async (_value:FormData) => {}});
export const Input = (p: {label: string; limit: number; placeholder: string}) => {
  const {name, setName} = useContext(Fields);
  return h('label', {}, p.label, h('input', {value:name, maxLength:p.limit, placeholder:p.placeholder,
    onChange:(e:any) => setName(e.target.value)}));
};

export function Portrait({mood}:{mood:string}) {
  return h('div', {className:'portrait', 'aria-hidden':true},
    h('i', {className:'ear left'}), h('i', {className:'ear right'}),
    h('div', {className:mood==='Dreaming'?'face sleeping':'face'}, h('i'), h('i'), h('b')));
}

export function Needs({state:w}:{state:World}) {
  const meters = [
    ['Happiness','♥',w.happiness,'happy'], ['Energy','✦',w.energy,'energy'],
    ['Fullness','●',w.fullness,'fullness'], ['Curiosity','◈',w.curiosity,'curiosity'],
  ] as const;
  return h('div',{className:'meters'}, meters.map(([label, icon, value, kind]) =>
    h('div',{key:kind,className:'meter '+kind},
      h('div',{className:'meterLabel'}, h('span', {'aria-hidden':true},icon),h('label',{htmlFor:'need-'+kind},label),h('span',{},Math.round(value)+'%')),
      h('progress',{id:'need-'+kind,value,max:100}))));
}

const artTitles = ['Sun over Clover Hill', 'Moonlit Pond', 'Wildflower Study', 'After the Rain'];
const buildTitles = ['Stone cairn', 'Meadow arch', 'Woodland tower', 'Clover stack'];

export function ArtThumb({item}:{item:Item}) {
  const el = useRef<HTMLCanvasElement>(null);
  useEffect(() => {
    const ctx = el.current?.getContext('2d');
    if (ctx) paint(ctx, item, 156, 108);
  }, [item.id, item.kind, item.style, item.level]);
  const title = (item.kind==='art'?artTitles:buildTitles)[item.style];
  return h('figure',{className:'thumb'},h('canvas',{ref:el,width:156,height:108,'aria-hidden':true}),
    h('figcaption',{},title));
}

export function HabitatView({state}:{state:World}) {
  const surface = useRef<HTMLDivElement>(null);
  const driver = useRef<ReturnType<typeof renderHabitat> | null>(null);
  const latest = useRef(state);
  latest.current = state;
  const [failed,setFailed] = useState(false);
  const [tracking,setTracking] = useState(true);
  const [touch,setTouch] = useState(() => navigator.maxTouchPoints > 0 || matchMedia('(pointer: coarse)').matches);
  useEffect(() => {
    const query = matchMedia('(pointer: coarse)');
    const changed = () => setTouch(navigator.maxTouchPoints > 0 || query.matches);
    query.addEventListener('change',changed);
    return () => query.removeEventListener('change',changed);
  }, []);
  useEffect(() => {
    try {
      const current = renderHabitat(surface.current!, () => latest.current, setTracking, () => setFailed(true));
      driver.current = current;
      return () => { driver.current = null; current.dispose(); };
    } catch {
      setFailed(true);
    }
  }, []);
  return h('div',{className:'viewport'+(failed?' failed':'')},
    h('div',{ref:surface,className:'surface'}),
    failed ? h('div',{className:'viewFallback'},h(Portrait,{mood:state.mood}),
      h('p',{},'3D view needs WebGL 2. Your pet is still active.'),h('p',{},'Follow along in Field notes.')) :
      h('div',{className:'overlay'},
        h('div',{className:'viewHeading'},h('span',{className:'liveDot'}),tracking?'Watching '+state.name:'First-person view'),
        h('button',{type:'button',className:'focusPet','aria-pressed':tracking,onClick:() => driver.current?.follow()},'◎  Find '+state.name),
        touch ? h('div',{className:'dpad',role:'group','aria-label':'Walk'},
          ...([['up','↑','Move forward'],['left','←','Move left'],['down','↓','Move back'],['right','→','Move right']] as const)
            .map(([dir, icon, label]) => h('button',{key:dir,type:'button',className:dir,'aria-label':label,
              onPointerDown:(e:any) => {e.preventDefault();e.currentTarget.setPointerCapture(e.pointerId);driver.current?.move(dir,true);},
              onLostPointerCapture:() => driver.current?.move(dir,false),
              onPointerUp:() => driver.current?.move(dir,false),
              onKeyDown:(e:any) => {if(e.key===' '||e.key==='Enter'){e.preventDefault();driver.current?.move(dir,true);}},
              onKeyUp:() => driver.current?.move(dir,false),
              onBlur:() => driver.current?.move(dir,false),
            },icon))) : null,
        h('span',{className:'viewHelp'},touch?'Drag to look':'Drag to look · WASD to walk')));
}

export function DecisionForm({state}:{state:SettingsState}) {
  const {save} = useContext(Fields);
  const [mode,setMode] = useState(state.mode);
  const [model,setModel] = useState(state.localModel);
  const [key,setKey] = useState('');
  const [clearKey,setClearKey] = useState(false);
  const [message,setMessage] = useState('');
  const [busy,setBusy] = useState(false);
  const wantsLocal = mode==='local' || mode==='both';
  const wantsShared = mode==='service' || mode==='both';
  const submit = async (event: any) => {
    event.preventDefault();
    setBusy(true);
    try {
      await save({mode,model,key,clearKey});
      setKey('');
      setClearKey(false);
      setMessage('Decision settings saved.');
    } catch (e: any) {
      setMessage(e.message);
    } finally {
      setBusy(false);
    }
  };
  return h('details', {},
    h('summary',{}, 'Decision settings'),
    h('form',{onSubmit:submit},
      h('p',{className:'muted'},'Models can suggest activities. Instincts handle food, sleep and outages.'),
      h('label',{},'Decision mode',h('select',{value:mode,onChange:(e:any)=>setMode(e.target.value)},
        h('option',{value:'instincts'},'Instincts'),
        h('option',{value:'local'},'Local Ollama model'),
        h('option',{value:'service',disabled:!state.hasService},'Workspace AI Service'),
        h('option',{value:'both',disabled:!state.hasService},'Local model + AI Service'))),
      wantsLocal && h('label',{},'Installed Ollama model',h('input',{value:model,placeholder:'qwen3.5:9b',maxLength:100,
        required:true,onChange:(e:any)=>setModel(e.target.value)})),
      wantsShared && h('label',{},state.paired?'Replace the Pet app key':'Pet app key',h('input',{value:key,type:'password',
        autoComplete:'new-password',maxLength:200,required:!state.paired,onChange:(e:any)=>setKey(e.target.value)})),
      wantsShared && h('p',{className:'muted'},'In AI Service, create a local app key with the ID pet. Paste it here.'),
      !wantsShared && state.paired && h('label',{className:'check'},h('input',{type:'checkbox',checked:clearKey,
        onChange:(e:any)=>setClearKey(e.target.checked)}),'Forget the saved AI Service key'),
      h('button',{className:'secondary',type:'submit',disabled:busy},'Save settings'),
      (message || state.modelStatus) && h('p',{role:'status',className:'muted'},[message,state.modelStatus].filter(Boolean).join(' '))));
}
