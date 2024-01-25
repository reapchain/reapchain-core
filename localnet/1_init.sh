#!/bin/bash
source "$1"/env.sh

if [ ! $2 ]; then
  echo "Please input server index"
  exit
fi

INDEX=`expr $2`
if [ $INDEX -lt $NUMBER_OF_STANDING_NODES ]; then
		NODE_NAME="standing"$((${INDEX}))
else
  NODE_NAME="steering"$((${INDEX}))
fi
CURRENT_DATA_DIR=$DATA_DIR/$NODE_NAME

echo "NODE_NAME: ${NODE_NAME}"
echo "CURRENT_DATA_DIR: ${CURRENT_DATA_DIR}"

echo "RPC_PORT: ${RPC_PORTS[$INDEX]}"
echo "REAPCHAIN_PORT: ${REAPCHAIN_PORTS[$INDEX]}"

echo "cp $1/genesis.json $CURRENT_DATA_DIR/config/genesis.json"
cp $1/genesis.json $CURRENT_DATA_DIR/config/genesis.json

sed -i.temp "s/laddr = \"tcp:\/\/127.0.0.1:26657\"/laddr = \"tcp:\/\/${HOST}:${RPC_PORTS[$INDEX]}\"/g" $CURRENT_DATA_DIR/config/config.toml
sed -i.temp "s/laddr = \"tcp:\/\/0.0.0.0:26656\"/laddr = \"tcp:\/\/${HOST}:${REAPCHAIN_PORTS[$INDEX]}\"/g" $CURRENT_DATA_DIR/config/config.toml 
sed -i.temp 's/allow_duplicate_ip = false/allow_duplicate_ip = true/g' $CURRENT_DATA_DIR/config/config.toml
# Mempool Size
sed -i.temp 's/size = 5000/size = 500000/g' $CURRENT_DATA_DIR/config/config.toml
# Create empty block
sed -i.temp 's/create_empty_blocks = true/create_empty_blocks = true/g' $CURRENT_DATA_DIR/config/config.toml


# Update peer connect type
if [ $USE_SEED_NODE == "enable" ]; then
  sed -i.temp "s/seeds = .*/seeds = \"$SEED_NODE_ID@$SEED_HOST:$SEED_P2P_PORT\"/g" $CURRENT_DATA_DIR/config/config.toml
  echo " > used seed mode"
else
  echo "Update persistent_peers"
  sed -i.temp "s/seeds = .*/seeds = \"\"/g" $CURRENT_DATA_DIR/config/config.toml

  PERSISTENT_PEERS=""
  i=0
  while [ $i -lt $NODE_COUNT ]; do
      PERSISTENT_PEERS=$PERSISTENT_PEERS${NODE_IDS[$i]}'@'${HOST}':'${REAPCHAIN_PORTS[$i]}','
      i=$((${i}+1))
  done  
  PERSISTENT_PEERS=${PERSISTENT_PEERS%,} #%, = ,를 제거하겠다는 의미
  echo "PERSISTENT_PEERS : "$PERSISTENT_PEERS

  sed -i.temp "s/persistent_peers = .*/persistent_peers = \"$PERSISTENT_PEERS\"/g" $CURRENT_DATA_DIR/config/config.toml

  echo " > used persistent_peers"
fi

rm -rf $CURRENT_DATA_DIR/config/app.toml.temp
rm -rf $CURRENT_DATA_DIR/config/config.toml.temp
rm -rf $CURRENT_DATA_DIR/config/genesis.json.temp 