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
$GOEXEC build -v -o $output_dir/bin/kitex_bencher $pb_dir/kitex/client
$GOEXEC build -v -o $output_dir/bin/kitex-mux_bencher $pb_dir/kitex-mux/client
$GOEXEC build -v -o $output_dir/bin/kitex_reciever $pb_dir/kitex
$GOEXEC build -v -o $output_dir/bin/kitex-mux_reciever $pb_dir/kitex-mux

# build others
$GOEXEC mod edit -replace=github.com/apache/thrift=github.com/apache/thrift@v0.14.2
$GOEXEC mod tidy
$GOEXEC build -v -o $output_dir/bin/grpc_bencher $pb_dir/grpc/client
$GOEXEC build -v -o $output_dir/bin/rpcx_bencher $pb_dir/rpcx/client
$GOEXEC build -v -o $output_dir/bin/arpc_bencher $pb_dir/arpc/client
$GOEXEC build -v -o $output_dir/bin/grpc_reciever $pb_dir/grpc
$GOEXEC build -v -o $output_dir/bin/rpcx_reciever $pb_dir/rpcx
$GOEXEC build -v -o $output_dir/bin/arpc_reciever $pb_dir/arpc
