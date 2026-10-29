"""Local One Workspace publisher. Called by publish-release.ps1; never runs in Actions."""
import argparse
import hashlib
import json
from pathlib import Path
import re
import shutil
import subprocess
import sys
import tempfile
import urllib.request
import zipfile

ROOT = Path(__file__).resolve().parents[1]
REPO = 'oneworkspacestudio/one-workspace'
INDEX_URL = 'https://raw.githubusercontent.com/' + REPO + '/main/updates/releases-v1.json'


def run(args, **kwargs):
    return subprocess.check_output([str(a) for a in args], text=True, **kwargs).strip()


def git(*args):
    return run(['git', '-c', 'safe.directory=' + ROOT.as_posix(), '-C', ROOT] + list(args))


def api(endpoint, payload=None):
    command = ['gh', 'api', 'repos/' + REPO + '/' + endpoint]
    if payload is None:
        return json.loads(run(command))
    # A structured JSON file preserves release-note newlines. Never use shell interpolation.
    with tempfile.TemporaryDirectory(prefix='workspace-publish-') as folder:
        request = Path(folder) / 'request.json'
        request.write_text(json.dumps(payload), encoding='utf-8')
        method = 'PATCH' if re.match(r'^releases/[0-9]+$', endpoint) else 'POST'
        return json.loads(run(command + ['--method', method, '--input', request]))


def sha(path):
    digest = hashlib.sha256()
    with path.open('rb') as source:
        for chunk in iter(lambda: source.read(1048576), b''):
            digest.update(chunk)
    return 'sha256:' + digest.hexdigest()


def verify_assets(release, files):
    remote = {a['name']: a for a in release['assets']}
    if set(remote) != {p.name for p in files}:
        raise RuntimeError('Draft assets differ from the prepared public asset directory')
    for local in files:
        asset = remote[local.name]
        if asset['state'] != 'uploaded' or asset['size'] != local.stat().st_size or asset.get('digest') != sha(local):
            raise RuntimeError('GitHub size/digest verification failed: ' + local.name)


def public_apk(url, expected):
    request = urllib.request.Request(url, headers={'User-Agent': 'OneWorkspacePublicReleaseCheck'})
    digest = hashlib.sha256()
    size = 0
    with urllib.request.urlopen(request, timeout=60) as response:
        for chunk in iter(lambda: response.read(1048576), b''):
            digest.update(chunk)
            size += len(chunk)
    if size != expected.stat().st_size or 'sha256:' + digest.hexdigest() != sha(expected):
        raise RuntimeError('Anonymous APK download did not match the signed local APK')


def prepare_package_assets(out):
    catalog = json.loads((out / 'catalog.json').read_text(encoding='utf-8'))
    with zipfile.ZipFile(out / 'workspace-bundle.zip') as bundle:
        for key, package in catalog['packages'].items():
            if not re.fullmatch(r'[A-Za-z0-9_.-]{1,180}', key):
                raise RuntimeError('Unsafe package asset identifier')
            dest = out / ('workspace-package-' + key + '.bin')
            with bundle.open(package['file']) as src, dest.open('wb') as target:
                shutil.copyfileobj(src, target)
            if dest.stat().st_size != package['size'] or sha(dest) != 'sha256:' + package['sha256']:
                raise RuntimeError('Independent package checksum mismatch: ' + key)

