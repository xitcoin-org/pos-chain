# Xitcoin Bridge Testnet Acceptance

## Environment

Acceptance testing uses an isolated route and test assets. Production token governance and production bridge credentials are outside the test environment.

## Required tests

1. A finalized Cronos lock mints the exact native test XTC amount.
2. A finalized native XTC burn unlocks the exact Cronos test-token amount.
3. Minting without verified collateral is rejected.
4. Unlocking without a finalized burn is rejected.
5. A transfer identifier cannot be processed twice.
6. A modified chain, route, transaction, destination, amount or nonce is rejected.
7. A signer below the configured threshold cannot authorize settlement.
8. The configured signer threshold authorizes a valid settlement.
9. An unknown or retired signer is rejected.
10. Pausing blocks new settlement and preserves accounting.
11. The guardian cannot mint, burn, release collateral or resume settlement.
12. Route closure blocks new transfers and preserves valid redemptions.
13. Per-transfer, daily and outstanding-mint limits are enforced.
14. Confirmation depth and source-chain reorganizations are handled correctly.
15. A reconciliation mismatch pauses settlement.
16. Native XTC minted minus native XTC burned equals represented collateral.
17. A dead-address transfer cannot satisfy a burn requirement.

## Evidence

Each test record includes:

- test identifier and software revision;
- source event and destination transaction;
- route, addresses and amount;
- signer approvals;
- supply and collateral before and after;
- expected and observed results;
- genesis hash.

## Exit criteria

The testnet phase is complete when all required tests pass, authorization and replay protection are independently reviewed, accounting reconciles without mismatch, route limits and pause controls are validated, and operational procedures are approved.
