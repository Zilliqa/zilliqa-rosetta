#!/bin/bash
# ========================================
# RUN ROSETTA START
# ========================================
function run_rosetta() {
cwd=$(pwd)
cd "/rosetta"
echo $BLOCKCHAIN_NETWORK
mv "$BLOCKCHAIN_NETWORK.config.local.yaml" "config.local.yaml"
nohup ./main &
cd "$cwd"

}


# ========================================
# RUN SEEDNODE START
# ========================================
function run_seednode() {
#if [ "$BLOCKCHAIN_NETWORK" = "testnet" ]
#then
#    seed_configuration_url="https://testnet-join.zilliqa.com/seed-configuration.tar.gz"
#elif [ "$BLOCKCHAIN_NETWORK" = "mainnet" ]
#then
#    seed_configuration_url="https://mainnet-join.zilliqa.com/seed-configuration.tar.gz"
#else
#    echo "Error, unknown $BLOCKCHAIN_NETWORK, terminating."
#    exit 1
#fi

seed_configuration_url="https://joel-refactor2-join.dev.z7a.xyz/seed-configuration.tar.gz"
echo $seed_configuration_url

NONINTERACTIVE="true"
ZILLIQA_PATH="/run/zilliqa"

curl -O "$seed_configuration_url"
tar -zxvf seed-configuration.tar.gz

export NONINTERACTIVE ZILLIQA_PATH && ./launch.sh

}


# ========================================
# SCRIPT START
# ========================================
if [ "$GENKEYPAIR" = "true" ]
then
    genkeypair
    exit 0
else
    run_rosetta
    run_seednode
fi
