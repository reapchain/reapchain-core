#!/bin/bash

source ./env.sh

i=0
while [ $i -lt $1 ]; do
    $TENDERMINT --home ./ gen-node-key
    cp ./config/node_key.json $SCRIPT_DIR/node_keys/node_$i.json
    $TENDERMINT --home ./ gen-validator > $SCRIPT_DIR/priv_validator_keys/node_$i.json
    
    cat $SCRIPT_DIR/priv_validator_keys/node_$i.json
    rm -rf config
    rm -rf data
    i=$((${i}+1))
done