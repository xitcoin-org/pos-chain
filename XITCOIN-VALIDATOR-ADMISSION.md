# Xitcoin — Validator Admission Policy

## Public identity

- Network name: Xitcoin Testnet
- Cosmos chain ID: `xitcoin-testnet-1`
- EVM chain ID: `101089` (`0x18ae1`)
- Native asset: XTC
- Atomic denomination: `axtc`
- Decimal precision: 18
- Cosmos addresses: `xtc1...`
- Validator addresses: `xtcvaloper1...`
- EVM addresses: `0x...`
- Node program: `xitcoind`

The public name is `Xitcoin Testnet`. `xitcoin-testnet-1` is its technical Cosmos chain ID, not a public version or server name.

## Admission rule

Xitcoin is a permissioned-validator network.

Holding, receiving, staking or delegating XTC does not grant validator-admission authority. A validator may be created only when its operator address has been explicitly approved by the canonical on-chain authority and all staking, security and operational requirements are satisfied.

## Canonical testnet parameters

- Maximum validator and admission capacity: 258
- Initially approved validators: 4
- Additional validators currently announced: 0
- Protocol minimum self-delegation: 5,000,000 XTC
- Initial self-delegation of each core validator: 5,000,000 XTC
- Initial core validators:
  - Xitcoin Atlas
  - Xitcoin Borealis
  - Xitcoin Meridian
  - Xitcoin Zenith

The maximum is capacity, not a target count. Any expansion beyond the four initial validators requires separate review and authorization.

## Participation framework

The planning model contains:

- 193 reserved Member-State positions;
- 65 general validator positions.

These positions are future capacity only. They do not transfer funds, approve operators, activate validators or create automatic entitlements.

The 39 territorial consolidations in the allocation dataset are statistical population mappings and do not create additional positions.

## Authority and governance boundary

The approval list and its authority are blockchain state initialized in genesis. They are not read from a private server-side allowlist.

Only the canonical validator-admission authority can execute the defined approval and revocation actions during the current launch phase. The authority is controlled through the project's authorized custody process.

Token balances, staking balances, delegation weight and ordinary governance voting do not approve or revoke validators and do not override the admission authority. The technical presence of Cosmos governance infrastructure does not change this boundary.

Production custody, recovery and signer-change procedures must be documented and verified before mainnet.

## Required actions

- Approve: add a validator address to the on-chain approval list.
- Revoke: remove an address and deactivate the associated validator according to protocol rules.
- Protect: a revoked validator must not recreate or unjail itself without renewed approval.
- Audit: approvals, revocations and parameter updates must remain visible in blockchain state and transaction history.

## Mandatory validation

1. Reject an unapproved `CreateValidator` transaction.
2. Accept an approved `CreateValidator` transaction.
3. Confirm that revocation deactivates the affected validator.
4. Confirm that a revoked validator cannot recreate or unjail itself.
5. Confirm identical state transitions across nodes using the same genesis.
6. Confirm maximum-capacity and minimum-self-delegation enforcement.
7. Retain transaction hashes and release-specific evidence.

## Current release status

Xitcoin Testnet is running as an isolated four-validator staging network. Its technical Cosmos chain ID is `xitcoin-testnet-1`. Consensus, peer connectivity, binary identity, genesis identity, Cosmos identity and local-only EVM JSON-RPC have been validated.

Cosmos and EVM transaction validation remains a public-release gate. Existing public endpoints continue serving the historical testnet until the coordinated canonical cutover.
