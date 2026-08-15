# Xitcoin

Xitcoin is an EVM-compatible Proof-of-Stake blockchain built with Cosmos SDK, CometBFT and Cosmos EVM.

## Network

| Property | Testnet | Mainnet |
| --- | --- | --- |
| Status | Active; coordinated reset candidate under review | Not launched |
| Cosmos Chain ID | `xitcoin-testnet` after reset | `xitcoin` |
| EVM Chain ID | `101089` | `101088` |
| Native asset | XTC | XTC |
| Atomic denomination | `axtc` | `axtc` |
| Decimals | 18 | 18 |

XTC is the native asset on both the Cosmos and EVM interfaces. The EVM interface uses the canonical native precompile; no separate wrapped-XTC token is required.

## Build

Prerequisites: Go 1.25.9, Node.js 24 and the toolchain documented by the project.

```bash
git clone https://github.com/xitcoin-org/pos-chain.git
cd pos-chain
make build
```

Run the core validation suites:

```bash
go test ./x/validatoradmission/...
go test ./x/bridge/...
go test ./evmd/config
go test ./evmd/cmd/evmd/cmd
go test -tags test ./evmd/tests/integration/precompiles/werc20
```

## Documentation

- [Network identity](XITCOIN-NETWORK-IDENTITY.md)
- [Testnet](docs/testnet.md)
- [Testnet operations](docs/testnet-operations.md)
- [Bridge architecture](docs/bridge-architecture.md)
- [Bridge acceptance tests](docs/bridge-testnet-acceptance.md)
- [Validator admission](docs/validator-admission.md)
- [Development](docs/development.md)

## Contributing

Read [CONTRIBUTING.md](CONTRIBUTING.md) before submitting changes. Security reports follow the process in [SECURITY.md](SECURITY.md).

## License

Licensed under the [Apache License 2.0](LICENSE).
