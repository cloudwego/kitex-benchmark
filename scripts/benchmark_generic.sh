#!/bin/bash
set -e
CURDIR=$(cd $(dirname $0); pwd)

echo "Checking whether the environment meets the requirements ..."
source $CURDIR/env.sh
echo "Check finished."

srepo=("generic_http" "generic_http" "generic_json_default" "generic_json_fallback" "generic_map" "generic_ordinary")
crepo=("generic_http_default" "generic_http_fallback" "generic_json_default" "generic_json_fallback" "generic_map" "generic_ordinary")
ports=(8001 8001 8002 8003 8004 8005)
data=("1KB" "5KB" "10KB")

function finish_cmd_generic() {
  # to csv report
    ./scripts/reports/to_csv.sh output/"$1.log" > output/"$1".csv
    echo "output reports: output/$1.log, output/$1.csv"
    cat ./output/"$1.csv"
}

echo "Building generic services by exec build_generic.sh ..."
source $CURDIR/build_generic.sh
echo "Build finished."

# benchmark
for d in ${data[@]}; do
  report="$(date +%F-%H-%M)_${d}"
  tee_cmd_generic="tee -a output/$report.log"
  echo "------------------------------"
  echo "Data size: $d"
  for b in ${body[@]}; do
    for c in ${concurrent[@]}; do
      for ((i = 0; i < ${#srepo[@]}; i++)); do
        srp=${srepo[i]}
        crp=${crepo[i]}
        addr="127.0.0.1:${ports[i]}"
        # server start
        echo "Starting server [$srp], if failed please check [output/log/nohup.log] for detail"
        nohup $cmd_server $output_dir/bin/${srp}_${d}_reciever >> $output_dir/log/nohup.log 2>&1 &
        sleep 1
        echo "Server [$srp] running with [$cmd_server]"

        # run client
        echo "Client [$crp] running with [$cmd_client]"
        $cmd_client $output_dir/bin/${crp}_${d}_bencher -addr="$addr" -b=$b -c=$c -n=$n --sleep=$sleep | $tee_cmd_generic

        # stop server
        pid=$(ps -ef | grep ${srp}_${d}_reciever | grep -v grep | awk '{print $2}')
        disown $pid
        kill -9 $pid
        echo "Server [$srp] stopped, pid [$pid]."
        sleep 1
      done
    done
  done
  finish_cmd_generic $report
done
