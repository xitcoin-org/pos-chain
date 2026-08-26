# Xitcoin Testnet Status

## Network identity

| Property | Canonical value |
| --- | --- |
| Public name | Xitcoin Testnet |
| Cosmos Chain ID | `xitcoin-testnet-1` |
| EVM Chain ID | `101089` (`0x18ae1`) |
| Native asset | XTC |
| Atomic denomination | `axtc` |
| Decimals | 18 |
| Genesis time | `2026-08-25T21:48:17.77229Z` |
| Genesis SHA-256 | `55c8756a212b9e92c0e8427ea61caff7fa9dca40e801e4b848f59d1aa5f6dae6` |
| Genesis supply | 457,000,000 XTC |
| Initial validators | Atlas, Borealis, Meridian and Zenith |
| Validator capacity | 258 |
| Minimum self-delegation | 5,000,000 XTC |
| Module authority in deployed genesis | `xtc1vza8zsgvrfwmve084ytd8xqdkkm7u9e5csctc2` |

The Cosmos Chain ID identifies the current genesis. The `-1` suffix is not a
public version label and does not represent a validator or server number.

## Certified deployment state

The canonical four-validator testnet is publicly active. On 2026-08-26 the
sentry, validators, public endpoints, explorers, Blockscout, faucet and
healthcheck were verified against the same live chain.

| Component | Status | Endpoint |
| --- | --- | --- |
| Cosmos RPC | Active | `https://rpc-testnet.xitcoin.org` |
| Cosmos REST API | Active | `https://api-testnet.xitcoin.org` |
| EVM JSON-RPC | Active | `https://evm-rpc-testnet.xitcoin.org` |
| Cosmos explorer | Active | `https://explorer-testnet.xitcoin.org` |
| EVM explorer | Active; normal asynchronous indexing lag | `https://evm-explorer-testnet.xitcoin.org` |
| Faucet | Active; 10 XTC per accepted request | `https://faucet-testnet.xitcoin.org` |
| Bridge route | Not configured; disabled | — |

Blockscout may trail the chain head by a small number of blocks while indexing.
A moving one-to-two-block lag is normal; a persistent or increasing lag is not.

## Genesis allocations

| Allocation | Amount |
| --- | ---: |
| Sovereign reserve | 386,000,000 XTC |
| Faucet | 50,000,000 XTC |
| Security reserve | 960,000 XTC |
| Validator self-delegation (four total) | 20,000,000 XTC |
| Validator liquid gas (four total) | 40,000 XTC |
| **Total** | **457,000,000 XTC** |

The faucet uses its allocation and does not mint automatically. Its deployed
limits are a 24-hour address window, a 24-hour IP window and at most three
accepted requests per IP during that window.

No separate active bridge allocation exists in this genesis. A future bridge
route requires an explicit configuration and approval process.

## Acceptance record

- four validators active, enabled, not catching up, with equal voting power;
- sentry active, enabled, not catching up, and connected to the four validators;
- public RPC, API, EVM RPC, explorers and faucet reachable;
- Blockscout reindexed from the deployed genesis and tracks the live chain;
- healthcheck timer repeatedly completes successfully;
- obsolete Blockscout containers and volumes removed;
- one verified rollback point retained for Blockscout, platform endpoints and
  sentry application configuration;
- legacy testnet service units inactive and disabled;
- bridge route disabled.

The previously documented source revision, binary checksum and Actions run were
removed because that historical record is no longer resolvable from the current
public repository history. This document does not claim a source-to-binary
provenance link that cannot be independently verified. A future binary release
must publish a reachable source revision and reproducible checksums.

No production key, mnemonic, password, node key, validator private key or private
backup is stored in this repository. Testnet and mainnet must use independent
keys, genesis files, state and operational directories.

Testnet XTC has no monetary value.
