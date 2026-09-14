# Security Policy

## Reporting a vulnerability

Report security vulnerabilities through GitHub Private Vulnerability Reporting for this repository. Include the affected component, impact, reproduction steps and any supporting evidence.

Do not open a public issue for an undisclosed vulnerability.

## Reviewed Go advisories

[SECURITY-ASSESSMENT.md](SECURITY-ASSESSMENT.md) is the current source of truth
for PR44's six dispositions, exact versions, evidence, expiry and residual risk.
[Provenance](docs/go-fork-provenance.md) records immutable sources and sums.
The 2026-09-14 technical assessment is by the PR author; it is not an independent
approval or release authorization.

- GO-2023-1821 and GO-2023-1881: x/crisis is absent from the examined production
  imports and application registration; these are distinct non-applicability
  conclusions, not blanket SDK exemptions.
- GO-2024-2584: the consumed SDK fork based on v0.54.4 contains the slashing
  fix; the open Go advisory range is reconciled against source.
- GO-2025-3442: CometBFT v0.39.4 contains the peer-height regression fix.
  Blocksync is active and the broad Go module findings remain recorded.
- GO-2026-4479: Geth uses STUN v3.1.5 over UDP4; DTLS v3.1.4 has the nonce
  fix. DTLS v2.2.12 remains transitively required by CometBFT/IBC but is absent
  from the examined production imports.
- GO-2026-5932: SDK armor uses ProtonMail/go-crypto v1.4.1 with CRC24
  compatibility checking. Legacy OpenPGP is a test oracle only; x/crypto
  remains for other primitives. CRC24 is not authentication.

The gates retain all six identities, exact dependency and assessment checks,
obsolete-import/raw-key RPC rejection, scanner errors, new findings and expiry.
Dependency security scans root and evmd; Test and build also scans its compiled
artifact. Historical symbol-bearing qualification is separately identified;
scanner silence under fork coordinates or in stripped binaries is not clearance.
Renewal requires a new justified assessment and tests, never just a date change.

## Scope

Reports may cover:

- consensus and state-transition logic;
- EVM execution and precompiles;
- staking, governance and validator admission;
- bridge authorization and accounting;
- build and release integrity;
- public RPC and API behavior implemented by this repository.

Third-party services and upstream dependencies should also be reported to their respective maintainers when appropriate.

## Disclosure

Please allow maintainers time to investigate and coordinate a fix before public disclosure. Receipt, assessment and remediation timelines depend on severity and reproducibility.
