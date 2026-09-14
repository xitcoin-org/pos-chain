'use strict';

const http = require('node:http');
const fs = require('node:fs/promises');
const path = require('node:path');
const { Journal, validAddress, clientIp } = require('./recovery');
const { readBody, runBounded } = require('./bounded');

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
let inFlight = 0;
let journal;
const trustedPeers = (process.env.TRUSTED_PROXY_IPS || '').split(',').map((s) => s.trim()).filter(Boolean);

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

async function faucetAddress() {
  return (await fs.readFile(path.join(cfg.faucetHome, 'state', 'address'), 'utf8')).trim();
}

async function balance(address) {
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), 5000);
  try {
  const response = await fetch(`${cfg.apiUrl}/cosmos/bank/v1beta1/balances/${address}`, { signal: controller.signal });
  if (!response.ok) throw new Error(`balance query HTTP ${response.status}`);
  let text = '';
  let bytes = 0;
  for await (const chunk of response.body) {
    bytes += chunk.length;
    if (bytes > 65536) { controller.abort(); throw new Error('balance_output_limit'); }
    text += Buffer.from(chunk).toString();
  }
  const data = JSON.parse(text);
  const coin = (data.balances || []).find((item) => item.denom === 'axtc');
  return BigInt(coin ? coin.amount : '0');
  } finally { clearTimeout(timer); }
}

async function runTx(to) {
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
  const stdout = await runBounded(cfg.bin, args, { cwd: cfg.faucetHome, env: { ...process.env, HOME: cfg.faucetHome } });
  const result = JSON.parse(stdout);
  if (Number(result.code || result.tx_response?.code || 0) !== 0) throw new Error('submission_unknown');
  return result.txhash || result.tx_response?.txhash || '';
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

    if (inFlight >= 16) return send(res, 503, { error: 'faucet_busy' });
    inFlight++;
    try {
      const raw = await readBody(req);
      const body = JSON.parse(raw || '{}');
      const address = body?.address;
      if (!validAddress(address)) return send(res, 400, { error: 'invalid_xitcoin_address' });
      const ip = clientIp(req, trustedPeers);
      const source = await faucetAddress();
      if (await balance(source) < cfg.amount + cfg.reserve) return send(res, 503, { error: 'faucet_not_funded' });
      const claim = await journal.claim(address, ip, runTx);
      return send(res, 202, { status: claim.state, request_id: claim.id, txhash: claim.txhash, reconciliation_required: true });
    } finally { inFlight--; }
  } catch (err) {
    const name = String(err.message);
    const status = ['address_limit', 'ip_limit'].includes(name) ? 429 :
      name === 'journal_unavailable' ? 503 :
      ['body_limit', 'body_timeout', 'body_aborted', 'invalid_ip'].includes(name) || err instanceof SyntaxError ? 400 : 500;
    console.error('faucet request failed', status);
    if (!res.headersSent) send(res, status, { error: status === 500 ? 'internal_error' : status === 400 ? 'invalid_request' : name });
  }
});

server.requestTimeout = 10000;
server.headersTimeout = 10000;
server.maxConnections = 32;

async function start() {
  journal = await Journal.open(path.join(cfg.faucetHome, 'state', 'recovery-v1'), {
    legacyFile: stateFile, sender: await faucetAddress(), chainId: cfg.chainId, amount: cfg.amount,
    addressWindow: cfg.addressWindow, ipWindow: cfg.ipWindow, ipLimit: cfg.ipLimit,
  });
  for (const signal of ['SIGTERM', 'SIGINT']) process.once(signal, () => {
    server.close(() => journal.close().catch(() => { process.exitCode = 1; }));
  });
  server.listen(cfg.port, cfg.host, () => console.log(`Xitcoin faucet listening on ${cfg.host}:${cfg.port}`));
}
start().catch(() => { console.error('faucet journal requires operator review'); process.exitCode = 1; });
