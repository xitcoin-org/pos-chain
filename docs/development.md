# Development

This guide describes how to build and test the Xitcoin node from source.

## Requirements

- Go version declared in `go.mod`
- A supported Linux, macOS or Windows development environment
- Git

## Build

From the repository root:

```bash
make build
```

The resulting daemon is written to:

```text
build/xitcoind
```

## Test

Run tests in both Go modules with the test build tag:

```bash
GOWORK=off go test -tags=test ./...
(cd evmd && GOWORK=off go test -tags=test ./...)
```

Some permission-oriented client tests must run as a non-privileged user.

## Scope

This repository contains the Xitcoin chain implementation. It does not contain production keys, node state, private genesis files, backups or production infrastructure configuration.

For upstream attribution, see [UPSTREAMS.md](../UPSTREAMS.md).

## PR44 security and qualification

[Current assessment](../SECURITY-ASSESSMENT.md) records the exact dependency
versions, six advisory dispositions and limits. Root tests alone do not cover
evmd. Versioned graph CI reuses digest-bound acquired application results and
runs targeted guard tests; it does not claim a new full suite or binary scan.
The ordinary dependency gates separately scan current source. Dated remediation
reports describe historical candidates and must not replace current provenance.
