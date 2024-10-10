#!/bin/bash
set -e
GOEXEC=${GOEXEC:-"go"}

# clean
rm -rf $output_dir/bin/ && mkdir -p $output_dir/bin/
rm -rf $output_dir/log/ && mkdir -p $output_dir/log/

# build grpc
$GOEXEC mod tidy
$GOEXEC build -v -o $output_dir/bin/grpc_bencher $streaming_dir/grpc/client
$GOEXEC build -v -o $output_dir/bin/grpc_reciever $streaming_dir/grpc

# build kitex_grpc
$GOEXEC mod tidy
$GOEXEC build -v -o $output_dir/bin/kitex_grpc_bencher $streaming_dir/kitex_grpc/client
$GOEXEC build -v -o $output_dir/bin/kitex_grpc_reciever $streaming_dir/kitex_grpc

# build kitex_tts_lconn
$GOEXEC mod tidy
$GOEXEC build -v -o $output_dir/bin/kitex_tts_lconn_bencher $streaming_dir/kitex_tts_lconn/client
$GOEXEC build -v -o $output_dir/bin/kitex_tts_lconn_reciever $streaming_dir/kitex_tts_lconn

# build kitex_tts_mux
$GOEXEC mod tidy
$GOEXEC build -v -o $output_dir/bin/kitex_tts_mux_bencher $streaming_dir/kitex_tts_mux/client
$GOEXEC build -v -o $output_dir/bin/kitex_tts_mux_reciever $streaming_dir/kitex_tts_mux

