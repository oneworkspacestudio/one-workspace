"""The native, portable Web Search service is built from the suite source.

Browser/model-server payloads are immutable program packages; the verified model
library is installed explicitly into private data and retained across updates.
"""
import hashlib,json,subprocess
from pathlib import Path

def websearch_component(root,go,out,version):
 root,out=Path(root),Path(out);out.mkdir(parents=True,exist_ok=True)
 subprocess.run(['npm.cmd','ci','--ignore-scripts'],cwd=root/'internal/websearch/web',check=True,shell=True)
 subprocess.run(['node','build.mjs'],cwd=root/'internal/websearch/web',check=True)
 subprocess.run([go,'build','-buildvcs=false','-trimpath','-ldflags','-s -w','-o',str(out/'WebSearch.exe'),'./cmd/websearch'],cwd=root,check=True)
 lock=json.loads((root/'internal/websearch/runtime-lock.json').read_text())
 runtime=out/'runtime'
 from search_runtime import prepare_runtime
 prepare_runtime(lock,runtime,out/'downloads')
 compatibility={'schema':1,'id':'websearch','version':'0.1.0-preview.1','workspaceInstallerIntegration':13,'minimumWorkspaceVersion':'0.3.0-preview.120','route':'/websearch','oidc':True,'executable':'WebSearch.exe','health':'/healthz','runtime':'native-portable-v1','modelSetup':'explicit-owner-verified-download','storage':'private-websearch-v1'}
 (out/'compatibility.json').write_text(json.dumps(compatibility,indent=2))
 inputs={'websearch':{'source':str(out/'WebSearch.exe'),'file':'WebSearch.exe'},'websearch-contract':{'source':str(out/'compatibility.json'),'file':'docs/WebSearch/compatibility.json'}}
 for path in sorted(runtime.rglob('*')):
  if not path.is_file():continue
  rel=path.relative_to(runtime).as_posix();key='websearch-runtime-'+hashlib.sha256(rel.encode()).hexdigest()[:20]
  inputs[key]={'source':str(path),'file':'websearch-runtime/'+rel}
 for name in ['LICENSE','NOTICE']:
  inputs['websearch-'+name]={'source':str(root/name),'file':'licenses/WebSearch-'+name}
 for rel in ['docs/WEB_SEARCH.md','docs/licenses/WebSearch-THIRD-PARTY.txt']:
  inputs['websearch-doc-'+Path(rel).name]={'source':str(root/rel),'file':rel}
 app={'id':'websearch','name':'Web Search','description':'Let your LLM research the web in a private browser','version':'0.1.0-preview.1','kind':'process','package':'websearch','executable':'WebSearch.exe','port':9043,'route':'/websearch','stripPrefix':True,'health':'/healthz','depends':['accounts'],'oidc':True,'callback':'/websearch/auth/callback','args':['--data','${DATA}/websearch','--listen','127.0.0.1:${PORT}','--origin','${ORIGIN}/websearch','--no-browser'],'source':'https://github.com/oneworkspacestudio/one-workspace/releases/tag/v'+version}
 return {'app':app,'inputs':inputs}
