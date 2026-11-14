"""Rebuild the manager and add pinned components to a verified published bundle.

This is the normal organization build when managedSeed is set. Existing app bytes
and contracts are preserved exactly; the complete managed registry is validated.
"""
import argparse,hashlib,json,os,shutil,subprocess,tempfile,urllib.request,zipfile
from pathlib import Path
from managed_component import remote_component, tickets_component, hosting_replacements
from voice_component import voice_component
from books_component import books_component
from commute_component import commute_component
from podcasts_component import podcasts_component
from survival_component import survival_component
from myai_component import myai_replacements
from websearch_component import websearch_component
from gamestudio_component import gamestudio_component
from sunnymoon_component import sunnymoon_component
from pet_component import pet_component
from froststock_component import froststock_component
from tube_component import tube_replacements
from slicer_component import slicer_replacements
from finance_component import finance_replacements
from news_component import news_component
from authenticator_component import authenticator_component

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

def integrated_build(args, pin, seed, cache, out, components):
 env = dict(os.environ, CGO_ENABLED='0', GOTOOLCHAIN='local')
 def go(*values):
  subprocess.run([args.go, *values], cwd=ROOT, env=env, check=True)
 inputs = {key: item for component in components for key, item in component['inputs'].items()}
 manager = out / 'Workspace.exe'
 go('build', '-buildvcs=false', '-trimpath', '-ldflags', f'-s -w -X main.version={args.version}',
    '-o', str(manager), './cmd/workspace')
 payload = ROOT / 'internal/payload/bundle/payload.zip'
 payload.parent.mkdir(parents=True, exist_ok=True)
 with zipfile.ZipFile(seed) as base:
  cat = json.loads(base.read('catalog.json'))
  cat.update(version=args.version, installerIntegration=18)
  apps = {app['id']: app for app in cat['apps']}
  apps.update((c['app']['id'], c['app']) for c in components)
  cat['apps'] = list(apps.values())
  replacements = {'manager': manager, 'notice-Workspace-release-status.md': ROOT / 'docs/RELEASE-STATUS.md'}
  for key, replacement in replacements.items():
   cat['packages'][key].update(size=replacement.stat().st_size, sha256=sha(replacement))
  for key, item in inputs.items():
   file = Path(item['source'])
   cat['packages'][key] = dict(file=item['file'], size=file.stat().st_size, sha256=sha(file))
  registry = json.loads((ROOT / 'internal/catalog/managed-apps.json').read_text())['apps']
  assert {app['id'] for app in registry} == {app['id'] for app in cat['apps']}
  catalog = json.dumps(cat, indent=2)
  (out / 'catalog.json').write_text(catalog, encoding='utf-8')
  go('run', '-buildvcs=false', './cmd/catalogcheck', str(out / 'catalog.json'))
  with zipfile.ZipFile(payload, 'w', compression=zipfile.ZIP_DEFLATED, compresslevel=6) as dest:
   dest.writestr('catalog.json', catalog)
   for key, package in cat['packages'].items():
    file = package['file']
    if key in replacements:
     dest.write(replacements[key], file)
    elif key in inputs:
     dest.write(inputs[key]['source'], file)
    else:
     h = hashlib.sha256()
     size = 0
     with base.open(file) as source, dest.open(file, 'w', force_zip64=True) as target:
      while block := source.read(4194304):
       size += len(block)
       h.update(block)
       target.write(block)
     assert size == package['size'] and h.hexdigest() == package['sha256'], key
 shutil.copyfile(payload, out / 'workspace-bundle.zip')
 go('build', '-buildvcs=false', '-trimpath', '-ldflags',
    '-s -w -X github.com/oneworkspacestudio/one-workspace/internal/distribution.Repository=oneworkspacestudio/one-workspace',
    '-o', str(out / 'WorkspaceSetup-windows-amd64.exe'), './cmd/setup')
 (out / 'managed-build.json').write_text(json.dumps({
  'schema': 1, 'version': args.version, 'seed': pin, 'candidate': False,
  'completeInventory': len(cat['apps']), 'builder': 'scripts/build-managed-seed.py',
 }, indent=2), encoding='utf-8')
 print('Verified managed bundle:', len(cat['apps']), 'apps;', len(cat['packages']), 'packages', flush=True)

