#!/bin/bash

repo=("grpc" "kitex" "kitex-mux" "rpcx" "arpc" "arpc_nbio")

# build
. ./scripts/build_pb.sh

# benchmark
. ./scripts/kill_servers.sh
core=0
for ((i = 0; i < ${#repo[@]}; i++)); do
  rp=${repo[i]}

  # server start
  nohup taskset -c $core-$(($core + 3)) ./output/bin/${rp}_reciever >> output/log/nohup.log 2>&1 &
  echo "server $rp running at cpu $core-$(($core + 3))"
  core=$(($core + 4))
  sleep 1
done
