# Xitcoin Testnet Status

## Current services

| Component | Status | Endpoint |
| --- | --- | --- |
| Cosmos RPC | Active | `https://rpc-testnet.xitcoin.org` |
| Cosmos REST API | Active | `https://api-testnet.xitcoin.org` |
| EVM JSON-RPC | Active | `https://evm-rpc-testnet.xitcoin.org` |
| Cosmos explorer | Active | `https://explorer-testnet.xitcoin.org` |
| EVM explorer | Active | `https://evm-explorer-testnet.xitcoin.org` |
| Faucet | Active | `https://faucet-testnet.xitcoin.org` |
| Bridge route | Disabled | — |

## Network transition

The running testnet and the coordinated-reset candidate are distinct states.

| State | Cosmos Chain ID | Genesis SHA-256 |
| --- | --- | --- |
| Current network | `xitcoin-testnet-2026-1` | `626f034f92f30cc6016b0175ea7b84e3f4b4b79543aea76352eeffc050d69a04` |
| Reset candidate | `xitcoin-testnet` | `818096564a458b68ddd56ac95592ec7bac64c88f6fcbb9742cd39114229884b0` |

The candidate includes the `axtc` atomic denomination, zero inflation, the declared supply ceiling and the native XTC EVM representation. It is not active until the coordinated reset is completed.

## Verification

Use the commands in [Testnet Operations](testnet-operations.md) to query the live network and verify the candidate genesis. Testnet XTC has no monetary value.
