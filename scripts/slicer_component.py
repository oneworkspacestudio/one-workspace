"""Replace only My Slicer's existing owned payloads in a managed-seed build."""
import hashlib,json,re
from pathlib import Path

def slicer_replacements(catalog,stage,pin):
 root=Path(stage).resolve()
 c=json.loads((root/'compatibility.json').read_text(encoding='utf-8-sig'))
 b=json.loads((root/'build-info.json').read_text(encoding='utf-8-sig'))
 assert pin['repository']=='oneworkspacestudio/my-slicer' and pin['tag']=='v'+c['version']
 assert c['schema']==1 and c['app']=='slicer' and c['route']=='/slicer' and c['executable']=='MySlicer.exe'
 assert c['callback']=='/slicer/auth/callback' and c['health']=='/healthz' and c['workspaceInstallerIntegration']==2
 assert b['version']==c['version'] and b['githubActionsUsed'] is False and re.fullmatch('[a-f0-9]{40}',b['sourceCommit'])
 app=next(x for x in catalog['apps'] if x['id']=='slicer')
 assert app['package']=='slicer' and app['executable']==c['executable'] and app['route']==c['route'] and app['oidc'] is True
 mapping={'slicer':('MySlicer.exe','MySlicer.exe')}
 for name in ['LICENSE','NOTICE','THIRD-PARTY-LICENSES.txt']:mapping['slicer-'+name]=('licenses/MySlicer-'+name,name)
 mapping['slicer-contract']=('docs/MySlicer/compatibility.json','compatibility.json')
 for name in ['SETUP.md','ACCEPTANCE.md','ENGINE.md']:mapping['slicer-doc-'+name]=('docs/MySlicer/'+name,'docs/'+name)
 assert {k for k in catalog['packages'] if k=='slicer' or k.startswith('slicer-')}==set(mapping),'Unexpected Slicer package ownership'
 replacements={}
 for key,(target,member) in mapping.items():
  payload=root/member
  assert payload.is_file() and not payload.is_symlink(),'Missing Slicer payload '+member
  assert catalog['packages'][key]['file']==target,('Unexpected Slicer path',key)
  replacements[target]=payload
 # Change the app version/source only after every owned payload is validated.
 for key,(target,member) in mapping.items():
  payload=replacements[target]
  catalog['packages'][key]={'file':target,'size':payload.stat().st_size,'sha256':hashlib.sha256(payload.read_bytes()).hexdigest()}
 app['version']=c['version'];app['source']='https://github.com/'+pin['repository']+'/releases/tag/'+pin['tag']
 return replacements
