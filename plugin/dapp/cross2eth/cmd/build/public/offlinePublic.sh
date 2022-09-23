#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck disable=SC2154
# shellcheck source=/dev/null
set -x
set +e

# 主要在平行链上测试

source "./publicTest.sh"

# $1 Key $2 addr $3 label $4 amount $5 evm amount
function ChainImportKey() {
    local key="${1}"
    local addr="${2}"
    local label="${3}"
    local amount="${4}"
    local evm_amount="${5}"
    # 转帐到 DeployAddr 需要手续费
    result=$(${ChainCli} account import_key -k "${key}" -l "${label}")
    check_addr "${result}" "${addr}"
    hash=$(${ChainCli} send coins transfer -a "${amount}" -n test -t "${addr}" -k 4257d8692ef7fe13c68b65d6a52f03933db2fa5ce8faf210b5b8b80c721ced01)
    check_tx "${ChainCli}" "${hash}"

    # 转账到 EVM  合约中
    hash=$(${ChainCli} send coins send_exec -e evm -a "${evm_amount}" -k "${addr}")
    check_tx "${ChainCli}" "${hash}"
    result=$(${ChainCli} account balance -a "${addr}" -e evm)
    #    balance_ret "${result}" "${evm_amount}.0000" # 平行链查询方式不一样 直接去掉金额匹配
}

# chian33 初始化准备
function InitChainValidator() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    # 转帐到 DeployAddr 需要手续费
    ChainImportKey "${chainDeployKey}" "${chainDeployAddr}" "DeployAddr" 2200 1000

    ChainImportKey "${chainTestAddrKey1}" "${chainTestAddr1}" "cross2ethAddr1" 2200 1000
    ChainImportKey "${chainTestAddrKey2}" "${chainTestAddr2}" "cross2ethAddr2" 2200 1000

    # 导入 chainValidators 私钥生成地址
    for name in a b c d p sp; do
        eval chainValidatorKey=\$chainValidatorKey${name}
        eval chainValidator=\$chainValidator${name}
        result=$(${ChainCli} account import_key -k "${chainValidatorKey}" -l validator$name)
        check_addr "${result}" "${chainValidator}"

        # chainValidator 要有手续费
        hash=$(${ChainCli} send coins transfer -a 100 -t "${chainValidator}" -k 4257d8692ef7fe13c68b65d6a52f03933db2fa5ce8faf210b5b8b80c721ced01)
        check_tx "${ChainCli}" "${hash}"
        result=$(${ChainCli} account balance -a "${chainValidator}" -e coins)
        #        balance_ret "${result}" "100.0000" # 平行链查询方式不一样 直接去掉金额匹配
    done

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function coins_cross_transfer() {
    local key="${1}"
    local addr="${2}"
    local amount="${3}"
    local para_amount="${4}"
    local evm_amount="${5}"
    # 先把 bty 转入到 paracross 合约中
    hash=$(${MainCli} send coins send_exec -e paracross -a "${amount}" -k "${key}")
    check_tx "${MainCli}" "${hash}"

    # 主链中的 bty 夸链到 平行链中
    hash=$(${Para8801Cli} send para cross_transfer -a "${para_amount}" -e coins -s bty -t "${addr}" -k "${key}")
    check_tx "${Para8801Cli}" "${hash}"
    check_tx "${Para8901Cli}" "${hash}"
    result=$(${Para8901Cli} asset balance -a "${addr}" --asset_exec paracross --asset_symbol coins.bty | jq -r .balance)
    is_equal "${result}" "${para_amount}.0000"

    # 把平行链中的 bty 转入 平行链中的 evm 合约
    hash=$(${Para8901Cli} send para transfer_exec -a "${evm_amount}" -e user.p.para.evm -s coins.bty -k "${key}")
    check_tx "${Para8901Cli}" "${hash}"
    result=$(${Para8901Cli} asset balance -a "${addr}" --asset_exec paracross --asset_symbol coins.bty -e user.p.para.evm | jq -r .balance)
    is_equal "${result}" "${evm_amount}.0000"
}

