#!/bin/bash
set -e

# benchmark params
srepo=("grpc" "kitex")
crepo=("grpc" "kitex")
ports=(8000 8001)

CURDIR=$(cd $(dirname $0); pwd)
echo "Checking whether the environment meets the requirements ..."
source $CURDIR/base.sh
echo "Check finished."

echo "Building streaming services by exec build_streaming.sh..."
source $CURDIR/build_streaming.sh
echo "Build finished."
# benchmark
for b in ${body[@]}; do
  for c in ${concurrent[@]}; do
    for q in ${qps[@]}; do
      for ((i = 0; i < ${#srepo[@]}; i++)); do
        srp=${srepo[i]}
        crp=${crepo[i]}
        addr="127.0.0.1:${ports[i]}"
        kill_pid_listening_on_port ${ports[i]}
        # server start
        echo "Starting server [$srp], if failed please check [output/log/nohup.log] for detail."
        nohup $cmd_server $output_dir/bin/${srp}_reciever >> $output_dir/log/nohup.log 2>&1 &
        sleep 1
        echo "Server [$srp] running with [$cmd_server]"

        # run client
        echo "Client [$crp] running with [$cmd_client]"
        $cmd_client $output_dir/bin/${crp}_bencher -addr="$addr" -b=$b -c=$c -qps=$q -n=$n --sleep=$sleep | $tee_cmd

        # stop server
        kill_pid_listening_on_port ${ports[i]}
      done
    done
  done
done

finish_cmd
