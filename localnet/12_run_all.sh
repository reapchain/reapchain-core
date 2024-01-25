#!/bin/bash

# ./12_run_all.sh [DIRECTORY]
if [[(-z "$1") ]]; then
	echo $'\nPlease input the working directory.'
	echo $'\n./12_run_all.sh [DIRECTORY]\n'
	exit
fi

echo "################ Run All Nodes"
source "./$1"/env.sh

i=0
while [ $i -lt $NODE_COUNT ]; do
  NODE_NUM=$((${i}))
  echo ""
  echo "Starting Node$NODE_NUM"
  echo "RPC - http://${HOST}:${RPC_PORTS[$NODE_NUM]}"
  echo "P2P - ${HOST}:${REAPCHAIN_PORTS[$NODE_NUM]}"
  
  $DIRECTORY/2_run.sh $DIRECTORY b $NODE_NUM
  i=$((${i}+1))
done


 
