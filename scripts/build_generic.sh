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
$GOEXEC build -v -o $output_dir/bin/generic_binary_bencher $generic_dir/binary/client
$GOEXEC build -v -o $output_dir/bin/generic_http_bencher $generic_dir/http/client
$GOEXEC build -v -o $output_dir/bin/generic_map_bencher $generic_dir/map/client
$GOEXEC build -v -o $output_dir/bin/generic_json_bencher $generic_dir/json/client
$GOEXEC build -v -o $output_dir/bin/generic_ordinary_bencher $generic_dir/ordinary/client

# build servers
$GOEXEC build -v -o $output_dir/bin/generic_binary_reciever $generic_dir/binary
$GOEXEC build -v -o $output_dir/bin/generic_http_reciever $generic_dir/http
$GOEXEC build -v -o $output_dir/bin/generic_map_reciever $generic_dir/map
$GOEXEC build -v -o $output_dir/bin/generic_json_reciever $generic_dir/json
$GOEXEC build -v -o $output_dir/bin/generic_ordinary_reciever $generic_dir/ordinary