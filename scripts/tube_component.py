"""Consume the public My Tube pin without changing other managed app payloads."""
import hashlib,json,re
from pathlib import Path

def tube_replacements(catalog,stage,pin):
 root=Path(stage).resolve()
 info=json.loads((root/'build-info.json').read_text(encoding='utf-8-sig'))
 assert pin['repository']=='oneworkspacestudio/my-tube'
 assert info['repository']==pin['repository'] and info['tag']==pin['tag']=='v'+info['version']
 assert info['sourceDirty'] is False and info['githubActionsUsed'] is False
 assert re.fullmatch('[a-f0-9]{40}',info['sourceCommit'])
 assert info['exeSHA256']==hashlib.sha256((root/'MyTube.exe').read_bytes()).hexdigest()
 app=next(a for a in catalog['apps'] if a['id']=='my-tube')
 assert app['package']=='my-tube' and app['executable']=='MyTube.exe'
 assert app['route']=='/my-tube' and app['oidc'] is True
 mapping={'my-tube':('MyTube.exe','MyTube.exe')}
 for name in ['LICENSE','NOTICE','THIRD-PARTY-LICENSES.txt']:
  mapping['my-tube-'+name]=('licenses/MyTube-'+name,name)
 assert {k for k in catalog['packages'] if k=='my-tube' or k.startswith('my-tube-')}==set(mapping),'Unexpected My Tube package ownership'
 replacements={}
 for key,(target,member) in mapping.items():
  payload=root/member
  assert payload.is_file() and not payload.is_symlink()
  assert catalog['packages'][key]['file']==target,('Unexpected My Tube destination',key)
  catalog['packages'][key]={'file':target,'size':payload.stat().st_size,'sha256':hashlib.sha256(payload.read_bytes()).hexdigest()}
  replacements[target]=payload
 app['version']=info['version']
 app['source']='https://github.com/'+pin['repository']+'/releases/tag/'+pin['tag']
 return replacements
