# Xitcoin Testnet Status

## Network naming

The public network name is **Xitcoin Testnet**.

| Identifier | Value |
| --- | --- |
| Cosmos Chain ID | `xitcoin-testnet-1` |
| EVM Chain ID | `101089` (`0x18ae1`) |

The Cosmos Chain ID identifies the current genesis. The `-1` suffix is not a
public version label and does not represent a validator or server number.

The four initial validators are Xitcoin Atlas, Xitcoin Borealis, Xitcoin
Meridian and Xitcoin Zenith.

## Current release state

The canonical four-validator network is active in staging. Its binary, genesis,
validator identities, consensus and local RPC interfaces have been validated.

The public endpoint cutover has not yet been completed. Existing public domains
continue serving the historical testnet until the coordinated canonical
cutover.

## Public services

| Component | Current status | Endpoint |
| --- | --- | --- |
| Cosmos RPC | Historical network; canonical cutover pending | `https://rpc-testnet.xitcoin.org` |
| Cosmos REST API | Historical network; canonical cutover pending | `https://api-testnet.xitcoin.org` |
| EVM JSON-RPC | Historical network; canonical cutover pending | `https://evm-rpc-testnet.xitcoin.org` |
| Cosmos explorer | Historical network; canonical cutover pending | `https://explorer-testnet.xitcoin.org` |
| EVM explorer | Historical network; canonical cutover pending | `https://evm-explorer-testnet.xitcoin.org` |
| Faucet | Historical network; canonical cutover pending | `https://faucet-testnet.xitcoin.org` |
| Bridge route | Not configured; disabled | — |

## Canonical configuration

| Parameter | Canonical value |
| --- | --- |
| Public name | Xitcoin Testnet |
| Cosmos Chain ID | `xitcoin-testnet-1` |
| EVM Chain ID | `101089` (`0x18ae1`) |
| Native asset | XTC |
| Atomic denomination | `axtc` |
| Decimals | 18 |
| Genesis SHA-256 | `7d13d7ed6a19ea48e2ce3c408f1f457e0961e72df6dd480d8200a6db5bae8414` |
| Initial validators | Atlas, Borealis, Meridian and Zenith |
| Validator capacity | 258 |
| Minimum self-delegation | 5,000,000 XTC |
| Administrative authority | `xtc1e3q4pm23ky0qetnep33j4yezq6c3lc7fcds4je` |

## Remaining public-release gates

- complete Cosmos and EVM transaction validation;
- validate admission, revocation and administrative signing;
- complete the public endpoint cutover;
- reconcile explorers and faucet with the canonical genesis;
- complete monitoring and rollback acceptance;
- validate the Cronos bridge independently before route activation.

Testnet XTC has no monetary value.
