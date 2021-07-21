#!/bin/bash

. ./scripts/env.sh
repo=("kitex" "kitex-mux")
ports=(8001 8002)
ip=${IP:-"127.0.0.1"}

. ./scripts/build_thrift.sh

# benchmark
for b in ${body[@]}; do
  for c in ${concurrent[@]}; do
    for ((i = 0; i < ${#repo[@]}; i++)); do
      rp=${repo[i]}
      addr="${ip}:${ports[i]}"

      # run client
      taskset -c 0-20 ./output/bin/${rp}_bencher -addr="$addr" -b=$b -c=$c -n=$n
    done
  done
done
