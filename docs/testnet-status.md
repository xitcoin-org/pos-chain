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

## Canonical network configuration

The public testnet release is defined by one canonical network identity. Deployment and activation are coordinated operational procedures and do not alter these identifiers.

| Parameter | Canonical value |
| --- | --- |
| Cosmos Chain ID | `xitcoin-testnet` |
| EVM Chain ID | `101089` (`0x18ae1`) |
| Genesis SHA-256 | `a0f835a5e12101382475ff50ba31f947ab703a072258d5ad2a0e62a8d9faf997` |
| Bridge route | Disabled until operational authorization |

The canonical genesis defines the `axtc` atomic denomination, zero inflation, the declared supply ceiling and the native XTC EVM representation. Operators must verify the complete genesis SHA-256 before installation or activation.

## Activation status

The canonical release becomes active only through the coordinated testnet release procedure. Publication of the configuration does not authorize validator replacement, bridge activation, contract deployment or asset transfers.

## Verification

Use the commands in [Testnet Operations](testnet-operations.md) to verify the active network and the canonical genesis independently. Testnet XTC has no monetary value.
