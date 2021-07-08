#!/bin/bash

source ./env.sh

echo http://${HOST[NODE_INDEX]}:$RPC_PORT'/validators?height=1&page=1&per_page=30'
curl http://${HOST[NODE_INDEX]}:$RPC_PORT'/validators?height=1&page=1&per_page=30'