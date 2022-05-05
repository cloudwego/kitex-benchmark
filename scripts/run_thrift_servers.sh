#!/bin/bash
set -e
CURDIR=$(cd $(dirname $0); pwd)
repo=("kitex" "kitex-mux")

# build
source $CURDIR/env.sh
source $CURDIR/build_thrift.sh

# benchmark
source $CURDIR/kill_servers.sh
core=0
for ((i = 0; i < ${#repo[@]}; i++)); do
  rp=${repo[i]}

  # server start
  nohup taskset -c $core-$(($core + 3)) $output_dir/bin/${rp}_reciever >> $output_dir/log/nohup.log 2>&1 &
  echo "Server [$rp] running at cpu $core-$(($core + 3)) ..."
  core=$(($core + 4))
  sleep 1
done
