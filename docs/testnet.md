# Xitcoin Testnet

Xitcoin Testnet provides public Cosmos and EVM interfaces for application, wallet, validator and protocol testing.

## Network identity

| Field | Value |
| --- | --- |
| Cosmos Chain ID | `xitcoin-testnet-1` after the coordinated reset |
| EVM Chain ID | `101089` (`0x18ae1`) |
| Native asset | XTC |
| Atomic denomination | `axtc` |
| Decimals | 18 |

## Public services

| Service | Endpoint |
| --- | --- |
| Cosmos RPC | `https://rpc-testnet.xitcoin.org` |
| Cosmos REST API | `https://api-testnet.xitcoin.org` |
| EVM JSON-RPC | `https://evm-rpc-testnet.xitcoin.org` |
| Cosmos explorer | `https://explorer-testnet.xitcoin.org` |
| EVM explorer | `https://evm-explorer-testnet.xitcoin.org` |
| Faucet | `https://faucet-testnet.xitcoin.org` |

## Genesis candidate

The candidate for the coordinated reset is published at [`networks/testnet/genesis.json`](../networks/testnet/genesis.json). Verify it before use:

```bash
cd networks/testnet
sha256sum -c genesis.sha256
```

Expected SHA-256:

```text
7d13d7ed6a19ea48e2ce3c408f1f457e0961e72df6dd480d8200a6db5bae8414
```

The candidate becomes authoritative only after the coordinated reset is announced and completed. Until then, public endpoints continue to report the active network state.

See [Testnet Operations](testnet-operations.md) for verification commands.
