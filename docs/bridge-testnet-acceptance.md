# Xitcoin Bridge Testnet Acceptance Plan

## Status

Test plan only. No bridge contract, signer, reserve, route, asset or deployment is created by this document.

## Test Environment

- Isolated testnet only.
- Test assets only; no canonical Cronos XTC and no real user funds.
- Separate test wallets and signer material.
- The production XTC proxy and its voter wallets are not used.

## Required Test Cases

1. A valid Cronos test-token vault deposit releases the same amount from the Xitcoin test reserve.
2. A valid return to the Xitcoin test reserve releases the same amount from the Cronos test vault.
3. A transfer identifier cannot be processed twice.
4. A conflicting destination, amount or nonce is rejected.
5. One bridge signer alone cannot authorize settlement.
6. Two authorized bridge signers can authorize settlement.
7. A guardian pause prevents new transfers.
8. Pausing does not erase valid recorded transfers or reserve accounting.
9. Route closure prevents new transfers while preserving valid pending redemptions.
10. Per-transfer and daily limits are enforced.
11. Testnet fees remain disabled.
12. Recovery rejects canonical bridge-reserve XTC and permits only unrelated test assets through the required approval path.
13. Source-chain confirmation and reorganization handling are tested before settlement.
14. Reserve reconciliation detects any accounting mismatch and blocks further settlement until reviewed.

## Evidence Required

Each test records:

- source transaction identifier;
- destination transaction identifier;
- amount and destination;
- signer approvals;
- reserve balances before and after;
- expected result;
- actual result.

## Exit Criteria

No testnet bridge implementation can advance to audit until all required tests pass, replay protection is independently reviewed, and reserve reconciliation produces no mismatch.
