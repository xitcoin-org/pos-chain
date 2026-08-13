# Xitcoin Public Testnet Operations

## Scope

This document is the public verification reference for Xitcoin Testnet. It
contains read-only checks only. It does not provide validator credentials,
operator access, or deployment instructions.

## Network identity

| Field | Value |
| --- | --- |
| Cosmos chain ID | `xitcoin-testnet-2026-1` |
| EVM chain ID | `101089` (`0x18ae1`) |
| Display denomination | `XTC` |
| Base denomination | `xtc` |
| Decimals | `18` |

## Public endpoints

| Service | Endpoint |
| --- | --- |
| Cosmos RPC | `https://rpc-testnet.xitcoin.org` |
| Cosmos REST API | `https://api-testnet.xitcoin.org` |
| EVM JSON-RPC | `https://evm-rpc-testnet.xitcoin.org` |
| Cosmos explorer | `https://explorer-testnet.xitcoin.org` |
| EVM explorer | `https://evm-explorer-testnet.xitcoin.org` |
| Faucet | `https://faucet-testnet.xitcoin.org` |

## Public verification

The following commands are safe to run locally. They submit no transaction and
do not require a wallet or private key.

### Verify Cosmos synchronization

```bash
curl -fsS https://rpc-testnet.xitcoin.org/status |
  jq '.result.sync_info | {latest_block_height, latest_block_time, catching_up}'
```

Expected result: `catching_up` is `false`.

### Verify Cosmos network identity

```bash
curl -fsS https://api-testnet.xitcoin.org/cosmos/base/tendermint/v1beta1/node_info |
  jq '.default_node_info | {network, moniker, version}'
```

Expected network: `xitcoin-testnet-2026-1`.

### Verify EVM network identity

```bash
curl -fsS -H 'content-type: application/json' \
  --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
  https://evm-rpc-testnet.xitcoin.org | jq .
```

Expected result: `0x18ae1`.

### Verify the official genesis file

```bash
sha256sum networks/testnet-101089/genesis.json
```

Expected SHA-256:
`626f034f92f30cc6016b0175ea7b84e3f4b4b79543aea76352eeffc050d69a04`.

## Active release provenance

The active public release is `v0.1.0-testnet.1`. Operational verification on
2026-08-13 recorded the running binary as built from source commit
`46111748fc5002dc20f6a1b6ab57622cb0cc0e71`.

Source commits do not deploy automatically. Any future upgrade must publish a
reviewed release record with the source revision, binary checksum, genesis
checksum, test evidence, rollout plan, and rollback plan.

## Change-control boundaries

- The public testnet genesis is immutable after launch.
- Testnet XTC has no market value; no fabricated price or market-cap data is
  displayed.
- The bridge is disabled by default: no route, relayer, settlement, custody,
  minting, vault, or real asset is enabled.
- A bridge code change is not a production bridge deployment.

## Current operational priorities

1. Maintain monitoring, faucet reserve, and tested backup/restore procedures.
2. Perform recurring wallet, transaction, staking, governance, and EVM user-path acceptance checks.
3. Add independent validator and RPC infrastructure before mainnet planning.
4. Enable state sync only after snapshot and recovery procedures are tested.
5. Audit supply accounting and define mainnet economics before any mainnet genesis.
