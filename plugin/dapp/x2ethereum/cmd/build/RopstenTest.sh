#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null

# 测试前请先阅读文档，对代码进行修改 http://note.youdao.com/noteshare?id=675ce9a1162b6639182206864985e935&sub=14A1CFAA14CB4318B207EAEC94991A93
set -x
set +e

source "./publicTest.sh"

CLIA="./ebcli_A"
CLIB="./ebcli_B"
CLIC="./ebcli_C"
CLID="./ebcli_D"

chainSenderAddr="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
# validatorsAddr=["0x92c8b16afd6d423652559c6e266cbe1c29bfd84f", "0x0df9a824699bc5878232c9e612fe1a5346a5a368", "0xcb074cb21cdddf3ce9c3c0a7ac4497d633c9d9f1", "0xd9dab021e74ecf475788ed7b61356056b2095830"]
ethValidatorAddrKeyA="3fa21584ae2e4fd74db9b58e2386f5481607dfa4d7ba0617aaa7858e5025dc1e"
ethValidatorAddrKeyB="a5f3063552f4483cfc20ac4f40f45b798791379862219de9e915c64722c1d400"
ethValidatorAddrKeyC="bbf5e65539e9af0eb0cfac30bad475111054b09c11d668fc0731d54ea777471e"
ethValidatorAddrKeyD="c9fa31d7984edf81b8ef3b40c761f1847f6fcd5711ab2462da97dc458f1f896b"
# 新增地址 chain 需要导入地址 转入 10 bty当收费费
chainValidator1="1GTxrmuWiXavhcvsaH5w9whgVxUrWsUMdV"
chainValidator2="155ooMPBTF8QQsGAknkK7ei5D78rwDEFe6"
chainValidator3="13zBdQwuyDh7cKN79oT2odkxYuDbgQiXFv"
chainValidator4="113ZzVamKfAtGt9dq45fX1mNsEoDiN95HG"
chainValidatorKey1="0xd627968e445f2a41c92173225791bae1ba42126ae96c32f28f97ff8f226e5c68"
chainValidatorKey2="0x9d539bc5fd084eb7fe86ad631dba9aa086dba38418725c38d9751459f567da66"
chainValidatorKey3="0x0a6671f101e30a2cc2d79d77436b62cdf2664ed33eb631a9c9e3f3dd348a23be"
chainValidatorKey4="0x3818b257b05ee75b6e43ee0e3cfc2d8502342cf67caed533e3756966690b62a5"
ethReceiverAddr1="0xa4ea64a583f6e51c3799335b28a8f0529570a635"
ethReceiverAddrKey1="355b876d7cbcb930d5dfab767f66336ce327e082cbaa1877210c1bae89b1df71"
ethReceiverAddr2="0x0c05ba5c230fdaa503b53702af1962e08d0c60bf"
ethReceiverAddrKey2="9dc6df3a8ab139a54d8a984f54958ae0661f880229bf3bdbb886b87d58b56a08"

maturityDegree=10
tokenAddrBty=""

function kill_ebrelayerC() {
    kill_ebrelayer "./C/ebrelayer"
}
function kill_ebrelayerD() {
    kill_ebrelayer "./D/ebrelayer"
}
function start_ebrelayerC() {
    start_ebrelayer_and_unlock C
    block_wait "${ChainCli}" $((maturityDegree + 3))
    eth_block_wait 3 https://ropsten-rpc.linkpool.io/
    sleep 1
}
function start_ebrelayerD() {
    start_ebrelayer_and_unlock D
    block_wait "${ChainCli}" $((maturityDegree + 3))
    eth_block_wait 3 https://ropsten-rpc.linkpool.io/
    sleep 1
}

