#!/bin/bash

. ./scripts/env.sh

repo=("grpc" "kitex" "kitex-mux" "rpcx" "arpc")
ports=(8000 8001 8002 8003 8004)

ip=${IP:-"127.0.0.1"}

# build
. ./scripts/build_pb.sh

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
