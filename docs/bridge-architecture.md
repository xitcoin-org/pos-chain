# Xitcoin Bridge Architecture

## Overview

The canonical route connects the Cronos XTC contract with native XTC on Xitcoin:

```text
Cronos XTC <-> Xitcoin native XTC
```

Cronos XTC contract:

```text
0xE45FE733BC8617FA6DAC8437FC44B5FFFA949991
```

The bridge route is disabled until its contracts, signer configuration, limits and monitoring are approved for the target environment.

## Transfer model

### Cronos to Xitcoin

1. Canonical XTC is locked in the Cronos vault.
2. The deposit reaches the required confirmation depth.
3. The signer quorum attests to the source event and destination.
4. Replay protection validates the transfer identifier.
5. The Xitcoin bridge module mints the corresponding native XTC.

### Xitcoin to Cronos

1. Native XTC is submitted to the Xitcoin bridge module.
2. The module burns the submitted amount.
3. The burn reaches finality and is attested by the signer quorum.
4. The Cronos vault releases the corresponding locked XTC.

## Accounting

For each route:

```text
outstanding native XTC = locked canonical XTC - released canonical XTC
```

Every mint requires finalized collateral. Every release requires a finalized native burn. Route-level limits may be lower than the chain supply ceiling.

## Message domain

A signed bridge action binds the source chain, destination chain, route identifier, transaction hash, log index, sender, recipient, amount, nonce, action, signer-set version and validity window.

Processed source events are immutable and cannot be executed twice.

## Control model

- A threshold signer set authorizes settlements.
- A guardian may pause settlement.
- Resuming, changing signers and changing limits require threshold approval.
- Per-transfer, daily and outstanding-mint limits are enforced.
- Accounting reconciliation runs independently of settlement.
- A reconciliation mismatch pauses the route.

## Burn semantics

A token burn invokes the token's burn mechanism and reduces `totalSupply`. A transfer to a dead address changes the holder balance but does not, by itself, prove a reduction in `totalSupply`.

Administrative Cronos supply reduction is separate from bridge settlement and must follow the verified contract ownership and governance controls.

## Read-only Cronos verification

```bash
export CRONOS_RPC='https://evm.cronos.org'
export CRONOS_XTC='0xE45FE733BC8617FA6DAC8437FC44B5FFFA949991'

cast code "$CRONOS_XTC" --rpc-url "$CRONOS_RPC"
cast call "$CRONOS_XTC" 'totalSupply()(uint256)' --rpc-url "$CRONOS_RPC"
cast call "$CRONOS_XTC" 'owner()(address)' --rpc-url "$CRONOS_RPC"
```

Contract implementation, ABI and ownership must be verified against the canonical explorer before an administrative action.

## Deployment stages

1. Unit and integration testing.
2. Isolated testnet deployment with test assets.
3. Independent security review.
4. Limited production rollout.
5. Continuous monitoring and reconciliation.

Additional networks require route-specific adapters, confirmation policies and threat models.
