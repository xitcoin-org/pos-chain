# Independent security review — versioned Go integration

Refs xitcoin-org/pos-chain#33 and #44. Review is required from an independent security reviewer designated by xitcoin-org. No review is supplied or fabricated by this mission. The existing 2026-09-05 expiry and the six original advisory records remain unchanged. Fork publication, qualification, integration, merge and deployment are separate states.

| Advisory | Original upstream identity | Evidence to assess | Required decision |
|---|---|---|---|
| GO-2023-1821 | github.com/cosmos/cosmos-sdk | Original advisory retained; deprecated x/crisis must remain absent from production imports and registration. Fork coordinates do not clear the finding. | Independent assessment on exact root/evmd graph and symbols. |
| GO-2023-1881 | github.com/cosmos/cosmos-sdk | Same x/crisis constraint; retain distinct advisory and original affected ranges. | Independent assessment; no blanket SDK exemption. |
| GO-2024-2584 | github.com/cosmos/cosmos-sdk v0.54.4 | Upstream source dedeb7c80a91c47ae83f5352e29c3dd34e4a3fc6 contains slashing and vesting fixes; original OSV open range and upstream advisory discrepancy are preserved in PR44 evidence. | Reconcile source fix with advisory metadata; fork rename or scanner absence is not clearance. |
| GO-2025-3442 | github.com/cometbft/cometbft v0.39.4 | Unchanged module remains present; original module-level findings retained. | Assess exact graph/symbols and compatible remediation or formally reviewed disposition. |
| GO-2026-4479 | github.com/pion/dtls/v2 v2.2.12 via Cosmos Geth v1.17.2-cosmos-0 | Patch moves NAT STUN to v3.1.5, preserving UDP4 and removing DTLS v2 from the graph; historical loopback tests pass. | Confirm absence from both versioned production graphs and binary, inspect STUN semantics. |
| GO-2026-5932 | golang.org/x/crypto/openpgp | SDK armor uses ProtonMail/go-crypto v1.4.1. Historical encryption and synthetic compatibility tests retained; old OpenPGP is a test oracle only. | Confirm production graph/symbol absence; retain x/crypto module-level advisory and stripped-binary scanner limitation. |

Original records and historical proofs: docs/security-review-20260913.md and docs/security-remediation-20260913 in signed PR44 head 52c0ea290645c4b53b5648178f7176dfdc5a401c. New versioned graph results must be added separately, never substituted for the unfavorable baseline results.

Review packet must contain fork base and signed commit SHAs, Go-generated versions and sums, corresponding source/license notices, exact patches, upstream prepublication CI results, source scans for both modules, symbol-bearing binary scan, module-identity mapping, root revision consumed by evmd, and required PR checks on the final head. No release binary is executed against a node. Public activation 11.2.8 and PR21 remain preserved; Cyberscope contract stays frozen.

## Versioned dependency sources

Both forks preserve upstream module declarations. The two application modules use remote, immutable Go replacements. The original module identities and all six advisory records above remain authoritative for review. Scanner silence under a fork name is not clearance.

