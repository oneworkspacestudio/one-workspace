"""Consume the verified public Finance pin while preserving its managed contract.

The portable archive has no embedded build-info/update manifest. The caller must
verify its public release digest and size before extraction (asset() does this).
A different runtime or installer contract needs explicit lifecycle acceptance.
"""
import hashlib
import json
from pathlib import Path

RUNTIME_SHA256 = 'f3c4675bf401042356310325a35de736de3c7b2f53aab2e74b611839f98da3ca'


def _sha(path):
    digest = hashlib.sha256()
    with path.open('rb') as stream:
        for chunk in iter(lambda: stream.read(4194304), b''):
            digest.update(chunk)
    return digest.hexdigest()


def finance_replacements(catalog, stage, pin):
    root = Path(stage).resolve()
    contract_file = root / 'compatibility.json'
    assert contract_file.is_file() and not contract_file.is_symlink(), 'Missing Finance compatibility contract'
    contract = json.loads(contract_file.read_text(encoding='utf-8-sig'))
    assert pin['repository'] == 'oneworkspacestudio/my-finance'
    assert pin['tag'] == 'v' + contract['version']
    assert contract['schema'] == 1 and contract['id'] == 'finance'
    assert contract['route'] == '/finance' and contract['callback'] == '/finance/auth/callback'
    assert contract['workspaceInstallerIntegration'] == 3
    assert contract['minimumWorkspaceVersion'] == '0.3.0-preview.100'
    assert contract['dataDirectory'] == 'finance' and contract['ledger'] == 'Actual Budget 26.10.0'
    apps = [app for app in catalog['apps'] if app['id'] == 'finance']
    assert len(apps) == 1, 'Expected one Finance app'
    app = apps[0]
    assert app['package'] == 'finance' and app['executable'] == 'MyFinance.exe'
    assert app['route'] == contract['route'] and app['callback'] == contract['callback']
    assert app['kind'] == 'process' and app['oidc'] is True
    assert app['health'] == '/healthz' and app['stripPrefix'] is True
    mapping = {'finance': ('MyFinance.exe', 'MyFinance.exe'),
               'finance-runtime': ('finance-runtime.zip', 'finance-runtime.zip')}
    for name in ['LICENSE', 'NOTICE', 'THIRD-PARTY-LICENSES.txt']:
        mapping['finance-' + name] = ('licenses/MyFinance-' + name, name)
    assert {key for key in catalog['packages'] if key == 'finance' or key.startswith('finance-')} == set(mapping), 'Unexpected Finance package ownership'
    replacements, packages = {}, {}
    for key, (target, member) in mapping.items():
        payload = root / member
        assert payload.is_file() and not payload.is_symlink(), 'Missing Finance payload ' + member
        assert catalog['packages'][key]['file'] == target, ('Unexpected Finance destination', key)
        packages[key] = {'file': target, 'size': payload.stat().st_size, 'sha256': _sha(payload)}
        replacements[target] = payload
    assert packages['finance-runtime']['sha256'] == RUNTIME_SHA256, 'Finance runtime changes need lifecycle acceptance'
    # Validate every payload before mutating the catalog; preserve all launch,
    # dependency, storage and unrelated application/package fields verbatim.
    catalog['packages'].update(packages)
    app['version'] = contract['version']
    app['source'] = 'https://github.com/' + pin['repository'] + '/releases/tag/' + pin['tag']
    return replacements
