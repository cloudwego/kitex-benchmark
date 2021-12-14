#!/bin/bash

GOEXEC=${GOEXEC:-"go"}

# clean
rm -rf output/bin/ && mkdir -p output/bin/
rm -rf output/log/ && mkdir -p output/log/

$GOEXEC mod edit -replace=github.com/apache/thrift=github.com/apache/thrift@v0.13.0
$GOEXEC mod tidy

# build clients
$GOEXEC build -v -o output/bin/kitex_bencher ./thrift/kitex/client
$GOEXEC build -v -o output/bin/kitex-mux_bencher ./thrift/kitex-mux/client

# build servers
$GOEXEC build -v -o output/bin/kitex_reciever ./thrift/kitex
$GOEXEC build -v -o output/bin/kitex-mux_reciever ./thrift/kitex-mux
