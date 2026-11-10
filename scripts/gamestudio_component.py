import json
import subprocess


def gamestudio_component(root, go, output, suite_version):
    output.mkdir(parents=True, exist_ok=True)
    web = root / 'internal/gamestudio/web'
    subprocess.run(['npm.cmd', 'ci', '--ignore-scripts'], cwd=web, check=True, shell=True)
    subprocess.run(['node', 'build.mjs'], cwd=web, check=True)
    binary = output / 'GameStudio.exe'
    subprocess.run([go, 'build', '-buildvcs=false', '-trimpath', '-ldflags', '-s -w',
                    '-o', str(binary), './cmd/gamestudio'], cwd=root, check=True)
    version = '0.1.0-preview.1'
    contract = output / 'compatibility.json'
    contract.write_text(json.dumps({
        'schema': 1, 'id': 'gamestudio', 'version': version, 'workspaceInstallerIntegration': 15,
        'minimumWorkspaceVersion': suite_version, 'route': '/gamestudio', 'oidc': True,
        'executable': 'GameStudio.exe', 'health': '/healthz', 'storage': 'private-gamestudio-v1',
    }, indent=2), encoding='utf-8')
    return {
        'app': {
            'id': 'gamestudio', 'name': 'Game Studio', 'version': version,
            'description': 'Describe a game, make it with DeepSeek and play it in a new tab',
            'kind': 'process', 'package': 'gamestudio', 'executable': 'GameStudio.exe',
            'port': 9045, 'route': '/gamestudio', 'stripPrefix': True, 'health': '/healthz',
            'depends': ['accounts'], 'oidc': True, 'callback': '/gamestudio/auth/callback',
            'args': ['--data', '${DATA}/gamestudio', '--listen', '127.0.0.1:${PORT}',
                     '--origin', '${ORIGIN}/gamestudio', '--no-browser'],
            'source': 'https://github.com/oneworkspacestudio/one-workspace/releases/tag/v' + suite_version,
        },
        'inputs': {
            'gamestudio': {'source': str(binary), 'file': 'GameStudio.exe'},
            'gamestudio-contract': {'source': str(contract), 'file': 'docs/GameStudio/compatibility.json'},
            'gamestudio-guide': {'source': str(root / 'docs/GAME_STUDIO.md'), 'file': 'docs/GameStudio/guide.md'},
            'gamestudio-license': {'source': str(root / 'LICENSE'), 'file': 'licenses/GameStudio-LICENSE.txt'},
            'gamestudio-notices': {'source': str(root / 'docs/GAME_STUDIO_NOTICES.md'), 'file': 'licenses/GameStudio-NOTICE.md'},
        },
    }
