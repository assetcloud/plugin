#!/usr/bin/env bash

set -e
set -o pipefail
#set -o verbose
#set -o xtrace

sedfix=""
if [ "$(uname)" == "Darwin" ]; then
    sedfix=".bak"
fi

OP="${1}"
PROJ="${2}"
DAPP="${3}"
EXTRA="${4}"

TESTCASEFILE="testcase.sh"
FORKTESTFILE="fork-test.sh"
DOCKER_COMPOSE_SH="docker-compose.sh"

function down_dapp() {
    app=$1

    if [ -d "${app}" ] && [ "${app}" != "system" ] && [ -d "${app}-ci" ]; then
        cd "${app}"-ci/ && pwd

        echo "============ down dapp=$app start ================="
        ./docker-compose-down.sh "${PROJ}" "${app}"
        echo "============ down dapp=$app end ================="

        cd .. && rm -rf ./"${app}"-ci
    fi
}

function run_dapp() {
    local app=$1
    local test=$2

    echo "============ run dapp=$app start ================="
    if [ "$app" == "metrics" ]; then
        cp ./ci/paracross/* ./metrics && echo $?
        cp -n ./* ./metrics/ && echo $?
        cp -r ci/dapptest/ metrics/ && echo $?
        cd metrics && pwd
        rm docker-compose-paracross.yml
        mv docker-compose-metrics.yml docker-compose-paracross.yml
        app="paracross"
    else
        rm -rf "${app}"-ci && mkdir -p "${app}"-ci && cp -r ./"${app}"/* ./"${app}"-ci && echo $?
        cp -n ./* ./"${app}"-ci/ && echo $?
        if [ "$app" == "paracross" ]; then
            cp -r dapptest/ "${app}"-ci/ && echo $?
        fi

        cd "${app}"-ci/ && pwd
    fi

    if [ "$test" == "$FORKTESTFILE" ]; then
        sed -i $sedfix 's/^system_coins_file=.*/system_coins_file="..\/system\/coins\/fork-test.sh"/g' system-fork-test.sh
        if ! ./system-fork-test.sh "${PROJ}" "${app}"; then
            exit 1
        fi
    elif [ "$test" == "$TESTCASEFILE" ]; then
        if ! ./${DOCKER_COMPOSE_SH} "${PROJ}" "${app}" "${EXTRA}"; then
            exit 1
        fi
    fi
    cd ..
    echo "============ run dapp=$app end ================="
}

function run_single_app() {
    app=$1
    testfile=$2

    if [ -d "${app}" ] && [ "${app}" != "system" ]; then
        if [ -e "$app/$testfile" ]; then
            run_dapp "${app}" "$testfile"
            if [ "$#" -gt 2 ]; then
                down_dapp "${app}"
            fi
        else
            echo "#### dapp=$app not exist the test file: $testfile ####"
        fi
    else
        echo "#### dapp=$app not exist or is system dir ####"
    fi
}

function main() {
    if [ "${OP}" == "run" ]; then
        #copy chain33 system-test-rpc.sh
        cp "$(go list -f "{{.Dir}}" github.com/assetcloud/chain)"/build/system-test-rpc.sh ./
        if [ "${DAPP}" == "all" ] || [ "${DAPP}" == "ALL" ]; then
            echo "============ run main start ================="
            if ! ./${DOCKER_COMPOSE_SH} "$PROJ"; then
                exit 1
            fi
            ./docker-compose-down.sh "$PROJ"
            echo "============ run main end ================="

            find . -maxdepth 1 -type d -name "*-ci" -exec rm -rf {} \;
            dir=$(find . -maxdepth 1 -type d ! -name system ! -name . | sed 's/^\.\///')
            for app in $dir; do
                run_single_app "${app}" "$TESTCASEFILE" "down"
            done
        elif [ -n "${DAPP}" ]; then
            run_single_app "${DAPP}" "$TESTCASEFILE"
        else
            ./${DOCKER_COMPOSE_SH} "${PROJ}"
        fi
    elif [ "${OP}" == "down" ]; then
        if [ "${DAPP}" == "all" ] || [ "${DAPP}" == "ALL" ]; then
            dir=$(find . -maxdepth 1 -type d ! -name system ! -name . | sed 's/^\.\///')
            for app in $dir; do
                down_dapp "${app}"
            done
            ./docker-compose-down.sh "${PROJ}"
        elif [ -n "${DAPP}" ]; then
            down_dapp "${DAPP}"
        else
            ./docker-compose-down.sh "${PROJ}"
        fi
    elif [ "${OP}" == "forktest" ]; then
        if [ "${DAPP}" == "all" ] || [ "${DAPP}" == "ALL" ]; then
            echo "============ run main start ================="
            if ! ./system-fork-test.sh "$PROJ"; then
                exit 1
            fi
            echo "============ down main test ================="
            ./docker-compose-down.sh "$PROJ"
            echo "============ run main end ================="
            # remove all the *-ci folders
            find . -maxdepth 1 -type d -name "*-ci" -exec rm -rf {} \;
            dir=$(find . -maxdepth 1 -type d ! -name system ! -name . | sed 's/^\.\///')
            for app in $dir; do
                run_single_app "${app}" "$FORKTESTFILE" "down"
            done
        elif [ -n "${DAPP}" ]; then
            run_single_app "${DAPP}" "$FORKTESTFILE"
        else
            ./system-fork-test.sh "${PROJ}"
        fi
    elif [ "${OP}" == "modify" ]; then
        sed -i $sedfix '/^useGithub=.*/a version=1' chain33.toml
    fi
}

# run script
main
