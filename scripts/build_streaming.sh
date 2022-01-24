#!/bin/bash
set -e
GOEXEC=${GOEXEC:-"go"}

# clean
rm -rf $output_dir/bin/ && mkdir -p $output_dir/bin/
rm -rf $output_dir/log/ && mkdir -p $output_dir/log/

# build kitex
$GOEXEC mod tidy
$GOEXEC build -v -o $output_dir/bin/kitex_bencher $streaming_dir/kitex/client
$GOEXEC build -v -o $output_dir/bin/kitex_reciever $streaming_dir/kitex

# build grpc
$GOEXEC mod tidy
$GOEXEC build -v -o $output_dir/bin/grpc_bencher $streaming_dir/grpc/client
$GOEXEC build -v -o $output_dir/bin/grpc_reciever $streaming_dir/grpc
