"""Strict public My Voice component; only package-owned program/document files."""
import json,hashlib,re
from pathlib import Path
def voice_component(root,allow_candidate=False):
 root=Path(root).resolve();c=json.loads((root/'compatibility.json').read_text(encoding='utf-8-sig'))
 assert c['schema']==1 and c['id']=='voice' and c['route']=='/voice' and c['oidc']
 assert c['workspaceInstallerIntegration']==7 and c['minimumWorkspaceVersion']=='0.3.0-preview.111'
 assert c['singleExecutable'] is True and c['capabilities']==['native-voice-v1']
 m=json.loads((root/'manifest.json').read_text(encoding='utf-8-sig'))
 assert m['schema']==1 and m['version']==c['version']
 assert m['executableSHA256']==hashlib.sha256((root/'MyVoice.exe').read_bytes()).hexdigest()
 if not allow_candidate:assert m.get('candidate') is False and m.get('builtNatively') is True and re.fullmatch('[a-f0-9]{40}',m['sourceCommit'])
 inputs={}
 for p in sorted(root.rglob('*')):
  if not p.is_file():continue
  rel=p.relative_to(root).as_posix()
  if rel=='MyVoice.exe':key,dest='voice',rel
  elif rel.startswith('licenses/'):key,dest='voice-license-'+hashlib.sha256(rel.encode()).hexdigest()[:16],'licenses/MyVoice/'+rel[9:]
  elif rel in ['LICENSE','NOTICE','THIRD-PARTY-NOTICES.md']:key,dest='voice-'+rel,'licenses/MyVoice-'+rel
  elif rel in ['compatibility.json','manifest.json']:key,dest='voice-'+('contract' if rel=='compatibility.json' else 'manifest'),'docs/MyVoice/'+rel
  elif rel in ['README.md','INSTALL-WINDOWS.md','MODEL-LICENSES.md','TEST-REPORT-WINDOWS.md','ARCHITECTURE.md','SECURITY.md','CHANGELOG.md','docs/API.md','docs/go-dependencies.json','docs/web-dependencies.json','docs/evidence/windows-native.json','docs/pocket-dependencies.json']:key,dest='voice-doc-'+hashlib.sha256(rel.encode()).hexdigest()[:16],'docs/MyVoice/'+rel
  else:raise ValueError('Unexpected My Voice package member '+rel)
  inputs[key]={'source':str(p),'file':dest}
 assert all(k in inputs for k in ['voice','voice-contract','voice-manifest','voice-LICENSE','voice-NOTICE'])
 app={'id':'voice','name':'My Voice','description':'Turn text into speech and create private voices on your PC','version':c['version'],'kind':'process','package':'voice','executable':'MyVoice.exe','port':9028,'route':'/voice','stripPrefix':True,'health':'/healthz','depends':['accounts'],'oidc':True,'callback':'/voice/auth/callback','args':['--data','${DATA}/voice','--listen','127.0.0.1:${PORT}','--origin','${ORIGIN}/voice','--no-browser'],'source':'https://github.com/oneworkspacestudio/my-voice/releases/tag/v'+c['version']}
 return {'app':app,'inputs':inputs}
