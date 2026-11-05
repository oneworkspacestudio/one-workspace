"""Finance managed-seed pin, ownership and atomic validation regressions."""
import copy
import hashlib
import json
import os
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch

from finance_component import finance_replacements


class FinanceComponentTests(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory(dir=os.environ.get('MYFINANCE_TEST_TEMP'))
        self.addCleanup(self.tmp.cleanup)
        self.root = Path(self.tmp.name)
        self.contract = dict(schema=1, id='finance', version='0.1.0-preview.4',
                             route='/finance', callback='/finance/auth/callback',
                             workspaceInstallerIntegration=3, minimumWorkspaceVersion='0.3.0-preview.100',
                             dataDirectory='finance', ledger='Actual Budget 26.10.0')
        self.write_contract()
        self.catalog = dict(apps=[dict(id='finance', package='finance', executable='MyFinance.exe',
                                      route='/finance', callback='/finance/auth/callback', kind='process',
                                      health='/healthz', stripPrefix=True, oidc=True, port=9024,
                                      args=['--data', '${DATA}/finance'], depends=['accounts'], version='old'),
                                 dict(id='other', version='keep')],
                            packages={'other': dict(file='Other.exe', size=17, sha256='a' * 64)})
        for name in ['MyFinance.exe', 'finance-runtime.zip', 'LICENSE', 'NOTICE', 'THIRD-PARTY-LICENSES.txt']:
            (self.root / name).write_bytes(('synthetic ' + name).encode())
            key = 'finance' if name == 'MyFinance.exe' else 'finance-runtime' if name == 'finance-runtime.zip' else 'finance-' + name
            target = name if key in ['finance', 'finance-runtime'] else 'licenses/MyFinance-' + name
            self.catalog['packages'][key] = dict(file=target, size=1, sha256='0' * 64)
        runtime = hashlib.sha256((self.root / 'finance-runtime.zip').read_bytes()).hexdigest()
        runtime_patch = patch('finance_component.RUNTIME_SHA256', runtime)
        runtime_patch.start()
        self.addCleanup(runtime_patch.stop)
        self.pin = dict(repository='oneworkspacestudio/my-finance', tag='v0.1.0-preview.4')

    def write_contract(self):
        (self.root / 'compatibility.json').write_text(json.dumps(self.contract), encoding='utf-8')

    def test_only_five_owned_payloads_and_version_change(self):
        before = copy.deepcopy(self.catalog)
        replacements = finance_replacements(self.catalog, self.root, self.pin)
        self.assertEqual(len(replacements), 5)
        self.assertEqual(self.catalog['apps'][1:], before['apps'][1:])
        self.assertEqual(self.catalog['packages']['other'], before['packages']['other'])
        expected_app = dict(before['apps'][0], version='0.1.0-preview.4',
                            source='https://github.com/oneworkspacestudio/my-finance/releases/tag/v0.1.0-preview.4')
        self.assertEqual(self.catalog['apps'][0], expected_app)
        for key, package in self.catalog['packages'].items():
            if key == 'other':
                continue
            payload = replacements[package['file']].read_bytes()
            self.assertEqual(package['size'], len(payload))
            self.assertEqual(package['sha256'], hashlib.sha256(payload).hexdigest())
            self.assertEqual(package['file'], before['packages'][key]['file'])

    def test_incompatible_contract_or_pin_rejected_without_catalog_changes(self):
        for key, value in [('schema', 2), ('id', 'other'), ('version', 'wrong'),
                           ('route', '/other'), ('callback', '/other/auth/callback'),
                           ('workspaceInstallerIntegration', 2), ('minimumWorkspaceVersion', '0.3.0-preview.99'),
                           ('dataDirectory', 'other'), ('ledger', 'Actual Budget new')]:
            with self.subTest(key=key):
                old = self.contract[key]
                self.contract[key] = value
                self.write_contract()
                before = copy.deepcopy(self.catalog)
                with self.assertRaises(AssertionError):
                    finance_replacements(self.catalog, self.root, self.pin)
                self.assertEqual(self.catalog, before)
                self.contract[key] = old
        self.write_contract()
        with self.assertRaises(AssertionError):
            finance_replacements(self.catalog, self.root, dict(self.pin, repository='elsewhere/finance'))

    def test_missing_or_changed_runtime_and_notice_are_atomic(self):
        for name in ['NOTICE', 'finance-runtime.zip']:
            with self.subTest(name=name):
                path = self.root / name
                original = path.read_bytes()
                path.unlink()
                before = copy.deepcopy(self.catalog)
                with self.assertRaises(AssertionError):
                    finance_replacements(self.catalog, self.root, self.pin)
                self.assertEqual(self.catalog, before)
                path.write_bytes(original)
        (self.root / 'finance-runtime.zip').write_bytes(b'unaccepted runtime')
        before = copy.deepcopy(self.catalog)
        with self.assertRaises(AssertionError):
            finance_replacements(self.catalog, self.root, self.pin)
        self.assertEqual(self.catalog, before)

    def test_foreign_ownership_destinations_and_auth_changes_fail_closed(self):
        mutations = [lambda c: c['packages'].update({'finance-extra': dict(file='Other.exe')}),
                     lambda c: c['packages']['finance-runtime'].update(file='Other.exe'),
                     lambda c: c['apps'][0].update(oidc=False),
                     lambda c: c['apps'].append(copy.deepcopy(c['apps'][0]))]
        for mutate in mutations:
            candidate = copy.deepcopy(self.catalog)
            mutate(candidate)
            before = copy.deepcopy(candidate)
            with self.assertRaises(AssertionError):
                finance_replacements(candidate, self.root, self.pin)
            self.assertEqual(candidate, before)


if __name__ == '__main__':
    unittest.main()
