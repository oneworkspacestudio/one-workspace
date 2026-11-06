import {createElement as h,createContext,useContext,useEffect,useRef} from 'react';

export const RemoteInput=createContext<any>({send:()=>{},warn:()=>{}});
const navigation=new Set(['Enter','Tab','Escape','Backspace','Delete','ArrowLeft','ArrowRight','ArrowUp','ArrowDown','Home','End','PageUp','PageDown']);

// Adapter for AUI Image variant=remote. Only private browser input is sent.
export function RemoteScreen(p:any){
 const el=useRef<HTMLDivElement|null>(null),image=useRef<HTMLImageElement|null>(null),down=useRef(false),lastMove=useRef(0),lastClick=useRef({at:0,x:0,y:0,count:1}),pos=useRef({x:0,y:0});
 const remote=useContext(RemoteInput);
 function point(e:any){const r=image.current!.getBoundingClientRect();const value={x:Math.max(0,Math.min(1279,(e.clientX-r.left)/r.width*1280)),y:Math.max(0,Math.min(799,(e.clientY-r.top)/r.height*800))};pos.current=value;return value}
 useEffect(()=>{
  const node=el.current!;
  const wheel=(e:WheelEvent)=>{e.preventDefault();remote.send({type:'wheel',...point(e),deltaX:Math.max(-1600,Math.min(1600,e.deltaX)),deltaY:Math.max(-1600,Math.min(1600,e.deltaY))})};
  const paste=(e:ClipboardEvent)=>{e.preventDefault();const text=e.clipboardData?.getData('text/plain');if(text)remote.send({type:'text',text})};
  node.addEventListener('wheel',wheel,{passive:false});node.addEventListener('paste',paste);
  return()=>{node.removeEventListener('wheel',wheel);node.removeEventListener('paste',paste)};
 },[remote]);
 function key(e:any){
  if(e.nativeEvent.isComposing)return;
  const ctrl=e.ctrlKey||e.metaKey,lower=e.key.toLowerCase();
  if(ctrl&&lower==='v')return;
  e.preventDefault();
  if(ctrl&&lower==='l'){(document.querySelector('input[aria-label="Browser address"]') as HTMLInputElement)?.focus();return}
  if(e.altKey||(ctrl&&lower!=='a'&&!navigation.has(e.key))){remote.warn('Use the browser controls above; clipboard and system shortcuts are unavailable.');return}
  if(navigation.has(e.key)||ctrl){const keys=[];if(ctrl)keys.push('CTRL');if(e.shiftKey)keys.push('SHIFT');keys.push(e.key);remote.send({type:'key',keys});return}
  if(Array.from(e.key).length===1)remote.send({type:'text',text:e.key});
 }
 function release(e:any){if(!down.current)return;down.current=false;e.preventDefault();remote.send({type:'pointerUp',...point(e),buttons:0,count:lastClick.current.count})}
 return h('div',{id:'remote-browser-screen',ref:el,className:'remote-screen',tabIndex:0,role:'application','aria-label':'Remote browser screen',onKeyDown:key,
  onPointerDown:(e:any)=>{if(e.button!==0)return;e.preventDefault();el.current!.focus();el.current!.setPointerCapture(e.pointerId);const pt=point(e),old=lastClick.current,count=Date.now()-old.at<400&&Math.hypot(pt.x-old.x,pt.y-old.y)<5?2:1;lastClick.current={at:Date.now(),...pt,count};down.current=true;remote.send({type:'pointerDown',...pt,buttons:1,count})},
  onPointerMove:(e:any)=>{if(!down.current||Date.now()-lastMove.current<70)return;lastMove.current=Date.now();remote.send({type:'pointerMove',...point(e),buttons:1})},onPointerUp:release,
  onPointerCancel:()=>{if(down.current){down.current=false;remote.send({type:'pointerUp',...pos.current,buttons:0,count:1})}},onBlur:()=>{if(down.current){down.current=false;remote.send({type:'pointerUp',...pos.current,buttons:0,count:1})}}
 },h('img',{ref:image,src:p.src,alt:p.alt,draggable:false}));
}
