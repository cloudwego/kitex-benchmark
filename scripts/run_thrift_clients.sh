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
      echo "client $rp running with $cmd_client"
      $cmd_client ./output/bin/${rp}_bencher -addr="$addr" -b=$b -c=$c -n=$n --sleep=$sleep | $tee_cmd
    done
  done
done

finish_cmd
