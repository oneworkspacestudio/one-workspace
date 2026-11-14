import {createElement as h, useEffect, useState} from 'react';

type History = {
  day:number; firstDay:number; lastDay:number; startedLate:boolean;
  events:{tick:number; kind:string; text:string}[];
};

const marks:Record<string,string> = {
  explore:'✧', home:'⌂', treat:'✿', sleep:'☾', play:'◈', art:'◒', build:'▦', eat:'●',
  gather:'◇', learn:'✦', plan:'↗', clock:'◷',
};

function timeAt(tick:number) {
  const minute = (tick*2+510)%1440;
  return String(Math.floor(minute/60)).padStart(2,'0')+':'+String(minute%60).padStart(2,'0');
}

export function DayHistory({account,currentDay,tick}:{account:string;currentDay:number;tick:number}) {
  const [selected,setSelected] = useState<number>();
  const [history,setHistory] = useState<History>();
  const [input,setInput] = useState<string>();
  const [error,setError] = useState('');
  const day = selected ?? currentDay;
  // Refresh the current day without sending an additional request every tick.
  const revision = day === currentDay ? Math.floor(tick/5) : 0;
  useEffect(() => {
    const control = new AbortController();
    fetch('./api/history?day='+day,{signal:control.signal,cache:'no-store',
      headers:{'X-Workspace-Expected-Account':account},
    }).then(response => {
      if(!response.ok) throw Error('The activity history could not be loaded.');
      return response.json();
    }).then((value:History)=>{
      if(!control.signal.aborted){setHistory(value);setError('');}
    }).catch(()=>{
      if(!control.signal.aborted) setError('The activity history could not be loaded.');
    });
    return ()=>control.abort();
  },[account,day,revision]);
  const ready = history?.day===day;
  const firstDay = history?.firstDay ?? day;
  function change(value:number) {
    if(Number.isInteger(value) && value>=firstDay && value<=currentDay){
      setSelected(value===currentDay ? undefined : value);
      setInput(undefined);
    }
  }
  return h('section',{className:'journal'},
    h('h2',{},'Activity by day'),
    h('div',{className:'dayControls'},
      h('button',{type:'button',className:'secondary',disabled:day<=firstDay,onClick:()=>change(day-1),'aria-label':'Previous day'},'‹'),
      h('form',{onSubmit:(event)=>{event.preventDefault();change(+(input??day));setInput(undefined);}},
        h('label',{},'Day ',h('input',{type:'number',min:firstDay,max:currentDay,value:input??day,
          onChange:e=>setInput(e.currentTarget.value)})),
        h('button',{type:'submit',className:'secondary'},'Go')),
      h('button',{type:'button',className:'secondary',disabled:day>=currentDay,onClick:()=>change(day+1),'aria-label':'Next day'},'›'),
      day!==currentDay && h('button',{type:'button',className:'secondary',onClick:()=>change(currentDay)},'Today')),
    ready && history.startedLate && day===firstDay &&
      h('p',{className:'muted'},'The earliest available day starts partway through.'),
    error ? h('p',{role:'status'},error) : !ready ? h('p',{role:'status'},'Loading activity…') :
      history.events.length===0 ? h('p',{className:'muted'},'No recorded activity on this day.') :
        h('ol',{className:'timeline'},history.events.map((event,i)=>
          h('li',{className:'event',key:event.tick+'-'+i},
            h('span',{className:'eventIcon','aria-hidden':true},marks[event.kind] || '·'),
            h('span',{},event.text),h('time',{className:'eventTime'},timeAt(event.tick))))));
}
