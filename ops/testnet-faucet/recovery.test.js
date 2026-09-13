'use strict';
const { test } = require('node:test');
const assert = require('node:assert/strict');
const fs = require('node:fs/promises');
const os = require('node:os');
const path = require('node:path');
const { PassThrough } = require('node:stream');
const { EventEmitter } = require('node:events');
const { Journal, validAddress, clientIp, normalizeIp } = require('./recovery');
const { readBody, runBounded } = require('./bounded');

// Test-only Bech32 encoder: 20 deterministic account bytes, no keys or signatures.
function address(seed = 0, prefix = 'xtc', length = 20) {
  const bytes = Buffer.alloc(length, seed), words = [];
  let acc = 0, bits = 0;
  for (const b of bytes) { acc = (acc << 8) | b; bits += 8; while (bits >= 5) { bits -= 5; words.push((acc >>> bits) & 31); } }
  if (bits) words.push((acc << (5 - bits)) & 31);
  const expanded = [...prefix].map((c) => c.charCodeAt(0) >> 5).concat([0], [...prefix].map((c) => c.charCodeAt(0) & 31));
  let chk = 1;
  for (const v of [...expanded, ...words, 0, 0, 0, 0, 0, 0]) {
    const top = chk >>> 25; chk = ((chk & 0x1ffffff) << 5) ^ v;
    for (const [i, g] of [0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3].entries()) if ((top >>> i) & 1) chk ^= g;
  }
  chk ^= 1;
  const alphabet = 'qpzry9x8gf2tvdw0s3jn54khce6mua7l';
  return prefix + '1' + [...words, ...Array.from({ length: 6 }, (_, i) => (chk >>> (5 * (5 - i))) & 31)].map((v) => alphabet[v]).join('');
}
const policy = { sender: address(99), chainId: 'fixture-chain', amount: '1', addressWindow: 86400000, ipWindow: 86400000, ipLimit: 3, clock: () => 1000 };
async function fixture(t, override = {}) {
  const dir = await fs.mkdtemp(path.join(os.tmpdir(), 'faucet-offline-'));
  t.after(() => fs.rm(dir, { recursive: true })); // Only this test's synthetic files.
  return { dir, journal: await Journal.open(dir, { ...policy, ...override }) };
}

test('checksum, prefix, decoded size and case validated before submission', async (t) => {
  assert(validAddress(address()));
  const { journal } = await fixture(t);
  for (const bad of [address().slice(0, -1) + 'q', address(0, 'cosmos'), address(0, 'xtc', 32), address().toUpperCase(), null]) {
    assert.equal(validAddress(bad), false);
    await assert.rejects(journal.claim(bad, '127.0.0.1', () => assert.fail('must not submit')), /invalid_xitcoin/);
  }
  assert.equal(journal.rows.length, 0);
});

test('untrusted forwarding headers ignored; trusted peer must supply a single valid IP', () => {
  const req = { socket: { remoteAddress: '::ffff:127.0.0.1' }, headers: { 'x-real-ip': '192.0.2.5' } };
  assert.equal(clientIp(req), '127.0.0.1');
  assert.equal(clientIp(req, ['127.0.0.1']), '192.0.2.5');
  assert.equal(normalizeIp('2001:0db8:0:0:0:0:0:1'), '2001:db8::1');
  req.headers['x-real-ip'] = '192.0.2.5, 192.0.2.6';
  assert.throws(() => clientIp(req, ['127.0.0.1']), /invalid_ip/);
});

test('reservation is durable before adapter; restart and old pending claims never resend', async (t) => {
  const { dir, journal } = await fixture(t);
  let calls = 0;
  const result = await journal.claim(address(), '192.0.2.1', async () => {
    calls++;
    const envelope = JSON.parse(await fs.readFile(path.join(dir, 'journal.json')));
    assert.equal(JSON.parse(envelope.payload)[0].state, 'unknown');
    throw new Error('accepted but response lost');
  });
  assert.equal(result.state, 'unknown');
  await journal.close();
  const recovered = await Journal.open(dir, { ...policy, clock: () => 999999999 });
  await assert.rejects(recovered.claim(address(), '192.0.2.2', () => calls++), /address_limit/);
  assert.equal(calls, 1);
});

