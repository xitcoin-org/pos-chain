# Go advisory reassessment — 2026-09-13

Related: #33, #41, #42, #43.

Status: incomplete; no renewed exception approval, release clearance or deployed-node vulnerability verdict.

## Executed evidence

PR43 commit 443acedf3e8f188dcf6d3c03fa186dbee7cb9c23, Dependency security run 34781675017, job 103789716965. Go 1.26.7, govulncheck 1.7.0, root ./... scan completed. It reports reachable GO-2026-5932, GO-2026-4479 and GO-2024-2584. Three other required-module findings were not classified as called. This root result does not cover the separate evmd module or the deployed binary. The expired review still fails the gate.

## Cosmos GO-2024-2584

The wrapper comment citing only <0.47.10 is incomplete. The upstream advisory explicitly lists fixes 0.47.10 and 0.50.5: https://github.com/cosmos/cosmos-sdk/security/advisories/GHSA-86h5-xcpx-cfqc

Source inspection of v0.54.4 confirms the undelegation-after-redelegation slashing logic from d1b5b0c5ae2c51206cc1849e09e4d59986742cc3: it retrieves destination unbonding delegations, checks creation height/maturity, subtracts the eligible slash amount and persists the updated entry. This is evidence of the upstream fix, not merely a numerical version comparison. https://github.com/cosmos/cosmos-sdk/blob/v0.54.4/x/staking/keeper/slash.go

The OSV JSON also references 7dbed2fc0c3ed7c285645e21cb1037d8810372ae, which addresses a separate blocked-recipient vesting issue. Its BlockedAddr guard is present in v0.54.4 CreatePeriodicVestingAccount: https://github.com/cosmos/cosmos-sdk/blob/v0.54.4/x/auth/vesting/msg_server.go

The OSV source retains an introduced 0.50.0 range without a following fixed event. This explains a metadata mismatch with the upstream advisory. Preserve the scanner result and obtain review before changing acceptance policy.

## OpenPGP GO-2026-5932

The root scan reaches Cosmos keyring import/export and armor. v0.54.4 crypto/armor.go imports x/crypto/openpgp/armor for encoding, while its encryption uses Argon2/ChaCha20Poly1305 and legacy bcrypt/xsalsa paths. This distinction narrows investigation but does not prove the unmaintained parser safe. No fixed x/crypto version is listed. A maintained replacement requires upstream/fork compatibility tests for existing armored keys; no real keys may be used. https://github.com/cosmos/cosmos-sdk/blob/v0.54.4/crypto/armor.go

## DTLS GO-2026-4479

The root scan reports dtls/v2 v2.2.12 through tracer initialization and indirect/interface paths. v3.1.4 also exists in go.mod; that does not replace v2. Determine the full import path and whether an actual DTLS session is reachable before selecting a compatible dependency fix. No v2 fixed release is listed in the official advisory. Do not replace a v2 module with v3 by version string alone. https://github.com/golang/vulndb/blob/master/data/osv/GO-2026-4479.json

## Next gate

Run evmd analysis even when the root scan fails; preserve both failures. Then inspect full dependency paths for DTLS and keyring. No date extension, removed advisory, suppressed scanner output, contract change, deployment, transaction, restart or upgrade is part of this PR.
