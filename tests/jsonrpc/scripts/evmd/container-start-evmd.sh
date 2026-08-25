#!/bin/bash

# Container-friendly evmd initialization and startup script
# This runs inside the evmd container, so no Docker commands

set -e

echo "🔧 Starting evmd container initialization..."

# Set up variables (same as start-evmd.sh)
KEYRING="test"
KEYALGO="eth_secp256k1"
CHAINDIR="/data"
GENESIS="$CHAINDIR/config/genesis.json"
TMP_GENESIS="$CHAINDIR/config/tmp_genesis.json"
CHAIN_ID="local-4221"
BASEFEE=10000000

# Deterministic local-only identities. No reusable credential is stored in Git.
local_test_private_key() {
    printf 'xitcoin-local-test-only:%s' "$1" | sha256sum | awk '{print $1}'
}

VAL_KEY="mykey"
USER1_KEY="dev0"
USER2_KEY="dev1"
USER3_KEY="dev2"
USER4_KEY="dev3"

# Initialize chain directly (no Docker wrapper)
echo "🔧 Initializing chain..."
evmd init localtestnet -o --chain-id "$CHAIN_ID" --home "$CHAINDIR"

# Set client config
evmd config set client chain-id "$CHAIN_ID" --home "$CHAINDIR"
evmd config set client keyring-backend "$KEYRING" --home "$CHAINDIR"

# Add keys
echo "🔧 Adding standard test keys..."
for keyname in "$VAL_KEY" "$USER1_KEY" "$USER2_KEY" "$USER3_KEY" "$USER4_KEY"; do
    private_key="$(local_test_private_key "$keyname")"
    printf '\n' | evmd keys unsafe-import-eth-key "$keyname" "$private_key" --keyring-backend "$KEYRING" --home "$CHAINDIR"
    unset private_key
done

# Configure genesis file
echo "🔧 Configuring genesis file..."
jq '.app_state["staking"]["params"]["bond_denom"]="atest"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="atest"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state["gov"]["params"]["min_deposit"][0]["denom"]="atest"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state["gov"]["params"]["expedited_min_deposit"][0]["denom"]="atest"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state["bank"]["denom_metadata"]=[{"description":"The native staking token for evmd.","denom_units":[{"denom":"atest","exponent":0,"aliases":["attotest"]},{"denom":"test","exponent":18,"aliases":[]}],"base":"atest","display":"test","name":"Test Token","symbol":"TEST","uri":"","uri_hash":""}]' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state["evm"]["params"]["evm_denom"]="atest"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state["mint"]["params"]["mint_denom"]="atest"' "$GENESIS" > "$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

# Add genesis accounts
echo "🔧 Setting up genesis accounts..."
evmd genesis add-genesis-account "$VAL_KEY" 100000000000000000000000000atest --keyring-backend "$KEYRING" --home "$CHAINDIR"
evmd genesis add-genesis-account "$USER1_KEY" 1000000000000000000000atest --keyring-backend "$KEYRING" --home "$CHAINDIR"
evmd genesis add-genesis-account "$USER2_KEY" 1000000000000000000000atest --keyring-backend "$KEYRING" --home "$CHAINDIR"
evmd genesis add-genesis-account "$USER3_KEY" 1000000000000000000000atest --keyring-backend "$KEYRING" --home "$CHAINDIR"
evmd genesis add-genesis-account "$USER4_KEY" 1000000000000000000000atest --keyring-backend "$KEYRING" --home "$CHAINDIR"

# Generate validator transaction
evmd genesis gentx "$VAL_KEY" 1000000000000000000000atest --gas-prices "${BASEFEE}atest" --keyring-backend "$KEYRING" --chain-id "$CHAIN_ID" --home "$CHAINDIR"
evmd genesis collect-gentxs --home "$CHAINDIR"
admission_authority="$(evmd keys show "$VAL_KEY" -a --keyring-backend "$KEYRING" --home "$CHAINDIR")"
admission_validator="$(evmd keys show "$VAL_KEY" -a --bech val --keyring-backend "$KEYRING" --home "$CHAINDIR")"
jq --arg authority "$admission_authority" --arg validator "$admission_validator" '.app_state.validator_admission={authority:$authority,approved_validators:[$validator],max_approved_validators:1,minimum_self_delegation:"1000000000000000000000atest"}' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
evmd genesis validate-genesis --home "$CHAINDIR"

# Reduce block time by adjusting consensus timeouts
CONFIG_TOML="$CHAINDIR/config/config.toml"
sed -i 's/timeout_commit = "5s"/timeout_commit = "500ms"/g' "$CONFIG_TOML"
sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' "$CONFIG_TOML"
sed -i 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "100ms"/g' "$CONFIG_TOML"
sed -i 's/timeout_prevote = "1s"/timeout_prevote = "300ms"/g' "$CONFIG_TOML"
sed -i 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "100ms"/g' "$CONFIG_TOML"
sed -i 's/timeout_precommit = "1s"/timeout_precommit = "300ms"/g' "$CONFIG_TOML"
sed -i 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "100ms"/g' "$CONFIG_TOML"

echo "🚀 Starting evmd..."
exec evmd start \
    --home "$CHAINDIR" \
    --minimum-gas-prices=0.0001atest \
    --json-rpc.enable \
    --json-rpc.api eth,txpool,personal,net,debug,web3 \
    --json-rpc.address 0.0.0.0:8545 \
    --json-rpc.ws-address 0.0.0.0:8546 \
    --json-rpc.enable-profiling \
    --keyring-backend test \
    --chain-id "$CHAIN_ID"
