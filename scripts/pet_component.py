"""Build the Pet executable and its managed release payloads."""
import json
import os
import subprocess


def pet_component(root, go, output):
    output.mkdir(parents=True, exist_ok=True)
    web = root / 'internal/pet/web'
    subprocess.run(['npm.cmd', 'ci', '--ignore-scripts'], cwd=web, check=True, shell=True)
    subprocess.run(['node', 'build.mjs'], cwd=web, check=True)
    binary = output / 'Pet.exe'
    subprocess.run([go, 'build', '-buildvcs=false', '-trimpath', '-ldflags', '-s -w',
                    '-o', str(binary), './cmd/pet'], cwd=root, check=True,
                   env={**os.environ, 'CGO_ENABLED': '0', 'GOTOOLCHAIN': 'local'})
    version = subprocess.check_output([str(binary), '--version'], text=True).strip()
    minimum = '0.3.0-preview.124'
    contract = output / 'compatibility.json'
    contract.write_text(json.dumps({
        'schema': 1, 'id': 'pet', 'version': version, 'workspaceInstallerIntegration': 16,
        'minimumWorkspaceVersion': minimum, 'route': '/pet', 'oidc': True,
        'executable': 'Pet.exe', 'health': '/healthz', 'storage': 'private-pet-v1',
    }, indent=2), encoding='utf-8')
    return {
        'app': {
            'id': 'pet', 'name': 'Pet', 'version': version,
            'description': 'Visit a virtual pet in a living woodland world',
            'kind': 'process', 'package': 'pet', 'executable': 'Pet.exe',
            'port': 9046, 'route': '/pet', 'stripPrefix': True, 'health': '/healthz',
            'depends': ['accounts'], 'oidc': True, 'callback': '/pet/auth/callback',
            'args': ['--data', '${DATA}/pet', '--listen', '127.0.0.1:${PORT}',
                     '--origin', '${ORIGIN}/pet', '--no-browser'],
            'source': 'https://github.com/oneworkspacestudio/one-workspace/releases/tag/v0.3.0-preview.127',
        },
        'inputs': {
            'pet': {'source': str(binary), 'file': 'Pet.exe'},
            'pet-contract': {'source': str(contract), 'file': 'docs/Pet/compatibility.json'},
            'pet-guide': {'source': str(root / 'docs/PET.md'), 'file': 'docs/Pet/guide.md'},
            'pet-license': {'source': str(root / 'LICENSE'), 'file': 'licenses/Pet-LICENSE.txt'},
            'pet-notices': {'source': str(web / 'THIRD-PARTY.txt'), 'file': 'licenses/Pet-THIRD-PARTY.txt'},
        },
    }
