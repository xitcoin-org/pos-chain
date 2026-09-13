'use strict';

const fs = require('node:fs/promises');
const path = require('node:path');
const { createHash, randomUUID } = require('node:crypto');
const { isIP } = require('node:net');
const MAX_BYTES = 8 * 1024 * 1024;
const PENDING = new Set(['reserved', 'prepared', 'submitted', 'unknown']);
const STATES = new Set([...PENDING, 'confirmed', 'failed_definite']);
const digest = (text) => createHash('sha256').update(text).digest('hex');

function validAddress(address) {
  if (typeof address !== 'string' || !/^xtc1[023456789acdefghjklmnpqrstuvwxyz]{38}$/.test(address)) return false;
  const alphabet = 'qpzry9x8gf2tvdw0s3jn54khce6mua7l';
  const values = [...address.slice(4)].map((c) => alphabet.indexOf(c));
  const prefix = [...'xtc'].map((c) => c.charCodeAt(0));
  let chk = 1;
  for (const value of [...prefix.map((c) => c >> 5), 0, ...prefix.map((c) => c & 31), ...values]) {
    const top = chk >>> 25;
    chk = ((chk & 0x1ffffff) << 5) ^ value;
    [0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3].forEach((g, i) => { if ((top >>> i) & 1) chk ^= g; });
  }
  // 32 data words encode exactly 20 account bytes, with no padding.
  return (chk >>> 0) === 1;
}

function normalizeIp(value) {
  if (typeof value !== 'string') throw new Error('invalid_ip');
  let ip = value.trim();
  if (ip.startsWith('::ffff:') && isIP(ip.slice(7)) === 4) ip = ip.slice(7);
  if (!isIP(ip) || ip.includes('%')) throw new Error('invalid_ip');
  return isIP(ip) === 6 ? new URL(`http://[${ip}]/`).hostname.slice(1, -1) : ip;
}

function clientIp(req, trustedPeers = []) {
  const peer = normalizeIp(req.socket.remoteAddress);
  const trusted = trustedPeers.map(normalizeIp).includes(peer);
  return trusted && req.headers['x-real-ip'] !== undefined ? normalizeIp(req.headers['x-real-ip']) : peer;
}

function validate(rows) {
  if (!Array.isArray(rows)) throw new Error('invalid_journal');
  const ids = new Set();
  for (const row of rows) {
    if (!row || typeof row.id !== 'string' || ids.has(row.id) || !STATES.has(row.state) ||
        !validAddress(row.address) || !validAddress(row.sender) || typeof row.ip !== 'string' || normalizeIp(row.ip) !== row.ip ||
        !Number.isSafeInteger(row.at) || row.at < 0 || typeof row.amount !== 'string' || typeof row.chainId !== 'string' || !row.chainId ||
        row.denom !== 'axtc' || !/^[1-9][0-9]*$/.test(row.amount) ||
        (row.txhash !== null && !/^[0-9A-F]{64}$/.test(row.txhash))) throw new Error('invalid_journal');
    ids.add(row.id);
  }
}

async function readBounded(file) {
  if ((await fs.stat(file)).size > MAX_BYTES) throw new Error('journal_capacity');
  return fs.readFile(file, 'utf8');
}

class Journal {
  static async open(directory, { legacyFile, chainId, sender, amount, addressWindow, ipWindow, ipLimit, clock = Date.now } = {}) {
    if (!validAddress(sender) || !chainId || !/^[1-9][0-9]*$/.test(String(amount)) ||
        ![addressWindow, ipWindow, ipLimit].every((x) => Number.isSafeInteger(x) && x > 0)) throw new Error('invalid_policy');
    await fs.mkdir(directory, { recursive: true, mode: 0o700 });
    const lock = path.join(directory, 'writer.lock');
    await fs.mkdir(lock, { mode: 0o700 }); // Never steal a lock, even after a crash.
    const journal = new Journal();
    Object.assign(journal, { directory, lock, chainId, sender, amount: String(amount), addressWindow, ipWindow, ipLimit, clock, rows: [], busy: false, poisoned: false });
    try {
      const file = path.join(directory, 'journal.json');
      let text;
      try { text = await readBounded(file); } catch (err) { if (err.code !== 'ENOENT') throw err; }
      if (text !== undefined) {
        const envelope = JSON.parse(text);
        if (envelope.version !== 1 || typeof envelope.payload !== 'string' || digest(envelope.payload) !== envelope.sha256) throw new Error('journal_checksum');
        journal.rows = JSON.parse(envelope.payload);
        validate(journal.rows);
      } else {
        const leftovers = (await fs.readdir(directory)).filter((name) => name !== 'writer.lock');
        if (leftovers.length) throw new Error('incomplete_journal_requires_review');
        if (legacyFile) {
          let legacy;
          try { legacy = JSON.parse(await readBounded(legacyFile)); } catch (err) { if (err.code !== 'ENOENT') throw err; }
          if (legacy !== undefined) {
            if (!Array.isArray(legacy.claims)) throw new Error('invalid_legacy');
            journal.rows = legacy.claims.map((row) => ({ id: randomUUID(), address: row.address, ip: normalizeIp(row.ip), at: row.at,
              chainId, sender, denom: 'axtc', amount: String(amount), state: 'unknown', txhash: typeof row.txhash === 'string' && /^[a-fA-F0-9]{64}$/.test(row.txhash) ? row.txhash.toUpperCase() : null }));
            validate(journal.rows); // Legacy broadcast hashes do not establish inclusion.
          }
        }
        await journal.persist();
      }
      if (journal.rows.some((row) => row.sender !== sender || row.chainId !== chainId || row.amount !== String(amount))) throw new Error('policy_changed_requires_review');
      return journal;
    } catch (err) {
      await fs.rmdir(lock);
      throw err;
    }
  }

