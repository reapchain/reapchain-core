#!/bin/bash

source ./env.sh

if [ $# -eq 0 ] ; then             
    echo "Please input node index (0~30)"
    exit
fi

CURRENT_NODE_INDEX=$1

CURRENT_DATA_DIR=$DATA_DIR/$CURRENT_NODE_INDEX
CURRENT_LOG_DIR=$LOG_DIR/$CURRENT_NODE_INDEX

PERSISTENT_PEERS=""
CURRENT_TENDERMINT_PORT=$TENDERMINT_PORT
i=0
j=0
while [ $i -lt $NODE_COUNT ]; do
    if [ $i -ne $NODE_INDEX ]; then
        PERSISTENT_PEERS=$PERSISTENT_PEERS${NODE_IDS[i]}"@"${HOSTS[j]}":"$CURRENT_TENDERMINT_PORT","
    fi

    i=$((${i}+1))
    CURRENT_TENDERMINT_PORT=$((${CURRENT_TENDERMINT_PORT}+1))

    if [ $j -ne `expr $HOST_COUNT - 1` ]; then
        check=`expr $i % $ALOCATED_NODE_COUNT`
        if [ $check -eq 0 ] ; then
            j=$((${j}+1))
            CURRENT_TENDERMINT_PORT=$TENDERMINT_PORT
        fi    
    fi    
done

if [ $CURRENT_NODE_INDEX -eq 0 ] ; then
    echo "$TENDERMINT start --proxy_app=kvstore --home $CURRENT_DATA_DIR"
    $TENDERMINT start --proxy_app=kvstore --home $CURRENT_DATA_DIR
fi
if [ $CURRENT_NODE_INDEX -ne 0 ] ; then
    echo "$TENDERMINT start --proxy_app=kvstore --home $CURRENT_DATA_DIR --p2p.persistent_peers d25dea4666e0894b7051582fb69bbac8e4040c46@147.46.240.248:27000"
    $TENDERMINT start --proxy_app=kvstore --home $CURRENT_DATA_DIR --p2p.persistent_peers d25dea4666e0894b7051582fb69bbac8e4040c46@147.46.240.248:27000
fi