# Xitcoin Testnet Status

## Network identity

| Property | Canonical value |
| --- | --- |
| Public name | Xitcoin Testnet |
| Cosmos Chain ID | `xitcoin-testnet-v2-1` |
| EVM Chain ID | `101089` (`0x18ae1`) |
| Native asset | XTC |
| Atomic denomination | `axtc` |
| Decimals | 18 |
| Genesis time | `2026-08-29T09:49:07Z` |
| Genesis SHA-256 | `5db34acf6496b2c76a6f516e0eb605caef6762552584ddbed7c8703239f33d72` |
| Genesis supply | 477,000,000 XTC |
| Initial validators | Atlas, Borealis, Meridian and Zenith |
| Validator capacity | 258 |
| Minimum self-delegation | 5,000,000 XTC |
| Module authority in deployed genesis | `xtc1vza8zsgvrfwmve084ytd8xqdkkm7u9e5csctc2` |

`xitcoin-testnet-v2-1` is the only active public testnet generation. `xitcoin-testnet-1` is retired and retained only as historical evidence.

## Verified deployment state

The four-validator V2 testnet is publicly active. On 2026-09-02, the sentry, validators, public endpoints, explorers, faucet, Cosmos transactions and EVM transactions were verified against the same live chain.

| Component | Status | Endpoint |
| --- | --- | --- |
| Cosmos RPC | Active | `https://rpc-testnet.xitcoin.org` |
| Cosmos REST API | Active | `https://api-testnet.xitcoin.org` |
| EVM JSON-RPC | Active | `https://evm-rpc-testnet.xitcoin.org` |
| Cosmos explorer | Active | `https://explorer-testnet.xitcoin.org` |
| EVM explorer | Active; indexed from height 1 | `https://evm-explorer-testnet.xitcoin.org` |
| Faucet | Active; 10 XTC per accepted request | `https://faucet-testnet.xitcoin.org` |
| Bridge route | Not configured; disabled | — |

## Genesis bank allocations

| Allocation | Amount |
| --- | ---: |
| Sovereign reserve | 386,000,000 XTC |
| Faucet | 50,000,000 XTC |
| Validator-incentives module | 20,000,000 XTC |
| Four validator accounts before genesis transactions | 21,000,000 XTC |
| **Total** | **477,000,000 XTC** |

The four genesis transactions self-delegate 20,000,000 XTC from the validator accounts. The remaining validator liquid balance is 1,000,000 XTC total, or 250,000 XTC per validator.

The faucet uses its finite allocation and does not mint automatically. Its deployed limits are a 24-hour address window, a 24-hour IP window and at most three accepted requests per IP during that window.

No bridge route configuration exists in this genesis. Activating a future route requires explicit configuration, authorization and a separate acceptance process.

## Acceptance record

- exact canonical genesis published with its SHA-256 checksum;
- four validators active, enabled, not catching up, with equal initial voting power;
- public RPC, API, EVM RPC, explorers and faucet reachable;
- Blockscout indexing complete from the first EVM height;
- native Cosmos and EVM transactions confirmed and visible in their respective explorers;
- bridge route absent;
- mainnet not launched and no Cronos transaction performed as part of this verification.

The active binary identifies source revision `e06e232b95b0c40e3c718da3b8f447eed0588972`. A production release still requires its own reproducible build, published checksums and independent acceptance record.

No production key, mnemonic, password, node key, validator private key or private backup is stored in this repository. Testnet and mainnet must use independent keys, genesis files, state and operational directories.

Testnet XTC has no monetary value.