  async persist() {
    validate(this.rows);
    const payload = JSON.stringify(this.rows);
    const text = JSON.stringify({ version: 1, sha256: digest(payload), payload });
    if (Buffer.byteLength(text) > MAX_BYTES) throw new Error('journal_capacity');
    // Keep immutable recovery snapshots; at the limit, stop instead of deleting history.
    const snapshots = (await fs.readdir(this.directory)).filter((name) => name.startsWith('snapshot-'));
    if (snapshots.length >= 64) throw new Error('journal_capacity_requires_archival');
    const snapshot = path.join(this.directory, `snapshot-${randomUUID()}.json`);
    const next = path.join(this.directory, `next-${randomUUID()}.json`);
    for (const file of [snapshot, next]) {
      const handle = await fs.open(file, 'wx', 0o600);
      try { await handle.writeFile(text); await handle.sync(); } finally { await handle.close(); }
    }
    await fs.rename(next, path.join(this.directory, 'journal.json'));
    const dir = await fs.open(this.directory, 'r');
    try { await dir.sync(); } finally { await dir.close(); }
  }

  async exclusive(action) {
    if (this.busy || this.poisoned) throw new Error('journal_unavailable');
    this.busy = true;
    try { return await action(); } finally { this.busy = false; }
  }

  async commit() {
    try { await this.persist(); } catch (err) { this.poisoned = true; throw err; }
  }

  async claim(address, ip, submit) {
    return this.exclusive(async () => {
      if (!validAddress(address)) throw new Error('invalid_xitcoin_address');
      ip = normalizeIp(ip);
      const now = this.clock();
      if (!Number.isSafeInteger(now) || now < 0) throw new Error('invalid_clock');
      if (this.rows.some((r) => r.address === address && (PENDING.has(r.state) || now - r.at < this.addressWindow))) throw new Error('address_limit');
      if (this.rows.filter((r) => r.ip === ip && (PENDING.has(r.state) || now - r.at < this.ipWindow)).length >= this.ipLimit) throw new Error('ip_limit');
      const row = { id: randomUUID(), address, ip, at: now, chainId: this.chainId, sender: this.sender, denom: 'axtc', amount: this.amount, state: 'reserved', txhash: null };
      this.rows.push(row);
      await this.commit(); // Nothing has reached the submission adapter yet.
      row.state = 'unknown';
      await this.commit(); // CLI has no prepare/broadcast split: fail closed before invoking it.
      try {
        const hash = await submit(address);
        if (typeof hash === 'string' && /^[a-fA-F0-9]{64}$/.test(hash)) {
          row.txhash = hash.toUpperCase();
          row.state = 'submitted';
        }
      } catch { /* Ambiguous child outcomes never release quota or trigger a resend. */ }
      await this.commit();
      return { ...row };
    });
  }

  // lookup must be a read-only receipt adapter; this method never submits.
  // The HTTP server deliberately does not expose reconciliation as a public API.
  async reconcile(id, lookup, timeout = 5000) {
    return this.exclusive(async () => {
      const row = this.rows.find((entry) => entry.id === id);
      if (!row) throw new Error('unknown_request');
      if (!PENDING.has(row.state) || !row.txhash) return { ...row };
      let timer, receipt;
      const controller = new AbortController();
      try {
        receipt = await Promise.race([
          Promise.resolve().then(() => lookup(row.txhash, controller.signal)),
          new Promise((_, reject) => { timer = setTimeout(() => { controller.abort(); reject(new Error('receipt_timeout')); }, timeout); }),
        ]);
      } catch { /* Absence, timeout and read failure are not non-submission. */ }
      finally { clearTimeout(timer); }
      const matches = receipt && receipt.txhash === row.txhash && receipt.chainId === row.chainId &&
        receipt.sender === row.sender && receipt.recipient === row.address && receipt.denom === row.denom &&
        receipt.amount === row.amount && Number.isSafeInteger(receipt.height) && receipt.height > 0 &&
        Number.isSafeInteger(receipt.code) && receipt.code >= 0;
      row.state = matches ? (receipt.code === 0 ? 'confirmed' : 'failed_definite') : 'unknown';
      await this.commit();
      return { ...row };
    });
  }

  async close() {
    if (this.busy) throw new Error('journal_busy');
    this.poisoned = true;
    await fs.rmdir(this.lock);
  }
}

module.exports = { Journal, validAddress, clientIp, normalizeIp };
