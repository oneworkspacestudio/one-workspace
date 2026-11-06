import unittest,tempfile,json,hashlib
from pathlib import Path
from podcasts_component import podcasts_component
class PodcastsContract(unittest.TestCase):
 def test_namespace_hash_and_candidate_boundary(self):
  with tempfile.TemporaryDirectory() as d:
   r=Path(d);(r/'MyPodcasts.exe').write_bytes(b'synthetic');(r/'LICENSE').write_text('MIT');(r/'NOTICE').write_text('Synthetic notices')
   (r/'compatibility.json').write_text(json.dumps(dict(schema=1,id='podcasts',route='/podcasts',oidc=True,version='0.1.0-preview.1',workspaceInstallerIntegration=14,minimumWorkspaceVersion='0.3.0-preview.122',singleExecutable=True,capabilities=['native-podcasts-v1'])))
   (r/'manifest.json').write_text(json.dumps(dict(version='0.1.0-preview.1',candidate=True,executableSHA256=hashlib.sha256(b'synthetic').hexdigest())))
   c=podcasts_component(r,True);self.assertEqual(c['app']['callback'],'/podcasts/auth/callback');self.assertEqual(c['app']['port'],9044);self.assertTrue(all(k=='podcasts' or k.startswith('podcasts-') for k in c['inputs']))
   with self.assertRaises(AssertionError):podcasts_component(r)
   (r/'private-data.json').write_text('{}')
   with self.assertRaises(ValueError):podcasts_component(r,True)
   (r/'private-data.json').unlink();(r/'MyPodcasts.exe').write_bytes(b'changed')
   with self.assertRaises(AssertionError):podcasts_component(r,True)
if __name__=='__main__':unittest.main()
