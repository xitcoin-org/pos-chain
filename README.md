# Xitcoin

Xitcoin is an EVM-compatible Proof-of-Stake blockchain built with the Cosmos SDK, CometBFT and Cosmos EVM.

## Status

Xitcoin is under active engineering and private-testnet validation.

- Public Testnet: planned, EVM Chain ID `101089`
- Mainnet: not launched
- Mainnet Chain ID `101088`: reserved for the future Xitcoin public network
- Native asset: `XTC`
- Technical base denomination: `xtc` (18 decimals)
- Inflation policy: 0%
- Maximum supply policy: 5,250,000,000 XTC
- Native EVM representation: XTC; no wrapped-XTC asset

No public RPC, faucet, explorer, bridge or production service is announced by this repository yet.

## Scope

This repository contains the Xitcoin chain source code and its public technical documentation.

It does not contain:

- validator keys or node keys;
- mnemonics, passwords, backup archives or private infrastructure;
- private testnet topology or server addresses;
- a bridge or migration system for the existing external Xitcoin asset on Cronos.

## Public Testnet

The future public testnet will use Chain ID `101089`. It will have a new genesis, separate validator keys and dedicated public infrastructure. The current private testnet is deliberately isolated and will never be exposed as the public testnet.

See [Public Testnet 101089](docs/public-testnet-101089.md).

## Development

The chain is derived from the Cosmos EVM stack. Upstream attribution, licensing and maintenance references are listed in [UPSTREAMS.md](UPSTREAMS.md).

Before contributing, read [CONTRIBUTING.md](CONTRIBUTING.md) and [SECURITY.md](SECURITY.md).

## License

The source code is distributed under the [Apache License 2.0](LICENSE).
