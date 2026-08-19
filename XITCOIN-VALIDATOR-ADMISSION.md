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

The public network name and the technical Cosmos chain ID serve different purposes. Version labels such as V1 or V2 are not public network names.

## Admission rule

Xitcoin is a permissioned-validator network.

Holding, receiving or staking XTC does not automatically grant the right to become a validator. A validator may be created only when its validator operator address is present in the on-chain approval list and all staking and operational requirements are satisfied.

## Candidate network parameters

- Maximum validator and admission capacity: 258
- Initial approved validators: 4
- Protocol minimum self-delegation: 1,000,000 XTC
- Initial self-delegation of each core validator: 5,000,000 XTC
- Initial core validators:
  - Xitcoin Atlas
  - Xitcoin Borealis
  - Xitcoin Meridian
  - Xitcoin Zenith

The protocol minimum is an admission floor. The larger initial core-validator self-delegation is a deployment value and does not change that floor.

## Participation framework

The 258-position planning model contains:

- 195 sovereign reference positions:
  - 193 United Nations Member States;
  - the Holy See;
  - the State of Palestine.
- 63 public positions.

The sovereign allocation index is a deterministic planning and reference record. It does not transfer funds, approve a validator, activate a validator or create an automatic entitlement.

The 39 territorial consolidations in the allocation dataset are statistical population mappings. They do not create additional validator positions.

## Authority and on-chain record

The approval list and its authority are blockchain state initialized in genesis. They are not read from a private server-side allowlist.

The admission authority can execute the defined on-chain approval and revocation actions. Production authority, multisignature custody and signer-change procedures must be documented and verified before mainnet.

## Required actions

- Approve: add a validator address to the on-chain approval list.
- Revoke: remove an address from the approval list and deactivate the associated validator according to protocol rules.
- Protect: a revoked validator must not be able to recreate or unjail itself without renewed approval.
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

The `xitcoin-testnet-1` release candidate is running as an isolated four-validator staging network. Consensus, peer connectivity, binary identity, genesis identity, Cosmos identity and local-only EVM JSON-RPC have been validated.

Cosmos and EVM transaction validation remains a release gate. The currently published testnet remains unchanged until coordinated activation.
