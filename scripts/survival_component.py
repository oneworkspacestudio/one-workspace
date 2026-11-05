"""Strict My Survival package ownership and offline size budget."""
import hashlib,json,re
from pathlib import Path

def survival_component(root,allow_candidate=False):
 root=Path(root).resolve();c=json.loads((root/'compatibility.json').read_text(encoding='utf-8-sig'))
 assert c['schema']==1 and c['id']=='survival' and c['route']=='/survival' and c['oidc']
 assert c['workspaceInstallerIntegration']==11 and c['minimumWorkspaceVersion']=='0.3.0-preview.117'
 assert c['singleExecutable'] and c['capabilities']==['offline-survival-v1'] and c['publicHandbook']
 m=json.loads((root/'manifest.json').read_text(encoding='utf-8-sig'))
 exe=root/'MySurvival.exe';assert m['version']==c['version'] and m['executableSHA256']==hashlib.sha256(exe.read_bytes()).hexdigest()
 assert 0<m['offlineBytes']<=20*1024*1024
 if not allow_candidate:assert m.get('candidate') is False and re.fullmatch('[a-f0-9]{40}',m['sourceCommit'])
 inputs={};size=0
 for p in sorted(root.rglob('*')):
  if not p.is_file():continue
  rel=p.relative_to(root).as_posix();size+=p.stat().st_size
  if rel=='MySurvival.exe':key,dest='survival',rel
  elif rel in ['LICENSE','NOTICE','THIRD-PARTY-LICENSES.txt']:key,dest='survival-'+rel,'licenses/MySurvival-'+rel
  elif rel in ['compatibility.json','manifest.json']:key,dest='survival-'+('contract' if rel=='compatibility.json' else 'manifest'),'docs/MySurvival/'+rel
  elif rel in ['README.md','My-Survival-offline.html','docs/INTEGRATION.md','docs/ACCEPTANCE.md','docs/SOURCES.md']:key,dest='survival-doc-'+hashlib.sha256(rel.encode()).hexdigest()[:16],'docs/MySurvival/'+rel
  else:raise ValueError('Unexpected My Survival payload '+rel)
  inputs[key]={'source':str(p),'file':dest}
 assert size<=150*1024*1024
 assert all(k in inputs for k in ['survival','survival-contract','survival-manifest','survival-LICENSE','survival-NOTICE'])
 app={'id':'survival','name':'My Survival','description':'Offline emergency guidance, household supplies and an encrypted survival plan','version':c['version'],'kind':'process','package':'survival','executable':'MySurvival.exe','port':9041,'route':'/survival','stripPrefix':True,'health':'/healthz','depends':['accounts'],'oidc':True,'callback':'/survival/auth/callback','args':['--data','${DATA}/survival','--listen','127.0.0.1:${PORT}','--origin','${ORIGIN}/survival','--no-browser'],'source':'https://github.com/oneworkspacestudio/my-survival/releases/tag/v'+c['version']}
 return {'app':app,'inputs':inputs}
