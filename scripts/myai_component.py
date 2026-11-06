"""Replace only the pinned My AI namespace in the ordinary managed-seed build."""
import hashlib,json,re
from pathlib import Path

def myai_replacements(catalog,stage,pin):
 root=Path(stage).resolve()/'MyAI'
 c=json.loads((root/'compatibility.json').read_text(encoding='utf-8-sig'))
 b=json.loads((root/'build-info.json').read_text(encoding='utf-8-sig'))
 assert pin['repository']=='oneworkspacestudio/my-ai' and pin['tag']=='v'+c['version']
 assert c['schema']==1 and c['id']=='myai' and c['route']=='/myai' and c['executable']=='MyAI.exe' and c['health']=='/healthz' and c['oidc'] is True and c['embeddedFrontend'] is True
 assert c['workspaceInstallerIntegration']==12 and c['minimumWorkspaceVersion']=='0.3.0-preview.118'
 assert b['schema']==1 and b['version']==c['version'] and b['sourceDirty'] is False and b['githubActionsUsed'] is False and re.fullmatch('[a-f0-9]{40}',b['commit'])
 assert b['exeSHA256']==hashlib.sha256((root/'MyAI.exe').read_bytes()).hexdigest()
 app=next(a for a in catalog['apps'] if a['id']=='myai')
 assert app['package']=='myai' and app['route']==c['route'] and app['executable']==c['executable'] and app['oidc'] is True
 mapping={'myai':('MyAI.exe','MyAI.exe'),'myai-contract':('docs/MyAI/compatibility.json','compatibility.json')}
 for name in ['LICENSE','NOTICE','THIRD-PARTY-LICENSES.txt']:mapping['myai-'+name]=('licenses/MyAI-'+name,name)
 for name in ['SETUP.md','SUPPORT.md','ADR-001-architecture.md','SECURITY-CHECKS.md']:mapping['myai-doc-'+name]=('docs/MyAI/'+name,'docs/'+name)
 assert {k for k in catalog['packages'] if k=='myai' or k.startswith('myai-')}<=set(mapping),'Unrecognized existing My AI payload ownership'
 service=json.loads((root/'aiservice-compatibility.json').read_text(encoding='utf-8-sig'))
 assert service['id']=='aiservice' and service['route']=='/aiservice' and service['executable']=='AIService.exe' and service['oidc'] is True and service['workspaceInstallerIntegration']==12
 assert b['serviceSHA256']==hashlib.sha256((root/'AIService.exe').read_bytes()).hexdigest()
 service_app=json.loads((root/'aiservice-app.json').read_text(encoding='utf-8-sig'));service_app.update(version=c['version'],source='https://github.com/'+pin['repository']+'/releases/tag/'+pin['tag'])
 existing=next((a for a in catalog['apps'] if a['id']=='aiservice'),None)
 if existing is None:catalog['apps'].append(service_app)
 else:
  assert {k:v for k,v in existing.items() if k not in ['version','source']}=={k:v for k,v in service_app.items() if k not in ['version','source']}
  existing.update(version=service_app['version'],source=service_app['source'])
 mapping.update({'aiservice':('AIService.exe','AIService.exe'),'aiservice-contract':('docs/AIService/compatibility.json','aiservice-compatibility.json')})
 for name in ['LICENSE','NOTICE','THIRD-PARTY-LICENSES.txt']:mapping['aiservice-'+name]=('licenses/AIService-'+name,name)
 replacements={}
 for key,(target,member) in mapping.items():
  payload=root/member;assert payload.is_file() and not payload.is_symlink()
  if key in catalog['packages']:assert catalog['packages'][key]['file']==target,('Unexpected My AI ownership',key)
  catalog['packages'][key]={'file':target,'size':payload.stat().st_size,'sha256':hashlib.sha256(payload.read_bytes()).hexdigest()}
  replacements[target]=payload
 app['version']=c['version'];app['source']='https://github.com/'+pin['repository']+'/releases/tag/'+pin['tag']
 return replacements
