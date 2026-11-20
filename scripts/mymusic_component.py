import json
import os
import subprocess


def mymusic_component(root, go, output, suite_version):
    output.mkdir(parents=True, exist_ok=True)
    for name in ('index.html', 'app.js', 'style.css'):
        if (root / 'cmd/mymusic/web' / name).stat().st_size == 0:
            raise ValueError('Vuno web asset is empty: ' + name)
    binary = output / 'MyMusic.exe'
    subprocess.run([go, 'build', '-buildvcs=false', '-trimpath', '-ldflags', '-s -w',
                    '-o', str(binary), './cmd/mymusic'], cwd=root,
                   env=dict(os.environ, CGO_ENABLED='0', GOTOOLCHAIN='local'), check=True)
    version = subprocess.check_output([binary, '--version'], text=True).strip()
    contract = output / 'compatibility.json'
    contract.write_text(json.dumps({
        'schema': 1, 'id': 'mymusic', 'version': version,
        'workspaceInstallerIntegration': 23,
        'minimumWorkspaceVersion': '0.3.0-preview.142', 'route': '/mymusic',
        'oidc': True, 'executable': 'MyMusic.exe', 'health': '/healthz',
        'storage': 'private-mymusic-v1', 'engine': 'external-loopback-ace-step-1.5',
    }, indent=2), encoding='utf-8')
    return {
        'app': {
            'id': 'mymusic', 'name': 'Vuno', 'version': version,
            'description': 'Generate songs, arrange tracks and export a mix',
            'kind': 'process', 'package': 'mymusic', 'executable': 'MyMusic.exe',
            'port': 9049, 'route': '/mymusic', 'stripPrefix': True, 'health': '/healthz',
            'depends': ['accounts'], 'oidc': True, 'callback': '/mymusic/auth/callback',
            'args': ['--data', '${DATA}/mymusic', '--listen', '127.0.0.1:${PORT}',
                     '--origin', '${ORIGIN}/mymusic', '--no-browser'],
            'source': 'https://github.com/oneworkspacestudio/one-workspace/releases/tag/v' + suite_version,
        },
        'inputs': {
            'mymusic': {'source': binary, 'file': 'MyMusic.exe'},
            'mymusic-contract': {'source': contract, 'file': 'docs/MyMusic/compatibility.json'},
            'mymusic-guide': {'source': root / 'cmd/mymusic/README.md', 'file': 'docs/MyMusic/guide.md'},
            'mymusic-license': {'source': root / 'LICENSE', 'file': 'licenses/MyMusic-LICENSE.txt'},
            'mymusic-notices': {'source': root / 'cmd/mymusic/NOTICE.md', 'file': 'licenses/MyMusic-NOTICE.md'},
        },
    }
