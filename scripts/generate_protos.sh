#!/usr/bin/env bash
set -euo pipefail

buf generate --template proto/buf.gen.gogo.yaml

test -d github.com/xitcoin-org/pos-chain
cp -r github.com/xitcoin-org/pos-chain/* ./
rm -rf github.com