function InitAndDeploy() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    result=$(${CLIA} relayer ethereum deploy)
    cli_ret "${result}" "deploy"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function StartRelayerAndDeploy() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    for name in A B C D; do
        local ebrelayer="./$name/ebrelayer"
        kill_ebrelayer "${ebrelayer}"
    done
    sleep 1

    rm -rf './A' './B' './C' './D' './datadir' './ebrelayer.log' './logs'
    mkdir './A' './B' './C' './D'
    cp './relayer.toml' './A/relayer.toml'
    cp './ebrelayer' './A/ebrelayer'

    start_ebrelayer_and_setpwd_unlock A

    # 部署合约
    InitAndDeploy

    # 获取 BridgeRegistry 地址
    result=$(${CLIA} relayer ethereum bridgeRegistry)
    BridgeRegistry=$(cli_ret "${result}" "bridgeRegistry" ".addr")
    #    BridgeRegistry="0xcA5E8FCE034888ea51eB568CCA83C413b9DE3F73"

    kill_ebrelayer "./A/ebrelayer"
    # 修改 relayer.toml 配置文件
    updata_relayer_toml_ropston "${BridgeRegistry}" ${maturityDegree} "./A/relayer.toml"
    updata_all_relayer_toml2

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# chian33 添加验证着及权重
function InitChainVilators() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    # 导入 chainValidators 私钥生成地址
    result=$(${ChainCli} account import_key -k ${chainValidatorKey1} -l validator1)
    check_addr "${result}" ${chainValidator1}
    result=$(${ChainCli} account import_key -k ${chainValidatorKey2} -l validator2)
    check_addr "${result}" ${chainValidator2}
    result=$(${ChainCli} account import_key -k ${chainValidatorKey3} -l validator3)
    check_addr "${result}" ${chainValidator3}
    result=$(${ChainCli} account import_key -k ${chainValidatorKey4} -l validator4)
    check_addr "${result}" ${chainValidator4}

    # SetConsensusThreshold
    hash=$(${ChainCli} send x2ethereum setconsensus -p 80 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${ChainCli}" "${hash}"

    # add a validator
    hash=$(${ChainCli} send x2ethereum add -a ${chainValidator1} -p 25 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${ChainCli}" "${hash}"
    hash=$(${ChainCli} send x2ethereum add -a ${chainValidator2} -p 25 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${ChainCli}" "${hash}"
    hash=$(${ChainCli} send x2ethereum add -a ${chainValidator3} -p 25 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${ChainCli}" "${hash}"
    hash=$(${ChainCli} send x2ethereum add -a ${chainValidator4} -p 25 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${ChainCli}" "${hash}"

    # query Validators
    totalPower=$(${ChainCli} send x2ethereum query totalpower -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq .totalPower | sed 's/\"//g')
    check_number 100 "${totalPower}"

    # cions 转帐到 x2ethereum 合约地址
    hash=$(${ChainCli} send coins send_exec -e x2ethereum -a 200 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)

    check_tx "${ChainCli}" "${hash}"
    result=$(${ChainCli} account balance -a 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -e x2ethereum)
    balance_ret "${result}" "200.0000"

    # chainValidator 要有手续费
    hash=$(${ChainCli} send coins transfer -a 10 -t "${chainValidator1}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${ChainCli}" "${hash}"
    result=$(${ChainCli} account balance -a "${chainValidator1}" -e coins)
    balance_ret "${result}" "10.0000"

    hash=$(${ChainCli} send coins transfer -a 10 -t "${chainValidator2}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${ChainCli}" "${hash}"
    result=$(${ChainCli} account balance -a "${chainValidator2}" -e coins)
    balance_ret "${result}" "10.0000"

    hash=$(${ChainCli} send coins transfer -a 10 -t "${chainValidator3}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${ChainCli}" "${hash}"
    result=$(${ChainCli} account balance -a "${chainValidator3}" -e coins)
    balance_ret "${result}" "10.0000"

    hash=$(${ChainCli} send coins transfer -a 10 -t "${chainValidator4}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${ChainCli}" "${hash}"
    result=$(${ChainCli} account balance -a "${chainValidator4}" -e coins)
    balance_ret "${result}" "10.0000"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function StartAllEbrelayer() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    start_ebrelayer_and_unlock A
    # 重启 ebrelayer 并解锁
    for name in B C D; do
        start_ebrelayer_and_setpwd_unlock $name
    done

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function EthImportKey() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    result=$(${CLIA} relayer ethereum import_chainprivatekey -k "${chainValidatorKey1}")
    cli_ret "${result}" "import_chainprivatekey"
    result=$(${CLIB} relayer ethereum import_chainprivatekey -k "${chainValidatorKey2}")
    cli_ret "${result}" "import_chainprivatekey"
    result=$(${CLIC} relayer ethereum import_chainprivatekey -k "${chainValidatorKey3}")
    cli_ret "${result}" "import_chainprivatekey"
    result=$(${CLID} relayer ethereum import_chainprivatekey -k "${chainValidatorKey4}")
    cli_ret "${result}" "import_chainprivatekey"

    result=$(${CLIA} relayer chain import_privatekey -k "${ethValidatorAddrKeyA}")
    cli_ret "${result}" "A relayer chain import_privatekey"
    result=$(${CLIB} relayer chain import_privatekey -k "${ethValidatorAddrKeyB}")
    cli_ret "${result}" "B relayer chain import_privatekey"
    result=$(${CLIC} relayer chain import_privatekey -k "${ethValidatorAddrKeyC}")
    cli_ret "${result}" "C relayer chain import_privatekey"
    result=$(${CLID} relayer chain import_privatekey -k "${ethValidatorAddrKeyD}")
    cli_ret "${result}" "D relayer chain import_privatekey"
    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestChainToEthAssets() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    # token4chain 在 以太坊 上先有 bty
    tokenSymbol="coins.bty"

    result=$(${CLIA} relayer ethereum token4chain -s "${tokenSymbol}")
    tokenAddrBty=$(cli_ret "${result}" "token4chain" ".addr")

    #tokenAddrBty="0xE79142B3171019fcfcA838f0792edB08d4F2a94F"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "0"

    # chain lock bty
    hash=$(${ChainCli} send x2ethereum lock -a 5 -t "${tokenSymbol}" -r ${ethReceiverAddr1} -q "${tokenAddrBty}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    block_wait "${ChainCli}" $((maturityDegree + 3))
    check_tx "${ChainCli}" "${hash}"

    result=$(${ChainCli} account balance -a 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -e x2ethereum)
    balance_ret "${result}" "195.0000"

    eth_block_wait $((maturityDegree + 3)) https://ropsten-rpc.linkpool.io/

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "5"

    # eth burn
    result=$(${CLIA} relayer ethereum burn -m 5 -k "${ethReceiverAddrKey1}" -r "${chainSenderAddr}" -t "${tokenAddrBty}")
    cli_ret "${result}" "burn"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "0"

    # eth 等待 10 个区块
    eth_block_wait $((maturityDegree + 3)) https://ropsten-rpc.linkpool.io/

    result=$(${ChainCli} account balance -a "${chainSenderAddr}" -e x2ethereum)
    balance_ret "${result}" "5"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# eth to chain
# 在以太坊上锁定资产,然后在 chain 上铸币,针对 eth 资产
function TestETH2ChainAssets() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    ${CLIA} relayer unlock -p 123456hzj

    result=$(${CLIA} relayer ethereum bridgeBankAddr)
    bridgeBankAddr=$(cli_ret "${result}" "bridgeBankAddr" ".addr")

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # eth lock 0.1
    result=$(${CLIA} relayer ethereum lock -m 0.1 -k "${ethReceiverAddrKey1}" -r 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    cli_ret "${result}" "lock"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "0.1"

    # eth 等待 10 个区块
    eth_block_wait $((maturityDegree + 3)) https://ropsten-rpc.linkpool.io/

    result=$(${ChainCli} x2ethereum balance -s 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -t eth | jq ".res" | jq ".[]")
    balance_ret "${result}" "0.1"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}")
    balance=$(cli_ret "${result}" "balance" ".balance")

    hash=$(${ChainCli} send x2ethereum burn -a 0.1 -t eth -r ${ethReceiverAddr2} -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    block_wait "${ChainCli}" $((maturityDegree + 3))
    check_tx "${ChainCli}" "${hash}"

    result=$(${ChainCli} x2ethereum balance -s 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -t eth | jq ".res" | jq ".[]")
    balance_ret "${result}" "0"

    eth_block_wait $((maturityDegree + 3)) https://ropsten-rpc.linkpool.io/

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}")
    cli_ret "${result}" "balance" ".balance" "$(echo "${balance}+0.1" | bc)"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestETH2ChainErc20() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    ${CLIA} relayer unlock -p 123456hzj

    # token4erc20 在 chain 上先有 token,同时 mint
    tokenSymbol="testc"
    result=$(${CLIA} relayer ethereum token4erc20 -s "${tokenSymbol}")
    tokenAddr=$(cli_ret "${result}" "token4erc20" ".addr")
    #tokenAddr="0x47F62ba65bCa4150BE98F31566DC559b9b04fc2D"

    # 先铸币 1000
    result=$(${CLIA} relayer ethereum mint -m 1000 -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "mint"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "1000"

    result=$(${CLIA} relayer ethereum bridgeBankAddr)
    bridgeBankAddr=$(cli_ret "${result}" "bridgeBankAddr" ".addr")

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # lock 100
    result=$(${CLIA} relayer ethereum lock -m 100 -k "${ethReceiverAddrKey1}" -r "${chainValidator1}" -t "${tokenAddr}")
    cli_ret "${result}" "lock"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "900"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "100"

    # eth 等待 10 个区块
    eth_block_wait $((maturityDegree + 3)) https://ropsten-rpc.linkpool.io/

    result=$(${ChainCli} x2ethereum balance -s "${chainValidator1}" -t "${tokenSymbol}" -a "${tokenAddr}" | jq ".res" | jq ".[]")
    balance_ret "${result}" "100"

    # chain burn 100
    hash=$(${ChainCli} send x2ethereum burn -a 100 -t "${tokenSymbol}" -r ${ethReceiverAddr2} -q "${tokenAddr}" -k "${chainValidator1}")
    block_wait "${ChainCli}" $((maturityDegree + 3))
    check_tx "${ChainCli}" "${hash}"

    result=$(${ChainCli} x2ethereum balance -s "${chainValidator1}" -t "${tokenSymbol}" -a "${tokenAddr}" | jq ".res" | jq ".[]")
    balance_ret "${result}" "0"

    eth_block_wait $((maturityDegree + 3)) https://ropsten-rpc.linkpool.io/

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "100"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestChainToEthAssetsKill() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    tokenSymbol="coins.bty"
    if [ "${tokenAddrBty}" == "" ]; then
        # token4chain 在 以太坊 上先有 bty
        result=$(${CLIA} relayer ethereum token4chain -s "${tokenSymbol}")
        tokenAddrBty=$(cli_ret "${result}" "token4chain" ".addr")
    fi

    #    tokenAddrBty="0xE79142B3171019fcfcA838f0792edB08d4F2a94F"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "0"

    kill_ebrelayerC
    kill_ebrelayerD

    # chain lock bty
    hash=$(${ChainCli} send x2ethereum lock -a 1.41 -t "${tokenSymbol}" -r ${ethReceiverAddr2} -q "${tokenAddrBty}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    block_wait "${ChainCli}" $((maturityDegree + 3))
    check_tx "${ChainCli}" "${hash}"

    eth_block_wait $((maturityDegree + 3)) https://ropsten-rpc.linkpool.io/

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "0"

    start_ebrelayerC

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "1.41"

    # eth burn
    result=$(${CLIA} relayer ethereum burn -m 1.41 -k "${ethReceiverAddrKey2}" -r "${chainValidator1}" -t "${tokenAddrBty}")
    cli_ret "${result}" "burn"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "0"

    # eth 等待 10 个区块
    eth_block_wait $((maturityDegree + 3)) https://ropsten-rpc.linkpool.io/

    result=$(${ChainCli} account balance -a "${chainValidator1}" -e x2ethereum)
    balance_ret "${result}" "0"

    start_ebrelayerD

    result=$(${ChainCli} account balance -a "${chainValidator1}" -e x2ethereum)
    balance_ret "${result}" "1.41"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# eth to chain
# 在以太坊上锁定资产,然后在 chain 上铸币,针对 eth 资产
function TestETH2ChainAssetsKill() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    ${CLIA} relayer unlock -p 123456hzj

    result=$(${CLIA} relayer ethereum bridgeBankAddr)
    bridgeBankAddr=$(cli_ret "${result}" "bridgeBankAddr" ".addr")

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}")
    balance=$(cli_ret "${result}" "balance" ".balance")

    kill_ebrelayerC
    kill_ebrelayerD

    # eth lock 0.1
    result=$(${CLIA} relayer ethereum lock -m 0.133 -k "${ethReceiverAddrKey1}" -r 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    cli_ret "${result}" "lock"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "$(echo "${balance}+0.133" | bc)"

    # eth 等待 10 个区块
    eth_block_wait $((maturityDegree + 3)) https://ropsten-rpc.linkpool.io/

    balance=$(${ChainCli} x2ethereum balance -s 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -t eth | jq ".res" | jq ".[]" | jq -r ".balance")
    balance_ret "${result}" "0"

    start_ebrelayerC
    start_ebrelayerD

    result=$(${ChainCli} x2ethereum balance -s 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -t eth | jq ".res" | jq ".[]" | jq -r ".balance")
    balance_ret "${result}" "0.133"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}")
    balance=$(cli_ret "${result}" "balance" ".balance")

    kill_ebrelayerC
    kill_ebrelayerD

    hash=$(${ChainCli} send x2ethereum burn -a 0.133 -t eth -r ${ethReceiverAddr2} -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    block_wait "${ChainCli}" $((maturityDegree + 3))
    check_tx "${ChainCli}" "${hash}"

    result=$(${ChainCli} x2ethereum balance -s 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -t eth | jq ".res" | jq ".[]")
    balance_ret "${result}" "0"

    eth_block_wait $((maturityDegree + 3)) https://ropsten-rpc.linkpool.io/

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "0.133"

    start_ebrelayerC
    start_ebrelayerD

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}")
    cli_ret "${result}" "balance" ".balance" "$(echo "${balance}+0.133" | bc)"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestETH2ChainErc20Kill() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    ${CLIA} relayer unlock -p 123456hzj

    # token4erc20 在 chain 上先有 token,同时 mint
    tokenSymbol="testcc"
    result=$(${CLIA} relayer ethereum token4erc20 -s "${tokenSymbol}")
    tokenAddr=$(cli_ret "${result}" "token4erc20" ".addr")

    # 先铸币 1000
    result=$(${CLIA} relayer ethereum mint -m 1000 -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "mint"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "1000"

    result=$(${CLIA} relayer ethereum bridgeBankAddr)
    bridgeBankAddr=$(cli_ret "${result}" "bridgeBankAddr" ".addr")

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    kill_ebrelayerC
    kill_ebrelayerD

    # lock 100
    result=$(${CLIA} relayer ethereum lock -m 100 -k "${ethReceiverAddrKey1}" -r "${chainValidator1}" -t "${tokenAddr}")
    cli_ret "${result}" "lock"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "900"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "100"

    # eth 等待 10 个区块
    eth_block_wait $((maturityDegree + 3)) https://ropsten-rpc.linkpool.io/

    result=$(${ChainCli} x2ethereum balance -s "${chainValidator1}" -t "${tokenSymbol}" -a "${tokenAddr}" | jq ".res" | jq ".[]")
    balance_ret "${result}" "0"

    start_ebrelayerC
    start_ebrelayerD

    result=$(${ChainCli} x2ethereum balance -s "${chainValidator1}" -t "${tokenSymbol}" -a "${tokenAddr}" | jq ".res" | jq ".[]")
    balance_ret "${result}" "100"

    kill_ebrelayerC
    kill_ebrelayerD

    # chain burn 100
    hash=$(${ChainCli} send x2ethereum burn -a 100 -t "${tokenSymbol}" -r ${ethReceiverAddr2} -q "${tokenAddr}" -k "${chainValidator1}")
    block_wait "${ChainCli}" $((maturityDegree + 3))
    check_tx "${ChainCli}" "${hash}"

    result=$(${ChainCli} x2ethereum balance -s "${chainValidator1}" -t "${tokenSymbol}" -a "${tokenAddr}" | jq ".res" | jq ".[]")
    balance_ret "${result}" "0"

    eth_block_wait $((maturityDegree + 3)) https://ropsten-rpc.linkpool.io/

    start_ebrelayerC

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "100"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function AllRelayerMainTest() {
    set +e
    #    docker cp ${GOPATH}/src/github.com/assetcloud/plugin/build/ci/x2ethereum build_chain_1:/root/x2ethereum
    #    docker cp ${GOPATH}/src/github.com/assetcloud/plugin/plugin/dapp/x2ethereum/cmd/build/RopstenTest.sh build_chain_1:/root/x2ethereum/RopstenTest.sh

    ChainCli="./../chain-cli"

    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    if [[ ${1} != "" ]]; then
        maturityDegree=${1}
        echo -e "${GRE}maturityDegree is ${maturityDegree} ${NOC}"
    fi

    # init
    StartRelayerAndDeploy
    InitChainVilators
    StartAllEbrelayer
    EthImportKey

    # test
    TestChainToEthAssets
    TestETH2ChainAssets
    TestETH2ChainErc20

    # kill relayer and start relayer
    TestChainToEthAssetsKill
    TestETH2ChainAssetsKill
    TestETH2ChainErc20Kill

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

AllRelayerMainTest 1
