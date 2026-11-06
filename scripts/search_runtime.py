"""Checked build-time portable browser/server preparation; no startup downloads."""
import hashlib,json,shutil,urllib.request,zipfile
from pathlib import Path
def digest(p):
 h=hashlib.sha256()
 with Path(p).open('rb') as f:
  for b in iter(lambda:f.read(4<<20),b''):h.update(b)
 return h.hexdigest()
def prepare_runtime(lock,out,cache):
 out,cache=Path(out),Path(cache);cache.mkdir(parents=True,exist_ok=True)
 for kind in ['browser','llama']:
  pin=lock[kind];archive=cache/(pin['sha256']+'.zip')
  if not archive.exists():
   with urllib.request.urlopen(pin['url'],timeout=180) as src,archive.open('wb') as dst:shutil.copyfileobj(src,dst,4<<20)
  if archive.stat().st_size!=pin['size'] or digest(archive)!=pin['sha256']:raise ValueError('Portable runtime digest mismatch: '+kind)
  with zipfile.ZipFile(archive) as z:
   for info in z.infolist():
    rel=Path(info.filename)
    if rel.is_absolute() or '..' in rel.parts or ':' in info.filename or '\\' in info.filename or (info.external_attr>>16)&0o170000==0o120000:raise ValueError('Unsafe runtime member')
    if kind=='llama' and not (info.filename.endswith('.dll') or info.filename=='llama-server.exe' or info.filename.startswith('LICENSE')):continue
    z.extract(info,out/kind)
 (out/'provenance.json').write_text(json.dumps(lock,indent=2))
 manifest={p.relative_to(out/'browser').as_posix():digest(p) for p in sorted((out/'browser').rglob('*')) if p.is_file()}
 (out/'browser-files.json').write_text(json.dumps(manifest,indent=2))
