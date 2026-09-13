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

## Binary and isolated remediation evidence (2026-09-13)

Run 34781993413 completed both root and evmd diagnostics at PR43 head
14c8b4569c8f1136e555766e97a0aa43dbb532ca. Both analyses reported
GO-2024-2584, GO-2026-4479 and GO-2026-5932. The review remains expired
on 2026-09-05; no exception is renewed by this investigation.

The actual CI artifact 9903565616 from run 33780690435 identifies baseline
5ec8692e8fc1813d0892ee535af1a73953a1c4fb. Its xitcoind SHA256 is
2cca8f6e1f1a1e3db9d33e3183fdc202707b82d743a6cebc68dd7937dd58afe9,
matching its packaged manifest. Comparing this baseline with PR43 shows only
three diagnostic/documentation files changed; application Go sources and module
locks are identical. This is analysis of a CI artifact, not an attestation of
any deployed node. The binary was inspected without executing it.

The Go build metadata contains DTLS v2.2.12, STUN v2.0.0 and DTLS v3.1.4.
The govulncheck 1.7.0 binary JSON reports GO-2026-4479 at module level but
no DTLS v2 package/function finding. Module inclusion alone therefore does not
prove a callable DTLS session. The source production graph was collected with
`go list -deps -json ./cmd/evmd` from evmd, without test imports. Geth's NAT
STUN implementation uses `stun/v2.Dial("udp4", server)`; the transport argument
must be considered when assessing the generic STUN/DTLS dependency path.
`go mod why` paths involving ethclient.test are not exposure evidence.

SDK v0.54.4 slashing and vesting source were retrieved again from their exact
upstream tag and inspected for the previously identified fixes. The slashing
logic handles destination unbonding entries and persists their reduced balances;
the periodic-vesting recipient check rejects blocked addresses. The upstream
slashing advisory explicitly names fixes in v0.47.10 and v0.50.5, whereas OSV
retains an open range introduced at v0.50.0. Preserve this discrepancy for
independent security review; do not suppress the finding.

SDK v0.55.0 crypto/armor.go still imports golang.org/x/crypto/openpgp/armor.
An isolated SDK v0.54.4 candidate switches only that production import to
ProtonMail/go-crypto v1.4.1. Synthetic tests compare historical armor encoders,
public versions 0.0.0/0.0.1, key info, bcrypt/XSalsa20 and Argon2/ChaCha formats,
wrong passwords and malformed input. An isolated Geth candidate changes the
STUN import to v3.1.5 while preserving the UDP4 transport. These are local
candidates, not integrated or released dependency fixes. Full two-module
compilation, tests and a corrected-binary scan are prerequisites to clearance.

A signed replacement branch replays PR43's diagnostic changes from the existing
main baseline. No main rewrite, forced update, contract change, protection
change or security-gate relaxation is part of this regularization. CI failure
on the expired review remains expected and must block merging.
