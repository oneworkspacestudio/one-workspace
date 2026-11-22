"""Build the embedded AI Fieldnotes service for the Workspace catalogue."""
import json
import os
import subprocess


def fieldnotes_component(root, go, output, suite_version):
    output.mkdir(parents=True, exist_ok=True)
    binary = output / 'AIFieldnotes.exe'
    subprocess.run([go, 'build', '-buildvcs=false', '-trimpath', '-ldflags', '-s -w',
                    '-o', str(binary), './cmd/fieldnotes'], cwd=root, check=True,
                   env={**os.environ, 'CGO_ENABLED': '0', 'GOTOOLCHAIN': 'local'})
    version = subprocess.check_output([str(binary), '--version'], text=True).strip()
    contract = output / 'compatibility.json'
    contract.write_text(json.dumps({
        'schema': 1, 'id': 'fieldnotes', 'version': version,
        'workspaceInstallerIntegration': 24, 'minimumWorkspaceVersion': '0.3.0-preview.145',
        'route': '/fieldnotes', 'oidc': True, 'executable': binary.name,
        'health': '/healthz', 'storage': 'private-fieldnotes-v1',
    }, indent=2), encoding='utf-8')
    return {
        'app': {
            'id': 'fieldnotes', 'name': 'AI Fieldnotes', 'version': version,
            'description': 'Research AI engineering and build a reading list for offline Android',
            'kind': 'process', 'package': 'fieldnotes', 'executable': binary.name,
            'port': 9050, 'route': '/fieldnotes', 'stripPrefix': True, 'health': '/healthz',
            'depends': ['accounts'], 'oidc': True, 'callback': '/fieldnotes/auth/callback',
            'args': ['--data', '${DATA}/fieldnotes', '--listen', '127.0.0.1:${PORT}',
                     '--origin', '${ORIGIN}/fieldnotes', '--no-browser'],
            'source': 'https://github.com/oneworkspacestudio/one-workspace/releases/tag/v' + suite_version,
        },
        'inputs': {
            'fieldnotes': {'source': str(binary), 'file': binary.name},
            'fieldnotes-contract': {'source': str(contract), 'file': 'docs/AIFieldnotes/compatibility.json'},
            'fieldnotes-guide': {'source': str(root / 'docs/AI_FIELDNOTES.md'), 'file': 'docs/AIFieldnotes/guide.md'},
            'fieldnotes-license': {'source': str(root / 'LICENSE'), 'file': 'licenses/AIFieldnotes-LICENSE.txt'},
            'fieldnotes-ui-notices': {'source': str(root / 'internal/fieldnotes/ui/dist/THIRD-PARTY-NOTICES.txt'),
                                     'file': 'licenses/AIFieldnotes-UI-NOTICES.txt'},
        },
    }