- `github.com/cosmos/cosmos-sdk` required `v0.54.4` → `github.com/xitcoin-org/cosmos-sdk v0.54.5-0.20260914091530-52ff14a25bee`. [Signed source](https://github.com/xitcoin-org/cosmos-sdk/commit/52ff14a25bee524f0e48ae7f00664443bcb22335), [upstream base](https://github.com/cosmos/cosmos-sdk/commit/dedeb7c80a91c47ae83f5352e29c3dd34e4a3fc6). Go module sum `h1:WU03GGJuL72vYRtn90DF79x51Y5TWGFbLa4XgTvKXuA=`; go.mod sum `h1:/1rNNn6uo2iw2ZEp9N8GK9PdWgH0OmxDly8MEAArYPI=`; qualified patch SHA256 `dc46d5e1e8497222040f0688277217599ce884e52ebdc820ba65f6a79a74902c`.
- `github.com/ethereum/go-ethereum` required `v1.16.9` → `github.com/xitcoin-org/go-ethereum v1.17.2-cosmos-0.0.20260913233706-e09f79643cd4`. [Signed source](https://github.com/xitcoin-org/go-ethereum/commit/e09f79643cd464404876f3565ce0a8fd8c0aeb16), [upstream base](https://github.com/cosmos/go-ethereum/commit/d99d6fa2c8d98b7cd653de4a9386d2da3db8f25c). Go module sum `h1:xDPGjArTeKQ8j19NP5e968xfZLaidTHxaGGk89JRIww=`; go.mod sum `h1:QtIPOaMKuz4zDZoQQ72nJJ9i5njr+MyRMl2tv7XteJQ=`; qualified patch SHA256 `8316073e371bdccaf788bd0f881d23094f1d486d1d7202d3edaa31fdf2416db4`.

Prepublication CI: [run 34788124256](https://github.com/xitcoin-org/pos-chain/actions/runs/34788124256), [run 34791392591](https://github.com/xitcoin-org/pos-chain/actions/runs/34791392591), [run 34793489922](https://github.com/xitcoin-org/pos-chain/actions/runs/34793489922), [run 34794203018](https://github.com/xitcoin-org/pos-chain/actions/runs/34794203018), [run 34788224704](https://github.com/xitcoin-org/pos-chain/actions/runs/34788224704), [run 34789811082](https://github.com/xitcoin-org/pos-chain/actions/runs/34789811082).

SDK production armor uses ProtonMail/go-crypto v1.4.1. Full upstream tests exposed its deliberate CRC24 omission; the fork restores Tendermint's checksum corruption check using the maintained encoder, with no copied or imported obsolete OpenPGP implementation. Existing encryption algorithms and upstream keyring expectations are preserved. Public golden vectors and CRC regression tests run in the SDK; the original dynamic legacy oracle remains separately available in the review packet. Original root/evmd compatibility tests are retained.

The separate historical CRC reader oracle retains 135 differential cases; its local checkpoint PASS is preserved and the published SDK revision will also be checked with the preserved oracle.

DTLS v2 remains in the module requirement graph through CometBFT v0.39.4 and IBC v11.2.0. Root `go mod why -m` reports it is not needed, and `go list -deps ./...` contains only DTLS v3. This residual module requirement and GO-2026-4479 remain review items. The provenance checker reports residual requirements and rejects effective production imports of DTLS v2 or STUN v2.

The Geth fork changes STUN to v3.1.5 with UDP4 preserved. Full root upstream tests passed; the initially failing keeper module required Go-generated tidy changes. Keeper build/tests, lint, generated-source and bad-dependency checks then passed. The first combined test-command failure is retained rather than relabeled a success.

## Licenses and corresponding source

The SDK core retains Apache-2.0 and its per-file notices. The enterprise modules retain their separate Cosmos Labs evaluation licenses; they were evaluated in upstream qualification and are not added to the application. ProtonMail armor retains its BSD-style license. Geth library files retain LGPL notices, and commands/components retain applicable GPL and per-file notices. Pion STUN retains its MIT license. The fork sources retain the original license files and explicit modification notices. No deployment binary is published by this integration; a symbol-bearing qualification binary is built and scanned on an isolated runner, then only its hash and build information are retained. A later distributor must assess applicable source and relinking obligations for its actual binary distribution.

## Review and merge status

No independent review is supplied by this mission. The exception expired on 2026-09-05; its gate and dependency acceptance locks remain unchanged. New graph diagnostics run separately and cannot grant acceptance. An independent security reviewer designated by xitcoin-org must assess the six cases against these exact revisions and source/binary findings. Repository maintainers must obtain the required approvals and successful mandatory checks before a normal merge. No administrative bypass, alert deletion, arbitrary expiry extension, force-push or deployment is part of this candidate.

Historical patch-only evidence remains under `security-remediation-20260913` and at signed head `52c0ea290645c4b53b5648178f7176dfdc5a401c`. It describes the earlier local graph; current locks and this provenance file describe the versioned candidate. The complete integration commit and the root library revision consumed by evmd are recorded separately to avoid a self-referential Go version.

Root library consumed by evmd: `v0.1.0-testnet.4.0.20260914092940-c82f40f042be`, signed commit `c82f40f042be57d6b8f82968a6d5a62a16b358c5`. Root Go sources and locks match that immutable library revision.
