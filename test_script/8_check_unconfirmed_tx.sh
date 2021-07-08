#!/bin/bash

source ./env.sh

curl http://${HOST[NODE_INDEX]}:$RPC_PORT'/unconfirmed_txs?limit=1'