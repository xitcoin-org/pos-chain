#!/usr/bin/env bash
set -euo pipefail

# Temporary, narrowly scoped exceptions for upstream advisories that cannot be
# remediated safely in this dependency line. These are dependency locks, not a
# generic advisory allow-list: any dependency change requires a fresh review.
review_by=2026-09-05
if [[ "$(date -u +%F)" > "$review_by" ]]; then
  echo "govulncheck exception review expired on $review_by" >&2
  exit 1
fi

accepted=(
  GO-2023-1821 # x/crisis is compiled upstream but not imported or registered by Xitcoin.
  GO-2023-1881 # x/crisis is compiled upstream but not imported or registered by Xitcoin.
  GO-2024-2584 # Cosmos SDK 0.54.4 is newer than the advisory's <0.47.10 range.
  GO-2025-3442 # CometBFT 0.39 has no compatible fixed release; exact version is locked.
  GO-2026-4479 # Pion DTLS v2 has no fixed release; v3 is not API-compatible.
  GO-2026-5932 # OpenPGP is pulled by Cosmos keyring; no x/crypto version fixes it.
)

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

require_module() {
  local module=$1 expected=$2 actual
  actual="$(cd "$repo_root" && go list -m -f '{{.Version}}' "$module")"
  if [[ "$actual" != "$expected" ]]; then
    printf 'Reviewed govulncheck exception invalidated: %s is %s, expected %s\n' \
      "$module" "$actual" "$expected" >&2
    exit 1
  fi
}

require_replacement() {
  local module=$1 expected_path=$2 expected_version=$3 actual
  actual="$(cd "$repo_root" && go list -m -f '{{if .Replace}}{{.Replace.Path}} {{.Replace.Version}}{{end}}' "$module")"
  if [[ "$actual" != "$expected_path $expected_version" ]]; then
    printf 'Reviewed govulncheck exception invalidated: replacement for %s is %s\n' \
      "$module" "${actual:-absent}" >&2
    exit 1
  fi
}

require_module github.com/cosmos/cosmos-sdk v0.54.4
require_module github.com/cometbft/cometbft v0.39.4
require_module github.com/pion/dtls/v2 v2.2.12
require_module golang.org/x/crypto v0.55.0
require_replacement github.com/ethereum/go-ethereum github.com/cosmos/go-ethereum v1.17.2-cosmos-0

# Do not reintroduce a JSON-RPC method accepting raw private keys. OpenPGP
# armor remains compiled through the upstream Cosmos CLI/keyring only.
if grep -R --line-number --include='*.go' -E 'ImportRawKey|personal_importRawKey' \
  "$repo_root/rpc"; then
  echo 'Raw private-key import must not be exposed by JSON-RPC' >&2
  exit 1
fi

# x/crisis is deprecated and has two no-fix advisories. It is present in the
# Cosmos SDK module archive, but must never be wired into this application.
if grep -R --line-number --include='*.go' 'cosmos-sdk/x/crisis' "$repo_root"; then
  echo 'The deprecated Cosmos x/crisis module must not be imported' >&2
  exit 1
fi

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
