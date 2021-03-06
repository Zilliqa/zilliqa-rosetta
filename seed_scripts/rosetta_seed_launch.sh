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
if [ "$BLOCKCHAIN_NETWORK" = "testnet" ]
then
    seed_configuration_url="https://testnet-join.zilliqa.com/seed-configuration.tar.gz"
elif [ "$BLOCKCHAIN_NETWORK" = "mainnet" ]
then
    seed_configuration_url="https://mainnet-join.zilliqa.com/seed-configuration.tar.gz"
else
    echo "Error, unknown $BLOCKCHAIN_NETWORK, terminating."
    exit 1
fi

echo $seed_configuration_url

NONINTERACTIVE="true"
ZILLIQA_PATH="/run/zilliqa"

curl -O "$seed_configuration_url"
tar -zxvf seed-configuration.tar.gz

export NONINTERACTIVE ZILLIQA_PATH && ./launch.sh

tail -f zilliqa-00001-log.txt

}


# ========================================
# SCRIPT START
# ========================================
if [ "$GENKEYPAIR" = "true" ]
then
    genkeypair
    exit 0
else
    echo "starting rosetta"
    run_rosetta
    echo "starting seed node"
    run_seednode
fi
