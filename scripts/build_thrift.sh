#!/bin/bash
set -e
CURDIR=$(cd $(dirname $0); pwd)

# clean
rm -rf output/ && mkdir -p output/bin/ && mkdir -p output/log/

go mod edit -replace=github.com/apache/thrift=github.com/apache/thrift@v0.13.0
go mod tidy

# build clients
go build -v -o output/bin/kitex_bencher $CURDIR/../thrift/kitex/client
go build -v -o output/bin/kitex-mux_bencher $CURDIR/../thrift/kitex-mux/client

# build servers
go build -v -o output/bin/kitex_reciever $CURDIR/../thrift/kitex
go build -v -o output/bin/kitex-mux_reciever $CURDIR/../thrift/kitex-mux
