#!/bin/bash

source ./env.sh

index=$1

echo "http://${HOST[index]}:$RPC_PORT/dump_consensus_state"
curl http://${HOST[index]}:$RPC_PORT/dump_consensus_state


