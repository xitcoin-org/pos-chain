#!/usr/bin/env python3
import hashlib
import json
from pathlib import Path

root = Path(__file__).resolve().parents[1]
chain = json.loads((root / "networks/testnet/chain.json").read_text())
genesis_path = root / "networks/testnet/genesis.json"
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

expected_hash = (root / "networks/testnet/genesis.sha256").read_text().split()[0]
actual_hash = hashlib.sha256(genesis_path.read_bytes()).hexdigest()
assert actual_hash == expected_hash == chain["genesis"]["sha256"]

identity = (root / "XITCOIN-NETWORK-IDENTITY.md").read_text()
normalizer = (root / "scripts/normalize-xitcoin-testnet-genesis.js").read_text()
readme = (root / "README.md").read_text()
status = (root / "docs/testnet-status.md").read_text()
testnet = (root / "docs/testnet.md").read_text()

for text in (identity, normalizer, readme, status, testnet):
    assert "xitcoin-testnet-1" in text

assert "genesis.chain_id !== 'xitcoin-testnet-1'" in normalizer
assert "`101089`" in identity
assert "`101088`" in identity
assert "| Status | Canonical public testnet active | Not launched |" in readme

for path in root.rglob("*"):
    if not path.is_file() or ".git" in path.parts:
        continue
    if path.suffix not in {".md", ".json", ".js", ".py", ".sh", ".yml", ".yaml"}:
        continue
    text = path.read_text(encoding="utf-8", errors="ignore")
    if "xitcoin-testnet" in text and "xitcoin-testnet-1" not in text:
        raise AssertionError(f"obsolete testnet identity in {path.relative_to(root)}")

print("canonical_identity=OK")
print("genesis_checksum=OK")
print("currency_identity=OK")
print("mainnet_boundary=OK")
print("release_state_validation=OK")
