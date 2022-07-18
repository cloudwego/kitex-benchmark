#!/bin/bash
set -e
CURDIR=$(cd $(dirname $0); pwd)
GOEXEC=${GOEXEC:-"go"}

# clean
if [ -z "$output_dir" ]; then
  echo "output_dir is empty"
  exit 1
fi
rm -rf $output_dir/bin/ && mkdir -p $output_dir/bin/
rm -rf $output_dir/log/ && mkdir -p $output_dir/log/

$GOEXEC mod edit -replace=github.com/apache/thrift=github.com/apache/thrift@v0.13.0
$GOEXEC mod tidy

# build clients
$GOEXEC build -v -o $output_dir/bin/kitex_bencher $thrift_dir/kitex/client
$GOEXEC build -v -o $output_dir/bin/kitex-mux_bencher $thrift_dir/kitex-mux/client

# build servers
$GOEXEC build -v -o $output_dir/bin/kitex_reciever $thrift_dir/kitex
$GOEXEC build -v -o $output_dir/bin/kitex-mux_reciever $thrift_dir/kitex-mux
