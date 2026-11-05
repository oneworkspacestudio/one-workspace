"""My News managed contract, with private reading-library storage."""
import hashlib,json,re
from pathlib import Path

def news_component(root,allow_candidate=False):
 root=Path(root).resolve();c=json.loads((root/'compatibility.json').read_text(encoding='utf-8-sig'))
 assert c['schema']==1 and c['id']=='news' and c['route']=='/news' and c['oidc']
 assert c['workspaceInstallerIntegration']==9 and c['minimumWorkspaceVersion']=='0.3.0-preview.114'
 assert c['singleExecutable'] is True and c['capabilities']==['native-news-v1']
 m=json.loads((root/'manifest.json').read_text(encoding='utf-8-sig'))
 assert m['version']==c['version'] and m['executableSHA256']==hashlib.sha256((root/'MyNews.exe').read_bytes()).hexdigest()
 if not allow_candidate:assert m.get('candidate') is False and re.fullmatch('[a-f0-9]{40}',m['sourceCommit'])
 inputs={}
 for p in sorted(root.rglob('*')):
  if not p.is_file():continue
  rel=p.relative_to(root).as_posix()
  if rel=='MyNews.exe':key,dest='news',rel
  elif rel in ['LICENSE','NOTICE','THIRD-PARTY-LICENSES.txt']:key,dest='news-'+rel,'licenses/MyNews-'+rel
  elif rel in ['compatibility.json','manifest.json']:key,dest='news-'+('contract' if rel=='compatibility.json' else 'manifest'),'docs/MyNews/'+rel
  elif rel.startswith('licenses/'):key,dest='news-license-'+hashlib.sha256(rel.encode()).hexdigest()[:16],'licenses/MyNews/'+rel[9:]
  elif rel in ['README.md','docs/RESEARCH.md','docs/ACCEPTANCE.md','docs/SECURITY.md']:key,dest='news-doc-'+hashlib.sha256(rel.encode()).hexdigest()[:16],'docs/MyNews/'+rel
  elif rel=='web/public/news.svg' or (rel.startswith('docs/screenshots/') and p.suffix=='.png') or (rel.startswith('docs/evidence/') and p.suffix=='.json'):key,dest='news-doc-'+hashlib.sha256(rel.encode()).hexdigest()[:16],'docs/MyNews/'+rel
  else:raise ValueError('Unexpected News payload '+rel)
  inputs[key]={'source':str(p),'file':dest}
 assert all(k in inputs for k in ['news','news-contract','news-manifest','news-LICENSE','news-NOTICE'])
 app={'id':'news','name':'My News','description':'News from around the world, with saved reading copies for offline','version':c['version'],'kind':'process','package':'news','executable':'MyNews.exe','port':9030,'route':'/news','stripPrefix':True,'health':'/healthz','depends':['accounts'],'oidc':True,'callback':'/news/auth/callback','args':['--data','${DATA}/news','--listen','127.0.0.1:${PORT}','--origin','${ORIGIN}/news','--no-browser'],'source':'https://github.com/oneworkspacestudio/my-news/releases/tag/v'+c['version']}
 return {'app':app,'inputs':inputs}
