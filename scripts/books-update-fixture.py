"""Synthetic app-only upgrade fixture. Never publish this archive."""
from pathlib import Path
import json,sys,hashlib,zipfile
config,base,exe,out=map(Path,sys.argv[1:]);c=json.loads(config.read_text(encoding='utf-8'))['catalog'];v='0.1.0-preview.2-test'
a=next(a for a in c['apps'] if a['id']=='books');a.update(version=v,source='https://github.com/oneworkspacestudio/my-books/releases/tag/v'+v)
with zipfile.ZipFile(base) as z:
 blobs={p['file']:z.read(p['file']) for p in c['packages'].values()}
blobs['MyBooks.exe']=exe.read_bytes()
for key in ['books-contract','books-manifest']:
 p=c['packages'][key]['file'];obj=json.loads(blobs[p]);obj['version']=v
 if key=='books-manifest':obj['executableSHA256']=hashlib.sha256(blobs['MyBooks.exe']).hexdigest();obj['candidate']=True
 blobs[p]=json.dumps(obj).encode()
for p in c['packages'].values():b=blobs[p['file']];p.update(size=len(b),sha256=hashlib.sha256(b).hexdigest())
with zipfile.ZipFile(out,'w',zipfile.ZIP_DEFLATED) as z:
 z.writestr('catalog.json',json.dumps(c));
 for name,b in blobs.items():z.writestr(name,b)
print('Built isolated test-only Books app upgrade fixture')
