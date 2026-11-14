"""Package the native Sunny & Moon service for managed Workspace installs."""
import json
import os
import subprocess


def sunnymoon_component(root, go, output, suite_version):
    output.mkdir(parents=True, exist_ok=True)
    binary = output / 'SunnyMoon.exe'
    subprocess.run([go, 'build', '-buildvcs=false', '-trimpath', '-ldflags', '-s -w',
                    '-o', str(binary), './cmd/sunnymoon'], cwd=root, check=True,
                   env={**os.environ, 'CGO_ENABLED': '0', 'GOTOOLCHAIN': 'local'})
    version = subprocess.check_output([str(binary), '--version'], text=True).strip()
    contract = output / 'compatibility.json'
    contract.write_text(json.dumps({
        'schema': 1, 'id': 'sunnymoon', 'version': version,
        'workspaceInstallerIntegration': 17, 'minimumWorkspaceVersion': '0.3.0-preview.125',
        'route': '/sunnymoon', 'oidc': True, 'executable': 'SunnyMoon.exe',
        'health': '/healthz', 'storage': 'private-sunnymoon-v1',
    }, indent=2), encoding='utf-8')
    return {
        'app': {
            'id': 'sunnymoon', 'name': 'Sunny & Moon', 'version': version,
            'description': 'A gentle daily routine for toilet trips and bedtime',
            'kind': 'process', 'package': 'sunnymoon', 'executable': 'SunnyMoon.exe',
            'port': 9047, 'route': '/sunnymoon', 'stripPrefix': True, 'health': '/healthz',
            'depends': ['accounts'], 'oidc': True, 'callback': '/sunnymoon/auth/callback',
            'args': ['--data', '${DATA}/sunnymoon', '--listen', '127.0.0.1:${PORT}',
                     '--origin', '${ORIGIN}/sunnymoon', '--no-browser'],
            'source': 'https://github.com/oneworkspacestudio/one-workspace/releases/tag/v' + suite_version,
        },
        'inputs': {
            'sunnymoon': {'source': str(binary), 'file': 'SunnyMoon.exe'},
            'sunnymoon-contract': {'source': str(contract), 'file': 'docs/SunnyMoon/compatibility.json'},
            'sunnymoon-guide': {'source': str(root / 'docs/SUNNY_MOON.md'), 'file': 'docs/SunnyMoon/guide.md'},
            'sunnymoon-license': {'source': str(root / 'LICENSE'), 'file': 'licenses/SunnyMoon-LICENSE.txt'},
            'sunnymoon-font-ofl': {'source': str(root / 'internal/sunnymoon/DM-Sans-OFL.txt'),
                                   'file': 'licenses/SunnyMoon-DM-Sans-OFL.txt'},
            'sunnymoon-font-source': {'source': str(root / 'internal/sunnymoon/DM-Sans-source.ttf'),
                                      'file': 'licenses/SunnyMoon-DM-Sans-source.ttf'},
        },
    }
