'use strict';

const http = require('node:http');
const fs = require('node:fs/promises');
const path = require('node:path');
const { spawn } = require('node:child_process');

const cfg = {
  host: process.env.HOST,
  port: Number(process.env.PORT),
  chainId: process.env.CHAIN_ID,
  rpcNode: process.env.RPC_NODE,
  apiUrl: process.env.API_URL.replace(/\/$/, ''),
  bin: process.env.BIN,
  faucetHome: process.env.FAUCET_HOME,
  keyName: process.env.KEY_NAME,
  keyringDir: process.env.KEYRING_DIR,
  amount: BigInt(process.env.AMOUNT_BASE),
  reserve: BigInt(process.env.RESERVE_BASE),
  addressWindow: Number(process.env.ADDRESS_WINDOW_SECONDS) * 1000,
  ipWindow: Number(process.env.IP_WINDOW_SECONDS) * 1000,
  ipLimit: Number(process.env.IP_LIMIT),
};

const stateFile = path.join(cfg.faucetHome, 'state', 'claims.json');
let lock = Promise.resolve();

function formatXtc(amount) {
  const unit = 10n ** 18n;
  const whole = amount / unit;
  const fraction = String(amount % unit).padStart(18, '0').replace(/0+$/, '');
  return fraction ? `${whole}.${fraction}` : String(whole);
}

function send(res, status, body) {
  res.writeHead(status, { 'content-type': 'application/json; charset=utf-8', 'cache-control': 'no-store' });
  res.end(JSON.stringify(body));
}

async function stateRead() {
  try {
    return JSON.parse(await fs.readFile(stateFile, 'utf8'));
  } catch (err) {
    if (err.code === 'ENOENT') return { claims: [] };
    throw err;
  }
}

async function stateWrite(state) {
  const next = `${stateFile}.next`;
  await fs.writeFile(next, JSON.stringify(state), { mode: 0o600 });
  await fs.rename(next, stateFile);
}

async function faucetAddress() {
  return (await fs.readFile(path.join(cfg.faucetHome, 'state', 'address'), 'utf8')).trim();
}

async function balance(address) {
  const response = await fetch(`${cfg.apiUrl}/cosmos/bank/v1beta1/balances/${address}`);
  if (!response.ok) throw new Error(`balance query HTTP ${response.status}`);
  const data = await response.json();
  const coin = (data.balances || []).find((item) => item.denom === 'axtc');
  return BigInt(coin ? coin.amount : '0');
}

function runTx(to) {
  return new Promise((resolve, reject) => {
    const args = [
      'tx', 'bank', 'send', cfg.keyName, to, `${cfg.amount}axtc`,
      '--home', cfg.faucetHome,
      '--keyring-dir', cfg.keyringDir,
      '--keyring-backend', 'test',
      '--chain-id', cfg.chainId,
      '--node', cfg.rpcNode,
      '--gas', 'auto',
      '--gas-adjustment', '1.3',
      '--gas-prices', '0.000000000000000007axtc',
      '--yes',
      '--output', 'json',
    ];
    const child = spawn(cfg.bin, args, { cwd: cfg.faucetHome, env: { ...process.env, HOME: cfg.faucetHome } });
    let stdout = '';
    let stderr = '';
    child.stdout.on('data', (chunk) => { stdout += chunk; });
    child.stderr.on('data', (chunk) => { stderr += chunk; });
    child.on('error', reject);
    child.on('close', (code) => {
      if (code !== 0) return reject(new Error(stderr || `tx exited ${code}`));
      try {
        const result = JSON.parse(stdout);
        if (result.code && Number(result.code) !== 0) return reject(new Error(result.raw_log || 'transaction rejected'));
        resolve(result.txhash || result.tx_response?.txhash || '');
      } catch {
        reject(new Error('invalid transaction response'));
      }
    });
  });
}

function clientIp(req) {
  return String(req.headers['x-real-ip'] || req.socket.remoteAddress || '').trim();
}

const server = http.createServer(async (req, res) => {
  try {
    if (req.method === 'GET' && req.url === '/healthz') {
      const address = await faucetAddress();
      const available = await balance(address);
      return send(res, 200, {
        status: 'ok',
        chain_id: cfg.chainId,
        faucet_address: address,
        claim_amount_xtc: formatXtc(cfg.amount),
        funded: available >= cfg.amount + cfg.reserve,
      });
    }

    if (req.method !== 'POST' || req.url !== '/claim') {
      return send(res, 404, { error: 'not_found' });
    }

    let raw = '';
    req.on('data', (chunk) => {
      raw += chunk;
      if (raw.length > 2048) req.destroy();
    });
    req.on('end', () => {
      lock = lock.then(async () => {
        const body = JSON.parse(raw || '{}');
        const address = String(body.address || '');
        if (!/^xtc1[023456789acdefghjklmnpqrstuvwxyz]{38,90}$/.test(address)) {
          return send(res, 400, { error: 'invalid_xitcoin_address' });
        }

        const now = Date.now();
        const ip = clientIp(req);
        const state = await stateRead();
        const recentAddress = state.claims.find((c) => c.address === address && now - c.at < cfg.addressWindow);
        const recentIp = state.claims.filter((c) => c.ip === ip && now - c.at < cfg.ipWindow);

        if (recentAddress) return send(res, 429, { error: 'address_limit', retry_after_seconds: Math.ceil((cfg.addressWindow - (now - recentAddress.at)) / 1000) });
        if (recentIp.length >= cfg.ipLimit) return send(res, 429, { error: 'ip_limit' });

        const source = await faucetAddress();
        const available = await balance(source);
        if (available < cfg.amount + cfg.reserve) {
          return send(res, 503, { error: 'faucet_not_funded' });
        }

        const txhash = await runTx(address);
        state.claims.push({ address, ip, at: now, txhash });
        state.claims = state.claims.filter((c) => now - c.at < Math.max(cfg.addressWindow, cfg.ipWindow));
        await stateWrite(state);

        return send(res, 200, { ok: true, amount_xtc: formatXtc(cfg.amount), txhash });
      }).catch((err) => {
        console.error(err);
        if (!res.headersSent) send(res, 500, { error: 'claim_failed' });
      });
    });
  } catch (err) {
    console.error(err);
    send(res, 500, { error: 'internal_error' });
  }
});

server.listen(cfg.port, cfg.host, () => {
  console.log(`Xitcoin faucet listening on ${cfg.host}:${cfg.port}`);
});
