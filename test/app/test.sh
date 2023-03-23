#! /bin/bash
set -ex

#- kvstore over socket, curl
#- counter over socket, curl
#- counter over grpc, curl
#- counter over grpc, grpc

# TODO: install everything

export PATH="$GOBIN:$PATH"
export TMHOME=$HOME/.reapchain-core_app

function kvstore_over_socket(){
    rm -rf $TMHOME
    reapchain-core init
    echo "Starting kvstore_over_socket"
    abci-cli kvstore > /dev/null &
    pid_kvstore=$!
    reapchain-core node > reapchain-core.log &
    pid_reapchain-core=$!
    sleep 5

    echo "running test"
    bash test/app/kvstore_test.sh "KVStore over Socket"

    kill -9 $pid_kvstore $pid_reapchain-core
}

# start reapchain-core first
function kvstore_over_socket_reorder(){
    rm -rf $TMHOME
    reapchain-core init
    echo "Starting kvstore_over_socket_reorder (ie. start reapchain-core first)"
    reapchain-core node > reapchain-core.log &
    pid_reapchain-core=$!
    sleep 2
    abci-cli kvstore > /dev/null &
    pid_kvstore=$!
    sleep 5

    echo "running test"
    bash test/app/kvstore_test.sh "KVStore over Socket"

    kill -9 $pid_kvstore $pid_reapchain-core
}


function counter_over_socket() {
    rm -rf $TMHOME
    reapchain-core init
    echo "Starting counter_over_socket"
    abci-cli counter --serial > /dev/null &
    pid_counter=$!
    reapchain-core node > reapchain-core.log &
    pid_reapchain-core=$!
    sleep 5

    echo "running test"
    bash test/app/counter_test.sh "Counter over Socket"

    kill -9 $pid_counter $pid_reapchain-core
}

function counter_over_grpc() {
    rm -rf $TMHOME
    reapchain-core init
    echo "Starting counter_over_grpc"
    abci-cli counter --serial --abci grpc > /dev/null &
    pid_counter=$!
    reapchain-core node --abci grpc > reapchain-core.log &
    pid_reapchain-core=$!
    sleep 5

    echo "running test"
    bash test/app/counter_test.sh "Counter over GRPC"

    kill -9 $pid_counter $pid_reapchain-core
}

function counter_over_grpc_grpc() {
    rm -rf $TMHOME
    reapchain-core init
    echo "Starting counter_over_grpc_grpc (ie. with grpc broadcast_tx)"
    abci-cli counter --serial --abci grpc > /dev/null &
    pid_counter=$!
    sleep 1
    GRPC_PORT=36656
    reapchain-core node --abci grpc --rpc.grpc_laddr tcp://localhost:$GRPC_PORT > reapchain-core.log &
    pid_reapchain-core=$!
    sleep 5

    echo "running test"
    GRPC_BROADCAST_TX=true bash test/app/counter_test.sh "Counter over GRPC via GRPC BroadcastTx"

    kill -9 $pid_counter $pid_reapchain-core
}

case "$1" in 
    "kvstore_over_socket")
    kvstore_over_socket
    ;;
"kvstore_over_socket_reorder")
    kvstore_over_socket_reorder
    ;;
    "counter_over_socket")
    counter_over_socket
    ;;
"counter_over_grpc")
    counter_over_grpc
    ;;
    "counter_over_grpc_grpc")
    counter_over_grpc_grpc
    ;;
*)
    echo "Running all"
    kvstore_over_socket
    echo ""
    kvstore_over_socket_reorder
    echo ""
    counter_over_socket
    echo ""
    counter_over_grpc
    echo ""
    counter_over_grpc_grpc
esac

