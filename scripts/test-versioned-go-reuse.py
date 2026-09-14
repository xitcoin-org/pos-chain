#!/usr/bin/env python3
"""Real filesystem/Git and archive mutations; no Go suite is rerun."""
import copy
import importlib.util
import json
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import unittest
import zipfile

sys.dont_write_bytecode = True
SOURCE = Path(__file__).resolve().parents[1]
spec = importlib.util.spec_from_file_location('reuse', SOURCE / 'scripts/verify-versioned-go-reuse.py')
guard = importlib.util.module_from_spec(spec)
spec.loader.exec_module(guard)


class ReuseTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.temp = tempfile.TemporaryDirectory(prefix='xitcoin-reuse-')
        cls.root = Path(cls.temp.name) / 'repo'
        subprocess.run(['git', 'clone', '--quiet', '--shared', '--no-checkout', str(SOURCE), str(cls.root)], check=True)
        subprocess.run(['git', '-C', str(cls.root), 'checkout', '--quiet', 'HEAD'], check=True)
        # GitHub checkout is shallow/detached. Unreferenced fetched ancestors
        # need an explicit local fetch; clone does not promise to copy them.
        subprocess.run(['git', '-C', str(cls.root), 'fetch', '--quiet', '--depth=2',
                        str(SOURCE), guard.REVIEWED_HEAD], check=True)
        cls.policy = guard.read_policy(SOURCE)
        for name in cls.policy['mission_files'].keys() | {guard.POLICY}:
            shutil.copy2(SOURCE / name, cls.root / name)
        cls.history = Path(cls.temp.name) / 'history'
        # CI already downloaded and checked all four archives before this test.
        source_history = SOURCE / 'versioned-evidence/historical-artifacts'
        shutil.copytree(source_history, cls.history)

    @classmethod
    def tearDownClass(cls):
        cls.temp.cleanup()

    def test_exact_delta(self):
        result = guard.validate(self.root, self.history)
        self.assertFalse(result['security_acceptance'])
        self.assertEqual(len(result['history']), 4)
        self.assertTrue(result['scans_are_historical'])

    def test_source_mutations(self):
        paths = ['evmd/app.go', 'go.mod', 'go.sum', 'evmd/go.mod', 'evmd/go.sum',
                 'tests/integration/wallets/test_ledger_suite.go',
                 'tests/integration/wallets/test_legder.go',
                 'tests/integration/x/erc20/test_ibc_callback.go',
                 'tests/integration/x/erc20/test_proposals.go',
                 'scripts/verify-go-fork-provenance.py',
                 'docs/go-fork-provenance.json', 'docs/versioned-go-qualification-resume.json',
                 *guard.GATE_FILES, *self.policy['mission_files']]
        for name in paths:
            with self.subTest(path=name):
                path = self.root / name
                original = path.read_bytes()
                try:
                    path.write_bytes(original + b'\n# unqualified mutation\n')
                    with self.assertRaises(ValueError):
                        guard.validate(self.root, self.history)
                finally:
                    path.write_bytes(original)

    def test_deleted_mode_symlink_and_new_inputs(self):
        path = self.root / 'scripts/run-govulncheck.sh'
        raw, mode = path.read_bytes(), path.stat().st_mode
        try:
            path.unlink()
            with self.assertRaises(ValueError): guard.validate(self.root, self.history)
            path.symlink_to(self.root / 'README.md')
            with self.assertRaises(ValueError): guard.validate(self.root, self.history)
            path.unlink(); path.write_bytes(raw); path.chmod(mode ^ 0o100)
            with self.assertRaises(ValueError): guard.validate(self.root, self.history)
        finally:
            if path.is_symlink(): path.unlink()
            path.write_bytes(raw); path.chmod(mode)
        for name in ['unqualified.go', 'scripts/unqualified.py', 'versioned-evidence/unqualified.go']:
            path = self.root / name
            path.parent.mkdir(exist_ok=True)
            try:
                path.write_text('unqualified\n')
                with self.assertRaises(ValueError): guard.validate(self.root, self.history)
            finally: path.unlink()

    def test_staged_source_mutation(self):
        path = self.root / 'go.mod'; original = path.read_bytes()
        try:
            path.write_bytes(original + b'\n// staged drift\n')
            subprocess.run(['git', '-C', str(self.root), 'add', 'go.mod'], check=True)
            path.write_bytes(original)
            with self.assertRaisesRegex(ValueError, 'staged'):
                guard.validate(self.root, self.history)
        finally:
            subprocess.run(['git', '-C', str(self.root), 'reset', '--quiet', 'HEAD', '--', 'go.mod'], check=True)
            path.write_bytes(original)

    def test_policy_mutations(self):
        path = self.root / guard.POLICY
        original = path.read_bytes()
        path.write_bytes(original + b' ')
        try:
            with self.assertRaises(ValueError): guard.validate(self.root, self.history)
        finally: path.write_bytes(original)
        mutations = []
        for key, value in [('qualified_head', guard.REVIEWED_HEAD), ('review_expired', '2099-01-01'),
                           ('security_acceptance', True), ('reason', 'different')]:
            p = copy.deepcopy(self.policy); p[key] = value; mutations.append(p)
        p = copy.deepcopy(self.policy); p['mission_files'][guard.VERIFIER]['sha256'] = '0' * 64; mutations.append(p)
        p = copy.deepcopy(self.policy); p['mission_files']['unqualified.go'] = {'sha256': '0'*64, 'mode': '100644'}; mutations.append(p)
        p = copy.deepcopy(self.policy); p['artifacts'] = []; mutations.append(p)
        try:
            for value in mutations:
                with self.subTest(value=value.get('reason')):
                    path.write_text(json.dumps(value))
                    with self.assertRaises(ValueError): guard.validate(self.root, self.history)
        finally: path.write_bytes(original)

    def test_artifact_absent_or_altered(self):
        for item in self.policy['artifacts']:
            path = self.history / item['file']; original = path.read_bytes()
            try:
                path.unlink()
                with self.assertRaisesRegex(ValueError, 'absent'): guard.validate(self.root, self.history)
                path.write_bytes(original + b'altered')
                with self.assertRaisesRegex(ValueError, 'altered'): guard.validate(self.root, self.history)
            finally: path.write_bytes(original)

    def test_historical_semantics_even_with_matching_archive_hash(self):
        # Test the result checks separately from the immutable policy pin.
        item = next(a for a in self.policy['artifacts'] if a['run_id'] == 34832506821 and a['module'] == 'evmd')
        path = self.history / item['file']; original = path.read_bytes()
        with zipfile.ZipFile(path) as archive:
            entries = {name: archive.read(name) for name in archive.namelist()}
        for kind in ['missing-results', 'empty-results', 'failed-tests', 'blocked-tests', 'wrong-command',
                     'missing-log', 'wrong-manifest', 'accepted', 'wrong-binary', 'wrong-head']:
            with self.subTest(kind=kind):
                mutated = dict(entries); policy = copy.deepcopy(self.policy)
                a = next(a for a in policy['artifacts'] if a['file'] == item['file'])
                if kind == 'missing-results': del mutated['results.json']
                elif kind == 'empty-results': mutated['results.json'] = b'{}'
                elif kind in ['failed-tests', 'blocked-tests', 'wrong-command']:
                    r = json.loads(mutated['results.json'])
                    field, value = {'failed-tests': ('exit_code', 1), 'blocked-tests': ('blocked', 'resource'),
                                    'wrong-command': ('command', ['true'])}[kind]
                    r['tests'][field] = value; mutated['results.json'] = json.dumps(r).encode()
                elif kind == 'missing-log': del mutated['tests.log']
                elif kind == 'wrong-manifest': mutated['preserved-checks.json'] = b'{}'
                elif kind == 'accepted':
                    q = json.loads(mutated['qualification.json']); q['security_acceptance'] = True
                    mutated['qualification.json'] = json.dumps(q).encode()
                elif kind == 'wrong-binary': mutated['buildinfo.log'] = b'(devel)'
                elif kind == 'wrong-head': a['head_sha'] = guard.REVIEWED_HEAD
                try:
                    with zipfile.ZipFile(path, 'w') as archive:
                        for name, raw in mutated.items(): archive.writestr(name, raw)
                    a['sha256'] = guard.digest(path.read_bytes())
                    with self.assertRaises((ValueError, KeyError)):
                        guard.validate_history(self.root, policy, self.history)
                finally: path.write_bytes(original)


if __name__ == '__main__':
    unittest.main(verbosity=2)
