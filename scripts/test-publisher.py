import importlib.util
from pathlib import Path
import subprocess
import unittest

spec = importlib.util.spec_from_file_location('publisher', Path(__file__).with_name('publish-github-release.py'))
publisher = importlib.util.module_from_spec(spec)
spec.loader.exec_module(publisher)

class PublisherWindowsArguments(unittest.TestCase):
    def test_large_catalog_uploads_every_asset_once_with_bounded_utf16_commands(self):
        paths = ['C:/Synthetic staging folder \U0001f4e6/' + 'x'*100 + '/workspace-package-' + str(n) + '.bin' for n in range(600)]
        commands = list(publisher.upload_commands('v0.3.0-preview.104', paths))
        self.assertGreater(len(commands), 1)
        self.assertEqual([p for command in commands for p in command[6:]], paths)
        for command in commands:
            self.assertLessEqual(len(subprocess.list2cmdline(command).encode('utf-16-le')) // 2, 24000)
    def test_empty_input_never_starts_an_upload(self):
        self.assertEqual(list(publisher.upload_commands('v1', [])), [])
    def test_single_overlong_path_is_refused_before_upload(self):
        with self.assertRaises(RuntimeError):
            list(publisher.upload_commands('v1', ['x'*30000]))

if __name__ == '__main__':
    unittest.main()