function initPara() {
    # para add
    hash=$(${Para8901Cli} send coins transfer -a 10000 -n test -t "${chainReceiverAddr}" -k CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944)
    check_tx "${Para8901Cli}" "${hash}"

    ChainCli=${Para8901Cli}
    InitChainValidator

    coins_cross_transfer "${chainDeployKey}" "${chainDeployAddr}" 1000 800 500
    coins_cross_transfer "${chainTestAddrKey1}" "${chainTestAddr1}" 1000 800 500
    coins_cross_transfer "${chainTestAddrKey2}" "${chainTestAddr2}" 1000 800 500

    # 平行链共识节点增加测试币
    ${MainCli} send coins transfer -a 1000 -n test -t "1KSBd17H7ZK8iT37aJztFB22XGwsPTdwE4" -k "${chainReceiverAddrKey}"
    ${MainCli} send coins transfer -a 1000 -n test -t "1JRNjdEqp4LJ5fqycUBm9ayCKSeeskgMKR" -k "${chainReceiverAddrKey}"
    ${MainCli} send coins transfer -a 1000 -n test -t "1NLHPEcbTWWxxU3dGUZBhayjrCHD3psX7k" -k "${chainReceiverAddrKey}"
    ${MainCli} send coins transfer -a 1000 -n test -t "1MCftFynyvG2F4ED5mdHYgziDxx6vDrScs" -k "${chainReceiverAddrKey}"
}

function initMultisignChainAddr() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    for name in A B C D; do
        eval chainMultisignKey=\$chainMultisignKey${name}
        eval chainMultisign=\$chainMultisign${name}
        result=$(${ChainCli} account import_key -k "${chainMultisignKey}" -l multisignAddr$name)
        check_addr "${result}" "${chainMultisign}"

        # chainMultisign 要有手续费
        hash=$(${ChainCli} send coins transfer -a 10 -t "${chainMultisign}" -k "${chainDeployAddr}")
        check_tx "${ChainCli}" "${hash}"
        result=$(${ChainCli} account balance -a "${chainMultisign}" -e coins)
        balance_ret "${result}" "10.0000"
    done

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function transferChainMultisignFee() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    # chainMultisignAddr 要有手续费
    hash=$(${ChainCli} send coins transfer -a 10 -t "${chainMultisignAddr}" -k "${chainDeployAddr}")
    check_tx "${ChainCli}" "${hash}"
    result=$(${ChainCli} account balance -a "${chainMultisignAddr}" -e coins)
    balance_ret "${result}" "10.0000"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# lock eth 判断是否转入多签地址金额是否正确
