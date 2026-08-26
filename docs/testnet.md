# Xitcoin Testnet

**Xitcoin Testnet** is the public network name.

The technical Cosmos Chain ID is `xitcoin-testnet-1`. The `-1` suffix
identifies this genesis generation; it is not a public version label or a
validator number.

## Canonical identity

| Field | Value |
| --- | --- |
| Public name | Xitcoin Testnet |
| Cosmos Chain ID | `xitcoin-testnet-1` |
| EVM Chain ID | `101089` (`0x18ae1`) |
| Native asset | XTC |
| Atomic denomination | `axtc` |
| Decimals | 18 |
| Genesis time | `2026-08-25T21:48:17.77229Z` |
| Genesis supply | 457,000,000 XTC |
| Genesis SHA-256 | `55c8756a212b9e92c0e8427ea61caff7fa9dca40e801e4b848f59d1aa5f6dae6` |

## Deployment state

The canonical four-validator testnet is publicly active. The sentry, all four
validators, public endpoints, explorers, Blockscout and faucet were certified
together on 2026-08-26.

## Public services

| Service | Endpoint |
| --- | --- |
| Cosmos RPC | `https://rpc-testnet.xitcoin.org` |
| Cosmos REST API | `https://api-testnet.xitcoin.org` |
| EVM JSON-RPC | `https://evm-rpc-testnet.xitcoin.org` |
| Cosmos explorer | `https://explorer-testnet.xitcoin.org` |
| EVM explorer | `https://evm-explorer-testnet.xitcoin.org` |
| Faucet | `https://faucet-testnet.xitcoin.org` |

The faucet sends 10 testnet XTC per accepted request. Address and IP windows are
24 hours, with at most three accepted requests per IP in that window. It spends
from a finite 50,000,000 XTC genesis allocation and does not mint automatically.

## Validator and allocation policy

- Four initial validators: Atlas, Borealis, Meridian and Zenith.
- 5,000,000 XTC self-delegated by each validator.
- 10,000 liquid XTC per validator account for operational gas.
- Validator capacity and admission cap: 258.
- Genesis allocations: 386,000,000 XTC sovereign reserve; 50,000,000 XTC
  faucet; 960,000 XTC security reserve; 20,000,000 XTC total validator
  self-delegation; 40,000 XTC total validator liquid gas.
- No separate bridge allocation or settlement route is active. The bridge
  remains disabled until explicitly configured and approved.

## Canonical genesis

The canonical genesis is published at
[`networks/xitcoin-testnet-1/genesis.json`](../networks/xitcoin-testnet-1/genesis.json).

Verify it with:

```bash
cd networks/xitcoin-testnet-1
sha256sum -c genesis.sha256
```

See [Testnet Status](testnet-status.md) for the current deployment boundary and
[Testnet Operations](testnet-operations.md) for verification commands.
