# Verified Go remediation candidate

Refs #33, #43, #44. This is a reviewable patch set, not an active dependency
release. The repository's dependency locks and expired-review gate remain intact.
No contract, network configuration, deployed node, key or service is changed.

## Scope and provenance

- SDK v0.54.4 (`dedeb7c80a91c47ae83f5352e29c3dd34e4a3fc6`): replace only the
  production ASCII-armor import with ProtonMail/go-crypto v1.4.1. Keep the
  Argon2/ChaCha20Poly1305 and legacy bcrypt/XSalsa20 encryption paths unchanged.
  Tests use the old encoder/reader only as a synthetic compatibility oracle.
- Cosmos Geth v1.17.2-cosmos-0
  (`d99d6fa2c8d98b7cd653de4a9386d2da3db8f25c`): migrate NAT STUN to v3.1.5,
  retaining the UDP4 call. Tidy removes DTLS v2 from the dependency graph.
  The response and missing-address behavior are tested on loopback only.
- Xitcoin root and evmd: local relative replacements plus regression tests for
  Ethereum key import/export through the actual codecs and keyring options.
  No private test material is embedded; all keys are generated synthetically.

## Reconstruct the isolated candidate

Use Go 1.26.7 and preserve at least 5 GiB free. Start from PR43 head
`14c8b4569c8f1136e555766e97a0aa43dbb532ca`. Source Go files and module locks
at signed diagnostic commit `5adc799e264bf3b013fe237bc213e577dde3e61a` are
identical to that baseline.

Create three sibling copies named `pos-chain`, `sdk-armor` and `geth-stun3`.
Obtain the two dependencies at the exact versions/commits above; never edit a
shared Go module cache. Apply `sdk-armor.patch` inside sdk-armor,
`geth-stun3.patch` inside geth-stun3, and `pos-chain-integration.patch` inside
pos-chain, using `patch --batch -p1`. All patch hashes are in validation.json.
The relative replacements resolve both modules to the same isolated sources.
The patches include Go-generated sums; do not edit those sums manually.

## Executed validation

Both modules passed `go build ./...`. The evmd command was built both with and
without symbol stripping. Fifteen targeted top-level test suites passed:
SDK armor and historical formats, Xitcoin Ethereum keyring import/export,
Geth STUN v3 loopback responses and parsing, and evmd encoding/configuration.
The exact suite names and successful stages are in validation.json. These are
not claims that all SDK, Geth or Xitcoin tests/lints have run.

Both production package graphs exclude golang.org/x/crypto/openpgp and
DTLS v2; ProtonMail armor and STUN v3 are present. Source scans for root and
evmd, plus both binary scans, were collected without advisory filtering.

## Scanner limits and acceptance

The stripped binary scan still reports OpenPGP package/symbol placeholders.
This is govulncheck's conservative module-level fallback when symbols are
missing; see [the scanner implementation](https://github.com/golang/vuln/blob/v1.7.0/internal/vulncheck/binary.go).
The unstripped binary scan contains only module-level findings for
GO-2026-5932 and GO-2025-3442, with no package/function findings for either.
No stripped-binary result has been erased or accepted as clean.

Local SDK/Geth replacements are recorded as `(devel)` and do not establish
upstream advisory/version matching. The disappearance of SDK advisories in
these scans must not be presented as a vulnerability clearance. The independent
source evidence that v0.54.4 contains the Cosmos slashing fix remains in
[the security review](../security-review-20260913.md). The open OSV range and
GO-2024-2584 finding from the original versioned build remain review items.
SDK v0.55.0 still imports the old OpenPGP armor package.

Before active integration, publish maintained, immutable SDK/Geth revisions
under their repository requirements, replace the local paths with exact
versions in both modules, qualify that exact graph and obtain independent
security review. Reassess the existing version locks and six advisory cases;
do not extend the expired 2026-09-05 date arbitrarily. Do not merge with red CI
or a missing mandatory review. These patches authorize no deployment.

The original licenses accompany this patch set. Fork maintenance and complete
upstream prepublication checks remain separate from these targeted validations.
