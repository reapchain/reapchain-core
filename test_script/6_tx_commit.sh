#!/bin/bash

source ./env.sh

curl http://${HOST[NODE_INDEX]}:$RPC_PORT'/broadcast_tx_commit?tx="345345"'
# curl http://${HOST[NODE_INDEX]}:$RPC_PORT'/broadcast_tx_commit?tx="name=abcd"'