#!/bin/bash

source ./env.sh

curl http://${HOST[NODE_INDEX]}:$RPC_PORT"/block?height=$1"