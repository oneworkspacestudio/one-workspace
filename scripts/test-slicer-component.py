"""Slicer managed-seed ownership and preservation regression checks."""
import copy,json,os,tempfile,unittest
from pathlib import Path
from slicer_component import slicer_replacements

class SlicerComponentTests(unittest.TestCase):
 def setUp(self):
  self.tmp=tempfile.TemporaryDirectory(dir=os.environ['MYSLICER_TEST_TEMP']);self.addCleanup(self.tmp.cleanup);self.root=Path(self.tmp.name)
  c=dict(schema=1,app='slicer',route='/slicer',callback='/slicer/auth/callback',health='/healthz',executable='MySlicer.exe',workspaceInstallerIntegration=2,version='0.1.0-preview.8')
  (self.root/'compatibility.json').write_text(json.dumps(c));(self.root/'build-info.json').write_text(json.dumps(dict(version=c['version'],githubActionsUsed=False,sourceCommit='1'*40)))
  self.catalog=dict(apps=[dict(id='slicer',package='slicer',executable='MySlicer.exe',route='/slicer',oidc=True,port=29123,args=['keep-existing'],version='0.1.0-preview.7'),dict(id='other',version='unchanged')],packages={'other':dict(file='Other.exe',size=1,sha256='f'*64)})
  for name in ['MySlicer.exe','LICENSE','NOTICE','THIRD-PARTY-LICENSES.txt']:
   (self.root/name).write_text('synthetic '+name)
   key='slicer' if name=='MySlicer.exe' else 'slicer-'+name
   self.catalog['packages'][key]=dict(file='MySlicer.exe' if key=='slicer' else 'licenses/MySlicer-'+name,size=1,sha256='0'*64)
  self.catalog['packages']['slicer-contract']=dict(file='docs/MySlicer/compatibility.json',size=1,sha256='0'*64)
  (self.root/'docs').mkdir()
  for name in ['SETUP.md','ACCEPTANCE.md','ENGINE.md']:
   (self.root/'docs'/name).write_text('synthetic guide')
   self.catalog['packages']['slicer-doc-'+name]=dict(file='docs/MySlicer/'+name,size=1,sha256='0'*64)
  self.pin=dict(repository='oneworkspacestudio/my-slicer',tag='v0.1.0-preview.8')
 def test_only_owned_bytes_and_version_change(self):
  before=copy.deepcopy(self.catalog);r=slicer_replacements(self.catalog,self.root,self.pin)
  self.assertEqual(len(r),8);self.assertEqual(self.catalog['apps'][1],before['apps'][1]);self.assertEqual(self.catalog['packages']['other'],before['packages']['other'])
  app=self.catalog['apps'][0];self.assertEqual(app['port'],29123);self.assertEqual(app['args'],['keep-existing']);self.assertEqual(app['version'],'0.1.0-preview.8')
 def test_missing_notice_leaves_catalog_unchanged(self):
  (self.root/'NOTICE').unlink();before=copy.deepcopy(self.catalog)
  with self.assertRaises(AssertionError):slicer_replacements(self.catalog,self.root,self.pin)
  self.assertEqual(self.catalog,before)
 def test_changed_path_or_unowned_payload_rejected(self):
  self.catalog['packages']['slicer']['file']='Other.exe'
  with self.assertRaises(AssertionError):slicer_replacements(self.catalog,self.root,self.pin)
  self.catalog['packages']['slicer']['file']='MySlicer.exe';self.catalog['packages']['slicer-extra']=dict(file='extra.exe')
  with self.assertRaises(AssertionError):slicer_replacements(self.catalog,self.root,self.pin)

if __name__=='__main__':unittest.main()
