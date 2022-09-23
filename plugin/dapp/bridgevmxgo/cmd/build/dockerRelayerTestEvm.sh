#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck disable=SC2219
# shellcheck source=/dev/null
set -x
set +e

# 主要在平行链上测试

source "./mainPubilcRelayerTest.sh"

# shellcheck disable=SC2034
{
    nonce=0
    chainBridgeBank=""
    chainBridgeRegistry=""
    chainMultisignAddr=""
    chainBtyERC20TokenAddr="1111111111111111111114oLvT2"

    chainUSDTBridgeTokenAddr=""
    chainUSDTBridgeTokenAddrOnETH=""
    chainUSDTBridgeTokenAddrOnBSC=""

    chainMainBridgeTokenAddr=""
    chainMainBridgeTokenAddrETH=""
    chainMainBridgeTokenAddrBNB=""

    ethereumBridgeBank=""
    ethereumBridgeRegistry=""
    ethereumMultisignAddr=""
    ethereumUSDTERC20TokenAddr=""
    ethereumBtyBridgeTokenAddr=""

    ethereumBridgeBankOnETH=""
    ethereumBridgeRegistryOnETH=""
    ethereumMultisignAddrOnETH=""
    ethereumUSDTERC20TokenAddrOnETH=""
    ethereumBtyBridgeTokenAddrOnETH=""

    ethereumBridgeBankOnBSC=""
    ethereumBridgeRegistryOnBSC=""
    ethereumMultisignAddrOnBSC=""
    ethereumUSDTERC20TokenAddrOnBSC=""
    ethereumBtyBridgeTokenAddrOnBSC=""

    # ETH 部署合约者的私钥 用于部署合约时签名使用
    ethDeployAddr="0x8AFDADFC88a1087c9A1D6c0F5Dd04634b87F303a"
    ethDeployKey="0x8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230"

    # chain 部署合约者的私钥 用于部署合约时签名使用
    chainDeployAddr="1JxhYLYsrscjTaQfaMoVUrnSdrejP7XRQD"
    chainDeployKey="0x9ef82623a5e9aac58d3a6b06392af66ec77289522b28896aed66abaaede66903"

    # eth 验证者私钥
    ethValidatorAddra="0x92C8b16aFD6d423652559C6E266cBE1c29Bfd84f"
    ethValidatorAddrb="0x0df9a824699bc5878232c9e612fe1a5346a5a368"
    ethValidatorAddrc="0xcb074cb21cdddf3ce9c3c0a7ac4497d633c9d9f1"
    ethValidatorAddrd="0xd9dab021e74ecf475788ed7b61356056b2095830"
    ethValidatorAddrKeya="3fa21584ae2e4fd74db9b58e2386f5481607dfa4d7ba0617aaa7858e5025dc1e"
    ethValidatorAddrKeyb="a5f3063552f4483cfc20ac4f40f45b798791379862219de9e915c64722c1d400"
    ethValidatorAddrKeyc="bbf5e65539e9af0eb0cfac30bad475111054b09c11d668fc0731d54ea777471e"
    ethValidatorAddrKeyd="c9fa31d7984edf81b8ef3b40c761f1847f6fcd5711ab2462da97dc458f1f896b"

    # 新增地址 chain 需要导入地址 转入 10 bty当收费费
    chainValidatora="1N6HstkyLFS8QCeVfdvYxx1xoryXoJtvvZ"
    chainValidatorb="155ooMPBTF8QQsGAknkK7ei5D78rwDEFe6"
    chainValidatorc="13zBdQwuyDh7cKN79oT2odkxYuDbgQiXFv"
    chainValidatord="113ZzVamKfAtGt9dq45fX1mNsEoDiN95HG"
    chainValidatorKeya="0x027ca96466c71c7e7c5d73b7e1f43cb889b3bd65ebd2413eefd31c6709c262ae"
    chainValidatorKeyb="0x9d539bc5fd084eb7fe86ad631dba9aa086dba38418725c38d9751459f567da66"
    chainValidatorKeyc="0x0a6671f101e30a2cc2d79d77436b62cdf2664ed33eb631a9c9e3f3dd348a23be"
    chainValidatorKeyd="0x3818b257b05ee75b6e43ee0e3cfc2d8502342cf67caed533e3756966690b62a5"

    ethTestAddr1=0xbc333839E37bc7fAAD0137aBaE2275030555101f
    ethTestAddrKey1=0x0c61f5a879d70807686e43eccc1f52987a15230ae0472902834af4d1933674f2
    ethTestAddr2=0x495953A743ef169EC5D4aC7b5F786BF2Bd56aFd5
    ethTestAddrKey2=0x2809477ede1261da21270096776ba7dc68b89c9df5f029965eaa5fe7f0b80697

    ethReceiverAddr1="0xa4ea64a583f6e51c3799335b28a8f0529570a635"
    #ethReceiverAddrKey1="355b876d7cbcb930d5dfab767f66336ce327e082cbaa1877210c1bae89b1df71"

    chainTestAddr1="1Cj1rqUenPmkeD6A8MGEzkBKQFN2H9yL3x"
    chainTestAddrKey1="0x7269a7a87d476310da37a9ca1ddc9333c9d7a0dfe1f2998b84758843a895433b"
    chainTestAddr2="1BCGLhdcdthNutQowV2YShuuN9fJRRGLxu"
    chainTestAddrKey2="0xb74acfd4eebbbd07bcae212baa7f094235ab8dc04f2f1d828681477b98b24008"

    chainReceiverAddr="12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv"
    chainReceiverAddrKey="4257d8692ef7fe13c68b65d6a52f03933db2fa5ce8faf210b5b8b80c721ced01"

    paraMainAddr="12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv"
    paraMainAddrKey="4257d8692ef7fe13c68b65d6a52f03933db2fa5ce8faf210b5b8b80c721ced01"
}

