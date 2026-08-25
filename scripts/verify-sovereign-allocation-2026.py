#!/usr/bin/env python3
import csv
import json
import sys
from decimal import Decimal
from pathlib import Path

repository = Path(sys.argv[1]).resolve() if len(sys.argv) == 2 else Path(__file__).resolve().parents[1]
network = repository / "networks" / "xitcoin-testnet-1"

with (network / "sovereign-allocation-index-2026.csv").open(encoding="utf-8", newline="") as handle:
    allocations = list(csv.DictReader(handle))
with (network / "territorial-consolidation-2026.csv").open(encoding="utf-8", newline="") as handle:
    mappings = list(csv.DictReader(handle))
with (network / "sovereign-allocation-index-2026.json").open(encoding="utf-8") as handle:
    metadata = json.load(handle)

assert isinstance(metadata, dict)
assert metadata["selection_basis"] == "united_nations_member_states"
assert len(allocations) == 193
assert len(mappings) == 39
assert len({row["iso3"] for row in allocations}) == 193
assert sum(int(row["allocation_atomic"]) for row in allocations) == 386_000_000 * 10**18
assert sum(Decimal(row["allocation_xtc"]) for row in allocations) == Decimal("386000000.000000000000000000")

china = {
    row["source_iso3"]: row["target_iso3"]
    for row in mappings
    if row["source_iso3"] in {"HKG", "MAC", "TWN"}
}
assert china == {"HKG": "CHN", "MAC": "CHN", "TWN": "CHN"}
assert not {"COK", "NIU", "ESH"}.intersection(row["source_iso3"] for row in mappings)

print("positions=193")
print("consolidated_records=39")
print("allocation_total_xtc=386000000.000000000000000000")
print("china_components=HKG,MAC,TWN")
print("validation=OK")
