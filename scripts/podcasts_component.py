"""My Podcasts package: one executable with embedded UI and private apphost state."""
import hashlib,json,re
from pathlib import Path

def podcasts_component(root,allow_candidate=False):
 root=Path(root).resolve();c=json.loads((root/'compatibility.json').read_text(encoding='utf-8-sig'));m=json.loads((root/'manifest.json').read_text(encoding='utf-8-sig'))
 assert c['schema']==1 and c['id']=='podcasts' and c['route']=='/podcasts' and c['oidc'] is True
 assert c['workspaceInstallerIntegration']==14 and c['minimumWorkspaceVersion']=='0.3.0-preview.122'
 assert c['singleExecutable'] is True and c['capabilities']==['native-podcasts-v1']
 assert m['version']==c['version'] and m['executableSHA256']==hashlib.sha256((root/'MyPodcasts.exe').read_bytes()).hexdigest()
 if not allow_candidate:assert m.get('candidate') is False and re.fullmatch('[a-f0-9]{40}',m.get('sourceCommit',''))
 inputs={}
 for p in sorted(root.rglob('*')):
  if not p.is_file():continue
  rel=p.relative_to(root).as_posix()
  if rel=='MyPodcasts.exe':key,dest='podcasts',rel
  elif rel in ['LICENSE','NOTICE','THIRD-PARTY-NOTICES.md']:key,dest='podcasts-'+rel,'licenses/MyPodcasts-'+rel
  elif rel in ['compatibility.json','manifest.json']:key,dest='podcasts-'+('contract' if rel=='compatibility.json' else 'manifest'),'docs/MyPodcasts/'+rel
  elif rel.startswith('licenses/'):key,dest='podcasts-license-'+hashlib.sha256(rel.encode()).hexdigest()[:16],'licenses/MyPodcasts/'+rel[9:]
  elif rel in ['README.md','docs/ACCEPTANCE.md','docs/DATA-SOURCES.md','docs/SECURITY.md','docs/PRODUCTION.md']:key,dest='podcasts-doc-'+hashlib.sha256(rel.encode()).hexdigest()[:16],'docs/MyPodcasts/'+rel
  else:raise ValueError('Unexpected Podcasts payload '+rel)
  inputs[key]={'source':str(p),'file':dest}
 assert all(k in inputs for k in ['podcasts','podcasts-contract','podcasts-manifest','podcasts-LICENSE','podcasts-NOTICE'])
 app={'id':'podcasts','name':'My Podcasts','description':c.get('description','Create researched technology podcasts with private background production'),'version':c['version'],'kind':'process','package':'podcasts','executable':'MyPodcasts.exe','port':9044,'route':'/podcasts','stripPrefix':True,'health':'/healthz','depends':['accounts'],'oidc':True,'callback':'/podcasts/auth/callback','args':['--data','${DATA}/podcasts','--listen','127.0.0.1:${PORT}','--origin','${ORIGIN}/podcasts','--no-browser'],'source':'https://github.com/oneworkspacestudio/my-podcasts/releases/tag/v'+c['version']}
 return {'app':app,'inputs':inputs}
