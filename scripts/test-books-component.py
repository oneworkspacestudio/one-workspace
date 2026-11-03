import unittest,tempfile,json,hashlib
from pathlib import Path
from books_component import books_component
class BooksContract(unittest.TestCase):
 def test_owned_files_and_candidate_boundary(self):
  with tempfile.TemporaryDirectory() as d:
   r=Path(d);(r/'MyBooks.exe').write_bytes(b'synthetic program');(r/'LICENSE').write_text('Apache2');(r/'NOTICE').write_text('Notices')
   (r/'compatibility.json').write_text(json.dumps(dict(schema=1,id='books',route='/books',oidc=True,version='0.1.0-preview.1',workspaceInstallerIntegration=8,minimumWorkspaceVersion='0.3.0-preview.113',singleExecutable=True,capabilities=['native-books-v1'])))
   (r/'manifest.json').write_text(json.dumps(dict(version='0.1.0-preview.1',candidate=True,executableSHA256=hashlib.sha256(b'synthetic program').hexdigest())))
   c=books_component(r,True);self.assertEqual(c['app']['env']['MYBOOKS_LIBRARY_ROOT'],'${FILES}/My Books');self.assertEqual(c['app']['callback'],'/books/auth/callback');self.assertTrue(all(k.startswith('books') for k in c['inputs']))
   with self.assertRaises(AssertionError):books_component(r)
   (r/'unexpected-secret.txt').write_text('fixture')
   with self.assertRaises(ValueError):books_component(r,True)
if __name__=='__main__':unittest.main()
