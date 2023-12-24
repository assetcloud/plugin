#!/bin/bash
# 官方ci集成脚本
strpwd=$(pwd)
strcmd=${strpwd##*dapp/}
strapp=${strcmd%/cmd*}

OUT_DIR="${1}/$strapp"
#FLAG=$2

mkdir -p "${OUT_DIR}"
cp ./ci/* "${OUT_DIR}"

CHAIN_PATH=$(go list -f "{{.Dir}}" github.com/assetcloud/chain)
PLUGIN_PATH=$(go list -f "{{.Dir}}" github.com/assetcloud/plugin)
# copy chain toml

cp "${CHAIN_PATH}/cmd/chain/chain.test.toml" "${OUT_DIR}"
cp "${PLUGIN_PATH}/chain.para.toml" "${OUT_DIR}"
