#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null
set -x
set +e

# 主要在平行链上测试

source "./mainPubilcRelayerTest.sh"
Yellow="\[\033[0;33m\]"
le8=100000000

function start_docker_ebrelayerProxy() {
    updata_toml proxy
    # 代理转账中继器中的标志位ProcessWithDraw设置为true
    sed -i 's/^ProcessWithDraw=.*/ProcessWithDraw=true/g' "./relayerproxy.toml"

    # shellcheck disable=SC2154
    docker cp "./relayerproxy.toml" "${dockerNamePrefix}_ebrelayerproxy_1":/root/relayer.toml
    start_docker_ebrelayer "${dockerNamePrefix}_ebrelayerproxy_1" "/root/ebrelayer" "./ebrelayerproxy.log"
    sleep 1

    # shellcheck disable=SC2154
    init_validator_relayer "${CLIP}" "${validatorPwd}" "${chain33ValidatorKeyp}" "${ethValidatorAddrKeyp}"
}

function setWithdraw() {
    result=$(${CLIP} ethereum cfgWithdraw -f 1 -s ETH -a 100 -d 18)
    cli_ret "${result}" "cfgWithdraw"
    result=$(${CLIP} ethereum cfgWithdraw -f 1 -s USDT -a 100 -d 6)
    cli_ret "${result}" "cfgWithdraw"

    # 在chain33上的bridgeBank合约中设置proxyReceiver
    # shellcheck disable=SC2154
    ${Boss4xCLI} chain33 offline set_withdraw_proxy -c "${chain33BridgeBank}" -a "${chain33Validatorsp}" -k "${chain33DeployKey}" -n "set_withdraw_proxy:${chain33Validatorsp}"
    chain33_offline_send "set_withdraw_proxy.txt"
}

