#!/bin/bash

# clean
rm -rf output/ && mkdir -p output/bin/ && mkdir -p output/log/

# build clients
go mod edit -replace=github.com/apache/thrift=github.com/apache/thrift@v0.13.0
go build -v -o output/bin/kitex_bencher ./protobuf/kitex/client
go build -v -o output/bin/kitex-mux_bencher ./protobuf/kitex-mux/client

go mod edit -replace=github.com/apache/thrift=github.com/apache/thrift@v0.14.2
go build -v -o output/bin/grpc_bencher ./protobuf/grpc/client
go build -v -o output/bin/rpcx_bencher ./protobuf/rpcx/client
go build -v -o output/bin/arpc_bencher ./protobuf/arpc/client

# build servers
go mod edit -replace=github.com/apache/thrift=github.com/apache/thrift@v0.13.0
go build -v -o output/bin/kitex_reciever ./protobuf/kitex
go build -v -o output/bin/kitex-mux_reciever ./protobuf/kitex-mux

go mod edit -replace=github.com/apache/thrift=github.com/apache/thrift@v0.14.2
go build -v -o output/bin/grpc_reciever ./protobuf/grpc
go build -v -o output/bin/rpcx_reciever ./protobuf/rpcx
go build -v -o output/bin/arpc_reciever ./protobuf/arpc