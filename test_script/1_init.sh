#!/bin/bash

source ./env.sh

if [ $# -eq 0 ] ; then
  echo "Please input server index (1~5)"
  exit
fi

HOST_INDEX=`expr $1 - 1`
LOOP_COUNT=$ALOCATED_NODE_COUNT

if [ $1 -eq $HOST_COUNT ] ; then
  LAST_NODE_COUNT=`expr $NODE_COUNT % $ALOCATED_NODE_COUNT`
  LOOP_COUNT=`expr $ALOCATED_NODE_COUNT + $LAST_NODE_COUNT`
fi              

i=0
while [ $i -lt $LOOP_COUNT ]; do
  CURRENT_HOST_INDEX=`expr $HOST_INDEX \* $ALOCATED_NODE_COUNT`
  CURRENT_HOST_INDEX=`expr $CURRENT_HOST_INDEX + $i`

  CURRENT_RPC_PORT=`expr $RPC_PORT + $i`
  CURRENT_TENDERMINT_PORT=`expr $TENDERMINT_PORT + $i`

  CURRENT_DATA_DIR=$DATA_DIR/$CURRENT_HOST_INDEX
  CURRENT_LOG_DIR=$LOG_DIR/$CURRENT_HOST_INDEX

  rm -rf $CURRENT_DATA_DIR
  rm -rf $CURRENT_LOG_DIR
  mkdir -p $CURRENT_LOG_DIR

  $TENDERMINT init --home $CURRENT_DATA_DIR
  # options
  # --home string        directory for config and data (default "/home/tendermint/.tendermint")
  # --log_level string   log level (default "info")
  # --trace              print out full stack trace on errors
  echo "RPC_PORT: $CURRENT_RPC_PORT"
  echo "TENDERMINT_PORT: $CURRENT_TENDERMINT_PORT"

  echo "cp ./genesis_block/genesis.json $CURRENT_DATA_DIR/config"
  cp ./genesis_block/genesis.json $CURRENT_DATA_DIR/config

  echo "cp ./node_keys/node_${CURRENT_HOST_INDEX}.json $CURRENT_DATA_DIR/config/node_key.json"
  cp ./node_keys/node_${CURRENT_HOST_INDEX}.json $CURRENT_DATA_DIR/config/node_key.json

  echo "cp ./priv_validator_keys/node_${CURRENT_HOST_INDEX}.json $CURRENT_DATA_DIR/config/priv_validator_key.json"
  cp ./priv_validator_keys/node_${CURRENT_HOST_INDEX}.json $CURRENT_DATA_DIR/config/priv_validator_key.json

  # [setting 관련] ------------
  # db 사용 -> cleveldb 
  sed -i 's/db_backend = "goleveldb"/db_backend = "cleveldb"/g' $CURRENT_DATA_DIR/config/config.toml

  # ---------------------------

  # [consensus 관련]-----------
  # How long we wait for a proposal block before prevoting nil
  sed -i 's/timeout_propose = "3s"/timeout_propose = "6s"/g' $CURRENT_DATA_DIR/config/config.toml

  # How much timeout_propose increases with each round
  sed -i 's/timeout_propose_delta = "1000ms"/timeout_propose_delta = "1000ms"/g' $CURRENT_DATA_DIR/config/config.toml

  # How long we wait after receiving +2/3 prevotes for “anything” (ie. not a single block or nil)
  sed -i 's/timeout_prevote = "1s"/timeout_prevote = "2s"/g' $CURRENT_DATA_DIR/config/config.toml

  # How much the timeout_prevote increases with each round
  sed -i 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "1000ms"/g' $CURRENT_DATA_DIR/config/config.toml

  # How long we wait after receiving +2/3 precommits for “anything” (ie. not a single block or nil)
  sed -i 's/timeout_precommit = "1s"/timeout_precommit = "2s"/g' $CURRENT_DATA_DIR/config/config.toml

  # How much the timeout_precommit increases with each round
  sed -i 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "1000ms"/g' $CURRENT_DATA_DIR/config/config.toml

  # How long we wait after committing a block, before starting on the new
  # height (this gives us a chance to receive some more precommits, even
  # though we already have +2/3).
  sed -i 's/timeout_commit = "1s"/timeout_commit = "2s"/g' $CURRENT_DATA_DIR/config/config.toml

  # 빈블록 생성 여부 -> False
  sed -i 's/create_empty_blocks = true/create_empty_blocks = false/g' $CURRENT_DATA_DIR/config/config.toml

  # --------------------------

  # [mempool 관련]------------
  # mempool의 최대 트랜잭션 수
  sed -i 's/size = 5000/size = 1100000/g' $CURRENT_DATA_DIR/config/config.toml

  # 트랜잭션의 캐시 크기
  sed -i 's/cache_size = 10000/cache_size = 100000/g' $CURRENT_DATA_DIR/config/config.toml

  # --------------------------
  sed -i 's/allow_duplicate_ip = false/allow_duplicate_ip = true/g' $CURRENT_DATA_DIR/config/config.toml
  sed -i 's/max_num_inbound_peers = 40/max_num_inbound_peers = 50/g' $CURRENT_DATA_DIR/config/config.toml
  sed -i 's/max_num_outbound_peers = 10/max_num_outbound_peers = 50/g' $CURRENT_DATA_DIR/config/config.toml
  
  
  sed -i 's/max_tx_bytes = 1048576/max_tx_bytes = 3900/g' $CURRENT_DATA_DIR/config/config.toml
  sed -i 's/max_txs_bytes = 1073741824/max_txs_bytes = 4294967296/g' $CURRENT_DATA_DIR/config/config.toml
  

  # 자신의 RCP 주소 설정
  sed -i "s/laddr = \"tcp:\/\/127.0.0.1:26657\"/laddr = \"tcp:\/\/${HOST}:${CURRENT_RPC_PORT}\"/g" $CURRENT_DATA_DIR/config/config.toml
  
  # TENDERMINT 주소 설정
  sed -i "s/laddr = \"tcp:\/\/0.0.0.0:26656\"/laddr = \"tcp:\/\/0.0.0.0:${CURRENT_TENDERMINT_PORT}\"/g" $CURRENT_DATA_DIR/config/config.toml 

  i=$((${i}+1))
done