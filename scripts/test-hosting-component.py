import copy,json,tempfile,unittest
from pathlib import Path
from managed_component import hosting_replacements

class HostingManagedBuild(unittest.TestCase):
 def test_all_existing_payloads_are_replaced_without_changing_other_apps(self):
  with tempfile.TemporaryDirectory() as folder:
   root=Path(folder)
   compatibility={'application':'my-hosting','version':'0.1.0-preview.2','minimumWorkspaceVersion':'0.3.0-preview.110','workspaceInstallerIntegration':6,'capabilities':['native-hosting-git-v1']}
   (root/'compatibility.json').write_text(json.dumps(compatibility))
   for name in ['MyHosting.exe','LICENSE','NOTICE','THIRD-PARTY-LICENSES.txt','docs/OPERATIONS.md']:
    p=root/name;p.parent.mkdir(exist_ok=True);p.write_text('synthetic release file')
   hosting={'id':'hosting','package':'hosting','executable':'MyHosting.exe','route':'/hosting','oidc':True,'version':'0.1.0-preview.1','port':9007,'args':['--data','${DATA}/my-hosting']}
   other={'id':'accounts','version':'unchanged'}
   cat={'apps':[hosting,other],'packages':{'hosting':{'file':'MyHosting.exe'},'hosting-LICENSE':{'file':'licenses/MyHosting-LICENSE'},'hosting-NOTICE':{'file':'licenses/MyHosting-NOTICE'},'hosting-THIRD-PARTY-LICENSES.txt':{'file':'licenses/MyHosting-THIRD-PARTY-LICENSES.txt'},'hosting-compatibility.json':{'file':'docs/MyHosting/compatibility.json'},'hosting-doc-OPERATIONS.md':{'file':'docs/MyHosting/OPERATIONS.md'},'accounts':{'file':'MyAccounts.exe'}}}
   before=copy.deepcopy(cat)
   replacements=hosting_replacements(cat,root,{'repository':'oneworkspacestudio/my-hosting','tag':'v0.1.0-preview.2'})
   self.assertEqual(len(replacements),6)
   self.assertEqual(cat['packages'],before['packages'])
   self.assertEqual(cat['apps'][1],before['apps'][1])
   self.assertEqual(cat['apps'][0]['args'],before['apps'][0]['args'])
   self.assertEqual(cat['apps'][0]['port'],9007)
   self.assertEqual(cat['apps'][0]['version'],'0.1.0-preview.2')
   (root/'docs/OPERATIONS.md').unlink()
   with self.assertRaises(AssertionError):hosting_replacements(cat,root,{'repository':'oneworkspacestudio/my-hosting','tag':'v0.1.0-preview.2'})

if __name__=='__main__':unittest.main()
