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

types=("thrift" "grpc" "pb")

function log_prefix() {
    echo -n "[`date '+%Y-%m-%d %H:%M:%S'`] "
}

function prepare_old() {
    go get -v github.com/cloudwego/kitex@$old
    go mod tidy
}

function prepare_new() {
    go get -v github.com/cloudwego/kitex@$new
    go mod tidy
}

function compare() {
    type=$1
    report_dir=`date '+%m%d-%H%M'`-$type
    mkdir -p $PROJECT_ROOT/output/$report_dir
    log_prefix; echo "Begin comparing $type..."

    # old
    log_prefix; echo "Benchmark $type @ $old (old)"
    export REPORT_PREFIX=$report_dir/old-
    prepare_old
    time bash $PROJECT_ROOT/scripts/benchmark_$type.sh

    # new
    log_prefix; echo "Benchmark $type @ $new (new)"
    export REPORT_PREFIX=$report_dir/new-
    prepare_new
    time bash $PROJECT_ROOT/scripts/benchmark_$type.sh

    # compare results
    $PROJECT_ROOT/scripts/compare_report.sh $PROJECT_ROOT/output/$report_dir

    log_prefix; echo "End comparing $type..."
}

for tp in ${types[@]}; do
    compare $tp
done


log_prefix; echo "All benchmark finished"
