import copy,hashlib,json,tempfile,unittest
from pathlib import Path
from myai_component import myai_replacements

class MyAIContract(unittest.TestCase):
 def fixture(self,directory):
  root=Path(directory)/'MyAI';(root/'docs').mkdir(parents=True)
  for name in ['MyAI.exe','LICENSE','NOTICE','THIRD-PARTY-LICENSES.txt','docs/SETUP.md','docs/SUPPORT.md','docs/ADR-001-architecture.md','docs/SECURITY-CHECKS.md','docs/MODEL-ADVISOR.md']:(root/name).write_text('synthetic '+name)
  c=dict(schema=1,id='myai',route='/myai',executable='MyAI.exe',health='/healthz',oidc=True,embeddedFrontend=True,version='0.1.0-preview.22',workspaceInstallerIntegration=3,minimumWorkspaceVersion='0.3.0-preview.93')
  (root/'compatibility.json').write_text(json.dumps(c));(root/'build-info.json').write_text(json.dumps(dict(schema=1,version=c['version'],sourceDirty=False,githubActionsUsed=False,commit='a'*40,exeSHA256=hashlib.sha256((root/'MyAI.exe').read_bytes()).hexdigest())))
  cat=dict(apps=[dict(id='myai',package='myai',route='/myai',executable='MyAI.exe',oidc=True,version='old',args=['keep']),dict(id='other',version='keep')],packages={'myai':{'file':'MyAI.exe','size':1,'sha256':'old'},'other':{'file':'Other.exe','size':4,'sha256':'unchanged'}})
  pin=dict(repository='oneworkspacestudio/my-ai',tag='v'+c['version'])
  return root,cat,pin
 def test_seed_uses_public_app_pin_and_retains_every_unrelated_contract(self):
  with tempfile.TemporaryDirectory() as d:
   root,cat,pin=self.fixture(d);before=copy.deepcopy(cat);paths=myai_replacements(cat,d,pin)
   self.assertEqual(cat['apps'][1],before['apps'][1]);self.assertEqual(cat['packages']['other'],before['packages']['other']);self.assertEqual(cat['apps'][0]['args'],['keep'])
   self.assertEqual(cat['apps'][0]['version'],'0.1.0-preview.22');self.assertEqual(len(paths),9)
   self.assertIn('docs/MyAI/SUPPORT.md',paths);self.assertEqual(cat['packages']['myai']['sha256'],hashlib.sha256((root/'MyAI.exe').read_bytes()).hexdigest())
 def test_dirty_or_tampered_program_and_foreign_ownership_fail_closed(self):
  for change in ['dirty','program','ownership']:
   with self.subTest(change=change),tempfile.TemporaryDirectory() as d:
    root,cat,pin=self.fixture(d)
    if change=='dirty':
     p=root/'build-info.json';b=json.loads(p.read_text());b['sourceDirty']=True;p.write_text(json.dumps(b))
    elif change=='program':(root/'MyAI.exe').write_text('changed')
    else:cat['packages']['myai']['file']='Other.exe'
    with self.assertRaises(AssertionError):myai_replacements(cat,d,pin)
if __name__=='__main__':unittest.main()
