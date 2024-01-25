#!/bin/bash
if [ ! $3 ]; then
  echo "Please input server index"
  exit
fi
SCRIPT_PATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
source "$SCRIPT_PATH"/env.sh

RUN_TYPE=$2
INDEX=`expr $3`

if [ $INDEX -lt $NUMBER_OF_STANDING_NODES ]; then
		NODE_NAME="standing"$((${INDEX}))
else
  NODE_NAME="steering"$((${INDEX}))
fi
LOG_NAME="${NODE_NAME}_log"
CURRENT_DATA_DIR=$DATA_DIR/$NODE_NAME
 
if [ $RUN_TYPE = "f" ]; then
  echo "$BINARY start --home $CURRENT_DATA_DIR "
  $BINARY start --home $CURRENT_DATA_DIR --proxy_app=kvstore
else
  echo "$BINARY start --home $CURRENT_DATA_DIR --proxy_app=kvstore &> $LOG_DIR/$LOG_NAME.log &"
  $BINARY start --home $CURRENT_DATA_DIR --proxy_app=kvstore &> $LOG_DIR/$LOG_NAME.log &
fi