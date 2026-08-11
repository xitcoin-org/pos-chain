# Validator Admission

Xitcoin operates a permissioned validator admission model.

## Authority

A single on-chain authority is defined in the Validator Admission genesis state.
Only this authority can:

- approve a validator operator address;
- revoke an approved validator;
- update Validator Admission policy parameters.

Approval is an explicit on-chain transaction. An operator cannot self-approve.

## Validator requirements

A validator creation transaction is accepted only when:

1. the validator operator address is currently approved;
2. the self-delegation meets the active minimum;
3. the validator follows the network staking and consensus rules.

The initial mainnet policy is:

| Parameter | Value |
| --- | ---: |
| Maximum approved validator addresses | 208 |
| Minimum self-delegation | 5,000,000 XTC |

Technical on-chain amounts use the native `xtc` denomination with 18 decimal places.

## Revocation

Revoking an operator removes its approval. If that validator already exists, the
network jails it. A revoked validator cannot unjail unless it is approved again.

## Policy updates

Policy updates are on-chain actions restricted to the configured Xitcoin
authority. A capacity update cannot be reduced below the number of already
approved validator addresses.

## Operator process

1. Submit the validator operator address for review.
2. Receive explicit approval from the Xitcoin authority.
3. Create the validator with at least the required self-delegation.
4. Maintain operational, security and network-policy compliance.

This module does not alter the public testnet genesis or activate any mainnet
configuration by itself.
