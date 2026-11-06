"""Synthetic app-only updater/rollback transport, never a public component."""
import hashlib,json,sys,zipfile
from pathlib import Path
cfg=json.loads(Path(sys.argv[1]).read_text(encoding='utf-8'));cat=cfg['catalog'];root=Path(cfg['programs'])
app=next(a for a in cat['apps'] if a['id']=='websearch');app.update(version='0.1.0-preview.2-test',source='https://github.com/oneworkspacestudio/one-workspace/releases/tag/websearch-synthetic-test')
with zipfile.ZipFile(sys.argv[2],'w',compression=zipfile.ZIP_DEFLATED,compresslevel=6) as z:
 for key,p in cat['packages'].items():
  data=(root/p['file']).read_bytes()
  if key=='websearch':data+=b'\nSynthetic update transport marker; not a released component.\n'
  p.update(size=len(data),sha256=hashlib.sha256(data).hexdigest());z.writestr(p['file'],data)
 z.writestr('catalog.json',json.dumps(cat,indent=2))
print('Prepared isolated Web Search update fixture')
