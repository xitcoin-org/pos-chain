# Xitcoin Testnet — Project Record

## Current status
- Stage: private testnet preparation
- Network: not launched
- Base technology: Cosmos EVM
- Reference source commit: 5d4a9c786b083ff35ff5dcd314967cdf8dda2280
- Reference binary: /opt/xitcoin-testnet/bin/evmd-reference
- Reference binary SHA-256: 00cbf5b31f142edd594abc4b3c3731548a3b8b0b8a5c65adcb68163a01656342

## Intended network identity
- Chain name: Xitcoin
- Native display symbol: XTC
- Native base denomination: xits
- Native display denomination: XTC
- Decimal precision: 18
- Cosmos address prefix: xitcoin
- EVM addresses: standard 0x addresses
- Testnet chain ID: xitcoin-testnet-1
- Temporary testnet EVM chain ID: 20260807, checked against chainid.network on 2026-08-07. Production requires a separate final allocation.

## Governance and validator policy
- Testnet validators are private and explicitly authorised.
- Production validator admission must be implemented as an allow-list control in code.
- A token holder must not obtain validator admission or governance control merely by buying XTC.
- Validator keys, seed phrases and production secrets never belong in this repository.

## Production gates
1. Four-node private testnet passes consensus, restart and recovery tests.
2. Genesis allocation and XTC migration plan are independently reviewed.
3. Validator allow-list is implemented and tested in code.
4. RPC, explorer, monitoring, backups and incident procedures are tested.
5. Production configuration is created separately from testnet.

## Next technical step
Create the Xitcoin application configuration from the Cosmos EVM reference:
binary name, chain identity, denominations, address prefix and private testnet genesis.

## Monetary policy
- Hard cap: 5,250,000,000 XTC (5,250,000,000,000,000,000,000,000,000 xits).
- Inflation: 0%.
- No Foundation free mint and no automatic validator minting.
- Validator rewards must come from transaction fees and a separately identified reserve.
- Burns permanently reduce the circulating and total supply.
- Mint parameters cannot be changed through governance proposals.
