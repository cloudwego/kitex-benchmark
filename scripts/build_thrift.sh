#!/bin/bash
set -e
CURDIR=$(cd $(dirname $0); pwd)

GOEXEC=${GOEXEC:-"go"}

# clean
rm -rf $CURDIR/../output/ && mkdir -p $CURDIR/../output/bin/ && mkdir -p $CURDIR/../output/log/

$GOEXEC mod edit -replace=github.com/apache/thrift=github.com/apache/thrift@v0.13.0
$GOEXEC mod tidy

# build clients
go build -v -o $CURDIR/../output/bin/kitex_bencher $CURDIR/../thrift/kitex/client
go build -v -o $CURDIR/../output/bin/kitex-mux_bencher $CURDIR/../thrift/kitex-mux/client

# build servers
go build -v -o $CURDIR/../output/bin/kitex_reciever $CURDIR/../thrift/kitex
go build -v -o $CURDIR/../output/bin/kitex-mux_reciever $CURDIR/../thrift/kitex-mux
