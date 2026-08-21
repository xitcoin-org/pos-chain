#!/usr/bin/env bash
set -euo pipefail

# Temporary, narrowly scoped exceptions for upstream advisories that cannot be
# remediated in this dependency line. Re-review before this date; the script
# deliberately fails after it so exceptions cannot become permanent silently.
review_by=2026-09-21
if [[ "$(date -u +%F)" > "$review_by" ]]; then
  echo "govulncheck exception review expired on $review_by" >&2
  exit 1
fi

accepted=(
  GO-2024-2584 # Cosmos SDK 0.54.3 is newer than the advisory's <0.47.10 range.
  GO-2026-4479 # Pion DTLS v2 has no fixed release; v3 is not API-compatible.
  GO-2026-5932 # OpenPGP is pulled by Cosmos keyring; no x/crypto version fixes it.
)

report="$(mktemp)"
trap 'rm -f "$report"' EXIT

set +e
govulncheck "$@" 2>&1 | tee "$report"
status=${PIPESTATUS[0]}
set -e

if (( status == 0 )); then
  exit 0
fi
if (( status != 3 )); then
  exit "$status"
fi

mapfile -t found < <(
  sed -nE 's/^Vulnerability #[0-9]+: (GO-[0-9]{4}-[0-9]+)$/\1/p' "$report" |
    sort -u
)
mapfile -t allowed < <(printf '%s\n' "${accepted[@]}" | sort -u)
mapfile -t unexpected < <(comm -23 <(printf '%s\n' "${found[@]}") <(printf '%s\n' "${allowed[@]}"))

if (( ${#found[@]} == 0 )); then
  echo "govulncheck failed without a parseable advisory ID" >&2
  exit 1
fi
if (( ${#unexpected[@]} > 0 )); then
  printf 'Unaccepted Go vulnerability: %s\n' "${unexpected[@]}" >&2
  exit 1
fi

printf 'Only reviewed upstream exceptions remain (review by %s): %s\n' \
  "$review_by" "${found[*]}"
