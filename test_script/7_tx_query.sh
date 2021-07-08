#!/bin/bash

source ./env.sh

curl http://${HOST[NODE_INDEX]}:$RPC_PORT'/abci_query?data="345345"'
curl http://${HOST[NODE_INDEX]}:$RPC_PORT'/abci_query?data="name"'

unconfirmed_txs?limit=1