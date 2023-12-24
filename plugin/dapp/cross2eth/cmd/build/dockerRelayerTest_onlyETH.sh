#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null
set -x
set +e

# 主要在平行链上测试

source "./mainPubilcRelayerTest.sh"
source "./proxyVerifyTest.sh"

# shellcheck disable=SC2154
function StartDockerRelayerDeploy_onlyETH() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    # 修改 relayer.toml
    up_relayer_toml

    # 删除行
    sed -i "16,23"'d' "./relayer.toml"

    # 启动 ebrelayer
    start_docker_ebrelayerA

    docker cp "./deploy_chain.toml" "${dockerNamePrefix}_ebrelayera_1":/root/deploy_chain.toml
    docker cp "./deploy_ethereum.toml" "${dockerNamePrefix}_ebrelayera_1":/root/deploy_ethereum.toml

    # 部署合约 设置 bridgeRegistry 地址
    OfflineDeploy_chain
    # 修改 relayer.toml 字段
    sed -i 's/^BridgeRegistryOnChain=.*/BridgeRegistryOnChain="'"${chainBridgeRegistry}"'"/g' "./relayer.toml"

    # shellcheck disable=SC2154
    # shellcheck disable=SC2034
    {
        Boss4xCLI=${Boss4xCLIeth}
        CLIA=${CLIAeth}
        OfflineDeploy_ethereum "./deploy_ethereum.toml"
        ethereumBridgeBankOnETH="${ethereumBridgeBank}"
        ethereumBridgeRegistryOnETH="${ethereumBridgeRegistry}"
        ethereumMultisignAddrOnETH="${ethereumMultisignAddr}"

        sed -i '12,18s/BridgeRegistry=.*/BridgeRegistry="'"${ethereumBridgeRegistryOnETH}"'"/g' "./relayer.toml"
    }

    # 向离线多签地址打点手续费
    ChainCli=${MainCli}
    initMultisignChainAddr
    transferChainMultisignFee
    ChainCli=${Para8901Cli}

    docker cp "./relayer.toml" "${dockerNamePrefix}_ebrelayera_1":/root/relayer.toml
    InitRelayerA

    # 设置 token 地址
    # shellcheck disable=SC2154
    # shellcheck disable=SC2034
    {

        Boss4xCLI=${Boss4xCLIeth}
        CLIA=${CLIAeth}
        ethereumBridgeBank="${ethereumBridgeBankOnETH}"
        offline_create_bridge_token_chain_symbol "USDT"
        chainUSDTBridgeTokenAddrOnETH="${chainMainBridgeTokenAddr}"
        offline_create_bridge_token_chain_symbol "ETH"
        chainMainBridgeTokenAddrETH="${chainMainBridgeTokenAddr}"
        offline_create_bridge_token_eth_BTY
        ethereumBtyBridgeTokenAddrOnETH="${ethereumBtyBridgeTokenAddr}"
        offline_deploy_erc20_create_tether_usdt_USDT "USDT"
        ethereumUSDTERC20TokenAddrOnETH="${ethereumUSDTERC20TokenAddr}"
    }

    # shellcheck disable=SC2086
    {
        docker cp "${chainBridgeBank}.abi" "${dockerNamePrefix}_ebrelayera_1":/root/${chainBridgeBank}.abi
        docker cp "${chainBridgeRegistry}.abi" "${dockerNamePrefix}_ebrelayera_1":/root/${chainBridgeRegistry}.abi
        docker cp "${chainUSDTBridgeTokenAddrOnETH}.abi" "${dockerNamePrefix}_ebrelayera_1":/root/${chainUSDTBridgeTokenAddrOnETH}.abi
        docker cp "${chainMainBridgeTokenAddrETH}.abi" "${dockerNamePrefix}_ebrelayera_1":/root/${chainMainBridgeTokenAddrETH}.abi
        docker cp "${ethereumBridgeBankOnETH}.abi" "${dockerNamePrefix}_ebrelayera_1":/root/${ethereumBridgeBankOnETH}.abi
        docker cp "${ethereumBridgeRegistryOnETH}.abi" "${dockerNamePrefix}_ebrelayera_1":/root/${ethereumBridgeRegistryOnETH}.abi
    }

    # start ebrelayer B C D
    updata_toml_start_bcd
    restart_ebrelayerA

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# shellcheck disable=SC2034
# shellcheck disable=SC2154
function AllRelayerMainTest() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    set +e

    if [[ ${1} != "" ]]; then
        maturityDegree=${1}
        echo -e "${GRE}maturityDegree is ${maturityDegree} ${NOC}"
    fi

    # shellcheck disable=SC2120
    if [[ $# -ge 2 ]]; then
        chainID="${2}"
    fi

    get_cli

    # init
    ChainCli=${MainCli}
    InitChainValidator
    # para add
    initPara

    StartDockerRelayerDeploy_onlyETH
    #  test_all_onlyETH
    Boss4xCLI=${Boss4xCLIeth}
    CLIA=${CLIAeth}
    ethereumBridgeBank="${ethereumBridgeBankOnETH}"
    ethereumMultisignAddr="${ethereumMultisignAddrOnETH}"
    chainMainBridgeTokenAddr="${chainMainBridgeTokenAddrETH}"
    ethereumBtyBridgeTokenAddr="${ethereumBtyBridgeTokenAddrOnETH}"
    ethereumUSDTERC20TokenAddr="${ethereumUSDTERC20TokenAddrOnETH}"
    chainUSDTBridgeTokenAddr="${chainUSDTBridgeTokenAddrOnETH}"
    test_lock_and_burn "ETH" "USDT"

    # TestRelayerProxy_onlyETH
    start_docker_ebrelayerProxy
    setWithdraw_ethereum
    TestProxy

    echo_addrs
    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}
