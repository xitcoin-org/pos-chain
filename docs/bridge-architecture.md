# Xitcoin Bridge Architecture

## Status

Design specification only.

No bridge contract, relayer, reserve, minting authority, key, Safe, or public bridge route is created by this document.

## First Route

Cronos EVM XTC proxy:

`0xE45FE733BC8617FA6DAC8437FC44B5FFFA949991`

Route:

Cronos EVM XTC <-> Xitcoin native XTC

The existing Cronos XTC proxy is not modified for the bridge.

## Route Decision

Xitcoin includes IBC and ERC-20 middleware for processing verified ICS-20 packets. That middleware cannot lock or verify an external CRC-20 deposit on Cronos EVM.

The first XTC route therefore requires dedicated bridge components:

- a Cronos Vault that accepts only canonical XTC;
- an Xitcoin bridge adapter that verifies approved source-chain transfers;
- a dedicated Xitcoin reserve that releases native XTC;
- independent bridge signers and monitoring.

## Accounting

- One XTC locked in the Cronos Vault corresponds to one XTC released from the Xitcoin reserve.
- One XTC returned to the Xitcoin reserve corresponds to one XTC released from the Cronos Vault.
- No additional XTC is minted.
- Every transfer has a unique source transaction, nonce, destination and amount.
- A processed transfer can never be processed again.
- Reserve reconciliation must detect any mismatch before transfers resume.

## Controls

- Existing XTC proxy upgrades remain controlled by the existing 3-of-3 voter mechanism.
- Bridge governance and transfer attestations use a separate 2-of-3 multisignature.
- One emergency guardian may pause new bridge transfers only.
- Resume, limits, signer changes, fee changes and route closure require 2-of-3 approval.
- Route closure must preserve pending and valid user redemptions.

## Vault and Reserve Safety

- The Cronos Vault holds bridge XTC only and cannot mint XTC.
- The Xitcoin reserve is used only for bridge settlement.
- Recovery must never withdraw canonical XTC bridge reserves.
- Recovery of unrelated assets requires multisignature approval, a public event and a delay.
- The bridge reserve cannot be used for unrelated spending.

## Fees

- Fees are disabled for testnet.
- No production fee value is set by this document.
- Any future fee requires a documented cap, a public delay and 2-of-3 approval.

## Launch Policy

- Testnet first, with fictitious amounts only.
- Low per-transfer and daily limits at initial production launch.
- Multiple confirmations before accepting a source-chain event.
- Public events, monitoring and daily reserve reconciliation.
- Independent audit before any mainnet route with real funds.

## Future Routes

Solana, Sonic and Cosmos are future adapters to the same accounting model. They are not part of the first deployment.
