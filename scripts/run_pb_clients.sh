#!/bin/bash
set -e
CURDIR=$(cd $(dirname $0); pwd)

source $CURDIR/env.sh

repo=("grpc" "kitex-pb" "kitex-mux" "rpcx" "arpc" "kitex-grpc")
ports=(8000 8001 8002 8003 8004 8005)

ip=${IP:-"127.0.0.1"}

# build
source $CURDIR/build_pb.sh

# benchmark
for b in ${body[@]}; do
  for c in ${concurrent[@]}; do
    for ((i = 0; i < ${#repo[@]}; i++)); do
      rp=${repo[i]}
      addr="${ip}:${ports[i]}"

      # run client
      echo "Client [$rp] running with [$taskset_client]"
      $cmd_client $CURDIR/../output/bin/${rp}_bencher -addr="$addr" -b=$b -c=$c -n=$n --sleep=$sleep
    done
  done
done
