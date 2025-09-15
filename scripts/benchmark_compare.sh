#!/bin/bash

# Usage:
#   `bash script/benchmark_compare.sh $old $new`
# Where old/new could be branch name, tag or commit hash. For example:
#   `bash script/benchmark_compare.sh v0.7.1 develop`
# It will invoke the script/benchmark_$type.sh to run on the two versions of Kitex.
# The $type variable could be "thrift", "grpc" or "pb", as is listed in the `$types` variable below.
#
# When it finishes, the reports are saved under `output/mmdd-HHMM-$type/`. Take `output/0927-1401-thrift`
# as example, you can now run:
#   `bash script/compare_result.sh output/0927-1401-thrift`
# to get the difference of two reports.

cd `dirname $0`/..
PROJECT_ROOT=`pwd`

source $PROJECT_ROOT/scripts/util.sh && check_supported_env

old=$1
new=$2

if [ -z "$old" -o -z "$new" ]; then
    echo "Usage: $0 <old> <new>"
    echo "  old/new could be branch name, tag or commit hash."
    echo "  e.g. $0 v1.13.1 develop"
    exit 1
fi

function log_prefix() {
    echo -n "[`date '+%Y-%m-%d %H:%M:%S'`] "
}

function prepare_old() {
    git checkout go.mod go.sum
    go get -v github.com/cloudwego/kitex@$old
    go mod tidy
}

function prepare_new() {
    git checkout go.mod go.sum
    go get -v github.com/cloudwego/kitex@$new
    go mod tidy
}

function benchmark() {
  ( # use subshell to remove dirty env/variables
    build=$1
    srp=$2
    crp=$2
    port=$3
    # base
    source $PROJECT_ROOT/scripts/base.sh
    # build
    source $build
    # benchmark
    for b in ${body[@]}; do
      for c in ${concurrent[@]}; do
        for q in ${qps[@]}; do
          addr="127.0.0.1:${port}"
          kill_pid_listening_on_port ${port}
          # server start
          echo "Starting server [$srp], if failed please check [output/log/nohup.log] for detail."
          nohup $cmd_server $output_dir/bin/${srp}_reciever >> $output_dir/log/nohup.log 2>&1 &
          sleep 1
          echo "Server [$srp] running with [$cmd_server]"

          # run client
          echo "Client [$crp] running with [$cmd_client]"
          $cmd_client $output_dir/bin/${crp}_bencher -addr="$addr" -b=$b -c=$c -qps=$q -n=$n | $tee_cmd

          # stop server
          kill_pid_listening_on_port ${port}
        done
      done
    done

    finish_cmd
  )
}

function compare() {
    type=$1
    build=$2
    repo=$3
    port=$4
    report_dir=`date '+%m%d-%H%M'`-$type
    mkdir -p $PROJECT_ROOT/output/$report_dir
    log_prefix; echo "Begin comparing $type..."

    # old
    log_prefix; echo "Benchmark $type @ $old (old)"
    export REPORT_PREFIX=$report_dir/old-
    prepare_old
    time benchmark $build $repo $port

    # new
    log_prefix; echo "Benchmark $type @ $new (new)"
    export REPORT_PREFIX=$report_dir/new-
    prepare_new
    time benchmark $build $repo $port

    # compare results
    $PROJECT_ROOT/scripts/compare_report.sh $PROJECT_ROOT/output/$report_dir

    log_prefix; echo "End comparing $type..."
}

compare "thrift" "$PROJECT_ROOT/scripts/build_thrift.sh" "kitex" 8001

compare "thrift-mux" "$PROJECT_ROOT/scripts/build_thrift.sh" "kitex-mux" 8002

compare "protobuf" "$PROJECT_ROOT/scripts/build_pb.sh" "kitex" 8001

compare "grpc-unary" "$PROJECT_ROOT/scripts/build_grpc.sh" "kitex" 8006

compare "grpc-bidistream" "$PROJECT_ROOT/scripts/build_streaming.sh" "kitex_grpc" 8001

compare "ttstream-bidistream" "$PROJECT_ROOT/scripts/build_streaming.sh" "kitex_tts_lconn" 8002

compare "generic-json" "$PROJECT_ROOT/scripts/build_generic.sh" "generic_json" 8002

compare "generic-map" "$PROJECT_ROOT/scripts/build_generic.sh" "generic_map" 8003

# compare "generic-binary" "$PROJECT_ROOT/scripts/build_generic.sh" "generic_binary" 8004

log_prefix; echo "All benchmark finished"
