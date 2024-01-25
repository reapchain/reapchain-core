#!/bin/bash

# ./11_init_all.sh [DIRECTORY]
if [[ (-z "$1") ]]; then
	echo $'\nPlease input the working directory.'
	echo $'\n./11_init_all.sh [DIRECTORY]\n'
	exit
fi

echo "################ Initialize Nodes"
source "./$1"/env.sh

if [ $USE_SEED_NODE == "enable" ]; then
  $DIRECTORY/31_seed_node_init.sh $DIRECTORY

  CURRENT_DATA_DIR="$DATA_DIR/seed_node"
  seed_node_id=`$BINARY --home $CURRENT_DATA_DIR show-node-id`
  sed -i "s/SEED_NODE_ID=\"\"/SEED_NODE_ID=\"$seed_node_id\"/g" "./$1"/env.sh
  echo " > used seed mode"
fi

i=0
while [ $i -lt $NODE_COUNT ]; do
  NODE_NUM=$((${i}))
  echo ""
  echo "Initialize Directory: $DIRECTORY/1_init.sh for Node$NODE_NUM" 
  echo ""
  $DIRECTORY/1_init.sh $DIRECTORY $NODE_NUM
  i=$((${i}+1))
done