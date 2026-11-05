import tempfile,json,hashlib,unittest
from pathlib import Path
from survival_component import survival_component
class SurvivalComponent(unittest.TestCase):
 def test_boundary(self):
  with tempfile.TemporaryDirectory() as d:
   r=Path(d);(r/'MySurvival.exe').write_bytes(b'synthetic exe');(r/'LICENSE').write_text('Apache2');(r/'NOTICE').write_text('Original summaries')
   (r/'compatibility.json').write_text(json.dumps(dict(schema=1,id='survival',route='/survival',oidc=True,version='0.1.0-preview.1',workspaceInstallerIntegration=11,minimumWorkspaceVersion='0.3.0-preview.117',singleExecutable=True,capabilities=['offline-survival-v1'],publicHandbook=True)))
   (r/'manifest.json').write_text(json.dumps(dict(version='0.1.0-preview.1',executableSHA256=hashlib.sha256(b'synthetic exe').hexdigest(),candidate=True,offlineBytes=1000)))
   c=survival_component(r,True);self.assertEqual(c['app']['callback'],'/survival/auth/callback');self.assertTrue(all(k.startswith('survival') for k in c['inputs']))
   with self.assertRaises(AssertionError):survival_component(r)
   (r/'secret.txt').write_text('unowned file')
   with self.assertRaises(ValueError):survival_component(r,True)
if __name__=='__main__':unittest.main()
