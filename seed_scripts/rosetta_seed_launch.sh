#!/bin/bash

prog=$(basename $0)
mykeyfile=mykey.txt
myaddrfile=myaddr.txt
cmd_log=last.log


# ========================================
# RUN FUNCTION START
# ========================================
function run() {

# Abstracted variables
ip="$IP_ADDRESS"
multiplier_sync="$MULTIPLIER_SYNC"
extseedprivk="$SEED_PRIVATE_KEY"
testnet_name="$TESTNET_NAME"
bucket_name="$BUCKET_NAME"
# abstracted outside of run function in preparation of 7.2.0 seed configuration launch.sh scripts
#seed_configuration_url="$SEED_CONFIGURATION_URL"

name="zilliqa"
port="33133"
zilliqa_path="/zilliqa"
storage_path="`dirname \"$0\"`"
storage_path="`( cd \"$storage_path\" && pwd )`"
base_path="$storage_path"
exclusion_txbodies_mb="false"
isSeed="true"

curl -O "$seed_configuration_url"
tar -zxvf seed-configuration.tar.gz

if [ -z "$zilliqa_path" -o ! -x "/usr/local/bin/zilliqa"  ]
then
    echo "Cannot find zilliqa binary on the path you specified"
    exit 1
fi

echo "$ip"

if [ ! -s $mykeyfile ]
then
    echo "Cannot find $mykeyfile, please mount the key pair file."
    exit 1
fi

prikey=$(cat $mykeyfile | awk '{ print $2 }')
pubkey=$(cat $mykeyfile | awk '{ print $1 }')

/usr/local/bin/getaddr --pubk $pubkey > $myaddrfile

cmd="cp /zilliqa/scripts/download_incr_DB.py /run/zilliqa/download_incr_DB.py && cp /zilliqa/scripts/download_static_DB.py /run/zilliqa/download_static_DB.py && sed -i \"/TESTNET_NAME=/c\TESTNET_NAME= '${testnet_name}'\" /run/zilliqa/download_incr_DB.py /run/zilliqa/download_static_DB.py && sed -i \"/BUCKET_NAME=/c\BUCKET_NAME= '${bucket_name}'\" /run/zilliqa/download_incr_DB.py /run/zilliqa/download_static_DB.py && o1=\$(grep -oPm1 '(?<=<NUM_FINAL_BLOCK_PER_POW>)[^<]+' /run/zilliqa/constants.xml) && [ ! -z \$o1 ] && sed -i \"/NUM_FINAL_BLOCK_PER_POW=/c\NUM_FINAL_BLOCK_PER_POW= \$o1\" /run/zilliqa/download_incr_DB.py && o1=\$(grep -oPm1 '(?<=<INCRDB_DSNUMS_WITH_STATEDELTAS>)[^<]+' /run/zilliqa/constants.xml) && [ ! -z \$o1 ] && sed -i \"/NUM_DSBLOCK=/c\NUM_DSBLOCK= \$o1\" /run/zilliqa/download_incr_DB.py && zilliqa --privk $prikey --pubk $pubkey --address $ip --port $port --synctype 6 --recovery"

if [ "$multiplier_sync" = "N" ] || [ "$multiplier_sync" = "n" ]
then
    cmd="$cmd --l2lsyncmode --extseedprivk $extseedprivk"
fi

eval ${cmd}
echo
echo "Use 'cat $cmd_log' to see the command output"
echo "Use 'tail -f zilliqa-00001-log.txt' to see the runtime log"
}


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
# SCRIPT START
# ========================================
if [ "$GENKEYPAIR" = "true" ]
then
    genkeypair
    exit 0
else
    run_rosetta
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
fi
echo $seed_configuration_url
run
