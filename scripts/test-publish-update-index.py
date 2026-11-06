import importlib.util
import io
import json
from pathlib import Path
import unittest

spec = importlib.util.spec_from_file_location('channel_publisher', Path(__file__).with_name('publish-update-index.py'))
publisher = importlib.util.module_from_spec(spec)
spec.loader.exec_module(publisher)


class ReleasePaginationTests(unittest.TestCase):
    def test_later_page_is_retained_even_when_drafts_fill_first_page(self):
        pages = [[{'tag_name': str(i), 'draft': i < 8} for i in range(100)],
                 [{'tag_name': 'historical-public-release', 'draft': False}]]
        calls = []

        def opener(req, timeout):
            calls.append(req.full_url)
            self.assertEqual(timeout, 60)
            return io.StringIO(json.dumps(pages[len(calls) - 1]))

        releases = publisher.repository_releases('example/repository', 'synthetic', opener)
        self.assertEqual(len(releases), 101)
        self.assertEqual(releases[-1]['tag_name'], 'historical-public-release')
        self.assertTrue(calls[1].endswith('&page=2'))

    def test_failed_later_page_is_not_a_successful_partial_index(self):
        calls = 0

        def opener(req, timeout):
            nonlocal calls
            calls += 1
            if calls == 2:
                raise OSError('synthetic network failure')
            return io.StringIO(json.dumps([{}] * 100))

        with self.assertRaises(OSError):
            publisher.repository_releases('example/repository', 'synthetic', opener)


if __name__ == '__main__':
    unittest.main()
