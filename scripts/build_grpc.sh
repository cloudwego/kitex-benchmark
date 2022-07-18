#!/bin/bash
set -e
GOEXEC=${GOEXEC:-"go"}

# clean
if [ -z "$output_dir" ]; then
  echo "output_dir is empty"
  exit 1
fi
rm -rf $output_dir/bin/ && mkdir -p $output_dir/bin/
rm -rf $output_dir/log/ && mkdir -p $output_dir/log/

# build kitex
$GOEXEC mod edit -replace=github.com/apache/thrift=github.com/apache/thrift@v0.13.0
$GOEXEC mod tidy
$GOEXEC build -v -o $output_dir/bin/kitex_bencher $grpc_dir/kitex/client
$GOEXEC build -v -o $output_dir/bin/kitex_reciever $grpc_dir/kitex

# build others
$GOEXEC mod edit -replace=github.com/apache/thrift=github.com/apache/thrift@v0.14.2
$GOEXEC mod tidy
$GOEXEC build -v -o $output_dir/bin/grpc_bencher $grpc_dir/grpc/client
$GOEXEC build -v -o $output_dir/bin/grpc_reciever $grpc_dir/grpc
