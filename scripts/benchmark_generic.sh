#!/bin/bash
set -e
CURDIR=$(cd $(dirname $0); pwd)

echo "Checking whether the environment meets the requirements ..."
source $CURDIR/base.sh
echo "Check finished."

srepo=("generic_http" "generic_json" "generic_map" "generic_ordinary")
crepo=("generic_http" "generic_json" "generic_map" "generic_ordinary")
ports=(8001 8002 8003 8004)

echo "Building generic services by exec build_generic.sh ..."
source $CURDIR/build_generic.sh
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
        echo "Starting server [$srp], if failed please check [output/log/nohup.log] for detail"
        nohup $cmd_server $output_dir/bin/${srp}_reciever >> $output_dir/log/nohup.log 2>&1 &
        sleep 1
        echo "Server [$srp] running with [$cmd_server]"

        # run client
        echo "Client [$crp] running with [$cmd_client]"
        $cmd_client $output_dir/bin/${crp}_bencher -addr="$addr" -b=$b -c=$c -qps=$q -t=$t --sleep=$sleep | $tee_cmd

        # stop server
        kill_pid_listening_on_port ${ports[i]}
      done
    done
  done
done

finish_cmd
