#!/usr/bin/env bash
#shellcheck disable=SC2128
#shellcheck source=/dev/null
set -x
source ../dapp-test-common.sh

MAIN_HTTP=""

function rpc_test() {
    chain_RpcTestBegin cross2eth
    MAIN_HTTP="$1"
    echo "main_ip=$MAIN_HTTP"

    chain_RpcTestRst cross2eth "$CASE_ERR"
}

chain_debug_function rpc_test "$1" "$2"