function chain_offline_send_evm() {
    # shellcheck disable=SC2154
    result=$(${EvmxgoBoss4xCLI} chain offline send -f "${1}")
    hash=$(echo "${result}" | jq -r ".[0].TxHash")
    # shellcheck disable=SC2154
    check_tx "${ChainCli}" "${hash}"
}

# $1 symbol $2 bridgeTokenAddr
function updateConfig() {
    local symbol=$1
    local bridgeTokenAddr=$2
    tx=$(curl -s --data-binary '{"jsonrpc":"2.0","id":2,"method":"Chain.CreateTransaction","params":[{"execer":"manage","actionName":"Modify","payload":{"key":"evmxgo-mint-'"${symbol}"'","value":"{\"address\":\"'"${bridgeTokenAddr}"'\",\"precision\":8,\"introduction\":\"symbol:'"${symbol}"', bridgeTokenAddr:'"${bridgeTokenAddr}"'\"}","op":"add","addr":""}}]}' -H 'content-type:text/plain;' "http://${docker_chain_ip}:8901" | jq -r ".result")
    if [ "${tx}" == "" ]; then
        echo -e "${RED}update config create tx 1${NOC}"
        exit 1
    fi

    sign=$(${ChainCli} wallet sign -k "$paraMainAddrKey" -d "${tx}")
    hash=$(${ChainCli} wallet send -d "${sign}")
    check_tx "${ChainCli}" "${hash}"
}

function configbridgevmxgoAddr() {
    local bridgevmxgoAddr=$1
    tx=$(curl -s --data-binary '{"jsonrpc":"2.0","id":2,"method":"Chain.CreateTransaction","params":[{"execer":"manage","actionName":"Modify","payload":{"key":"bridgevmxgo-contract-addr","value":"{\"address\":\"'"${bridgevmxgoAddr}"'\"}","op":"add","addr":""}}]}' -H 'content-type:text/plain;' "http://${docker_chain_ip}:8901" | jq -r ".result")
    if [ "${tx}" == "" ]; then
        echo -e "${RED}update config create tx 1${NOC}"
        exit 1
    fi

    sign=$(${ChainCli} wallet sign -k "$paraMainAddrKey" -d "${tx}")
    hash=$(${ChainCli} wallet send -d "${sign}")
    check_tx "${ChainCli}" "${hash}"
}

function DeployEvmxgo() {
    # 在 chain 上部署合约
    # shellcheck disable=SC2154
    ${EvmxgoBoss4xCLI} chain offline create -f 1 -k "${chainDeployKey}" -n "deploy crossx to chain" -r "${chainDeployAddr}, [${chainValidatora}, ${chainValidatorb}, ${chainValidatorc}, ${chainValidatord}], [96, 1, 1, 1]"
    result=$(${EvmxgoBoss4xCLI} chain offline send -f "deployBridgevmxgo2Chain.txt")

    for i in {0..6}; do
        hash=$(echo "${result}" | jq -r ".[$i].TxHash")
        check_tx "${ChainCli}" "${hash}"
    done
    XgoBridgeRegistryOnChain=$(echo "${result}" | jq -r ".[6].ContractAddr")
    XgoChainOracle=$(echo "${result}" | jq -r ".[2].ContractAddr")

    # 拷贝 BridgeRegistry.abi 和 BridgeBank.abi
    cp XgoBridgeRegistryOnChain.abi "${XgoBridgeRegistryOnChain}.abi"
    XgoChainBridgeBank=$(${ChainCli} evm query -c "${chainDeployAddr}" -b "bridgeBank()" -a "${XgoBridgeRegistryOnChain}")
    cp XgoChainBridgeBank.abi "${XgoChainBridgeBank}.abi"
}