# eth to chain33 在以太坊上锁定 ETH 资产,然后在 chain33 上 burn
function TestETH2Chain33Assets_proxy() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    echo -e "${GRE}=========== eth to chain33 在以太坊上锁定 ETH 资产,然后在 chain33 上 burn ===========${NOC}"

    echo -e "${Yellow} lockAmount1 ${NOC}"
    local lockAmount1=$1

    echo -e "${Yellow} ethBridgeBank 初始金额 ${NOC}"
    # shellcheck disable=SC2154
    ethBridgeBankBalancebf=$(${CLIA} ethereum balance -o "${ethBridgeBank}" | jq -r ".balance")

    echo -e "${Yellow} chain33ReceiverAddr chain33 端 lock 后接收地址初始金额 ${NOC}"
    # shellcheck disable=SC2154
    chain33RBalancebf=$(${Chain33Cli} evm query -a "${chain33EthBridgeTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")

    echo -e "${Yellow} chain33Validatorsp chain33 代理地址初始金额 ${NOC}"
    chain33VspBalancebf=$(${Chain33Cli} evm query -a "${chain33EthBridgeTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33Validatorsp})")

    echo -e "${Yellow} lock ${NOC}"
    # shellcheck disable=SC2154
    result=$(${CLIA} ethereum lock -m "${lockAmount1}" -k "${ethTestAddrKey1}" -r "${chain33ReceiverAddr}")
    cli_ret "${result}" "lock"

    # eth 等待 2 个区块
    sleep 4

    echo -e "${Yellow} ethBridgeBank lock 后金额 ${NOC}"
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}")
    # shellcheck disable=SC2219
    let ethBridgeBankBalanceEnd=${ethBridgeBankBalancebf}+${lockAmount1}
    cli_ret "${result}" "balance" ".balance" "${ethBridgeBankBalanceEnd}"

    # shellcheck disable=SC2086
    sleep "${maturityDegree}"

    # chain33 chain33EthBridgeTokenAddr（ETH合约中）查询 lock 金额
    echo -e "${Yellow} chain33ReceiverAddr chain33 端 lock 后接收地址 lock 后金额 ${NOC}"
    # shellcheck disable=SC2154
    result=$(${Chain33Cli} evm query -a "${chain33EthBridgeTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    # shellcheck disable=SC2219
    let chain33RBalancelock=${lockAmount1}*${le8}+${chain33RBalancebf}
#    is_equal "${result}" "${chain33RBalancelock}"

    echo -e "${Yellow} chain33Validatorsp chain33 代理地址 lock 后金额 ${NOC}"
    result=$(${Chain33Cli} evm query -a "${chain33EthBridgeTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33Validatorsp})")
#    is_equal "${result}" "${chain33VspBalancebf}"

    echo -e "${Yellow} ethTestAddr2 ethereum withdraw 接收地址初始金额 ${NOC}"
    # shellcheck disable=SC2154
    ethT2Balancebf=$(${CLIA} ethereum balance -o "${ethTestAddr2}" | jq -r ".balance")

    echo -e "${Yellow} ethValidatorAddrp ethereum 代理地址初始金额 ${NOC}"
    # shellcheck disable=SC2154
    ethPBalancebf=$(${CLIA} ethereum balance -o "${ethValidatorAddrp}" | jq -r ".balance")

    echo -e "${Yellow} withdraw ${NOC}"
    # shellcheck disable=SC2154
    result=$(${CLIA} chain33 withdraw -m "${lockAmount1}" -k "${chain33ReceiverAddrKey}" -r "${ethTestAddr2}" -t "${chain33EthBridgeTokenAddr}")
    cli_ret "${result}" "withdraw"

    sleep "${maturityDegree}"

    # 查询 ETH 这端 bridgeBank 地址 0
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}")
    cli_ret "${result}" "balance" ".balance" "${ethBridgeBankBalanceEnd}"

    echo -e "${Yellow} chain33ReceiverAddr chain33 端 lock 后接收地址 withdraw 后金额 ${NOC}"
    result=$(${Chain33Cli} evm query -a "${chain33EthBridgeTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
#    is_equal "${result}" "${chain33RBalancebf}"

    echo -e "${Yellow} chain33Validatorsp chain33 代理地址 withdraw 后金额 ${NOC}"
    result=$(${Chain33Cli} evm query -a "${chain33EthBridgeTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33Validatorsp})")
    # shellcheck disable=SC2219
    let chain33VspBalancewithdraw=${lockAmount1}*${le8}+${chain33VspBalancebf}
#    is_equal "${result}" "${chain33VspBalancewithdraw}"

    echo -e "${Yellow} ethTestAddr2 ethereum withdraw 接收地址 withdraw 后金额 ${NOC}"
    result=$(${CLIA} ethereum balance -o "${ethTestAddr2}" | jq -r ".balance")
    # shellcheck disable=SC2219
    let ethT2BalanceEnd=${ethT2Balancebf}+${lockAmount1}-1
#    is_equal "${result}" "${ethT2BalanceEnd}"

    echo -e "${Yellow} ethValidatorAddrp ethereum 代理地址 withdraw 后金额 ${NOC}"
    result=$(${CLIA} ethereum balance -o "${ethValidatorAddrp}" | jq -r ".balance")
    # shellcheck disable=SC2219
    let ethPBalanceEnd=${ethPBalancebf}-${lockAmount1}+1-${result}
    if [ $ethPBalanceEnd -gt 1 ]
    then echo "error $ethPBalanceEnd 大于 1"
    else echo "ok"
    fi

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# eth to chain33 在以太坊上锁定 ETH 资产,然后在 chain33 上 burn
function TestETH2Chain33Assets_proxy_excess() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    echo -e "${GRE}=========== eth to chain33 在以太坊上锁定 ETH 资产,然后在 chain33 上 burn ===========${NOC}"
    # shellcheck disable=SC2154
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}")
    cli_ret "${result}" "balance" ".balance" "0"

    local lockAmount1=20
    # shellcheck disable=SC2154
    result=$(${CLIA} ethereum lock -m ${lockAmount1} -k "${ethTestAddrKey1}" -r "${chain33ReceiverAddr}")
    cli_ret "${result}" "lock"

    # eth 等待 2 个区块
    sleep 4

    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}")
    cli_ret "${result}" "balance" ".balance" "${lockAmount1}"

    # shellcheck disable=SC2086
    sleep "${maturityDegree}"

    # chain33 chain33EthBridgeTokenAddr（ETH合约中）查询 lock 金额
    # shellcheck disable=SC2154
    result=$(${Chain33Cli} evm query -a "${chain33EthBridgeTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    is_equal "${result}" "${lockAmount1}00000000"

    result=$(${Chain33Cli} evm query -a "${chain33EthBridgeTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33Validatorsp})")
    is_equal "${result}" "0"
    # 原来的数额
    # shellcheck disable=SC2154
    ethT2Balance=$(${CLIA} ethereum balance -o "${ethTestAddr2}" | jq -r ".balance")

    # shellcheck disable=SC2154
    result=$(${CLIA} ethereum balance -o "${ethValidatorAddrp}")

    echo '#5.burn ETH from Chain33 ETH(Chain33)-----> Ethereum'
    # shellcheck disable=SC2154
    result=$(${CLIA} chain33 withdraw -m ${lockAmount1} -k "${chain33ReceiverAddrKey}" -r "${ethTestAddr2}" -t "${chain33EthBridgeTokenAddr}")
    cli_ret "${result}" "withdraw"

    sleep "${maturityDegree}"

    # 查询 ETH 这端 bridgeBank 地址 0
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}")
    cli_ret "${result}" "balance" ".balance" "${lockAmount1}"

    echo "check the balance on chain33"
    result=$(${Chain33Cli} evm query -a "${chain33EthBridgeTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    is_equal "${result}" "0"

    result=$(${Chain33Cli} evm query -a "${chain33EthBridgeTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33Validatorsp})")
    is_equal "${result}" "${lockAmount1}00000000"

    result=$(${CLIA} ethereum balance -o "${ethTestAddr2}")
    let ethT2BalanceEnd=${ethT2Balance}+${lockAmount1}-1
    cli_ret "${result}" "balance" ".balance" "${ethT2BalanceEnd}"

    result=$(${CLIA} ethereum balance -o "${ethValidatorAddrp}")

    echo -e "${GRE}=========== $FUNCNAME 超额 ===========${NOC}"
    # shellcheck disable=SC2154
    result=$(${CLIA} ethereum lock -m 120 -k "${ethTestAddrKey1}" -r "${chain33ReceiverAddr}")
    cli_ret "${result}" "lock"

    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}")
    cli_ret "${result}" "balance" ".balance" "140"
    sleep "${maturityDegree}"
    result=$(${Chain33Cli} evm query -a "${chain33EthBridgeTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    is_equal "${result}" "12000000000"

    result=$(${CLIA} chain33 withdraw -m 120 -k "${chain33ReceiverAddrKey}" -r "${ethTestAddr2}" -t "${chain33EthBridgeTokenAddr}")
    cli_ret "${result}" "withdraw"

    sleep "${maturityDegree}"

    # 查询 ETH 这端 bridgeBank 地址 0
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}")
    cli_ret "${result}" "balance" ".balance" "140"

    echo "check the balance on chain33"
    result=$(${Chain33Cli} evm query -a "${chain33EthBridgeTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    is_equal "${result}" "0"

    result=$(${Chain33Cli} evm query -a "${chain33EthBridgeTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33Validatorsp})")
    is_equal "${result}" "14000000000"

    result=$(${CLIA} ethereum balance -o "${ethTestAddr2}")
    cli_ret "${result}" "balance" ".balance" "1019"

    result=$(${CLIA} ethereum balance -o "${ethValidatorAddrp}")

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestETH2Chain33USDT_proxy() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    echo -e "${GRE}=========== eth to chain33 在以太坊上锁定 USDT 资产,然后在 chain33 上 burn ===========${NOC}"
    # shellcheck disable=SC2154
    ${CLIA} ethereum token token_transfer -k "${ethTestAddrKey1}" -m 200 -r "${ethValidatorAddrp}" -t "${ethereumUSDTERC20TokenAddr}"

    # 查询 ETH 这端 bridgeBank 地址原来是 0
    # shellcheck disable=SC2154
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumUSDTERC20TokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # ETH 这端 lock 12个 USDT
    result=$(${CLIA} ethereum lock -m 12 -k "${ethTestAddrKey1}" -r "${chain33ReceiverAddr}" -t "${ethereumUSDTERC20TokenAddr}")
    cli_ret "${result}" "lock"

    # eth 等待 2 个区块
    sleep 4

    # 查询 ETH 这端 bridgeBank 地址 12 USDT
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumUSDTERC20TokenAddr}")
    cli_ret "${result}" "balance" ".balance" "12"

    sleep "${maturityDegree}"

    # chain33 chain33EthBridgeTokenAddr（ETH合约中）查询 lock 金额
    # shellcheck disable=SC2154
    result=$(${Chain33Cli} evm query -a "${chain33USDTBridgeTokenAddr}" -c "${chain33TestAddr1}" -b "balanceOf(${chain33ReceiverAddr})")
    # 结果是 12 * le8
    is_equal "${result}" "1200000000"

     result=$(${Chain33Cli} evm query -a "${chain33USDTBridgeTokenAddr}" -c "${chain33TestAddr1}" -b "balanceOf(${chain33Validatorsp})")
     is_equal "${result}" "0"

    # 原来的数额 0
    # shellcheck disable=SC2154
    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}" -t "${ethereumUSDTERC20TokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    result=$(${CLIA} ethereum balance -o "${ethValidatorAddrp}" -t "${ethereumUSDTERC20TokenAddr}")
    cli_ret "${result}" "balance" ".balance" "200"

    echo '#5.burn YCC from Chain33 YCC(Chain33)-----> Ethereum'
    result=$(${CLIA} chain33 withdraw -m 12 -k "${chain33ReceiverAddrKey}" -r "${ethReceiverAddr1}" -t "${chain33USDTBridgeTokenAddr}")
    cli_ret "${result}" "withdraw"

    sleep "${maturityDegree}"

    echo "check the balance on chain33"
    result=$(${Chain33Cli} evm query -a "${chain33USDTBridgeTokenAddr}" -c "${chain33TestAddr1}" -b "balanceOf(${chain33ReceiverAddr})")
    is_equal "${result}" "0"
    result=$(${Chain33Cli} evm query -a "${chain33USDTBridgeTokenAddr}" -c "${chain33TestAddr1}" -b "balanceOf(${chain33Validatorsp})")
     is_equal "${result}" "1200000000"

    # 查询 ETH 这端 bridgeBank 地址 0
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumUSDTERC20TokenAddr}")
    cli_ret "${result}" "balance" ".balance" "12"
    # 更新后的金额 12
    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}" -t "${ethereumUSDTERC20TokenAddr}")
    cli_ret "${result}" "balance" ".balance" "11"

    result=$(${CLIA} ethereum balance -o "${ethValidatorAddrp}" -t "${ethereumUSDTERC20TokenAddr}")
    cli_ret "${result}" "balance" ".balance" "189"

    echo -e "${GRE}=========== $FUNCNAME 超额 ===========${NOC}"
    result=$(${CLIA} ethereum lock -m 100 -k "${ethTestAddrKey1}" -r "${chain33ReceiverAddr}" -t "${ethereumUSDTERC20TokenAddr}")
    cli_ret "${result}" "lock"

    sleep "${maturityDegree}"

     result=$(${Chain33Cli} evm query -a "${chain33USDTBridgeTokenAddr}" -c "${chain33TestAddr1}" -b "balanceOf(${chain33ReceiverAddr})")
    # 结果是 12 * le8
