"""Publisher-side contract for My Remote; no user-machine paths enter a catalog."""
import hashlib,json
from pathlib import Path

def remote_component(root):
 root=Path(root).resolve()
 c=json.loads((root/'compatibility.json').read_text(encoding='utf-8-sig'))
 assert c['schema']==1 and c['id']=='remote' and c['route']=='/remote' and c['oidc']
 assert c['ownerStreamRevocation']=='manager-v1'
 assert c['workspaceInstallerIntegration']==3 and c['minimumWorkspaceVersion']=='0.3.0-preview.104'
 inputs={}
 for p in sorted(root.rglob('*')):
  if not p.is_file():continue
  rel=p.relative_to(root).as_posix()
  if rel=='MyRemote.exe':key,dest='remote',rel
  elif rel.startswith('remote-runtime/'):key,dest='remote-runtime-'+hashlib.sha256(rel.encode()).hexdigest()[:16],rel
  elif rel.startswith('licenses/'):key,dest='remote-license-'+hashlib.sha256(rel.encode()).hexdigest()[:16],'licenses/MyRemote/'+rel[9:]
  elif rel in ['LICENSE','NOTICE','THIRD-PARTY-LICENSES.txt']:key,dest='remote-'+rel,'licenses/MyRemote-'+rel
  elif rel in ['compatibility.json','manifest.json']:key,dest='remote-'+('contract' if rel=='compatibility.json' else 'manifest'),'docs/MyRemote/'+rel
  else:raise ValueError('Unexpected managed package member '+rel)
  inputs[key]={'source':str(p),'file':dest}
 assert 'remote' in inputs and 'remote-contract' in inputs
 assert all(any(x['file']==f for x in inputs.values()) for f in ['remote-runtime/node.exe','remote-runtime/engine/local-server.cjs','remote-runtime/engine/rfb.cjs','remote-runtime/engine/security.cjs','remote-runtime/engine/lifecycle.cjs','remote-runtime/node_modules/ws/index.js','remote-runtime/host/tvnserver.exe','remote-runtime/scripts/configure-host.ps1','remote-runtime/scripts/verify-host.ps1','remote-runtime/scripts/create-host-secret.cjs','remote-runtime/scripts/migrate-launchers.ps1'])
 app={'id':'remote','name':'My Remote','description':'Access this PC with owner-approved desktop sharing','version':c['version'],'kind':'process','package':'remote','executable':'MyRemote.exe','port':9025,'route':'/remote','stripPrefix':True,'health':'/healthz','depends':['accounts'],'oidc':True,'callback':'/remote/auth/callback','args':['--data','${DATA}/remote','--listen','127.0.0.1:${PORT}','--origin','${ORIGIN}/remote','--no-browser'],'source':'https://github.com/oneworkspacestudio/my-remote/releases/tag/v'+c['version']}
 return {'app':app,'inputs':inputs}

if __name__=='__main__':
 import sys
 Path(sys.argv[2]).write_text(json.dumps(remote_component(sys.argv[1]),indent=2),encoding='utf-8')
