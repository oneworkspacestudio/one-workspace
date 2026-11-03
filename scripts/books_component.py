"""My Books managed contract, including independent Files storage ownership."""
import hashlib,json,re
from pathlib import Path

def books_component(root,allow_candidate=False):
 root=Path(root).resolve();c=json.loads((root/'compatibility.json').read_text(encoding='utf-8-sig'))
 assert c['schema']==1 and c['id']=='books' and c['route']=='/books' and c['oidc']
 assert c['workspaceInstallerIntegration']==8 and c['minimumWorkspaceVersion']=='0.3.0-preview.113'
 assert c['singleExecutable'] is True and c['capabilities']==['native-books-v1']
 m=json.loads((root/'manifest.json').read_text(encoding='utf-8-sig'))
 assert m['version']==c['version'] and m['executableSHA256']==hashlib.sha256((root/'MyBooks.exe').read_bytes()).hexdigest()
 if not allow_candidate:assert m.get('candidate') is False and re.fullmatch('[a-f0-9]{40}',m['sourceCommit'])
 inputs={}
 for p in sorted(root.rglob('*')):
  if not p.is_file():continue
  rel=p.relative_to(root).as_posix()
  if rel=='MyBooks.exe':key,dest='books',rel
  elif rel in ['LICENSE','NOTICE','THIRD-PARTY-LICENSES.txt']:key,dest='books-'+rel,'licenses/MyBooks-'+rel
  elif rel in ['compatibility.json','manifest.json']:key,dest='books-'+('contract' if rel=='compatibility.json' else 'manifest'),'docs/MyBooks/'+rel
  elif rel.startswith('licenses/'):key,dest='books-license-'+hashlib.sha256(rel.encode()).hexdigest()[:16],'licenses/MyBooks/'+rel[9:]
  elif rel in ['README.md','docs/RESEARCH.md','docs/ACCEPTANCE.md','docs/SECURITY.md']:key,dest='books-doc-'+hashlib.sha256(rel.encode()).hexdigest()[:16],'docs/MyBooks/'+rel
  elif rel=='web/public/books.svg' or (rel.startswith('docs/screenshots/') and p.suffix=='.png') or (rel.startswith('docs/evidence/') and p.suffix=='.json'):key,dest='books-doc-'+hashlib.sha256(rel.encode()).hexdigest()[:16],'docs/MyBooks/'+rel
  else:raise ValueError('Unexpected Books payload '+rel)
  inputs[key]={'source':str(p),'file':dest}
 assert all(k in inputs for k in ['books','books-contract','books-manifest','books-LICENSE','books-NOTICE'])
 app={'id':'books','name':'My Books','description':'Your personal bookshelf, with comfortable reading on every screen','version':c['version'],'kind':'process','package':'books','executable':'MyBooks.exe','port':9029,'route':'/books','stripPrefix':True,'health':'/healthz','depends':['accounts'],'oidc':True,'callback':'/books/auth/callback','args':['--data','${DATA}/books','--listen','127.0.0.1:${PORT}','--origin','${ORIGIN}/books','--no-browser'],'env':{'MYBOOKS_LIBRARY_ROOT':'${FILES}/My Books'},'source':'https://github.com/oneworkspacestudio/my-books/releases/tag/v'+c['version']}
 return {'app':app,'inputs':inputs}
