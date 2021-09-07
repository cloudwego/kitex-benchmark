#!/bin/bash
set -e
CURDIR=$(cd $(dirname $0); pwd)

repo=("grpc" "kitex-pb" "kitex-grpc" "kitex-mux" "rpcx" "arpc" "arpc-nbio")

# build
source $CURDIR/build_pb.sh

# benchmark
source $CURDIR/kill_servers.sh
core=0
for ((i = 0; i < ${#repo[@]}; i++)); do
  rp=${repo[i]}

  # server start
  nohup taskset -c $core-$(($core + 3)) $CURDIR/../output/bin/${rp}_reciever >> $CURDIR/../output/log/nohup.log 2>&1 &
  echo "Server [$rp] running at cpu $core-$(($core + 3)) ..."
  core=$(($core + 4))
  sleep 1
done
