#!/bin/bash

# clean
rm -rf output/ && mkdir -p output/bin/ && mkdir -p output/log/

# build clients
go mod edit -replace=github.com/apache/thrift=github.com/apache/thrift@v0.13.0
go build -v -o output/bin/kitex_bencher ./thrift/kitex/client
go build -v -o output/bin/kitex-mux_bencher ./thrift/kitex-mux/client

# build servers
go mod edit -replace=github.com/apache/thrift=github.com/apache/thrift@v0.13.0
go build -v -o output/bin/kitex_reciever ./thrift/kitex
go build -v -o output/bin/kitex-mux_reciever ./thrift/kitex-mux
