import importlib.util
import datetime
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

    def test_workspace_pages_keep_asset_heavy_responses_small(self):
        calls = []

        def opener(req, timeout):
            calls.append(req.full_url)
            return io.StringIO(json.dumps([{}] * (20 if len(calls) == 1 else 1)))

        releases = publisher.repository_releases('oneworkspacestudio/one-workspace', 'synthetic', opener)
        self.assertEqual(len(releases), 21)
        self.assertTrue(all('per_page=20&' in url for url in calls))

    def test_channel_caps_newest_public_releases_after_filtering(self):
        repo = 'example/repository'
        start = datetime.date(2025, 1, 1)
        releases = [dict(tag_name=f'v{i}', draft=i == 101, assets=[],
                         published_at=(start + datetime.timedelta(days=i)).isoformat() + 'T00:00:00Z')
                    for i in range(103)]
        public = publisher.channel_releases(repo, releases, {(repo, 'v100')})
        self.assertEqual(len(public), 100)
        self.assertEqual([r['tag_name'] for r in public[:3]], ['v102', 'v99', 'v98'])
        self.assertEqual(public[-1]['tag_name'], 'v1')

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
