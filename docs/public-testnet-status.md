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