#    is_equal "${result}" "1200000000"

     result=$(${Chain33Cli} evm query -a "${chain33USDTBridgeTokenAddr}" -c "${chain33TestAddr1}" -b "balanceOf(${chain33Validatorsp})")
#     is_equal "${result}" "0"

    # 原来的数额 0
    # shellcheck disable=SC2154
    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}" -t "${ethereumUSDTERC20TokenAddr}")
#    cli_ret "${result}" "balance" ".balance" "0"

    result=$(${CLIA} ethereum balance -o "${ethValidatorAddrp}" -t "${ethereumUSDTERC20TokenAddr}")
#    cli_ret "${result}" "balance" ".balance" "200"

    echo '#5.burn YCC from Chain33 YCC(Chain33)-----> Ethereum'
    result=$(${CLIA} chain33 withdraw -m 100 -k "${chain33ReceiverAddrKey}" -r "${ethReceiverAddr1}" -t "${chain33USDTBridgeTokenAddr}")
    cli_ret "${result}" "withdraw"

    sleep "${maturityDegree}"

    echo "check the balance on chain33"
    result=$(${Chain33Cli} evm query -a "${chain33USDTBridgeTokenAddr}" -c "${chain33TestAddr1}" -b "balanceOf(${chain33ReceiverAddr})")