function set_config_ethereum() {
    ${EvmxgoBoss4xCLI} chain offline create_add_lock_list -s ETH -t "${chainMainBridgeTokenAddr}" -c "${XgoChainBridgeBank}" -k "${chainDeployKey}" -f 1
    chain_offline_send_evm "create_add_lock_list.txt"

    ${EvmxgoBoss4xCLI} chain offline create_add_lock_list -s USDT -t "${chainUSDTBridgeTokenAddr}" -c "${XgoChainBridgeBank}" -k "${chainDeployKey}" -f 1
    chain_offline_send_evm "create_add_lock_list.txt"

    updateConfig "ETH" "${chainMainBridgeTokenAddr}"
    updateConfig "USDT" "${chainUSDTBridgeTokenAddr}"
    configbridgevmxgoAddr "${XgoChainBridgeBank}"
}

function set_config_bsc() {
    ${EvmxgoBoss4xCLI} chain offline create_add_lock_list -s BNB -t "${chainMainBridgeTokenAddr}" -c "${XgoChainBridgeBank}" -k "${chainDeployKey}" -f 1
    chain_offline_send_evm "create_add_lock_list.txt"

    ${EvmxgoBoss4xCLI} chain offline create_add_lock_list -s BUSDT -t "${chainUSDTBridgeTokenAddr}" -c "${XgoChainBridgeBank}" -k "${chainDeployKey}" -f 1
    chain_offline_send_evm "create_add_lock_list.txt"

    updateConfig "BNB" "${chainMainBridgeTokenAddr}"
    updateConfig "BUSDT" "${chainUSDTBridgeTokenAddr}"
    configbridgevmxgoAddr "${XgoChainBridgeBank}"
}

