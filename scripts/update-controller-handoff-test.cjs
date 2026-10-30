const fs=require('node:fs'),vm=require('node:vm'),assert=require('node:assert/strict');
const text=fs.readFileSync(require('node:path').join(__dirname,'../internal/maintenance/page.go'),'utf8');
const refresh=text.slice(text.indexOf('async function refresh()'),text.indexOf('\nasync function action('));
async function check(remote,busy){
 const calls=[],nodes=new Map();let lost=0,cleared=0;
 const context={remote,calls,clearInterval(){cleared++},updateProgress(){},buttons(){},connectionLost(){lost++},document:{createElement(){return {append(){}}}},$:id=>{if(!nodes.has(id))nodes.set(id,{replaceChildren(){},append(){}});return nodes.get(id)},api:async(path)=>{calls.push(path);return {busy,controllerChanged:true,newController:true,plan:null,backups:[],error:'',version:'new',phase:'ready'};}};
 vm.createContext(context);
 vm.runInContext("let refreshing=false,stateEpoch=0,awaitingResult=false,pending=false,installer=null,currentChoices=[],newController=false,requestError='',lastProblem='',requesting=false,reopening=false,lastPlan='null',lastBackups='[]';const timer=1;let busy=false;"+refresh,context);
 await vm.runInContext('refresh()',context);
 assert.equal(lost,0,'refresh failed');
 assert.deepEqual(calls,!remote&&!busy?['state','reopen']:['state']);
 assert.equal(cleared,!remote&&!busy?1:0);
 if(!remote&&!busy){await vm.runInContext('refresh()',context);assert.equal(calls.filter(x=>x==='reopen').length,1,'duplicate reopen');assert.equal(lost,0);}
}
(async()=>{await check(false,false);await check(false,true);await check(true,false);console.log('PASS: idle local handoff once; active and remote sessions never launch duplicate controllers.');})().catch(e=>{console.error(e);process.exit(1)});
