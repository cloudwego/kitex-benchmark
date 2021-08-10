#!/bin/bash

. ./scripts/env.sh

repo=("grpc" "kitex" "kitex-mux" "rpcx" "arpc" "arpc_nbio")
ports=(8000 8001 8002 8003 8004 8005)

. ./scripts/build_pb.sh

# benchmark
for b in ${body[@]}; do
  for c in ${concurrent[@]}; do
    for ((i = 0; i < ${#repo[@]}; i++)); do
      rp=${repo[i]}
      addr="127.0.0.1:${ports[i]}"
      # server start
      nohup $taskset_server ./output/bin/${rp}_reciever >> output/log/nohup.log 2>&1 &
      sleep 1
      echo "server $rp running with $taskset_server"

      # run client
      $taskset_client ./output/bin/${rp}_bencher -addr="$addr" -b=$b -c=$c -n=$n

      # stop server
      pid=$(ps -ef | grep ${rp}_reciever | grep -v grep | awk '{print $2}')
      disown $pid
      kill -9 $pid
      sleep 1
    done
  done
done