test('hash is submission evidence, never confirmation; ambiguous output retains quota', async (t) => {
  for (const response of ['', '{broken', 'x'.repeat(64), 'a'.repeat(64)]) {
    const { journal } = await fixture(t);
    const row = await journal.claim(address(), '192.0.2.1', async () => response);
    assert.equal(row.state, response === 'a'.repeat(64) ? 'submitted' : 'unknown');
    await assert.rejects(journal.claim(address(), '192.0.2.2', async () => assert.fail()), /address_limit/);
  }
});

test('one owner across processes; overlapping claims cannot enter adapter', async (t) => {
  const { dir, journal } = await fixture(t);
  await assert.rejects(Journal.open(dir, policy), /EEXIST/);
  let release, entered;
  const started = new Promise((resolve) => { entered = resolve; });
  const first = journal.claim(address(), '192.0.2.1', () => { entered(); return new Promise((resolve) => { release = resolve; }); });
  await started;
  await assert.rejects(journal.claim(address(1), '192.0.2.2', async () => assert.fail()), /journal_unavailable/);
  release('a'.repeat(64)); await first;
});

test('IP quota remains across restart and beyond the window for unknown entries', async (t) => {
  const { dir, journal } = await fixture(t);
  for (let i = 0; i < 3; i++) await journal.claim(address(i), '192.0.2.1', async () => '');
  await journal.close();
  const recovered = await Journal.open(dir, { ...policy, clock: () => 999999999 });
  await assert.rejects(recovered.claim(address(4), '192.0.2.1', async () => assert.fail()), /ip_limit/);
});

test('write failure before reservation prevents submission and poisons owner', async (t) => {
  const { journal } = await fixture(t);
  journal.persist = async () => { throw new Error('disk_failed'); };
  await assert.rejects(journal.claim(address(), '192.0.2.1', async () => assert.fail()), /disk_failed/);
  await assert.rejects(journal.claim(address(1), '192.0.2.2', async () => assert.fail()), /journal_unavailable/);
});

test('crash between reservation and adapter leaves durable reservation blocking retry', async (t) => {
  const { dir, journal } = await fixture(t);
  const persist = journal.persist.bind(journal); let writes = 0;
  journal.persist = async () => { if (++writes === 2) throw new Error('crash'); await persist(); };
  await assert.rejects(journal.claim(address(), '192.0.2.1', async () => assert.fail()), /crash/);
  await journal.close();
  const recovered = await Journal.open(dir, policy);
  assert.equal(recovered.rows[0].state, 'reserved');
  await assert.rejects(recovered.claim(address(), '192.0.2.1', async () => assert.fail()), /address_limit/);
});

test('corrupt journal fails closed without changing evidence; missing journal with snapshots fails', async (t) => {
  const { dir, journal } = await fixture(t); await journal.close();
  const file = path.join(dir, 'journal.json');
  const envelope = JSON.parse(await fs.readFile(file)); envelope.sha256 = 'bad';
  const corrupt = JSON.stringify(envelope); await fs.writeFile(file, corrupt);
  await assert.rejects(Journal.open(dir, policy), /journal_checksum/);
  assert.equal(await fs.readFile(file, 'utf8'), corrupt);
  await fs.rename(file, path.join(dir, 'corrupt-evidence.json'));
  await assert.rejects(Journal.open(dir, policy), /incomplete_journal/);
});

test('legacy file remains byte-identical; unproven claims stay unknown', async (t) => {
  const parent = await fs.mkdtemp(path.join(os.tmpdir(), 'faucet-legacy-'));
  t.after(() => fs.rm(parent, { recursive: true }));
  const legacyFile = path.join(parent, 'claims.json');
  const text = JSON.stringify({ claims: [{ address: address(), ip: '192.0.2.1', at: 1, txhash: 'a'.repeat(64) }] });
  await fs.writeFile(legacyFile, text);
  const journal = await Journal.open(path.join(parent, 'new'), { ...policy, legacyFile });
  assert.equal(await fs.readFile(legacyFile, 'utf8'), text);
  assert.equal(journal.rows[0].state, 'unknown');
  await assert.rejects(journal.claim(address(), '192.0.2.1', async () => assert.fail()), /address_limit/);
});