function TestETH2EVMToChain() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    # 查询 ETH 这端 bridgeBank 地址原来是 0
    result=$(${CLIA} ethereum balance -o "${ethereumBridgeBank}")
    cli_ret "${result}" "balance" ".balance" "0"

    # ETH 这端 lock 11个
    result=$(${CLIA} ethereum lock -m 11 -k "${ethTestAddrKey1}" -r "${chainReceiverAddr}")
    cli_ret "${result}" "lock"

    # eth 等待 2 个区块
    sleep 4

    # 查询 ETH 这端 bridgeBank 地址 11
    result=$(${CLIA} ethereum balance -o "${ethereumBridgeBank}")
    cli_ret "${result}" "balance" ".balance" "11"

    sleep "${maturityDegree}"

    # chain chainMainBridgeTokenAddr（ETH合约中）查询 lock 金额
    result=$(${ChainCli} evm query -a "${chainMainBridgeTokenAddr}" -c "${chainDeployAddr}" -b "balanceOf(${chainReceiverAddr})")
    # 结果是 11 * le8
    is_equal "${result}" "1100000000"

    ${EvmxgoBoss4xCLI} chain offline approve_erc20 -a 330000000000 -s "${XgoChainBridgeBank}" -c "${chainMainBridgeTokenAddr}" -k "${chainReceiverAddrKey}" -f 1
    chain_offline_send_evm "approve_erc20.txt"

    hash=$(${ChainCli} send evm call -f 1 -k "${chainReceiverAddr}" -e "${XgoChainBridgeBank}" -p "lock(${chainTestAddr2}, ${chainMainBridgeTokenAddr}, 500000000)")
    check_tx "${ChainCli}" "${hash}"

    result=$(${ChainCli} evm query -a "${chainMainBridgeTokenAddr}" -c "${chainDeployAddr}" -b "balanceOf(${chainReceiverAddr})")
    is_equal "${result}" "600000000"

    result=$(${ChainCli} evm query -a "${chainMainBridgeTokenAddr}" -c "${chainDeployAddr}" -b "balanceOf(${XgoChainBridgeBank})")
    is_equal "${result}" "500000000"

    let nonce=nonce+1
    hash=$(${EvmxgoBoss4xCLI} chain burn_xgo -m "300000000" -f "${chainTestAddr2}" -r "${chainTestAddr2}" -o "${XgoChainOracle}" -n "${nonce}" -s "$1" -t "${chainMainBridgeTokenAddr}" -k "${chainValidatorKeya}")
    check_tx "${ChainCli}" "${hash}"

    result=$(${ChainCli} evm query -a "${chainMainBridgeTokenAddr}" -c "${chainDeployAddr}" -b "balanceOf(${XgoChainBridgeBank})")
    is_equal "${result}" "200000000"

    result=$(${ChainCli} evm query -a "${chainMainBridgeTokenAddr}" -c "${chainDeployAddr}" -b "balanceOf(${chainTestAddr2})")
    is_equal "${result}" "300000000"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function Testethereum2EVMToChain_usdt() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    # 查询 ETH 这端 bridgeBank 地址原来是
    result=$(${CLIA} ethereum balance -o "${ethereumBridgeBank}" -t "${ethereumUSDTERC20TokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # ETH 这端 lock 12个
    result=$(${CLIA} ethereum lock -m 12 -k "${ethTestAddrKey1}" -r "${chainReceiverAddr}" -t "${ethereumUSDTERC20TokenAddr}")
    cli_ret "${result}" "lock"

    # eth 等待 2 个区块
    sleep 4

    # 查询 ETH 这端 bridgeBank 地址 12
    result=$(${CLIA} ethereum balance -o "${ethereumBridgeBank}" -t "${ethereumUSDTERC20TokenAddr}")
    cli_ret "${result}" "balance" ".balance" "12"

    sleep "${maturityDegree}"

    # chain chainMainBridgeTokenAddr（ETH合约中）查询 lock 金额
    result=$(${ChainCli} evm query -a "${chainUSDTBridgeTokenAddr}" -c "${chainTestAddr1}" -b "balanceOf(${chainReceiverAddr})")
    is_equal "${result}" "1200000000"

    ${EvmxgoBoss4xCLI} chain offline approve_erc20 -a 330000000000 -s "${XgoChainBridgeBank}" -c "${chainUSDTBridgeTokenAddr}" -k "${chainReceiverAddrKey}" -f 1
    chain_offline_send_evm "approve_erc20.txt"

    hash=$(${ChainCli} send evm call -f 1 -k "${chainReceiverAddr}" -e "${XgoChainBridgeBank}" -p "lock(${chainTestAddr2}, ${chainUSDTBridgeTokenAddr}, 500000000)")
    check_tx "${ChainCli}" "${hash}"

    result=$(${ChainCli} evm query -a "${chainUSDTBridgeTokenAddr}" -c "${chainDeployAddr}" -b "balanceOf(${chainReceiverAddr})")
    is_equal "${result}" "700000000"

    result=$(${ChainCli} evm query -a "${chainUSDTBridgeTokenAddr}" -c "${chainDeployAddr}" -b "balanceOf(${XgoChainBridgeBank})")
    is_equal "${result}" "500000000"

    let nonce=nonce+1
    hash=$(${EvmxgoBoss4xCLI} chain burn_xgo -m "300000000" -f "${chainTestAddr2}" -r "${chainTestAddr2}" -o "${XgoChainOracle}" -n "${nonce}" -s "$1" -t "${chainUSDTBridgeTokenAddr}" -k "${chainValidatorKeya}")
    check_tx "${ChainCli}" "${hash}"

    result=$(${ChainCli} evm query -a "${chainUSDTBridgeTokenAddr}" -c "${chainDeployAddr}" -b "balanceOf(${XgoChainBridgeBank})")
    is_equal "${result}" "200000000"

    result=$(${ChainCli} evm query -a "${chainUSDTBridgeTokenAddr}" -c "${chainDeployAddr}" -b "balanceOf(${chainTestAddr2})")
    is_equal "${result}" "300000000"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function get_evm_cli() {
    # shellcheck disable=SC2034
    {
        paraName="user.p.para."
        # shellcheck disable=SC2154
        docker_chain_ip=$(get_docker_addr "${dockerNamePrefix}_chain_1")
        MainCli="./chain-cli --rpc_laddr http://${docker_chain_ip}:8801"
        Para8801Cli="./chain-cli --rpc_laddr http://${docker_chain_ip}:8901 --paraName ${paraName}"
        Para8901Cli="./chain-cli --rpc_laddr http://${docker_chain_ip}:8901 --paraName ${paraName}"

        docker_ebrelayera_ip=$(get_docker_addr "${dockerNamePrefix}_ebrelayera_1")
        CLIA="docker exec ${dockerNamePrefix}_ebrelayera_1 /root/ebcli_A"
        CLIB="docker exec ${dockerNamePrefix}_ebrelayerb_1 /root/ebcli_A"
        CLIC="docker exec ${dockerNamePrefix}_ebrelayerc_1 /root/ebcli_A"
        CLID="docker exec ${dockerNamePrefix}_ebrelayerd_1 /root/ebcli_A"

        docker_ganachetesteth_ip=$(get_docker_addr "${dockerNamePrefix}_ganachetesteth_1")
        docker_ganachetestbsc_ip=$(get_docker_addr "${dockerNamePrefix}_ganachetestbsc_1")
        Boss4xCLI="docker exec ${dockerNamePrefix}_ebrelayera_1 /root/boss4x --rpc_laddr http://${docker_chain_ip}:8901 --rpc_laddr_ethereum http://${docker_ganachetesteth_ip}:8545 --paraName ${paraName} --chainID ${chainID} --chainEthId 1337"

        Boss4xCLIeth="docker exec ${dockerNamePrefix}_ebrelayera_1 /root/boss4x --rpc_laddr http://${docker_chain_ip}:8901 --rpc_laddr_ethereum http://${docker_ganachetesteth_ip}:8545 --paraName ${paraName} --chainID ${chainID} --chainEthId 1337"
        Boss4xCLIbsc="docker exec ${dockerNamePrefix}_ebrelayera_1 /root/boss4x --rpc_laddr http://${docker_chain_ip}:8901 --rpc_laddr_ethereum http://${docker_ganachetestbsc_ip}:8545 --paraName ${paraName} --chainID ${chainID} --chainEthId 1337"

        CLIAeth="docker exec ${dockerNamePrefix}_ebrelayera_1 /root/ebcli_A --node_addr http://${docker_ganachetesteth_ip}:8545 --eth_chain_name Ethereum"
        CLIAbsc="docker exec ${dockerNamePrefix}_ebrelayera_1 /root/ebcli_A --node_addr http://${docker_ganachetestbsc_ip}:8545 --eth_chain_name Binance"

        EvmxgoBoss4xCLI="./evmxgoboss4x --rpc_laddr http://${docker_chain_ip}:8901 --paraName user.p.para. --chainID ${chainID}"
    }
}

function test_xgo() {
    TestETH2ChainAssets
    TestETH2ChainUSDT

    TestETH2EVMToChain "$1"
    Testethereum2EVMToChain_usdt "$2"
}

# shellcheck disable=SC2034
function test_evm_all() {
    # test
    Boss4xCLI=${Boss4xCLIeth}
    CLIA=${CLIAeth}
    ethereumBridgeBank="${ethereumBridgeBankOnETH}"
    chainMainBridgeTokenAddr="${chainMainBridgeTokenAddrETH}"
    ethereumBtyBridgeTokenAddr="${ethereumBtyBridgeTokenAddrOnETH}"
    ethereumUSDTERC20TokenAddr="${ethereumUSDTERC20TokenAddrOnETH}"
    chainUSDTBridgeTokenAddr="${chainUSDTBridgeTokenAddrOnETH}"
    set_config_ethereum
    test_xgo "ETH" "USDT"

    Boss4xCLI=${Boss4xCLIbsc}
    CLIA=${CLIAbsc}
    ethereumBridgeBank="${ethereumBridgeBankOnBSC}"
    chainMainBridgeTokenAddr="${chainMainBridgeTokenAddrBNB}"
    ethereumBtyBridgeTokenAddr="${ethereumBtyBridgeTokenAddrOnBSC}"
    ethereumUSDTERC20TokenAddr="${ethereumUSDTERC20TokenAddrOnBSC}"
    chainUSDTBridgeTokenAddr="${chainUSDTBridgeTokenAddrOnBSC}"
    set_config_bsc
    test_xgo "BNB" "BUSDT"
}

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

    get_evm_cli

    # init
    ChainCli=${MainCli}
    InitChainValidator
    # para add
    initPara

    ChainCli=${Para8901Cli}
    StartDockerRelayerDeploy

    DeployEvmxgo
    test_evm_all

    echo_addrs
    echo -e "${GRE}XgoChainOracle: ${XgoChainOracle} ${NOC}"
    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}
