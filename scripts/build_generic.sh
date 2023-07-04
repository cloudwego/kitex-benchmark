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
$GOEXEC build -v -o $output_dir/bin/generic_http_1KB_bencher $generic_dir/http/client/default/small
$GOEXEC build -v -o $output_dir/bin/generic_http_5KB_bencher $generic_dir/http/client/default/medium
$GOEXEC build -v -o $output_dir/bin/generic_http_10KB_bencher $generic_dir/http/client/default/large
$GOEXEC build -v -o $output_dir/bin/generic_http_dynamicgo_1KB_bencher $generic_dir/http/client/dynamicgo/small
$GOEXEC build -v -o $output_dir/bin/generic_http_dynamicgo_5KB_bencher $generic_dir/http/client/dynamicgo/medium
$GOEXEC build -v -o $output_dir/bin/generic_http_dynamicgo_10KB_bencher $generic_dir/http/client/dynamicgo/large
$GOEXEC build -v -o $output_dir/bin/generic_map_1KB_bencher $generic_dir/map/client/small
$GOEXEC build -v -o $output_dir/bin/generic_map_5KB_bencher $generic_dir/map/client/medium
$GOEXEC build -v -o $output_dir/bin/generic_map_10KB_bencher $generic_dir/map/client/large
$GOEXEC build -v -o $output_dir/bin/generic_json_1KB_bencher $generic_dir/json/client/default/small
$GOEXEC build -v -o $output_dir/bin/generic_json_5KB_bencher $generic_dir/json/client/default/medium
$GOEXEC build -v -o $output_dir/bin/generic_json_10KB_bencher $generic_dir/json/client/default/large
$GOEXEC build -v -o $output_dir/bin/generic_json_dynamicgo_1KB_bencher $generic_dir/json/client/dynamicgo/small
$GOEXEC build -v -o $output_dir/bin/generic_json_dynamicgo_5KB_bencher $generic_dir/json/client/dynamicgo/medium
$GOEXEC build -v -o $output_dir/bin/generic_json_dynamicgo_10KB_bencher $generic_dir/json/client/dynamicgo/large
$GOEXEC build -v -o $output_dir/bin/generic_ordinary_1KB_bencher $generic_dir/ordinary/client/small
$GOEXEC build -v -o $output_dir/bin/generic_ordinary_5KB_bencher $generic_dir/ordinary/client/medium
$GOEXEC build -v -o $output_dir/bin/generic_ordinary_10KB_bencher $generic_dir/ordinary/client/large

# build servers
$GOEXEC build -v -o $output_dir/bin/generic_http_1KB_reciever $generic_dir/http
$GOEXEC build -v -o $output_dir/bin/generic_http_5KB_reciever $generic_dir/http
$GOEXEC build -v -o $output_dir/bin/generic_http_10KB_reciever $generic_dir/http
$GOEXEC build -v -o $output_dir/bin/generic_map_1KB_reciever $generic_dir/map
$GOEXEC build -v -o $output_dir/bin/generic_map_5KB_reciever $generic_dir/map
$GOEXEC build -v -o $output_dir/bin/generic_map_10KB_reciever $generic_dir/map
$GOEXEC build -v -o $output_dir/bin/generic_json_1KB_reciever $generic_dir/json/server/default/small
$GOEXEC build -v -o $output_dir/bin/generic_json_5KB_reciever $generic_dir/json/server/default/medium
$GOEXEC build -v -o $output_dir/bin/generic_json_10KB_reciever $generic_dir/json/server/default/large
$GOEXEC build -v -o $output_dir/bin/generic_json_dynamicgo_1KB_reciever $generic_dir/json/server/dynamicgo/small
$GOEXEC build -v -o $output_dir/bin/generic_json_dynamicgo_5KB_reciever $generic_dir/json/server/dynamicgo/medium
$GOEXEC build -v -o $output_dir/bin/generic_json_dynamicgo_10KB_reciever $generic_dir/json/server/dynamicgo/large
$GOEXEC build -v -o $output_dir/bin/generic_ordinary_1KB_reciever $generic_dir/ordinary
$GOEXEC build -v -o $output_dir/bin/generic_ordinary_5KB_reciever $generic_dir/ordinary
$GOEXEC build -v -o $output_dir/bin/generic_ordinary_10KB_reciever $generic_dir/ordinary