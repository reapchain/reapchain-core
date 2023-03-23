#!/usr/bin/env bash

set -eo pipefail

proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  protoc \
  -I "proto" \
  -I "third_party/proto" \
  --gogofaster_out=\
Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/duration.proto=github.com/golang/protobuf/ptypes/duration,\
plugins=grpc,paths=source_relative:. \
  $(find "${dir}" -maxdepth 1 -name '*.proto')
done

cp -r ./reapchain-core/* ./proto/*
rm -rf reapchain-core

mv ./proto/reapchain-core/abci/types.pb.go ./abci/types

mv ./proto/reapchain-core/rpc/grpc/types.pb.go ./rpc/grpc
