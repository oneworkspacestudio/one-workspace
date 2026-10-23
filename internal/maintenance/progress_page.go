package maintenance

const progressScript = `
let operation=null,receivedAt=0,disconnected=true;
const panel=document.createElement('section');panel.id='operation';panel.hidden=true;panel.setAttribute('aria-labelledby','operation-title');
panel.innerHTML='<div class="progress-head"><h2 id="operation-title">Updating Workspace</h2><span id="step-count" class="muted"></span></div><progress id="overall" aria-label="Update stages completed"></progress><p id="current-detail" role="status" aria-live="polite"></p><progress id="meter" aria-label="Current stage progress"></progress><p id="transfer"></p><p id="timing" class="muted"></p><p id="heartbeat" role="status"></p><ol id="steps"></ol><details><summary>Activity details</summary><ol id="activity"></ol></details>';
$('error').after(panel);
function duration(seconds){seconds=Math.max(0,Math.round(seconds));if(seconds<60)return seconds+'s';if(seconds<3600)return Math.floor(seconds/60)+'m '+seconds%60+'s';return Math.floor(seconds/3600)+'h '+Math.floor(seconds%3600/60)+'m'}
function bytes(n){const units=['B','KB','MB','GB','TB'];let i=0;while(n>=1024&&i<units.length-1){n/=1024;i++}return n.toFixed(i?1:0)+' '+units[i]}
function updateProgress(p){operation=p;receivedAt=Date.now();disconnected=false;renderProgress()}
function renderProgress(){
 const p=operation;panel.hidden=!p;$('phase').hidden=!!p;if(!p)return;
 const running=p.status==='running',complete=p.status==='complete',failed=p.status==='failed';
 const titles={check:'Checking app releases',apply:'Updating Workspace',rollback:'Restoring your backup',recover:'Recovering Workspace'};
 $('operation-title').textContent=complete?(p.kind==='check'?'Release check complete':p.kind==='apply'?'Update complete':'Installation restored'):failed?'Operation needs attention':p.recovering?'Restoring the previous installation':titles[p.kind]||'Working';
 $('operation-title').className=complete?'complete':failed?'failed':'';
 const steps=p.steps||[],active=steps.findIndex(x=>x.state==='running'||x.state==='failed'),done=steps.filter(x=>x.state==='done'||x.state==='skipped').length;
 $('overall').hidden=!running;$('overall').max=steps.length||1;$('overall').value=done;
 $('step-count').textContent=complete?'Finished':active>=0?'Stage '+(active+1)+' of '+steps.length:'Preparing';
 $('current-detail').textContent=complete&&p.kind==='apply'?(p.detail||'').replace(/^Update complete\.\s*/,''):p.detail||'Preparing your request…';
 const quiet=(Date.now()-new Date(p.lastActivityAt).getTime())/1000,fresh=!disconnected&&quiet<15;
 const measured=p.total>0;$('meter').hidden=!running;
 if(measured){$('meter').max=p.total;$('meter').value=Math.min(p.done,p.total)}else $('meter').removeAttribute('value');
 let transfer='';if(measured){transfer=p.unit==='bytes'?bytes(p.done)+' of '+bytes(p.total):p.done+' of '+p.total+' '+p.unit;transfer+=' · '+Math.min(100,Math.floor(p.done/p.total*100))+'% of this stage'}
 if(running&&fresh&&p.bytesPerSecond>0)transfer+=' · '+bytes(p.bytesPerSecond)+'/s';
 $('transfer').textContent=transfer;$('transfer').hidden=!running;
 const since=running?(Date.now()-receivedAt)/1000:0,elapsed=p.elapsedSeconds+since;
 let timing=(complete?'Completed in ':failed?'Stopped after ':'Elapsed: ')+duration(elapsed);
 if(running){timing+=fresh&&p.remainingSeconds>=0?' · About '+duration(Math.max(0,p.remainingSeconds-since))+' remaining in this stage':' · Estimating this stage as work progresses';timing+='. Later stages have their own timing.'}
 $('timing').textContent=timing;
 $('heartbeat').textContent=disconnected&&running?'Connection to the updater was lost. The update may still be running; this page will reconnect automatically.':running&&quiet>20?'No new progress for '+duration(quiet)+'. '+(p.stage==='download'?'Waiting for the release server or network.':'Windows may still be copying, checking or starting an app.')+' Keep this window open.':'';
 $('steps').replaceChildren();for(const step of steps){const li=document.createElement('li');li.dataset.state=step.state;const marks={done:'✓ ',running:'› ',failed:'! ',skipped:'– ',pending:'○ '};li.textContent=(marks[step.state]||'')+step.label+(step.state==='skipped'?' (not required)':'');if(step.state==='running')li.setAttribute('aria-current','step');$('steps').append(li)}
 $('activity').replaceChildren();for(const line of p.history||[]){const li=document.createElement('li');li.textContent=new Date(line.at).toLocaleTimeString()+' · '+line.text;$('activity').append(li)}
}
setInterval(renderProgress,1000);
`
