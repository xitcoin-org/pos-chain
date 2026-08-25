#!/usr/bin/env python3
import hashlib
import json
from pathlib import Path

root = Path(__file__).resolve().parents[1]
chain = json.loads((root / "networks/xitcoin-testnet-1/chain.json").read_text())
genesis_path = root / "networks/xitcoin-testnet-1/genesis.json"
genesis = json.loads(genesis_path.read_text())

assert chain["status"] == "active"
assert chain["cosmos_chain_id"] == "xitcoin-testnet-1"
assert chain["evm_chain_id"] == 101089
assert chain["evm_chain_id_hex"] == "0x18ae1"
assert chain["native_currency"] == {
    "name": "Xitcoin",
    "symbol": "XTC",
    "decimals": 18,
    "base_denom": "axtc",
}
assert genesis["chain_id"] == "xitcoin-testnet-1"

expected_hash = (root / "networks/xitcoin-testnet-1/genesis.sha256").read_text().split()[0]
actual_hash = hashlib.sha256(genesis_path.read_bytes()).hexdigest()
assert actual_hash == expected_hash == chain["genesis"]["sha256"]

identity = (root / "XITCOIN-NETWORK-IDENTITY.md").read_text()
normalizer = (root / "scripts/normalize-xitcoin-testnet-genesis.js").read_text()
readme = (root / "README.md").read_text()
status = (root / "docs/testnet-status.md").read_text()
testnet = (root / "docs/testnet.md").read_text()

prohibited_test_credentials = (
    "gesture inject test cycle",
    "copper push brief egg",
    "maximum display century economy",
    "will wear settle write",
    "doll midnight silk carpet",
    "aunt imitate maximum student",
    "88cbead91aee890d",
    "741de4f8988ea941",
    "3b7955d25189c99a",
    "8a36c69d940a92fc",
)

for text in (identity, normalizer, readme, status, testnet):
    assert "xitcoin-testnet-1" in text

assert "genesis.chain_id !== 'xitcoin-testnet-1'" in normalizer
assert "`101089`" in identity
assert "`101088`" in identity
assert "| Status | Canonical public testnet active | Not launched |" in readme

for path in root.rglob("*"):
    if not path.is_file() or ".git" in path.parts:
        continue
    if path.resolve() == Path(__file__).resolve():
        continue
    if path.suffix not in {".md", ".json", ".js", ".py", ".sh", ".yml", ".yaml"}:
        continue
    text = path.read_text(encoding="utf-8", errors="ignore")
    for credential in prohibited_test_credentials:
        if credential in text:
            raise AssertionError(
                f"stored legacy test credential in {path.relative_to(root)}"
            )
    if "xitcoin-testnet" in text and "xitcoin-testnet-1" not in text:
        raise AssertionError(f"obsolete testnet identity in {path.relative_to(root)}")

for path in (
    root / "local_node.sh",
    root / "tests/jsonrpc/docker-compose.yml",
    root / "tests/jsonrpc/scripts/evmd/start-evmd.sh",
):
    assert "--privileged" not in path.read_text()

policy = (root / "x/validatoradmission/types/policy.go").read_text()
assert 'coin.Denom != "axtc"' in policy
assert "ValidateGenesisPolicy" in policy

print("canonical_identity=OK")
print("genesis_checksum=OK")
print("currency_identity=OK")
print("mainnet_boundary=OK")
print("release_state_validation=OK")
print("legacy_test_credentials=ABSENT")
print("local_validator_policy=OK")