def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--tag', required=True)
    parser.add_argument('--assets', type=Path, required=True)
    parser.add_argument('--notes', type=Path, required=True)
    parser.add_argument('--title', default='')
    parser.add_argument('--sdk', type=Path, required=True)
    args = parser.parse_args()
    if not re.fullmatch(r'(android-)?v[0-9]+\.[0-9]+\.[0-9]+(?:-[A-Za-z0-9.-]+)?', args.tag):
        raise RuntimeError('Use an Android or Workspace version tag')
    if git('status', '--porcelain'):
        raise RuntimeError('Release source must be committed and clean')
    source = git('rev-parse', 'HEAD')
    tree = git('rev-parse', 'HEAD:android')
    tag_ref = 'refs/tags/' + args.tag
    refs = dict((line.split()[1], line.split()[0]) for line in git('ls-remote', 'origin', tag_ref, tag_ref + '^{}').splitlines())
    if refs and refs.get(tag_ref + '^{}', refs.get(tag_ref)) != source:
        raise RuntimeError('Existing tag points to different source; use a new release tag')
    policy = json.loads((ROOT / 'android/release-policy.json').read_text(encoding='utf-8'))
    assert policy['schema'] == 1 and policy['requireAPKForEveryRelease'] is True
    version_file = (ROOT / 'android/app/build.gradle').read_text(encoding='utf-8')
    version = re.search(r"versionName\s+'([^']+)'", version_file).group(1)
    code = int(re.search(r'versionCode\s+(\d+)', version_file).group(1))
    if args.tag.startswith('android-') and args.tag != 'android-v' + version:
        raise RuntimeError('Android release tag must match versionName')
    out = args.assets.resolve()
    apk = out / ('OneWorkspace-Android-' + version + '.apk')
    if not apk.is_file():
        raise RuntimeError('Required signed companion APK is missing')
    if any(p.suffix == '.apk' and p != apk for p in out.iterdir()):
        raise RuntimeError('Remove stale/different APKs from this public staging directory')
    tools = sorted((p for p in (args.sdk / 'build-tools').iterdir() if re.fullmatch(r'\d+\.\d+\.\d+', p.name)), key=lambda p: tuple(map(int, p.name.split('.'))))[-1]
    signature = run([tools / ('apksigner.bat' if sys.platform == 'win32' else 'apksigner'), 'verify', '--verbose', '--print-certs', apk])
    certs = re.findall(r'Signer #\d+ certificate SHA-256 digest:\s*([a-f0-9]{64})', signature)
    if certs != [policy['certificateSHA256']]:
        raise RuntimeError('APK is not signed with the approved permanent Android key')
    badging = run([tools / ('aapt2.exe' if sys.platform == 'win32' else 'aapt2'), 'dump', 'badging', apk])
    if not re.search(r"package: name='" + re.escape(policy['applicationId']) + r"' versionCode='" + str(code) + r"' versionName='" + re.escape(version) + "'", badging):
        raise RuntimeError('APK package/version does not match this source')
    with zipfile.ZipFile(apk) as archive:
        provenance = json.loads(archive.read('assets/workspace-release.json'))
        if provenance != {'schema': 1, 'sourceCommit': source, 'androidTree': tree}:
            raise RuntimeError('APK provenance does not match the committed release source')
        for name in ['assets/licenses/APACHE-2.0.txt', 'assets/licenses/THIRD-PARTY-NOTICES.txt']:
            archive.getinfo(name)
    if not args.tag.startswith('android-'):
        for name in ['WorkspaceSetup-windows-amd64.exe', 'Workspace.exe', 'catalog.json', 'workspace-bundle.zip']:
            if not (out / name).is_file():
                raise RuntimeError('Workspace release is missing ' + name)
        if json.loads((out / 'catalog.json').read_text(encoding='utf-8'))['version'] != args.tag[1:]:
            raise RuntimeError('Workspace catalog version does not match release tag')
    if api('actions/permissions')['enabled'] is not False:
        raise RuntimeError('Keep GitHub Actions disabled before publishing locally')
    if api('commits/' + source)['sha'] != source:
        raise RuntimeError('Push the intended source before publishing')
    notes = args.notes.read_text(encoding='utf-8').strip()
    if not notes:
        raise RuntimeError('Release notes are required')
    releases = api('releases?per_page=100')
    existing = [r for r in releases if r['tag_name'] == args.tag]
    if existing and (not existing[0]['draft'] or existing[0]['target_commitish'] != source):
        raise RuntimeError('Published releases are immutable; choose a new tag. A draft must match this source.')
    run(['git', '-c', 'safe.directory=' + ROOT.as_posix(), '-C', ROOT, 'archive', '--format=zip', '--prefix=one-workspace-android/', '--output=' + str(out / 'OneWorkspace-Android-source.zip'), source, 'android', 'scripts/build-android.ps1', 'LICENSE', 'NOTICE'])
    for src, name in [('android/README.md', 'ANDROID.md'), ('android/app/src/main/assets/licenses/THIRD-PARTY-NOTICES.txt', 'ANDROID-NOTICES.txt'), ('LICENSE', 'ANDROID-LICENSE.txt')]:
        shutil.copyfile(ROOT / src, out / name)
    (out / 'ANDROID-SIGNATURE.txt').write_text(signature + '\n', encoding='utf-8')
    manifest = {'schema': 1, 'applicationId': policy['applicationId'], 'versionName': version, 'versionCode': code, 'sourceCommit': source, 'androidTree': tree, 'apk': apk.name, 'apkSHA256': sha(apk)[7:], 'certificateSHA256': certs[0], 'builtLocally': True, 'githubActionsUsed': False, 'checks': ['signature', 'package/version', 'embedded-source-provenance', 'embedded-licenses']}
    (out / 'ANDROID-BUILD.json').write_text(json.dumps(manifest, indent=2) + '\n', encoding='utf-8')
    if not args.tag.startswith('android-'):
        prepare_package_assets(out)
    files = sorted(p for p in out.iterdir() if p.is_file() and p.name not in ['SHA256SUMS', 'SHA256SUMS.txt'])
    for path in files:
        if path.suffix.lower() in ['.p12', '.jks', '.keystore', '.dpapi', '.db', '.log'] or path.name.startswith('.'):
            raise RuntimeError('Keep private material out of the public staging directory: ' + path.name)
    sums = out / 'SHA256SUMS.txt'
    sums.write_text(''.join(sha(p)[7:] + '  ' + p.name + '\n' for p in files), encoding='utf-8')
    files.append(sums)
    version_tag = args.tag[len('android-'):] if args.tag.startswith('android-') else args.tag
    release = existing[0] if existing else api('releases', {'tag_name': args.tag, 'target_commitish': source, 'name': args.title or args.tag, 'body': notes, 'draft': True, 'prerelease': '-' in version_tag, 'make_latest': 'false'})
    remote = {a['name']: a for a in release['assets']}
    missing = []
    for path in files:
        asset = remote.get(path.name)
        if asset:
            if asset.get('digest') != sha(path) or asset['size'] != path.stat().st_size:
                raise RuntimeError('Existing draft asset differs; review it before retrying: ' + path.name)
        else:
            missing.append(path)
    if missing:
        subprocess.run(['gh', 'release', 'upload', args.tag, '--repo', REPO] + [str(p) for p in missing], check=True)
    verify_assets(api('releases/' + str(release['id'])), files)
    release = api('releases/' + str(release['id']), {'draft': False})
    if api('commits/' + args.tag)['sha'] != source:
        raise RuntimeError('Published tag does not resolve to the intended source')
    public_apk('https://github.com/' + REPO + '/releases/download/' + args.tag + '/' + apk.name, apk)
    # Promotion is part of the supported publication flow. Do not move the source tag.
    subprocess.run([sys.executable, str(ROOT / 'scripts/publish-update-index.py')], cwd=str(ROOT), check=True)
    git('add', 'updates/releases-v1.json')
    changed = subprocess.run(['git', '-c', 'safe.directory=' + ROOT.as_posix(), '-C', str(ROOT), 'diff', '--cached', '--quiet', '--', 'updates/releases-v1.json']).returncode
    if changed == 1:
        git('commit', '-m', 'Promote ' + args.tag + ' with required signed Android APK [skip ci]', '--', 'updates/releases-v1.json')
        print(git('push', 'origin', 'HEAD:main'), flush=True)
    elif changed != 0:
        raise RuntimeError('Cannot inspect channel promotion')
    # Exact URL used by installers; fail clearly if CDN propagation still needs a retry.
    with urllib.request.urlopen(urllib.request.Request(INDEX_URL, headers={'User-Agent': 'OneWorkspacePublicReleaseCheck'}), timeout=60) as response:
        promoted = json.load(response)
    rows = promoted['repositories'][REPO]
    if not any(r['tag_name'] == args.tag and any(a['name'] == apk.name and a['digest'] == sha(apk) for a in r['assets']) for r in rows):
        raise RuntimeError('Release is published and index pushed; public CDN is still propagating. Verify the exact public index URL before announcing delivery.')
    print(json.dumps({'url': release['html_url'], 'apk': apk.name, 'assetsVerified': len(files), 'anonymousAPKVerified': True, 'publicChannelVerified': True}))


if __name__ == '__main__':
    main()
