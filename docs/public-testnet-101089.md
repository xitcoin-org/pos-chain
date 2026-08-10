# Xitcoin Public Testnet

## Status

Genesis released. P2P bootstrap endpoints are active. Public RPC, explorer and faucet are not announced yet.

## Genesis

The official genesis file is [`genesis.json`](../networks/testnet-101089/genesis.json). Its SHA-256 hash is [`626f034f92f30cc6016b0175ea7b84e3f4b4b79543aea76352eeffc050d69a04`](../networks/testnet-101089/genesis.sha256).

## Network identity

| Field | Value |
| --- | --- |
| Network name | Xitcoin Testnet |
| Cosmos Chain ID | `xitcoin-testnet-2026-1` |
| EVM Chain ID | `101089` (`0x18ae1`) |
| Native asset | XTC |
| Base denomination | `xtc` |
| Decimals | 18 |

The Cosmos Chain ID is intentionally different from the private validation network. The public testnet will use a fresh genesis and independent validator keys.

## Launch prerequisites

- Independent public-testnet infrastructure and validator operators
- Reproducible genesis ceremony and published SHA-256 hash
- Public RPC with rate limiting and monitoring
- Explorer and faucet
- Security review of supply accounting and future external-chain migration design
- Public validator, RPC and incident-response documentation

## P2P bootstrap

The initial Xitcoin Testnet validators are reachable at:

- `c135bf79b66db802f93c46170a48a166d24c5167@51.68.54.120:27656`
- `a494fded0411327150b57dabb1145d44c72cbd2a@51.68.54.120:27666`
- `e9bf23858046dfa6a0f74f312765cc3ca43f697b@51.68.54.120:27676`
- `f26c4c36a4a6049c4f3d709be3ccccacd28057bc@51.68.54.120:27686`

These are P2P endpoints only. No public RPC endpoint is announced yet.

## Endpoint policy

Endpoints will be published only after launch validation. Do not use unofficial RPC URLs or send assets to addresses claiming to be Xitcoin Testnet before the official launch announcement.

The machine-readable manifest is available at [`networks/testnet-101089/chain.json`](../networks/testnet-101089/chain.json).