function lock_eth_multisign() {
    local lockAmount=$1
    result=$(${CLIA} ethereum lock -m "${lockAmount}" -k "${ethTestAddrKey1}" -r "${chainReceiverAddr}")
    cli_ret "${result}" "lock"

    if [[ $# -eq 3 ]]; then
        local bridgeBankBalance=$2
        local multisignBalance=$3
        # eth 等待 2 个区块
        sleep 10
        #        eth_block_wait 2

        result=$(${CLIA} ethereum balance -o "${ethereumBridgeBank}")
        cli_ret "${result}" "balance" ".balance" "${bridgeBankBalance}"
        result=$(${CLIA} ethereum balance -o "${ethereumMultisignAddr}")
        cli_ret "${result}" "balance" ".balance" "${multisignBalance}"
    fi
}

function lock_ethereum_usdt_multisign() {
    local lockAmount=$1
    result=$(${CLIA} ethereum lock -m "${lockAmount}" -k "${ethTestAddrKey1}" -r "${chainReceiverAddr}" -t "${ethereumUSDTERC20TokenAddr}")
    cli_ret "${result}" "lock"

    if [[ $# -eq 3 ]]; then
        local bridgeBankBalance=$2
        local multisignBalance=$3

        # eth 等待 2 个区块
        sleep 10

        result=$(${CLIA} ethereum balance -o "${ethereumBridgeBank}" -t "${ethereumUSDTERC20TokenAddr}")
        cli_ret "${result}" "balance" ".balance" "${bridgeBankBalance}"
        result=$(${CLIA} ethereum balance -o "${ethereumMultisignAddr}" -t "${ethereumUSDTERC20TokenAddr}")
        cli_ret "${result}" "balance" ".balance" "${multisignBalance}"
    fi
}

# 检查交易是否执行成功 $1:交易hash
function check_eth_tx() {
    local tx=${1}
    ty=$(${CLIA} ethereum receipt -s "${tx}" | jq .status | sed 's/\"//g')
    if [[ ${ty} != 0x1 ]]; then
        echo -e "${RED}check eth tx error, hash is ${tx}${NOC}"
        exit_cp_file
    fi
}

# $1 send file
function chain_offline_send() {
    result=$(${Boss4xCLI} chain offline send -f "${1}")
    hash=$(echo "${result}" | jq -r ".[0].TxHash")
    check_tx "${ChainCli}" "${hash}"
}

# $1 send file $2 key
function ethereum_offline_sign_send() {
    key="${ethDeployKey}"
    if [[ ${2} != "" ]]; then
        key="${2}"
    fi
    ${Boss4xCLI} ethereum offline sign -f "${1}" -k "${key}"
    result=$(${Boss4xCLI} ethereum offline send -f "deploysigntxs.txt")
    hash=$(echo "${result}" | jq -r ".[0].TxHash")
    check_eth_tx "${hash}"
}

function OfflineDeploy_chain() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    # 在 chain 上部署合约
    #    ${Boss4xCLI} chain offline create -f 1 -k "${chainDeployKey}" -n "deploy crossx to chain" -r "${chainDeployAddr}, [${chainValidatora}, ${chainValidatorb}, ${chainValidatorc}, ${chainValidatord}], [25, 25, 25, 25]" -m "${chainMultisignA},${chainMultisignB},${chainMultisignC},${chainMultisignD}"
    ${Boss4xCLI} chain offline create_file -f 1 -k "${chainDeployKey}" -n "deploy crossx to chain" -c "./deploy_chain.toml"
    result=$(${Boss4xCLI} chain offline send -f "deployCrossX2Chain.txt")

    for i in {0..9}; do
        hash=$(echo "${result}" | jq -r ".[$i].TxHash")
        check_tx "${ChainCli}" "${hash}"
    done
    chainBridgeRegistry=$(echo "${result}" | jq -r ".[6].ContractAddr")
    # shellcheck disable=SC2034
    chainMultisignAddr=$(echo "${result}" | jq -r ".[7].ContractAddr")

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function OfflineDeploy_ethereum() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    local deployfile=$1
    # 在 Eth 上部署合约
    ${Boss4xCLI} ethereum offline create_file -f "${deployfile}"
    ${Boss4xCLI} ethereum offline sign -k "${ethDeployKey}"
    result=$(${Boss4xCLI} ethereum offline send -f "deploysigntxs.txt")
    for i in {0..7}; do
        hash=$(echo "${result}" | jq -r ".[$i].TxHash")
        check_eth_tx "${hash}"
    done
    ethereumBridgeBank=$(echo "${result}" | jq -r ".[3].ContractAddr")
    ethereumBridgeRegistry=$(echo "${result}" | jq -r ".[7].ContractAddr")
    # shellcheck disable=SC2034
    ethereumMultisignAddr=$(echo "${result}" | jq -r ".[8].ContractAddr")

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function OfflineDeploy() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    OfflineDeploy_chain
    # 修改 relayer.toml 字段
    sed -i 's/^BridgeRegistryOnChain=.*/BridgeRegistryOnChain="'"${chainBridgeRegistry}"'"/g' "./relayer.toml"

    {
        Boss4xCLI=${Boss4xCLIeth}
        CLIA=${CLIAeth}
        OfflineDeploy_ethereum "./deploy_ethereum.toml"
        ethereumBridgeBankOnETH="${ethereumBridgeBank}"
        ethereumBridgeRegistryOnETH="${ethereumBridgeRegistry}"
        ethereumMultisignAddrOnETH="${ethereumMultisignAddr}"

        sed -i '14,21s/BridgeRegistry=.*/BridgeRegistry="'"${ethereumBridgeRegistryOnETH}"'"/g' "./relayer.toml"

        Boss4xCLI=${Boss4xCLIbsc}
        CLIA=${CLIAbsc}
        cp "./deploy_ethereum.toml" "./deploy_bsc.toml"
        sed -i 's/^symbol=.*/symbol="BNB"/g' "./deploy_bsc.toml"
        docker cp "./deploy_bsc.toml" "${dockerNamePrefix}_ebrelayera_1":/root/deploy_bsc.toml
        OfflineDeploy_ethereum "./deploy_bsc.toml"
        ethereumBridgeBankOnBSC="${ethereumBridgeBank}"
        ethereumBridgeRegistryOnBSC="${ethereumBridgeRegistry}"
        ethereumMultisignAddrOnBSC="${ethereumMultisignAddr}"

        sed -i '23,30s/BridgeRegistry=.*/BridgeRegistry="'"${ethereumBridgeRegistryOnBSC}"'"/g' "./relayer.toml"
    }

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# init $1 CLI $2 pwd $3 chainValidatorKey $4 ethValidatorAddrKey
function init_validator_relayer() {
    local CLI=$1
    local pwd=$2
    local chainValidatorKey=$3
    local ethValidatorAddrKey=$4
    result=$(${CLI} set_pwd -p "${pwd}")
    cli_ret "${result}" "set_pwd"

    result=$(${CLI} unlock -p "${pwd}")
    cli_ret "${result}" "unlock"

    sleep 20

    result=$(${CLI} chain import_privatekey -k "${chainValidatorKey}")
    cli_ret "${result}" "chain import_privatekey"

    result=$(${CLI} ethereum import_privatekey -k "${ethValidatorAddrKey}")
    cli_ret "${result}" "ethereum import_privatekey"
}

# shellcheck disable=SC2120
function InitRelayerA() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    init_validator_relayer "${CLIA}" "${validatorPwd}" "${chainValidatorKeya}" "${ethValidatorAddrKeya}"

    ${CLIA} chain multisign set_multiSign -a "${chainMultisignAddr}"

    # 拷贝 BridgeRegistry.abi 和 BridgeBank.abi
    cp BridgeRegistry.abi "${chainBridgeRegistry}.abi"
    chainBridgeBank=$(${ChainCli} evm query -c "${chainDeployAddr}" -b "bridgeBank()" -a "${chainBridgeRegistry}")
    cp ChainBridgeBank.abi "${chainBridgeBank}.abi"

    ${CLIAeth} ethereum multisign set_multiSign -a "${ethereumMultisignAddrOnETH}"
    ${CLIAbsc} ethereum multisign set_multiSign -a "${ethereumMultisignAddrOnBSC}"

    # 拷贝 BridgeRegistry.abi 和 BridgeBank.abi
    cp BridgeRegistry.abi "${ethereumBridgeRegistryOnETH}.abi"
    cp EthBridgeBank.abi "${ethereumBridgeBankOnETH}.abi"
    cp BridgeRegistry.abi "${ethereumBridgeRegistryOnBSC}.abi"
    cp EthBridgeBank.abi "${ethereumBridgeBankOnBSC}.abi"
    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function offline_deploy_erc20_create_tether_usdt_USDT() {
    # eth 上 铸币 USDT
    local name=$1
    echo -e "${GRE}======= 在 ethereum 上创建 ERC20 ${name} ======${NOC}"
    ${Boss4xCLI} ethereum offline create_tether_usdt -m 33000000000000000000 -s "${name}" -d "${ethTestAddr1}"
    ${Boss4xCLI} ethereum offline sign -f "deployTetherUSDT.txt" -k "${ethTestAddrKey1}"
    sleep 10
    result=$(${Boss4xCLI} ethereum offline send -f "deploysigntxs.txt")
    hash=$(echo "${result}" | jq -r ".[0].TxHash")
    check_eth_tx "${hash}"
    ethereumUSDTERC20TokenAddr=$(echo "${result}" | jq -r ".[0].ContractAddr")
    ${Boss4xCLI} ethereum offline create_add_lock_list -s "${name}" -t "${ethereumUSDTERC20TokenAddr}" -c "${ethereumBridgeBank}" -d "${ethDeployAddr}"
    ethereum_offline_sign_send "create_add_lock_list.txt"
}

function offline_create_bridge_token_eth_BTY() {
    # 在 Ethereum 上创建 bridgeToken BTY
    echo -e "${GRE}======= 在 Ethereum 上创建 bridgeToken BTY ======${NOC}"
    ${Boss4xCLI} ethereum offline create_bridge_token -s BTY -c "${ethereumBridgeBank}" -d "${ethDeployAddr}"
    ethereum_offline_sign_send "create_bridge_token.txt"

    # shellcheck disable=SC2034
    ethereumBtyBridgeTokenAddr=$(${CLIA} ethereum receipt -s "${hash}" | jq -r .logs[0].address)
}

function offline_create_bridge_token_chain_symbol() {
    # 在 chain 上创建 bridgeToken ETH
    echo -e "${GRE}======= 在 chain 上创建 bridgeToken $1 ======${NOC}"
    local symbolName="$1"
    ${Boss4xCLI} chain offline create_bridge_token -c "${chainBridgeBank}" -s "${symbolName}" -k "${chainDeployKey}" -n "create_bridge_token:${symbolName}"
    chain_offline_send "create_bridge_token.txt"

    chainMainBridgeTokenAddr=$(${ChainCli} evm query -a "${chainBridgeBank}" -c "${chainDeployAddr}" -b "getToken2address(${symbolName})")
    echo "${symbolName} Token Addr= ${chainMainBridgeTokenAddr}"
    cp BridgeToken.abi "${chainMainBridgeTokenAddr}.abi"

    result=$(${ChainCli} evm query -a "${chainMainBridgeTokenAddr}" -c "${chainMainBridgeTokenAddr}" -b "symbol()")
    is_equal "${result}" "${symbolName}"
}

function offline_transfer_multisign_Eth_test() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    # transfer
    ${Boss4xCLI} ethereum offline multisign_transfer_prepare -a 3 -r "${ethereumBridgeBank}" -c "${ethereumMultisignAddr}" -d "${ethTestAddr1}"
    ${Boss4xCLI} ethereum offline sign_multisign_tx -k "${ethMultisignKeyB},${ethMultisignKeyC},${ethMultisignKeyD}"
    ${Boss4xCLI} ethereum offline send_multisign_tx -f sign_multisign_tx.txt -k "${ethTestAddrKey1}"
    sleep 10

    result=$(${CLIA} ethereum balance -o "${ethereumBridgeBank}")
    cli_ret "${result}" "balance" ".balance" "16"
    result=$(${CLIA} ethereum balance -o "${ethereumMultisignAddr}")
    cli_ret "${result}" "balance" ".balance" "20"

    # transfer
    ${Boss4xCLI} ethereum offline multisign_transfer_prepare -a 5 -r "${ethMultisignA}" -c "${ethereumMultisignAddr}" -d "${ethTestAddr1}"
    ${Boss4xCLI} ethereum offline sign_multisign_tx -k "${ethMultisignKeyA},${ethMultisignKeyC},${ethMultisignKeyD}"
    ${Boss4xCLI} ethereum offline send_multisign_tx -f sign_multisign_tx.txt -k "${ethTestAddrKey1}"

    result=$(${CLIA} ethereum balance -o "${ethMultisignA}")
    cli_ret "${result}" "balance" ".balance" "1005"
    result=$(${CLIA} ethereum balance -o "${ethereumMultisignAddr}")
    cli_ret "${result}" "balance" ".balance" "15"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function offline_transfer_multisign_EthUSDT() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    # transfer
    ${Boss4xCLI} ethereum offline multisign_transfer_prepare -a 8 -r "${ethereumBridgeBank}" -c "${ethereumMultisignAddr}" -d "${ethTestAddr1}" -t "${ethereumUSDTERC20TokenAddr}"
    ${Boss4xCLI} ethereum offline sign_multisign_tx -k "${ethMultisignKeyA},${ethMultisignKeyB},${ethMultisignKeyD}"
    ${Boss4xCLI} ethereum offline send_multisign_tx -f sign_multisign_tx.txt -k "${ethTestAddrKey1}"

    result=$(${CLIA} ethereum balance -o "${ethereumBridgeBank}" -t "${ethereumUSDTERC20TokenAddr}")
    cli_ret "${result}" "balance" ".balance" "80"
    result=$(${CLIA} ethereum balance -o "${ethereumMultisignAddr}" -t "${ethereumUSDTERC20TokenAddr}")
    cli_ret "${result}" "balance" ".balance" "80"

    # transfer
    ${Boss4xCLI} ethereum offline multisign_transfer_prepare -a 10 -r "${ethMultisignA}" -c "${ethereumMultisignAddr}" -d "${ethTestAddr1}" -t "${ethereumUSDTERC20TokenAddr}"
    ${Boss4xCLI} ethereum offline sign_multisign_tx -k "${ethMultisignKeyA},${ethMultisignKeyB},${ethMultisignKeyC}"
    ${Boss4xCLI} ethereum offline send_multisign_tx -f sign_multisign_tx.txt -k "${ethTestAddrKey1}"

    result=$(${CLIA} ethereum balance -o "${ethMultisignA}" -t "${ethereumUSDTERC20TokenAddr}")
    cli_ret "${result}" "balance" ".balance" "10"
    result=$(${CLIA} ethereum balance -o "${ethereumMultisignAddr}" -t "${ethereumUSDTERC20TokenAddr}")
    cli_ret "${result}" "balance" ".balance" "70"
    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}
