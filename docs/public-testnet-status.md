# Xitcoin public testnet — operational status

## Current public network

- Cosmos chain ID: `xitcoin-testnet-2026-1`
- EVM chain ID: `101089`
- Current public release: `v0.1.0-testnet.1`
- Cosmos RPC: `https://rpc-testnet.xitcoin.org`
- Cosmos REST API: `https://api-testnet.xitcoin.org`
- EVM JSON-RPC: `https://evm-rpc-testnet.xitcoin.org`
- Cosmos explorer: `https://explorer-testnet.xitcoin.org`

The testnet currently has four active public validators.

## Verified release metadata

- Application name: `Xitcoin`
- Server binary: `xitcoind`
- Cosmos SDK: `v0.54.3`
- Release source commit: `46111748fc5002dc20f6a1b6ab57622cb0cc0e71`

## Public functionality

- Cosmos RPC: active
- Cosmos REST API: active
- EVM JSON-RPC: active
- IBC API: active
- Cosmos explorer: active
- gRPC: active at `grpc-testnet.xitcoin.org:443`
- EVM explorer: pending Blockscout deployment
- Faucet: pending real faucet service
- CosmWasm explorer: not available from the current chain API
- State sync: not enabled yet

## Important testnet policy

Testnet XTC is not a market asset. No USD price, swap value, circulating-market value, or fabricated economic data must be displayed.

## Remaining testnet work

1. Build and publish a real rate-limited faucet.
2. Deploy Blockscout for the EVM explorer.
3. Finish the Cosmos explorer branding and keep only real modules/data.
4. Add monitoring, alerting and operational documentation.
5. Add independent validators on separate infrastructure.
6. Add independent validators on separate infrastructure.
7. Run wallet, transaction, staking, governance and EVM acceptance tests.
8. Audit the chain and define final mainnet economics before any mainnet genesis.

## Faucet public — 11 August 2026

<!-- FAUCET_PUBLIC_2026_08_11 -->

The Xitcoin public testnet faucet is active.

- Public page: https://faucet-testnet.xitcoin.org/
- Health endpoint: https://faucet-testnet.xitcoin.org/healthz
- Claim endpoint: `POST https://faucet-testnet.xitcoin.org/claim`
- Distribution: 100 XTC testnet per successful request
- Limits: one request per address per 24 hours; three requests per IP per 24 hours
- Faucet balance after the on-chain test: approximately 9,900 XTC testnet
- Test transaction: `97772B9FA283A4B982D8F59790DAAE4220306F417FF95100FE4CB52BD49F43CD`

Testnet tokens have no monetary value. No USD price or market-cap claim is displayed.

### Next operational work

1. Monitor faucet balance, failed requests, rate-limit events and certificate renewal.
2. Finish the Cosmos explorer presentation using only live chain data.
3. Deploy Blockscout for the EVM explorer at `evm-explorer-testnet.xitcoin.org`.
4. Test wallet connection, transfers, staking, governance and EVM transactions end to end.
5. Add independent validator infrastructure before mainnet planning.

## Explorer, faucet and monitoring — 2026-08-11

Marker: PUBLIC_MONITORING_EXPLORER_2026_08_11

- Cosmos explorer: `https://explorer-testnet.xitcoin.org`
- Faucet: `https://faucet-testnet.xitcoin.org`
- Public gRPC: `grpc-testnet.xitcoin.org:443`
- Explorer is configured only for Xitcoin Testnet. No third-party blockchain is administered by this instance.
- Explorer pages use live chain data: dashboard, governance, staking, blocks, transactions, uptime, IBC, supply, parameters, consensus and faucet.
- Transaction history uses the Cosmos SDK query form required by the API.
- XTC CoinGecko data is enabled. Market listings are displayed only from CoinGecko; no market is inserted manually.
- Public healthcheck runs every minute and verifies the four public validators, RPC, API, explorer, faucet, gRPC and TLS certificate validity.
- EVM explorer remains intentionally unavailable pending a separate Blockscout deployment.

### Remaining priorities

1. Perform and document the full public user-path test.
2. Maintain faucet funding and review claim limits.
3. Deploy Blockscout for the EVM explorer.
4. Prepare independent validator infrastructure before mainnet.

## Public end-to-end verification and faucet IP protection — 2026-08-11

Marker: PUBLIC_E2E_SECURITY_2026_08_11

- Faucet requests now use the original visitor IP supplied by Cloudflare only; client-provided forwarding headers are not trusted.
- A public faucet request was tested through `https://faucet-testnet.xitcoin.org/claim`.
- The faucet transaction was confirmed on-chain with code `0`.
- A testnet delegation of 10 XTC to a bonded public validator was confirmed on-chain with code `0`.
- The resulting transaction routes are available from the Cosmos explorer.
- These checks used only XTC testnet. No mainnet asset or private validator was involved.

### Remaining priorities

1. Manual wallet-extension connection check in the public explorer.
2. Deploy Blockscout for the EVM explorer.
3. Maintain faucet reserve and monitor operational alerts.
4. Prepare independent validator infrastructure before mainnet.
