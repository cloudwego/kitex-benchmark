#!/bin/bash
set -e
CURDIR=$(cd $(dirname $0); pwd)
repo=("grpc" "kitex" "kitex-mux" "rpcx" "arpc" "hessian" "triple")
ports=(8000 8001 8002 8003 8004 8005 8006)
ip=${IP:-"127.0.0.1"}

# build
source $CURDIR/env.sh
source $CURDIR/build_pb.sh

# benchmark
for b in ${body[@]}; do
  for c in ${concurrent[@]}; do
    for ((i = 0; i < ${#repo[@]}; i++)); do
      rp=${repo[i]}
      addr="${ip}:${ports[i]}"

      # run client
      echo "Client [$rp] running with [$taskset_client]"
      $cmd_client $output_dir/bin/${rp}_bencher -addr="$addr" -b=$b -c=$c -n=$n --sleep=$sleep | $tee_cmd
    done
  done
done

finish_cmd
