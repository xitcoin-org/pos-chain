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

The canonical four-validator testnet is publicly active. Its binary, genesis,
validator identities, consensus, transactions and public interfaces have been
validated.

The coordinated public endpoint cutover was completed on 2026-08-21. The
canonical release was deployed across the sentry and all four validators on
2026-08-22 using a rolling, rollback-protected procedure.

| Release property | Canonical value |
| --- | --- |
| Source revision | `2aa39b8c2ce7ac06278d58f2970225fd450e2c2c` |
| Linux AMD64 binary SHA-256 | `1958ca411353e79d3ff62a262960378141afb0343ab22689a759dffc610a9ecc` |
| Genesis SHA-256 | `7d13d7ed6a19ea48e2ce3c408f1f457e0961e72df6dd480d8200a6db5bae8414` |
| Release build record | [GitHub Actions run 32563535992](https://github.com/xitcoin-org/pos-chain/actions/runs/32563535992) |

The release artifact is retained temporarily by GitHub Actions as a build and
checksum record. It is not a node key, wallet secret, genesis authority or
active network dependency.

## Public services

| Component | Current status | Endpoint |
| --- | --- | --- |
| Cosmos RPC | Canonical network active | `https://rpc-testnet.xitcoin.org` |
| Cosmos REST API | Canonical network active | `https://api-testnet.xitcoin.org` |
| EVM JSON-RPC | Canonical network active | `https://evm-rpc-testnet.xitcoin.org` |
| Cosmos explorer | Canonical network active | `https://explorer-testnet.xitcoin.org` |
| EVM explorer | Canonical network active | `https://evm-explorer-testnet.xitcoin.org` |
| Faucet | Canonical network active | `https://faucet-testnet.xitcoin.org` |
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

## Completed public-testnet acceptance

- Cosmos and EVM transactions validated;
- admission, revocation and administrative signing validated;
- public RPC, API, explorer and faucet cutover completed;
- monitoring and rollback acceptance completed;
- the same verified release binary deployed on the sentry and four validators;
- validator identities, voting power, genesis and blockchain state preserved;
- obsolete deployment directories, temporary artifacts and disabled legacy
  service units removed from every testnet node;
- one verified rollback binary retained per node;
- system journals compressed, rotated and bounded to a 30-day retention policy;
- Cronos bridge remains inactive and separately gated.

No production key, mnemonic, password, node key, validator private key or private
backup is stored in this repository. Testnet and mainnet must use independent
keys, genesis files, state and operational directories.

Testnet XTC has no monetary value.
