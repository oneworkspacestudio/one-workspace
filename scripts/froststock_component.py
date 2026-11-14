"""Build the managed Froststock service and its owned packages."""
import json
import os
import subprocess


def froststock_component(root, go, output, suite_version):
    output.mkdir(parents=True, exist_ok=True)
    binary = output / 'Froststock.exe'
    subprocess.run([go, 'build', '-buildvcs=false', '-trimpath', '-ldflags', '-s -w',
                    '-o', str(binary), './cmd/froststock'], cwd=root, check=True,
                   env={**os.environ, 'CGO_ENABLED': '0', 'GOTOOLCHAIN': 'local'})
    version = subprocess.check_output([str(binary), '--version'], text=True).strip()
    contract = output / 'compatibility.json'
    contract.write_text(json.dumps({
        'schema': 1, 'id': 'froststock', 'version': version, 'workspaceInstallerIntegration': 18,
        'minimumWorkspaceVersion': '0.3.0-preview.128', 'route': '/froststock', 'oidc': True,
        'executable': 'Froststock.exe', 'health': '/healthz', 'storage': 'my-drive-froststock-v1',
    }, indent=2), encoding='utf-8')
    return {
        'app': {
            'id': 'froststock', 'name': 'Froststock', 'version': version,
            'description': 'Scan freezer items and expiry labels into My Drive',
            'kind': 'process', 'package': 'froststock', 'executable': 'Froststock.exe',
            'port': 9048, 'route': '/froststock', 'stripPrefix': True, 'health': '/healthz',
            'depends': ['accounts', 'drive'], 'oidc': True, 'callback': '/froststock/auth/callback',
            'args': ['--data', '${DATA}/froststock', '--listen', '127.0.0.1:${PORT}',
                     '--origin', '${ORIGIN}/froststock', '--no-browser'],
            'source': 'https://github.com/oneworkspacestudio/one-workspace/releases/tag/v' + suite_version,
        },
        'inputs': {
            'froststock': {'source': str(binary), 'file': 'Froststock.exe'},
            'froststock-contract': {'source': str(contract), 'file': 'docs/Froststock/compatibility.json'},
            'froststock-guide': {'source': str(root / 'docs/FROSTSTOCK.md'), 'file': 'docs/Froststock/guide.md'},
            'froststock-license': {'source': str(root / 'LICENSE'), 'file': 'licenses/Froststock-LICENSE.txt'},
            'froststock-notices': {'source': str(root / 'internal/froststock/NOTICE'), 'file': 'licenses/Froststock-NOTICE.txt'},
        },
    }
