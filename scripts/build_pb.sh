#!/bin/bash
set -e
CURDIR=$(cd $(dirname $0); pwd)

GOEXEC=${GOEXEC:-"go"}

# clean
rm -rf $CURDIR/../output/ && mkdir -p $CURDIR/../output/bin/ && mkdir -p $CURDIR/../output/log/

# build kitex
$GOEXEC mod edit -replace=github.com/apache/thrift=github.com/apache/thrift@v0.13.0
$GOEXEC mod tidy
go build -v -o $CURDIR/../output/bin/kitex-pb_bencher $CURDIR/../protobuf/kitex-pb/client
go build -v -o $CURDIR/../output/bin/kitex-grpc_bencher $CURDIR/../protobuf/kitex-grpc/client
go build -v -o $CURDIR/../output/bin/kitex-mux_bencher $CURDIR/../protobuf/kitex-mux/client
go build -v -o $CURDIR/../output/bin/kitex-pb_reciever $CURDIR/../protobuf/kitex-pb
go build -v -o $CURDIR/../output/bin/kitex-grpc_reciever $CURDIR/../protobuf/kitex-grpc
go build -v -o $CURDIR/../output/bin/kitex-mux_reciever $CURDIR/../protobuf/kitex-mux

# build others
$GOEXEC mod edit -replace=github.com/apache/thrift=github.com/apache/thrift@v0.14.2
$GOEXEC mod tidy
go build -v -o $CURDIR/../output/bin/grpc_bencher $CURDIR/../protobuf/grpc/client
go build -v -o $CURDIR/../output/bin/rpcx_bencher $CURDIR/../protobuf/rpcx/client
go build -v -o $CURDIR/../output/bin/arpc_bencher $CURDIR/../protobuf/arpc/client
go build -v -o $CURDIR/../output/bin/grpc_reciever $CURDIR/../protobuf/grpc
go build -v -o $CURDIR/../output/bin/rpcx_reciever $CURDIR/../protobuf/rpcx
go build -v -o $CURDIR/../output/bin/arpc_reciever $CURDIR/../protobuf/arpc
