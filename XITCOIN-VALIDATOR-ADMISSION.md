# Xitcoin — Validator Admission Policy

## Public identity
- Chain: Xitcoin
- Public token ticker: $XTC
- Cosmos addresses: xitcoin1...
- EVM addresses: 0x...
- Node program: xitcoind

## Rule
Xitcoin is a permissioned-validator network.

Holding or staking $XTC never grants the right to become a validator.
A validator may only be created when its validator address is in the
on-chain Xitcoin approval list.

## Authority
The approval list is blockchain state, initialized in genesis.
It is never read from a private file on one server.

The authority for future changes will be a Xitcoin company multisignature
address selected before the permanent testnet genesis is created.

## Actions required
- Approve: add a validator address to the on-chain approval list.
- Revoke: remove the address from the list and deactivate any existing
  validator using that address.
- Protect: a revoked validator must not be able to create or unjail itself.
- Audit: every approval and revocation is recorded on-chain.

## Mandatory tests before network launch
1. An unapproved CreateValidator transaction is rejected.
2. An approved CreateValidator transaction is accepted.
3. Revocation deactivates an existing validator.
4. A revoked validator cannot create or unjail itself.
5. All nodes reach the same result from the same genesis.

## Current status
- Policy recorded.
- No permanent genesis created.
- No validator key created.
- No Xitcoin node launched.
