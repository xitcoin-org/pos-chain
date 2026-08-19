# Xitcoin

Xitcoin is an EVM-compatible Proof-of-Stake blockchain built with Cosmos SDK, CometBFT and Cosmos EVM.

It is conceived as a common and evolving infrastructure layer through which public-sector participants, sovereign operators, companies, institutions, DeFi applications, communities and developers can interact using compatible transaction and application standards.

The network provides a shared technical language while each participant retains its own identity, authority, legal framework and operational responsibilities. Access to the protocol is based on published rules rather than institutional hierarchy or political alignment.

## Architecture and interoperability

Xitcoin combines:

- CometBFT consensus and Cosmos SDK protocol modules;
- native EVM execution and Ethereum-compatible tooling;
- Cosmos RPC, REST and gRPC interfaces;
- one canonical native-XTC economy across Cosmos and EVM execution;
- on-chain validator admission separated from staking, governance and operations;
- public source, release checksums and verifiable network state.

Cosmos and EVM activity execute on the same sovereign chain and use the same native XTC accounting. A future Cronos bridge is a separate security boundary using auditable lock/mint and burn/unlock accounting. It is not implied to be active by the existence of bridge source code.

## Network

| Property | Testnet | Mainnet |
| --- | --- | --- |
| Public name | Xitcoin Testnet | Xitcoin |
| Status | Active; coordinated reset candidate under review | Not launched |
| Cosmos Chain ID | `xitcoin-testnet-1` after reset | `xitcoin` |
| EVM Chain ID | `101089` | `101088` |
| Native asset | XTC | XTC |
| Atomic denomination | `axtc` | `axtc` |
| Decimals | 18 | 18 |

XTC is the native asset on both the Cosmos and EVM interfaces. The EVM interface uses the canonical native precompile; no separate wrapped-XTC token is required inside the Xitcoin network.

Technical service suffixes used during staging are not public network version names.

## Participation model

The candidate network aligns staking and validator-admission capacity at 258 positions:

- 195 sovereign reference positions;
- 63 public positions;
- 4 initially approved core validators;
- 1,000,000 XTC protocol minimum self-delegation;
- 5,000,000 XTC initial self-delegation for each core validator.

The sovereign framework keeps a defined participation pathway available for every reference position. Activation remains voluntary and begins only through the relevant participant's own initiative, an authorized operator and the same published admission, security and operational standards that apply across the network.

A reference or reserved position does not itself transfer assets or activate a validator. It preserves future access to the framework while the network continues to develop its public infrastructure, applications, developer ecosystem and community participation.

The sovereign reference methodology totals 390,000,000 XTC: 292,500,000 XTC in an equal component and 97,500,000 XTC in a square-root demographic component. This is a planning reference, not an executed transfer or automatic economic entitlement.

## Candidate testnet supply

The verified candidate genesis contains 1,250,000,000 testnet XTC:

- 20,000,000 XTC bonded by the four initial validators;
- 1,230,000,000 XTC remaining in five non-zero liquid genesis balances;
- zero mint inflation in the candidate genesis.

These quantities are testnet release accounting. They do not define the planned 5,250,000,000 XTC mainnet maximum supply or assign mainnet ownership. Exact account roles are published only after formal verification; the canonical genesis remains authoritative.

## Project repositories

| Repository | Responsibility |
| --- | --- |
| [`pos-chain`](https://github.com/xitcoin-org/pos-chain) | Consensus, Cosmos and EVM execution, native XTC, genesis and network documentation |
| [`contracts`](https://github.com/xitcoin-org/contracts) | Canonical Cronos V1, V2 and migration sources, deployments, audit scope and ecosystem references |
| [`migration-v1-to-v2`](https://github.com/xitcoin-org/migration-v1-to-v2) | Migration interface, reproducible build and continuity documentation |
| [`explorer-cosmos-testnet`](https://github.com/xitcoin-org/explorer-cosmos-testnet) | Cosmos testnet explorer source and network configuration |
| [`explorer-evm-testnet`](https://github.com/xitcoin-org/explorer-evm-testnet) | Blockscout deployment configuration and EVM explorer branding |
| [`brand`](https://github.com/xitcoin-org/brand) | Approved standalone Xitcoin token artwork and color references |

Contract addresses and audit claims are maintained only in `contracts`. Network identity and native-asset rules are maintained only in `pos-chain`. Other repositories link to those canonical records instead of duplicating them.

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
