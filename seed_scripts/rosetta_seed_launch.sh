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



name="zilliqa"
port="33133"
zilliqa_path="/zilliqa"
storage_path="`dirname \"$0\"`"
storage_path="`( cd \"$MY_PATH\" && pwd )`"
base_path="$storage_path"
exclusion_txbodies_mb="false"
isSeed="true"

if [ -z "$zilliqa_path" -o ! -x $zilliqa_path/build/bin/zilliqa ]
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

$zilliqa_path/build/bin/getaddr --pubk $pubkey > $myaddrfile

cmd="cp ${zilliqa_path}/scripts/download_incr_DB.py ${base_path}/download_incr_DB.py; sed -i \"/TESTNET_NAME=/c\TESTNET_NAME= '${testnet_name}'\" ${base_path}/download_incr_DB.py; sed -i \"/BUCKET_NAME=/c\BUCKET_NAME= '${bucket_name}'\" ${base_path}/download_incr_DB.py; o1=\$(grep -oPm1 '(?<=<NUM_FINAL_BLOCK_PER_POW>)[^<]+' ${base_path}/constants.xml); [ ! -z \$o1 ] && sed -i \"/NUM_FINAL_BLOCK_PER_POW=/c\NUM_FINAL_BLOCK_PER_POW= \$o1\" ${base_path}/download_incr_DB.py; o1=\$(grep -oPm1 '(?<=<INCRDB_DSNUMS_WITH_STATEDELTAS>)[^<]+' ${base_path}/constants.xml); [ ! -z \$o1 ] && sed -i \"/NUM_DSBLOCK=/c\NUM_DSBLOCK= \$o1\" ${base_path}/download_incr_DB.py"
eval ${cmd}

if [ ! -z "$storage_path" ]; then
 # patch constant for STORAGE_PATH
 tag="STORAGE_PATH"
 tag_value=$(grep "<$tag>.*<.$tag>" constants.xml | sed -e "s/^.*<$tag/<$tag/" | cut -f2 -d">"| cut -f1 -d"<")
 # Replacing element value with new storage path
 sed -i -e "s|<$tag>$tag_value</$tag>|<$tag>$storage_path</$tag>|g" constants.xml
fi

cmd="zilliqa --privk $prikey --pubk $pubkey --address $ip --port $port --synctype 6 --recovery"

if [ "$multiplier_sync" = "N" ] || [ "$multiplier_sync" = "n" ]
then
    cmd="$cmd --l2lsyncmode --extseedprivk $extseedprivk"
fi

$zilliqa_path/build/bin/$cmd > $cmd_log 2>&1 &

echo
echo "Use 'cat $cmd_log' to see the command output"
echo "Use 'tail -f zilliqa-00001-log.txt' to see the runtime log"
}




# ========================================
# SCRIPT START
# ========================================
run
