"""Rebuild the manager and add pinned components to a verified published bundle.

This is the normal organization build when managedSeed is set. Existing app bytes
and contracts are preserved exactly; the complete managed registry is validated.
"""
import argparse,hashlib,json,os,shutil,subprocess,tempfile,urllib.request,zipfile
from pathlib import Path
from managed_component import remote_component, tickets_component

ROOT=Path(__file__).resolve().parents[1]
def sha(p):
 h=hashlib.sha256()
 with Path(p).open('rb') as f:
  for b in iter(lambda:f.read(4194304),b''):h.update(b)
 return h.hexdigest()
def asset(pin,cache,provided=None):
 release=json.loads(subprocess.check_output(['gh','api',f"repos/{pin['repository']}/releases/tags/{pin['tag']}"]))
 assert not release['draft']
 a=next(a for a in release['assets'] if a['name']==pin['asset'])
 assert a['digest']=='sha256:'+pin['sha256']
 target=Path(provided) if provided else cache/(pin['sha256']+'-'+pin['asset'])
 if not target.exists():
  with urllib.request.urlopen(a['browser_download_url'],timeout=120) as src,target.open('wb') as dst:shutil.copyfileobj(src,dst,4194304)
 assert target.stat().st_size==a['size'] and sha(target)==pin['sha256'],str(target)
 return target
def extract(z,target):
 with zipfile.ZipFile(z) as src:
  for m in src.infolist():
   assert not m.filename.startswith(('/', '\\')) and '\\' not in m.filename and ':' not in m.filename and '..' not in Path(m.filename).parts
   assert (m.external_attr>>16)&0o170000!=0o120000
  src.extractall(target)
def main():
 p=argparse.ArgumentParser();p.add_argument('--go',default='go');p.add_argument('--version',required=True);p.add_argument('--cache',required=True);p.add_argument('--out',required=True);p.add_argument('--seed-archive');p.add_argument('--tickets-candidate',help='Local acceptance only; publisher rejects candidate builds');a=p.parse_args()
 cache=Path(a.cache).resolve();out=Path(a.out).resolve();cache.mkdir(parents=True,exist_ok=True);out.mkdir(parents=True,exist_ok=True)
 lock=json.loads((ROOT/'organization-components.lock.json').read_text(encoding='utf-8-sig'))
 seed=asset(lock['managedSeed'],cache,a.seed_archive);remote=asset(lock['myRemote'],cache)
 if a.tickets_candidate:
  tickets=Path(a.tickets_candidate).resolve();assert sha(tickets)==lock['myTickets']['sha256']
 else:tickets=asset(lock['myTickets'],cache)
 stage=Path(tempfile.mkdtemp(prefix='remote-',dir=cache));extract(remote,stage)
 contract=remote_component(stage)
 ticket_stage=Path(tempfile.mkdtemp(prefix='tickets-',dir=cache));extract(tickets,ticket_stage)
 ticket_contract=tickets_component(ticket_stage)
 contracts=[contract,ticket_contract]
 env=os.environ.copy();env.update(CGO_ENABLED='0',GOTOOLCHAIN='local')
 def go(*args):subprocess.run([a.go,*args],cwd=ROOT,env=env,check=True)
 go('build','-buildvcs=false','-trimpath','-ldflags',f'-s -w -X main.version={a.version}','-o',str(out/'Workspace.exe'),'./cmd/workspace')
 with zipfile.ZipFile(seed) as base:
  cat=json.loads(base.read('catalog.json'));old=json.loads(json.dumps(cat));cat['version']=a.version;cat['installerIntegration']=5
  for key,item in old['packages'].items():
   h=hashlib.sha256();size=0
   with base.open(item['file']) as src:
    for b in iter(lambda:src.read(4194304),b''):h.update(b);size+=len(b)
   assert size==item['size'] and h.hexdigest()==item['sha256'],key
  replacements={'Workspace.exe':out/'Workspace.exe','licenses/components/Workspace-release-status.md':ROOT/'docs/RELEASE-STATUS.md'}
  for k,pkg in cat['packages'].items():
   if pkg['file'] in replacements:
    f=replacements[pkg['file']];pkg.update(size=f.stat().st_size,sha256=sha(f))
  for contract in contracts:
   assert not any(x['id']==contract['app']['id'] for x in cat['apps'])
   cat['apps'].append(contract['app'])
   for k,item in contract['inputs'].items():
    assert k not in cat['packages'];f=Path(item['source']);cat['packages'][k]={'file':item['file'],'size':f.stat().st_size,'sha256':sha(f)}
  registry=json.loads((ROOT/'internal/catalog/managed-apps.json').read_text())['apps'];assert {x['id'] for x in registry}=={x['id'] for x in cat['apps']}
  (out/'catalog.json').write_text(json.dumps(cat,indent=2),encoding='utf-8')
  go('run','-buildvcs=false','./cmd/catalogcheck',str(out/'catalog.json'))
  payload=ROOT/'internal/payload/bundle/payload.zip';payload.parent.mkdir(parents=True,exist_ok=True)
  with zipfile.ZipFile(payload,'w',compression=zipfile.ZIP_DEFLATED,compresslevel=6) as dst:
   dst.writestr('catalog.json',json.dumps(cat,indent=2))
   for pkg in old['packages'].values():
    rel=pkg['file']
    if rel in replacements:dst.write(replacements[rel],rel)
    else:
     with base.open(rel) as src,dst.open(rel,'w',force_zip64=True) as target:shutil.copyfileobj(src,target,4194304)
   for contract in contracts:
    for item in contract['inputs'].values():dst.write(item['source'],item['file'])
 shutil.copyfile(payload,out/'workspace-bundle.zip')
 go('build','-buildvcs=false','-trimpath','-ldflags','-s -w -X github.com/oneworkspacestudio/one-workspace/internal/distribution.Repository=oneworkspacestudio/one-workspace','-o',str(out/'WorkspaceSetup-windows-amd64.exe'),'./cmd/setup')
 (out/'managed-build.json').write_text(json.dumps({'schema':1,'version':a.version,'seed':lock['managedSeed'],'remote':lock['myRemote'],'tickets':lock['myTickets'],'candidate':bool(a.tickets_candidate),'priorAppContractsRetained':True,'priorPackagesRetainedExceptManagerAndReleaseNotes':True,'completeInventory':len(cat['apps']),'builder':'scripts/build-managed-seed.py'},indent=2),encoding='utf-8')
 print('Verified managed bundle:',len(cat['apps']),'apps;',len(cat['packages']),'packages',flush=True)
if __name__=='__main__':main()