#    is_equal "${result}" "0"
    result=$(${Chain33Cli} evm query -a "${chain33USDTBridgeTokenAddr}" -c "${chain33TestAddr1}" -b "balanceOf(${chain33Validatorsp})")
#     is_equal "${result}" "1200000000"

    # 查询 ETH 这端 bridgeBank 地址 0
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumUSDTERC20TokenAddr}")
#    cli_ret "${result}" "balance" ".balance" "12"
    # 更新后的金额 12
    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}" -t "${ethereumUSDTERC20TokenAddr}")
#    cli_ret "${result}" "balance" ".balance" "11"

    result=$(${CLIA} ethereum balance -o "${ethValidatorAddrp}" -t "${ethereumUSDTERC20TokenAddr}")
#    cli_ret "${result}" "balance" ".balance" "189"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestRelayerProxy() {
    start_docker_ebrelayerProxy
    setWithdraw

    TestETH2Chain33Assets_proxy 20
#    TestETH2Chain33Assets_proxy 30
#    TestETH2Chain33Assets_proxy_excess 60
#    TestETH2Chain33USDT_proxy
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
        # shellcheck disable=SC2034
        chain33ID="${2}"
    fi

    get_cli

    # init
    # shellcheck disable=SC2154
    # shellcheck disable=SC2034
    Chain33Cli=${MainCli}
    InitChain33Validator
    # para add
    initPara

    StartDockerRelayerDeploy
#    test_all

    TestRelayerProxy

    echo_addrs
    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}
