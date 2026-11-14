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
# Released v1 clients reject repositories with more than this many entries.
MAX_PUBLIC_RELEASES = 100


def repository_releases(repo, token, opener=urllib.request.urlopen):
    releases = []
    page_size = 20 if repo == ORG + '/one-workspace' else 100
    for page in range(1, 101):
        req = urllib.request.Request(
            f'https://api.github.com/repos/{repo}/releases?per_page={page_size}&page={page}',
            headers={'Authorization': 'Bearer ' + token, 'User-Agent': 'OneWorkspaceChannelPublisher',
                     'Accept': 'application/vnd.github+json', 'X-GitHub-Api-Version': '2022-11-28'})
        with opener(req, timeout=60) as response:
            batch = json.load(response)
        if not isinstance(batch, list):
            raise ValueError('GitHub release list was not an array')
        releases.extend(batch)
        if len(batch) < page_size:
            return releases
    raise ValueError('Release pagination exceeded its safety bound; channel was not written')

def withheld_releases():
    data = json.loads((ROOT / 'updates/withheld-releases.json').read_text(encoding='utf-8'))
    assert data['schema'] == 1
    result = set()
    for row in data['releases']:
        assert row['repository'].startswith(ORG + '/') and row['tag'] and row['reason']
        result.add((row['repository'], row['tag']))
    return result



def channel_releases(repo, releases, withheld):
    public = [release for release in releases
              if not release['draft'] and release['published_at']
              and (repo, release['tag_name']) not in withheld]
    public.sort(key=lambda release: release['published_at'], reverse=True)
    return [{
        'tag_name': release['tag_name'], 'draft': False, 'published_at': release['published_at'],
        'assets': [{key: a.get(key) for key in
            ('name', 'browser_download_url', 'url', 'state', 'digest', 'size')}
            for a in release['assets']]
    } for release in public[:MAX_PUBLIC_RELEASES]]


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
        # A private app may be distributed in a public Workspace bundle without
        # publishing its private release metadata or unusable download links.
        metadata_req = urllib.request.Request('https://api.github.com/repos/' + repo,
            headers={'Authorization': 'Bearer ' + token, 'User-Agent': 'OneWorkspaceChannelPublisher',
                'Accept': 'application/vnd.github+json'})
        with urllib.request.urlopen(metadata_req, timeout=60) as response:
            metadata = json.load(response)
        if metadata.get('private', True):
            repositories[repo] = []
            print(repo, 'private source; distributed through Workspace bundles', flush=True)
            continue
        releases = repository_releases(repo, token)
        public = channel_releases(repo, releases, withheld)
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
