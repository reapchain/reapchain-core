#!/bin/bash
SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

# ./launch_local_testnet.sh [NUMBER_OF_STANDING_NODES] [NUMBER_OF_STEERING_CANDIDATE_NODES] [DIRECTORY]
if [[ ($1 -eq 0) || ($2 -eq 0) || (-z "$3") ]]; then
	echo $'\nPlease input the number of standing and steering nodes and also a valid directory.'
	echo $'\n./launch_local_testnet.sh [NUMBER_OF_STANDING_NODES] [NUMBER_OF_STEERING_CANDIDATE_NODES] [DIRECTORY]\n'
	exit
fi

NUMBER_OF_STANDING_NODES=$1
NUMBER_OF_STEERING_CANDIDATE_NODES=$2
DIRECTORY=$3

./10_setup.sh $NUMBER_OF_STANDING_NODES $NUMBER_OF_STEERING_CANDIDATE_NODES $DIRECTORY
./11_init_all.sh $DIRECTORY
./12_run_all.sh $DIRECTORY
