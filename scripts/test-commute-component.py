import unittest,tempfile,json,hashlib
from pathlib import Path
from commute_component import commute_component
class CommuteContract(unittest.TestCase):
 def test_namespace_hash_and_candidate_boundary(self):
  with tempfile.TemporaryDirectory() as d:
   r=Path(d);(r/'MyCommute.exe').write_bytes(b'synthetic');(r/'LICENSE').write_text('MIT');(r/'NOTICE').write_text('Synthetic notices')
   (r/'compatibility.json').write_text(json.dumps(dict(schema=1,id='commute',route='/commute',oidc=True,version='0.1.0-preview.1',workspaceInstallerIntegration=10,minimumWorkspaceVersion='0.3.0-preview.116',singleExecutable=True,capabilities=['native-commute-v1'])))
   (r/'manifest.json').write_text(json.dumps(dict(version='0.1.0-preview.1',candidate=True,executableSHA256=hashlib.sha256(b'synthetic').hexdigest())))
   c=commute_component(r,True);self.assertEqual(c['app']['callback'],'/commute/auth/callback');self.assertEqual(c['app']['port'],9040);self.assertTrue(all(k=='commute' or k.startswith('commute-') for k in c['inputs']))
   with self.assertRaises(AssertionError):commute_component(r)
   (r/'private-data.json').write_text('{}')
   with self.assertRaises(ValueError):commute_component(r,True)
   (r/'private-data.json').unlink();(r/'MyCommute.exe').write_bytes(b'changed')
   with self.assertRaises(AssertionError):commute_component(r,True)
if __name__=='__main__':unittest.main()
