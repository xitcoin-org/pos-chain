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

Run the complete test suite with the test build tag:

```bash
go test -tags=test ./...
```

Some permission-oriented client tests must run as a non-privileged user.

## Scope

This repository contains the Xitcoin chain implementation. It does not contain production keys, node state, private genesis files, backups or production infrastructure configuration.

For upstream attribution, see [UPSTREAMS.md](../UPSTREAMS.md).
