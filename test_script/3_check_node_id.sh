#!/bin/bash

source ./env.sh


CURRENT_NODE_INDEX=$1

CURRENT_DATA_DIR=$DATA_DIR/$CURRENT_NODE_INDEX

echo "NODE_ID"
$TENDERMINT --home $CURRENT_DATA_DIR show-node-id