# Xitcoin Testnet

Xitcoin Testnet provides public Cosmos and EVM interfaces for application, wallet, validator and protocol testing.

## Network identity

| Field | Value |
| --- | --- |
| Cosmos Chain ID | `xitcoin-testnet` after the coordinated reset |
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
818096564a458b68ddd56ac95592ec7bac64c88f6fcbb9742cd39114229884b0
```

The candidate becomes authoritative only after the coordinated reset is announced and completed. Until then, public endpoints continue to report the active network state.

See [Testnet Operations](testnet-operations.md) for verification commands.
