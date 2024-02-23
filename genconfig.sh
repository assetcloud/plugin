#!/bin/bash

set -x
config="chain.para.toml"

function genConfig() {
    echo -e "start auto modify config"

    peerArray=$1
    num=$2

    validatorNodes=""
    seeds=""
    for((i=0;i<$num;i++));
    do
	validatorNodes=${validatorNodes},\""${peerArray[$i]}:33001"\"
        seeds=${seeds},"\"${peerArray[$i]}:13801\""
    done
    validatorNodes=$(echo $validatorNodes|sed '/^,/s///g')
    seeds=$(echo $seeds|sed '/^,/s///g')

    sed -i "s/^validatorNodes=.*/validatorNodes=[${validatorNodes}]/g" ${config}
    sed -i "s/^seeds=.*/seeds=[${seeds}]/g" ${config}
    for((i=0;i<$num;i++));
    do
	rm -rf ${peerArray[$i]}
        mkdir -p ${peerArray[$i]}
        cp ${config} ${peerArray[$i]}/${config}
    done

}

function prepareChainPkg() {
    echo -e "start prepare chain33 files"
    
    peerArray=$1
    num=$2

    ./chain-cli qbft  gen_file -n $num -t bls

    for((i=0;i<$num;i++));
    do
        cp priv_validator_$i.json  ${peerArray[$i]}/priv_validator.json
        cp genesis_file.json  ${peerArray[$i]}/genesis.json
	cp chain ${peerArray[$i]}/chain
	cp chain-cli ${peerArray[$i]}/chain-cli
    done

}

function main() {
    peers=$1

    peerArray=(${peers//,/ })

    peerNum=${#peerArray[*]}

    genConfig $peerArray $peerNum

    prepareChainPkg $peerArray $peerNum
}

main $1

