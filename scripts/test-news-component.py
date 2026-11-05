import unittest,tempfile,json,hashlib
from pathlib import Path
from news_component import news_component
class NewsContract(unittest.TestCase):
 def test_owned_files_and_candidate_boundary(self):
  with tempfile.TemporaryDirectory() as d:
   r=Path(d);(r/'MyNews.exe').write_bytes(b'synthetic program');(r/'LICENSE').write_text('Apache2');(r/'NOTICE').write_text('Notices')
   (r/'compatibility.json').write_text(json.dumps(dict(schema=1,id='news',route='/news',oidc=True,version='0.1.0-preview.1',workspaceInstallerIntegration=9,minimumWorkspaceVersion='0.3.0-preview.115',singleExecutable=True,capabilities=['native-news-v1'])))
   (r/'manifest.json').write_text(json.dumps(dict(version='0.1.0-preview.1',candidate=True,executableSHA256=hashlib.sha256(b'synthetic program').hexdigest())))
   c=news_component(r,True);self.assertEqual(c['app']['callback'],'/news/auth/callback');self.assertTrue(all(k.startswith('news') for k in c['inputs']))
   with self.assertRaises(AssertionError):news_component(r)
   (r/'unexpected-secret.txt').write_text('fixture')
   with self.assertRaises(ValueError):news_component(r,True)
if __name__=='__main__':unittest.main()
