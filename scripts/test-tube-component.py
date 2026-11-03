import copy,hashlib,json,tempfile,unittest
from pathlib import Path
from tube_component import tube_replacements

class TubeComponent(unittest.TestCase):
 def fixture(self,d):
  root=Path(d)
  names=['MyTube.exe','LICENSE','NOTICE','THIRD-PARTY-LICENSES.txt']
  for name in names:(root/name).write_text('synthetic '+name)
  info=dict(repository='oneworkspacestudio/my-tube',tag='v0.1.0-preview.6',version='0.1.0-preview.6',sourceDirty=False,githubActionsUsed=False,sourceCommit='a'*40,exeSHA256=hashlib.sha256((root/'MyTube.exe').read_bytes()).hexdigest())
  (root/'build-info.json').write_text(json.dumps(info))
  cat=dict(apps=[dict(id='my-tube',package='my-tube',executable='MyTube.exe',route='/my-tube',oidc=True,version='old',args=['keep']),dict(id='other',version='keep')],packages={'other':dict(file='Other.exe',size=1,sha256='unchanged')})
  for name in names:
   key='my-tube' if name=='MyTube.exe' else 'my-tube-'+name
   cat['packages'][key]=dict(file=name if name=='MyTube.exe' else 'licenses/MyTube-'+name,size=1,sha256='old')
  return root,cat,dict(repository=info['repository'],tag=info['tag'])
 def test_replaces_only_owned_payloads_and_keeps_launch_contract(self):
  with tempfile.TemporaryDirectory() as d:
   root,cat,pin=self.fixture(d);before=copy.deepcopy(cat)
   replacements=tube_replacements(cat,d,pin)
   self.assertEqual(len(replacements),4)
   self.assertEqual(cat['apps'][0]['args'],['keep'])
   self.assertEqual(cat['apps'][0]['version'],'0.1.0-preview.6')
   self.assertEqual(cat['apps'][1],before['apps'][1])
   self.assertEqual(cat['packages']['other'],before['packages']['other'])
   self.assertEqual(cat['packages']['my-tube']['sha256'],hashlib.sha256((root/'MyTube.exe').read_bytes()).hexdigest())
 def test_rejects_dirty_or_tampered_source_and_foreign_ownership(self):
  for mutation in ['dirty','binary','extra-package','destination']:
   with self.subTest(mutation=mutation),tempfile.TemporaryDirectory() as d:
    root,cat,pin=self.fixture(d)
    if mutation=='dirty':
     p=root/'build-info.json';info=json.loads(p.read_text());info['sourceDirty']=True;p.write_text(json.dumps(info))
    elif mutation=='binary':(root/'MyTube.exe').write_text('tampered')
    elif mutation=='extra-package':cat['packages']['my-tube-foreign']=dict(file='Other.exe')
    else:cat['packages']['my-tube']['file']='Other.exe'
    with self.assertRaises(AssertionError):tube_replacements(cat,d,pin)

if __name__=='__main__':unittest.main()
