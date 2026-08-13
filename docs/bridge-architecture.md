# Xitcoin bridge architecture

## Status

Design specification only.

No bridge contract, relayer, reserve, minting authority, key, Safe, or public bridge route is created by this document.

## First route

Cronos EVM XTC proxy:
`0xE45FE733BC8617FA6DAC8437FC44B5FFFA949991`

First route:

Cronos EVM <-> Xitcoin PoS native XTC

The existing Cronos XTC proxy is not modified for the bridge.

## Accounting

- One XTC locked on Cronos corresponds to one XTC released from the sole Xitcoin bridge reserve.
- One XTC locked on Xitcoin corresponds to one XTC released from the Cronos Vault.
- No additional XTC is minted.
- Every transfer has a unique source transaction, nonce, destination and amount.
- A processed transfer can never be processed again.

## Controls

- Existing XTC proxy upgrades remain controlled by its existing 3-of-3 voter mechanism.
- The bridge uses a separate 2-of-3 multisig with three independent keys.
- One emergency guardian may pause new bridge transfers only.
- Resume, limits, relayer changes, fee changes, upgrades and route closure require 2-of-3 approval.
- Route closure must preserve pending and valid user redemptions.

## Vault safety

- The Cronos Vault holds bridge XTC only.
- It cannot mint XTC.
- Its recovery function must never withdraw canonical XTC reserves.
- Recovery of unrelated assets requires multisig approval, a public event and a delay.
- The bridge reserve cannot be used for unrelated spending.

## Launch policy

- Testnet first, with no real XTC.
- Low per-transfer and daily limits at launch.
- Multiple confirmations before accepting a source-chain event.
- Public events, monitoring and daily reserve reconciliation.
- Independent audit before any mainnet route with real funds.
- No protocol fee at initial launch. Any later fee requires multisig approval, a public delay and a documented cap.

## Future routes

Solana, Sonic and Cosmos are future adapters to the same accounting model. They are not part of the first deployment.
