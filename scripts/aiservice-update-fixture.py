"""Synthetic app-only update archive; never publish."""
from pathlib import Path
import json,sys,hashlib,zipfile
config,exe,myai,out=map(Path,sys.argv[1:]);c=json.loads(config.read_text(encoding='utf-8'));cat=c['catalog'];v='0.1.0-preview.27-test'
for app in cat['apps']:
 if app['id'] in ['aiservice','myai']:app.update(version=v,source='https://github.com/oneworkspacestudio/my-ai/releases/tag/v'+v)
blobs={}
for key,p in cat['packages'].items():
 data=(Path(c['programs'])/p['file']).read_bytes()
 if key=='aiservice':data=exe.read_bytes()
 if key=='myai':data=myai.read_bytes()
 if key in ['aiservice-contract','myai-contract']:
  obj=json.loads(data);obj['version']=v;data=json.dumps(obj).encode()
 p.update(size=len(data),sha256=hashlib.sha256(data).hexdigest());blobs[p['file']]=data
with zipfile.ZipFile(out,'w',zipfile.ZIP_DEFLATED) as z:
 z.writestr('catalog.json',json.dumps(cat))
 for name,data in blobs.items():z.writestr(name,data)
print('Prepared synthetic AI Service update fixture')