test('body byte limit and absolute deadline', async () => {
  let stream = new PassThrough(); const result = readBody(stream, 4, 100);
  stream.end(Buffer.from('ééé')); await assert.rejects(result, /body_limit/);
  stream = new PassThrough(); await assert.rejects(readBody(stream, 4, 5), /body_timeout/);
});

test('fake child timeout, output cap, nonzero exit and success; no process spawned', async () => {
  for (const kind of ['timeout', 'output', 'exit', 'success']) {
    const child = new EventEmitter(); child.stdout = new PassThrough(); child.stderr = new PassThrough(); let kills = 0; child.kill = () => { kills++; };
    const result = runBounded('never-executed', [], {}, { spawnChild: () => child, timeout: 5, limit: 4 });
    if (kind === 'output') child.stderr.write('12345');
    if (kind === 'exit') child.emit('close', 1);
    if (kind === 'success') { child.stdout.write('ok'); child.emit('close', 0); }
    if (kind === 'success') { assert.equal(await result, 'ok'); assert.equal(kills, 0); }
    else { await assert.rejects(result, /unknown/); assert.equal(kills, 1); }
  }
});

test('failure after accepted adapter result retains unknown on disk', async (t) => {
  const { dir, journal } = await fixture(t);
  const persist = journal.persist.bind(journal); let writes = 0;
  journal.persist = async () => { if (++writes === 3) throw new Error('crash_after_submit'); await persist(); };
  let calls = 0;
  await assert.rejects(journal.claim(address(), '192.0.2.1', async () => { calls++; return 'a'.repeat(64); }), /crash_after_submit/);
  await journal.close();
  const recovered = await Journal.open(dir, policy);
  assert.equal(recovered.rows[0].state, 'unknown');
  await assert.rejects(recovered.claim(address(), '192.0.2.1', async () => calls++), /address_limit/);
  assert.equal(calls, 1);
});

test('malformed legacy history is preserved and prevents initialization', async (t) => {
  const parent = await fs.mkdtemp(path.join(os.tmpdir(), 'faucet-bad-legacy-'));
  t.after(() => fs.rm(parent, { recursive: true }));
  const legacyFile = path.join(parent, 'claims.json');
  await fs.writeFile(legacyFile, '{truncated');
  await assert.rejects(Journal.open(path.join(parent, 'new'), { ...policy, legacyFile }), SyntaxError);
  assert.equal(await fs.readFile(legacyFile, 'utf8'), '{truncated');
});

test('receipt reconciliation requires exact transaction context and inclusion', async (t) => {
  for (const changed of ['txhash', 'chainId', 'sender', 'recipient', 'denom', 'amount', 'height', 'code', 'absent', 'timeout', 'confirmed', 'failed']) {
    const { journal } = await fixture(t);
    const row = await journal.claim(address(), '192.0.2.1', async () => 'a'.repeat(64));
    const receipt = { txhash: row.txhash, chainId: row.chainId, sender: row.sender, recipient: row.address, denom: row.denom, amount: row.amount, height: 1, code: changed === 'failed' ? 5 : 0 };
    if (Object.hasOwn(receipt, changed)) receipt[changed] = changed === 'height' ? 0 : 'mismatch';
    let aborted = false;
    const result = await journal.reconcile(row.id, async (_, signal) => {
      if (changed === 'absent') return null;
      if (changed === 'timeout') { signal.addEventListener('abort', () => { aborted = true; }); return new Promise(() => {}); }
      return receipt;
    }, 5);
    assert.equal(result.state, changed === 'confirmed' ? 'confirmed' : changed === 'failed' ? 'failed_definite' : 'unknown');
    if (changed === 'timeout') assert(aborted);
    await assert.rejects(journal.claim(address(), '192.0.2.2', async () => assert.fail()), /address_limit/);
  }
});

test('request without hash never queries or infers non-submission', async (t) => {
  const { journal } = await fixture(t);
  const row = await journal.claim(address(), '192.0.2.1', async () => '');
  assert.equal((await journal.reconcile(row.id, async () => assert.fail())).state, 'unknown');
});
