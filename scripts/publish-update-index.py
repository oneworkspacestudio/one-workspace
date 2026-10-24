"""Refresh the public update channel after publishing any managed app/Workspace.

Run locally with GH_TOKEN/GITHUB_TOKEN or existing Git Credential Manager access.
Credentials are read into memory only. Review and commit the generated JSON.
This script does not publish releases or push commits and must not run in Actions.
"""
import datetime
import json
import os
from pathlib import Path
import re
import subprocess
import urllib.request

ROOT = Path(__file__).resolve().parents[1]
ORG = 'oneworkspacestudio'

def withheld_releases():
    data = json.loads((ROOT / 'updates/withheld-releases.json').read_text(encoding='utf-8'))
    assert data['schema'] == 1
    result = set()
    for row in data['releases']:
        assert row['repository'].startswith(ORG + '/') and row['tag'] and row['reason']
        result.add((row['repository'], row['tag']))
    return result



def main():
    token = os.environ.get('GH_TOKEN') or os.environ.get('GITHUB_TOKEN')
    if not token:
        result = subprocess.run(['git', 'credential', 'fill'],
            input='protocol=https\nhost=github.com\n\n', capture_output=True,
            text=True, check=True)
        fields = dict(line.split('=', 1) for line in result.stdout.splitlines() if '=' in line)
        token = fields.get('password')
    if not token:
        raise SystemExit('Publisher GitHub access is required; no credentials are written to the index')
    source = (ROOT / 'internal/distribution/updates.go').read_text(encoding='utf-8')
    groups = source.split('var updateGroups = []updateGroup{', 1)[1].split('\n}', 1)[0]
    repos = sorted(set(re.findall(r'^\s*\{"[^"\n]+", "[^"\n]+", "([^"\n]+)"', groups, re.MULTILINE)))
    assert 'one-workspace' in repos and len(repos) >= 20
    repositories = {}
    withheld = withheld_releases()
    for name in repos:
        repo = ORG + '/' + name
        req = urllib.request.Request('https://api.github.com/repos/' + repo + '/releases?per_page=100',
            headers={'Authorization': 'Bearer ' + token, 'User-Agent': 'OneWorkspaceChannelPublisher',
                'Accept': 'application/vnd.github+json', 'X-GitHub-Api-Version': '2022-11-28'})
        with urllib.request.urlopen(req, timeout=60) as response:
            releases = json.load(response)
        public = []
        for release in releases:
            if release['draft'] or not release['published_at'] or (repo, release['tag_name']) in withheld:
                continue
            public.append({
                'tag_name': release['tag_name'], 'draft': False, 'published_at': release['published_at'],
                'assets': [{key: a.get(key) for key in
                    ('name', 'browser_download_url', 'url', 'state', 'digest', 'size')}
                    for a in release['assets']]
            })
        repositories[repo] = public
        print(repo, len(public), 'published releases', flush=True)
    output = {'schema': 1,
        'publishedAt': datetime.datetime.now(datetime.timezone.utc).isoformat().replace('+00:00', 'Z'),
        'repositories': repositories}
    data = (json.dumps(output, separators=(',', ':'), ensure_ascii=True) + '\n').encode('utf-8')
    assert len(data) <= 8 << 20
    dest = ROOT / 'updates/releases-v1.json'
    dest.parent.mkdir(exist_ok=True)
    temporary = dest.with_suffix('.tmp')
    temporary.write_bytes(data)
    temporary.replace(dest)
    print('Wrote public update channel:', len(data), 'bytes. Review and commit it before delivery.')


if __name__ == '__main__':
    main()
