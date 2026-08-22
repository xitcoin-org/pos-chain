# Security Policy

## Reporting a vulnerability

Report security vulnerabilities through GitHub Private Vulnerability Reporting for this repository. Include the affected component, impact, reproduction steps and any supporting evidence.

Do not open a public issue for an undisclosed vulnerability.

## Reviewed upstream Go advisories

The release CI scans both source and the compiled `xitcoind` artifact. It
currently recognizes upstream-only advisories under exact dependency
and code-surface locks in `scripts/run-govulncheck.sh`:

- `GO-2024-2584`: the database range is stale for Cosmos SDK `v0.54.3`; the
  published fix is present from Cosmos SDK `v0.47.10`.
- `GO-2023-1821` and `GO-2023-1881`: the deprecated Cosmos `x/crisis`
  package is present in the upstream module archive but is neither imported
  nor registered by Xitcoin. CI prohibits adding an import.
- `GO-2025-3442`: the compatible CometBFT `v0.39` line has no published fixed
  release. The exact version is locked pending a Cosmos-compatible release;
  the binary scan prevents this advisory from being overlooked.
- `GO-2026-4479`: Pion DTLS v2 has no fixed v2 release. It is pulled only by
  the STUN/NAT path of the pinned Cosmos go-ethereum fork. Migration to DTLS
  v3 requires the corresponding upstream geth change.
- `GO-2026-5932`: Cosmos SDK still compiles the legacy OpenPGP armor helper.
  Xitcoin does not expose the former `personal_importRawKey` JSON-RPC method.

These are not permanent suppressions. A new advisory, a dependency or
replacement change, reintroduction of raw-key import, or the review deadline
fails CI. The exceptions must be removed as soon as compatible upstream
releases are available.

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
