# Xitcoin Testnet

**Xitcoin Testnet** is the public network name.

The technical Cosmos Chain ID of the current genesis is
`xitcoin-testnet-1`. The `-1` suffix identifies the genesis generation. It is
not a public version label and does not represent a validator or server number.

## Canonical identity

| Field | Value |
| --- | --- |
| Public name | Xitcoin Testnet |
| Cosmos Chain ID | `xitcoin-testnet-1` |
| EVM Chain ID | `101089` (`0x18ae1`) |
| Native asset | XTC |
| Atomic denomination | `axtc` |
| Decimals | 18 |
| Genesis SHA-256 | `7d13d7ed6a19ea48e2ce3c408f1f457e0961e72df6dd480d8200a6db5bae8414` |

## Deployment state

The canonical network is active as a four-validator staging network. Public
endpoint cutover remains pending.

Existing public domains continue serving the historical testnet and must not
yet be used as evidence that the canonical genesis is publicly active.

## Public endpoint transition

| Service | Current endpoint |
| --- | --- |
| Cosmos RPC | `https://rpc-testnet.xitcoin.org` |
| Cosmos REST API | `https://api-testnet.xitcoin.org` |
| EVM JSON-RPC | `https://evm-rpc-testnet.xitcoin.org` |
| Cosmos explorer | `https://explorer-testnet.xitcoin.org` |
| EVM explorer | `https://evm-explorer-testnet.xitcoin.org` |
| Faucet | `https://faucet-testnet.xitcoin.org` |

These endpoints remain associated with the historical public environment until
the coordinated canonical cutover is completed.

## Canonical genesis

The canonical genesis is published at
[`networks/testnet/genesis.json`](../networks/testnet/genesis.json).

Expected SHA-256:

`7d13d7ed6a19ea48e2ce3c408f1f457e0961e72df6dd480d8200a6db5bae8414`

See [Testnet Status](testnet-status.md) for the current deployment boundary and
[Testnet Operations](testnet-operations.md) for verification commands.
