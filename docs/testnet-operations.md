# Xitcoin Testnet Operations

## Network verification

Check CometBFT synchronization:

```bash
curl -fsS https://rpc-testnet.xitcoin.org/status |
  jq '.result.sync_info | {latest_block_height, latest_block_time, catching_up}'
```

Check the active Cosmos Chain ID:

```bash
curl -fsS https://rpc-testnet.xitcoin.org/status |
  jq -r '.result.node_info.network'
```

Check the EVM Chain ID:

```bash
curl -fsS   -H 'content-type: application/json'   --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}'   https://evm-rpc-testnet.xitcoin.org |
  jq -r '.result'
```

Expected EVM result: `0x18ae1`.

## Genesis verification

```bash
git clone https://github.com/xitcoin-org/pos-chain.git
cd pos-chain/networks/testnet
sha256sum -c genesis.sha256
```

## Native asset

```text
symbol: XTC
atomic denomination: axtc
decimals: 18
```

## Release verification

A release record must identify the source revision, binary checksum, genesis checksum and compatibility notes. Source changes do not activate a genesis or network upgrade automatically.

## Operational principles

- Publish release artifacts with reproducible checksums.
- Validate RPC identity before signing or broadcasting a transaction.
- Apply rate limiting and monitoring to public endpoints.
- Test backup, restoration and rollback procedures before upgrades.
- Coordinate genesis resets and consensus upgrades with every validator.