def main():
 p=argparse.ArgumentParser();p.add_argument('--go',default='go');p.add_argument('--version',required=True);p.add_argument('--cache',required=True);p.add_argument('--out',required=True);p.add_argument('--seed-archive');p.add_argument('--myai-candidate',help='Local acceptance only; publisher rejects candidate builds');p.add_argument('--tickets-candidate',help='Local acceptance only; publisher rejects candidate builds');p.add_argument('--authenticator-candidate',help='Local acceptance only');p.add_argument('--accounts-candidate',help='Local acceptance only');p.add_argument('--books-candidate',help='Local My Books acceptance archive; not a public pin');p.add_argument('--voice-candidate',help='Local acceptance only; publisher rejects candidate builds');p.add_argument('--news-candidate',help='Local My News acceptance archive; not a public pin');p.add_argument('--commute-candidate',help='Local My Commute acceptance only; never publish');p.add_argument('--survival-candidate',help='Local My Survival acceptance only');p.add_argument('--podcasts-candidate',help='Local My Podcasts acceptance only; never publish');a=p.parse_args()
 cache=Path(a.cache).resolve();out=Path(a.out).resolve();cache.mkdir(parents=True,exist_ok=True);out.mkdir(parents=True,exist_ok=True)
 lock=json.loads((ROOT/'organization-components.lock.json').read_text(encoding='utf-8-sig'))
 candidate = any(value for key, value in vars(a).items() if key.endswith('_candidate'))
 if lock.get('petSeed') and a.version != '0.3.0-preview.128' and not candidate:
  pin = lock['petSeed']
  base = asset(pin, cache, a.seed_archive)
  return integrated_build(a, pin, base, cache, out, [pet_component(ROOT,a.go,cache/'pet')])
 if lock.get('integratedSeed') and a.version == '0.3.0-preview.128' and not candidate:
  pin = lock['integratedSeed']
  base = asset(pin, cache, a.seed_archive)
  return integrated_build(a, pin, base, cache, out, [froststock_component(ROOT,a.go,cache/'froststock',a.version)])
 seed=asset(lock['managedSeed'],cache,a.seed_archive)
 remote=asset(lock['myRemote'],cache)
 hosting=asset(lock['myHosting'],cache)
 hosting_stage=Path(tempfile.mkdtemp(prefix='hosting-',dir=cache));extract(hosting,hosting_stage)
 myai=Path(a.myai_candidate).resolve() if a.myai_candidate else asset(lock['myAI'],cache)
 myai_stage=Path(tempfile.mkdtemp(prefix='myai-',dir=cache));extract(myai,myai_stage)
 if a.myai_candidate:lock['myAI']={'repository':'oneworkspacestudio/my-ai','tag':'v'+json.loads((myai_stage/'MyAI/compatibility.json').read_text())['version'],'asset':myai.name,'sha256':sha(myai)}
 tube=asset(lock['myTube'],cache)
 tube_stage=Path(tempfile.mkdtemp(prefix='tube-',dir=cache));extract(tube,tube_stage)
 slicer=asset(lock['mySlicer'],cache)
 slicer_stage=Path(tempfile.mkdtemp(prefix='slicer-',dir=cache));extract(slicer,slicer_stage)
 finance=asset(lock['myFinance'],cache)
 finance_stage=Path(tempfile.mkdtemp(prefix='finance-',dir=cache));extract(finance,finance_stage)
 if a.tickets_candidate:
  tickets=Path(a.tickets_candidate).resolve();assert sha(tickets)==lock['myTickets']['sha256']
 else:tickets=asset(lock['myTickets'],cache)
 stage=Path(tempfile.mkdtemp(prefix='remote-',dir=cache));extract(remote,stage)
 contract=remote_component(stage)
 ticket_stage=Path(tempfile.mkdtemp(prefix='tickets-',dir=cache));extract(tickets,ticket_stage)
 ticket_contract=tickets_component(ticket_stage)
 if a.authenticator_candidate:
  auth=Path(a.authenticator_candidate).resolve();assert sha(auth)==lock['oneAuthenticator']['sha256']
 else:auth=asset(lock['oneAuthenticator'],cache)
 if a.accounts_candidate:
  accounts=Path(a.accounts_candidate).resolve();assert sha(accounts)==lock['accounts']['sha256']
 else:accounts=asset(lock['accounts'],cache)
 auth_stage=Path(tempfile.mkdtemp(prefix='authenticator-',dir=cache));extract(auth,auth_stage)
 if a.voice_candidate:
  voice=Path(a.voice_candidate).resolve();assert sha(voice)==lock['myVoice']['sha256']
 else:voice=asset(lock['myVoice'],cache)
 voice_stage=Path(tempfile.mkdtemp(prefix='voice-',dir=cache));extract(voice,voice_stage)
 if a.books_candidate:books=Path(a.books_candidate).resolve()
 elif lock.get('myBooks'):books=asset(lock['myBooks'],cache)
 else:raise ValueError('Publish and pin My Books before a release, or use --books-candidate for local acceptance')
 books_stage=Path(tempfile.mkdtemp(prefix='books-',dir=cache));extract(books,books_stage)
 books_contract=books_component(books_stage,allow_candidate=bool(a.books_candidate))
 if a.news_candidate:news=Path(a.news_candidate).resolve()
 elif lock.get('myNews'):news=asset(lock['myNews'],cache)
 else:raise ValueError('Publish and pin My News before a release, or use --news-candidate for local acceptance')
 news_stage=Path(tempfile.mkdtemp(prefix='news-',dir=cache));extract(news,news_stage)
 news_contract=news_component(news_stage,allow_candidate=bool(a.news_candidate))
 if a.commute_candidate:commute=Path(a.commute_candidate).resolve()
 elif lock.get('myCommute'):commute=asset(lock['myCommute'],cache)
 else:raise ValueError('Publish and pin My Commute or pass --commute-candidate for local acceptance')
 commute_stage=Path(tempfile.mkdtemp(prefix='commute-',dir=cache));extract(commute,commute_stage)
 commute_contract=commute_component(commute_stage,allow_candidate=bool(a.commute_candidate))
 if a.survival_candidate:survival=Path(a.survival_candidate).resolve()
 elif lock.get('mySurvival'):survival=asset(lock['mySurvival'],cache)
 else:raise ValueError('Publish and pin My Survival before a release, or use --survival-candidate for local acceptance')
 survival_stage=Path(tempfile.mkdtemp(prefix='survival-',dir=cache));extract(survival,survival_stage)
 if a.podcasts_candidate:podcasts=Path(a.podcasts_candidate).resolve()
 elif lock.get('myPodcasts'):podcasts=asset(lock['myPodcasts'],cache)
 else:raise ValueError('Publish and pin My Podcasts or pass --podcasts-candidate for local acceptance')
 podcasts_stage=Path(tempfile.mkdtemp(prefix='podcasts-',dir=cache));extract(podcasts,podcasts_stage)
 podcasts_contract=podcasts_component(podcasts_stage,allow_candidate=bool(a.podcasts_candidate))
 contracts=[sunnymoon_component(ROOT,a.go,cache/'sunnymoon','0.3.0-preview.126'),pet_component(ROOT,a.go,cache/'pet'),gamestudio_component(ROOT,a.go,cache/'gamestudio','0.3.0-preview.123'),podcasts_contract,websearch_component(ROOT,a.go,cache/'websearch','0.3.0-preview.126'),survival_component(survival_stage,allow_candidate=bool(a.survival_candidate)),news_contract,commute_contract,books_contract,contract,ticket_contract,authenticator_component(auth_stage),voice_component(voice_stage,allow_candidate=bool(a.voice_candidate))]
 contracts.append(froststock_component(ROOT,a.go,cache/'froststock',a.version))
 env=os.environ.copy();env.update(CGO_ENABLED='0',GOTOOLCHAIN='local')
 def go(*args):subprocess.run([a.go,*args],cwd=ROOT,env=env,check=True)
 go('build','-buildvcs=false','-trimpath','-ldflags',f'-s -w -X main.version={a.version}','-o',str(out/'Workspace.exe'),'./cmd/workspace')
 with zipfile.ZipFile(seed) as base:
  cat=json.loads(base.read('catalog.json'));old=json.loads(json.dumps(cat));cat['version']=a.version;cat['installerIntegration']=18
  for key,item in old['packages'].items():
   h=hashlib.sha256();size=0
   with base.open(item['file']) as src:
    for b in iter(lambda:src.read(4194304),b''):h.update(b);size+=len(b)
   assert size==item['size'] and h.hexdigest()==item['sha256'],key
  account_app=next(x for x in cat['apps'] if x['id']=='accounts');account_file=cat['packages'][account_app['package']]['file'];account_app['version']=lock['accounts']['tag'];account_app['source']='https://github.com/'+lock['accounts']['repository']+'/releases/tag/'+lock['accounts']['tag']
  replacements={account_file:accounts,'Workspace.exe':out/'Workspace.exe','licenses/components/Workspace-release-status.md':ROOT/'docs/RELEASE-STATUS.md'}
  replacements.update(hosting_replacements(cat,hosting_stage,lock['myHosting']))
  replacements.update(myai_replacements(cat,myai_stage,lock['myAI']))
  replacements.update(tube_replacements(cat,tube_stage,lock['myTube']))
  replacements.update(slicer_replacements(cat,slicer_stage,lock['mySlicer']))
  replacements.update(finance_replacements(cat,finance_stage,lock['myFinance']))
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
   old_files={pkg['file'] for pkg in old['packages'].values()}
   for rel,source in replacements.items():
    if rel not in old_files:dst.write(source,rel)
   for contract in contracts:
    for item in contract['inputs'].values():dst.write(item['source'],item['file'])
 shutil.copyfile(payload,out/'workspace-bundle.zip')
 go('build','-buildvcs=false','-trimpath','-ldflags','-s -w -X github.com/oneworkspacestudio/one-workspace/internal/distribution.Repository=oneworkspacestudio/one-workspace','-o',str(out/'WorkspaceSetup-windows-amd64.exe'),'./cmd/setup')
 (out/'managed-build.json').write_text(json.dumps({'schema':1,'version':a.version,'seed':lock['managedSeed'],'remote':lock['myRemote'],'tickets':lock['myTickets'],'accounts':lock['accounts'],'hosting':lock['myHosting'],'voice':lock['myVoice'],'authenticator':lock['oneAuthenticator'],'books':lock.get('myBooks'),'news':lock.get('myNews'),'commute':lock.get('myCommute'),'myAI':lock['myAI'],'myTube':lock['myTube'],'mySlicer':lock['mySlicer'],'myFinance':lock['myFinance'],'survival':lock.get('mySurvival'),'podcasts':lock.get('myPodcasts'),'candidate':bool(a.podcasts_candidate or a.myai_candidate or a.survival_candidate or a.commute_candidate or a.news_candidate or a.books_candidate or a.tickets_candidate or a.authenticator_candidate or a.accounts_candidate or a.voice_candidate),'priorAppContractsRetainedExceptAccountsHostingMyAIMyTubeSlicerAndFinanceVersions':True,'priorPackagesRetainedExceptManagerAccountsHostingMyAIMyTubeSlicerFinanceAndReleaseNotes':True,'completeInventory':len(cat['apps']),'builder':'scripts/build-managed-seed.py'},indent=2),encoding='utf-8')
 print('Verified managed bundle:',len(cat['apps']),'apps;',len(cat['packages']),'packages',flush=True)
if __name__=='__main__':main()
