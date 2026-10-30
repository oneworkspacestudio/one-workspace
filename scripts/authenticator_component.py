"""Validated managed contract for the independent public OneAccount code generator."""
import hashlib,json
from pathlib import Path

def authenticator_component(root):
 root=Path(root).resolve();c=json.loads((root/'compatibility.json').read_text(encoding='utf-8-sig'))
 assert c['schema']==1 and c['id']=='authenticator' and c['route']=='/authenticator' and c['oidc'] is False
 assert c['singleExecutable'] is True and c['workspaceInstallerIntegration']==6
 assert c['authentication']=='independent-local-encrypted-vault' and c['publicBootstrap']=='static-allowlist-v1'
 inputs={}
 for p in sorted(root.rglob('*')):
  if not p.is_file():continue
  rel=p.relative_to(root).as_posix()
  if rel=='one-authenticator.exe':key,dest='authenticator',rel
  elif rel.startswith('licenses/'):key,dest='authenticator-license-'+hashlib.sha256(rel.encode()).hexdigest()[:16],'licenses/OneAuthenticator/'+rel[9:]
  elif rel in ['LICENSE','NOTICE','THIRD-PARTY-LICENSES.txt']:key,dest='authenticator-'+rel,'licenses/OneAuthenticator-'+rel
  elif rel in ['compatibility.json','manifest.json']:key,dest='authenticator-'+('contract' if rel=='compatibility.json' else 'manifest'),'docs/OneAuthenticator/'+rel
  else:raise ValueError('Unexpected managed Authenticator package member '+rel)
  inputs[key]={'source':str(p),'file':dest}
 assert all(k in inputs for k in ['authenticator','authenticator-contract','authenticator-manifest','authenticator-LICENSE'])
 app={'id':'authenticator','name':'One Authenticator','description':'Keep your OneAccount security codes on your phone, even offline','version':c['version'],'kind':'process','package':'authenticator','executable':'one-authenticator.exe','port':9027,'route':'/authenticator','stripPrefix':True,'health':'/healthz','oidc':False,'args':['--listen','127.0.0.1:${PORT}','--public-origin','${ORIGIN}'],'source':'https://github.com/oneworkspacestudio/one-authenticator/releases/tag/v'+c['version']}
 return {'app':app,'inputs':inputs}
