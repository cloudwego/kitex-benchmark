#!/bin/bash
#go mod edit --replace github.com/cloudwego/kitex=github.com/sinnera/kitex@perf/grpc_streaming
#go mod tidy
#go mod edit --replace github.com/cloudwego/netpoll-http2=github.com/sinnera/netpoll-http2@perf/grpc_streaming
#go mod tidy

#go mod edit --replace github.com/cloudwego/kitex=github.com/sinnera/kitex@perf/remove_nhttp2
#go mod edit --replace github.com/cloudwego/kitex=github.com/sinnera/kitex@hotfix/validate_grpc_option #最早的优化commit
#go mod tidy

go mod edit --replace github.com/cloudwego/kitex=github.com/cloudwego/kitex@bb1260b48c482c2a1012ddc88589ff7f233deef4
#go mod edit --replace github.com/cloudwego/kitex=github.com/cloudwego/kitex@6c033d991c18de59e106ca5d94e0543f6de94008
#go mod edit --replace github.com/cloudwego/kitex=github.com/cloudwego/kitex@4b5390de7460a5e0d7a670e2d79e7b196bc93478
#go mod edit --replace github.com/cloudwego/kitex=github.com/cloudwego/kitex@706250705e679be1afcbb49f6c369fdd5106432f
#go mod edit --replace github.com/cloudwego/kitex=github.com/sinnera/kitex@perf/merge_all_prs
#go mod edit --replace github.com/cloudwego/kitex=github.com/cloudwego/kitex@main
#go mod edit --replace github.com/cloudwego/kitex=github.com/cloudwego/kitex@v0.2.1

#go mod edit --replace github.com/cloudwego/kitex=github.com/sinnera/kitex@fix/fix_framer_rw
#go mod edit --replace github.com/cloudwego/kitex=github.com/sinnera/kitex@fix/fix_framer_rw_without_nhttp2
#go mod edit -replace github.com/apache/thrift=github.com/apache/thrift@v0.13.0
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
