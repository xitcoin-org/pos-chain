# Xitcoin Public Testnet

## Status

Planned. No public RPC, explorer, faucet, validator endpoint or genesis file is announced yet.

## Network identity

| Field | Value |
| --- | --- |
| Network name | Xitcoin Testnet |
| Cosmos Chain ID | `xitcoin-testnet-2026-1` |
| EVM Chain ID | `101089` (`0x18ae1`) |
| Native asset | XTC |
| Base denomination | `xits` |
| Decimals | 18 |

The Cosmos Chain ID is intentionally different from the private validation network. The public testnet will use a fresh genesis and independent validator keys.

## Launch prerequisites

- Independent public-testnet infrastructure and validator operators
- Reproducible genesis ceremony and published SHA-256 hash
- Public RPC with rate limiting and monitoring
- Explorer and faucet
- Security review of supply accounting and future external-chain migration design
- Public validator, RPC and incident-response documentation

## Endpoint policy

Endpoints will be published only after launch validation. Do not use unofficial RPC URLs or send assets to addresses claiming to be Xitcoin Testnet before the official launch announcement.

The machine-readable manifest is available at [`networks/testnet-101089/chain.json`](../networks/testnet-101089/chain.json).
