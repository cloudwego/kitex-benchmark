#!/bin/bash
#go mod edit --replace github.com/cloudwego/kitex=github.com/sinnera/kitex@perf/grpc_streaming
#go mod tidy
#go mod edit --replace github.com/cloudwego/netpoll-http2=github.com/sinnera/netpoll-http2@perf/grpc_streaming
#go mod tidy

#go mod edit --replace github.com/cloudwego/kitex=github.com/sinnera/kitex@perf/remove_nhttp2
#go mod tidy


go mod edit --replace github.com/cloudwego/kitex=github.com/sinnera/kitex@perf/merge_all_prs
go mod tidy

set -e
CURDIR=$(cd $(dirname $0); pwd)

echo "Checking whether the environment meets the requirements ..."
source $CURDIR/env.sh
echo "Check finished."

repo=("kitex")
cli_repo=("kitex")
ports=(8001)

echo "Building streaming services by exec build_streaming.sh..."
source $CURDIR/build_streaming.sh
echo "Build finished."
# benchmark
for b in ${body[@]}; do
  for c in ${concurrent[@]}; do
    for ((i = 0; i < ${#repo[@]}; i++)); do
      rp=${repo[i]}
      crp=${cli_repo[i]}
      addr="127.0.0.1:${ports[i]}"
      # server start
      echo "Starting server [$rp], if failed please check [output/log/nohup.log] for detail."
      nohup $cmd_server $output_dir/bin/${rp}_reciever >> $output_dir/log/nohup.log 2>&1 &
      sleep 1
      echo "Server [$rp] running with [$cmd_server]"

      # run client
      echo "Client [$crp] running with [$cmd_client]"
      $cmd_client $output_dir/bin/${crp}_bencher -addr="$addr" -b=$b -c=$c -n=$n --sleep=$sleep | $tee_cmd

      # stop server
      pid=$(ps -ef | grep ${rp}_reciever | grep -v grep | awk '{print $2}')
      disown $pid
      kill -9 $pid
      echo "Server [$rp] stopped, pid [$pid]."
      sleep 1
    done
  done
done

finish_cmd
