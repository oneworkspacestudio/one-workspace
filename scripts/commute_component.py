"""My Commute package: one executable with embedded UI and private apphost state."""
import hashlib,json,re
from pathlib import Path

def commute_component(root,allow_candidate=False):
 root=Path(root).resolve();c=json.loads((root/'compatibility.json').read_text(encoding='utf-8-sig'));m=json.loads((root/'manifest.json').read_text(encoding='utf-8-sig'))
 assert c['schema']==1 and c['id']=='commute' and c['route']=='/commute' and c['oidc'] is True
 assert c['workspaceInstallerIntegration']==10 and c['minimumWorkspaceVersion']=='0.3.0-preview.116'
 assert c['singleExecutable'] is True and c['capabilities']==['native-commute-v1']
 assert m['version']==c['version'] and m['executableSHA256']==hashlib.sha256((root/'MyCommute.exe').read_bytes()).hexdigest()
 if not allow_candidate:assert m.get('candidate') is False and re.fullmatch('[a-f0-9]{40}',m.get('sourceCommit',''))
 inputs={}
 for p in sorted(root.rglob('*')):
  if not p.is_file():continue
  rel=p.relative_to(root).as_posix()
  if rel=='MyCommute.exe':key,dest='commute',rel
  elif rel in ['LICENSE','NOTICE','THIRD-PARTY-NOTICES.md']:key,dest='commute-'+rel,'licenses/MyCommute-'+rel
  elif rel in ['compatibility.json','manifest.json']:key,dest='commute-'+('contract' if rel=='compatibility.json' else 'manifest'),'docs/MyCommute/'+rel
  elif rel.startswith('licenses/'):key,dest='commute-license-'+hashlib.sha256(rel.encode()).hexdigest()[:16],'licenses/MyCommute/'+rel[9:]
  elif rel in ['README.md','docs/ACCEPTANCE.md','docs/DATA-SOURCES.md','docs/SECURITY.md']:key,dest='commute-doc-'+hashlib.sha256(rel.encode()).hexdigest()[:16],'docs/MyCommute/'+rel
  else:raise ValueError('Unexpected Commute payload '+rel)
  inputs[key]={'source':str(p),'file':dest}
 assert all(k in inputs for k in ['commute','commute-contract','commute-manifest','commute-LICENSE','commute-NOTICE'])
 app={'id':'commute','name':'My Commute','description':'Follow live London buses, Tube and trains and plan your next connection','version':c['version'],'kind':'process','package':'commute','executable':'MyCommute.exe','port':9040,'route':'/commute','stripPrefix':True,'health':'/healthz','depends':['accounts'],'oidc':True,'callback':'/commute/auth/callback','args':['--data','${DATA}/commute','--listen','127.0.0.1:${PORT}','--origin','${ORIGIN}/commute','--no-browser'],'source':'https://github.com/oneworkspacestudio/my-commute/releases/tag/v'+c['version']}
 return {'app':app,'inputs':inputs}
