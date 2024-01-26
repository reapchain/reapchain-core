#!/bin/bash
OSTYPE="`uname | tr '[:upper:]' '[:lower:]'`"

SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

NUMBER_OF_STANDING_NODES=$1
NUMBER_OF_STEERING_CANDIDATE_NODES=$2
DIRECTORY="$SCRIPTPATH/$3"

if [[ $DIRECTORY =~ /$ ]]; then
    len=${#DIRECTORY}
    DIRECTORY=${DIRECTORY:0:len-1}    
fi

# Make sure Binaries are present in the bin directory
./bin/manage_binaries.sh

echo ""
echo "NUMBER_OF_STANDING_NODES: $NUMBER_OF_STANDING_NODES"
echo "NUMBER_OF_STEERING_CANDIDATE_NODES: $NUMBER_OF_STEERING_CANDIDATE_NODES"
echo "DIRECTORY: $DIRECTORY"
echo "" 

rm -rf $DIRECTORY
mkdir -p $DIRECTORY 


NODE_COUNT=$((${NUMBER_OF_STANDING_NODES}+${NUMBER_OF_STEERING_CANDIDATE_NODES}))
echo "Total the number of nodes : " $NODE_COUNT 

cp -r ./env.sh $DIRECTORY/env.sh

if [[ "$OSTYPE" == "darwin"* ]]; then
	BINARY_PATH=./bin/mac/podc
elif [[ "$OSTYPE" = "linux"* ]]; then
	BINARY_PATH=./bin/linux/podc
fi

BINARY=podc
sed -i.temp  "s/BINARY=./BINARY=${BINARY//\//\\/}/g" $DIRECTORY/env.sh 
sed -i.temp  "s/DIRECTORY=./DIRECTORY=${DIRECTORY//\//\\/}/g" $DIRECTORY/env.sh 
sed -i.temp  "s/NODE_COUNT=30/NODE_COUNT=${NODE_COUNT}/g" $DIRECTORY/env.sh

sed -i.temp  "s/NUMBER_OF_STANDING_NODES=14/NUMBER_OF_STANDING_NODES=${NUMBER_OF_STANDING_NODES}/g" $DIRECTORY/env.sh
sed -i.temp  "s/NUMBER_OF_STEERING_CANDIDATE_NODES=15/NUMBER_OF_STEERING_CANDIDATE_NODES=${NUMBER_OF_STEERING_CANDIDATE_NODES}/g" $DIRECTORY/env.sh

cp ./2_run.sh $DIRECTORY
cp ./1_init.sh $DIRECTORY

sed -i.temp  "s/Please input server index/Please input server index \(1 ~ ${NODE_COUNT}\)/g" $DIRECTORY/1_init.sh
sed -i.temp  "s/Please input server index/Please input server index \(1 ~ ${NODE_COUNT}\)/g" $DIRECTORY/2_run.sh

source $DIRECTORY/env.sh 
mkdir -p $LOG_DIR 

# generate node keys
i=0
current_steering_member_candidates=0
while [ $i -lt $NODE_COUNT ]; do
	ACCOUNT_NAME="account"$((${i}))
	if [ $i -lt $NUMBER_OF_STANDING_NODES ]; then
		NODE_NAME="standing"$((${i}))
	else
		NODE_NAME="steering"$((${i}))
	fi

  echo $'################ Initialize Node\n'
	echo "$BINARY init --home $DIRECTORY/$NODE_NAME"
	echo "" 

  $BINARY init --home $DIRECTORY/$NODE_NAME
  i=$((${i}+1))
done

echo "################ Make genesis.json"

VALIDATORS=""
STANDING_MEMBERS=""
QRNS=""
STEERING_MEMBER_CANDIDATES=""

i=0
current_steering_member_candidates=0
while [ $i -lt $NODE_COUNT ]; do
	if [ $i -lt $NUMBER_OF_STANDING_NODES ]; then
		NODE_NAME="standing"$((${i}))
	else
		NODE_NAME="steering"$((${i}))
	fi

  if [ $i -lt $NUMBER_OF_STANDING_NODES ]; then
    validator=`cat $DIRECTORY/$NODE_NAME/config/genesis.json | jq .validators[] | jq '.power="10"' | jq '.name="'${NODE_NAME}'"'`
		VALIDATORS=$VALIDATORS$validator","

    standing_member=`cat $DIRECTORY/$NODE_NAME/config/genesis.json | jq .standing_members[] | jq '.power="10"' | jq '.name="'${NODE_NAME}'"'`
		STANDING_MEMBERS=$STANDING_MEMBERS$standing_member","

    qrn=`cat $DIRECTORY/$NODE_NAME/config/genesis.json | jq .qrns[] | jq '.standing_member_index='${i}`
		QRNS=$QRNS$qrn","

  else
		if [ $current_steering_member_candidates -lt 15 ]; then
      validator=`cat $DIRECTORY/$NODE_NAME/config/genesis.json | jq .validators[] | jq '.power="10"' | jq '.name="'${NODE_NAME}'"'`
			VALIDATORS=$VALIDATORS$validator","
    fi

    steering_member_candidate=`cat $DIRECTORY/$NODE_NAME/config/genesis.json | jq .standing_members[] | jq '.power="10"' | jq '.name="'${NODE_NAME}'"'`
		STEERING_MEMBER_CANDIDATES=$STEERING_MEMBER_CANDIDATES$steering_member_candidate","

		current_steering_member_candidates=`expr $current_steering_member_candidates + 1`
  fi

	i=$((${i}+1))
done

SAMPLE_GENESIS_FILE=$SCRIPTPATH/resources/sample_genesis_new.json

VALIDATORS=${VALIDATORS%,}
STANDING_MEMBERS=${STANDING_MEMBERS%,}
QRNS=${QRNS%,}
STEERING_MEMBER_CANDIDATES=${STEERING_MEMBER_CANDIDATES%,}

if [[ "$OSTYPE" == "darwin"* ]]; then
	DATETIME=`gdate -u "+%Y-%m-%dT%H:%M:%S.%8NZ"`
elif [[ "$OSTYPE" = "linux"* ]]; then
	DATETIME=`date -u "+%Y-%m-%dT%H:%M:%S.%8NZ"`
fi

cat $SAMPLE_GENESIS_FILE | jq '.genesis_time="'$DATETIME'"' | jq ".validators=[$VALIDATORS]" | jq ".standing_members=[$STANDING_MEMBERS]" | jq ".qrns=[$QRNS]" | jq ".steering_member_candidates=[$STEERING_MEMBER_CANDIDATES]" > $DIRECTORY/genesis.json

i=0
while [ $i -lt $NODE_COUNT ]; do
	if [ $i -lt $NUMBER_OF_STANDING_NODES ]; then
		NODE_NAME="standing"$((${i}))
	else
		NODE_NAME="steering"$((${i}))
	fi

	NODE_IDS[$i]=`$BINARY show-node-id --home $DIRECTORY/$NODE_NAME`
  echo ${NODE_IDS[$i]}
	sed -i.temp "s/NODE_IDS\[${i}\]=\"\"/NODE_IDS\[${i}\]=\"${NODE_IDS[$i]}\"/g" $DIRECTORY/env.sh
 
	i=$((${i}+1))
done

echo "################ Done"
echo ""

rm -rf $DIRECTORY/env.sh.temp
rm -rf $DIRECTORY/1_init.sh.temp
rm -rf $DIRECTORY/2_run.sh.temp