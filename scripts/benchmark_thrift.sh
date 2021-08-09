#!/bin/bash

. ./scripts/env.sh
repo=("kitex" "kitex-mux")
ports=(8001 8002)

. ./scripts/build_thrift.sh

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
      echo "benchmark $rp: echo size=$b, concurrent=$c, n=$n at:$(date "+%Y-%m-%d %H:%M:%S")"
      $taskset_client ./output/bin/${rp}_bencher -addr="$addr" -b=$b -c=$c -n=$n

      # stop server
      pid=$(ps -ef | grep ${rp}_reciever | grep -v grep | awk '{print $2}')
      disown $pid
      kill -9 $pid
      sleep 1
    done
  done
done
