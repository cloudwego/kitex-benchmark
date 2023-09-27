#!/bin/bash

# This tool is designed to be used with `scripts/benchmark_compare.sh` which generates reports with
# names like `old-xxx.csv` and `new-xxx.csv` in the `output/mmdd-HHMM-$type` directory.
# Example:
#   `bash script/compare_result.sh output/0927-1401-thrift`

cwd=`pwd`
cd `dirname $0`/..
PROJECT_ROOT=`pwd`
cd "$cwd"

dir=${1}
if [ -z "$dir" ];then
    echo "Usage: $0 <directory>"
    echo "  directory can be absolute path or relative path"
    exit 1
fi

for new in $dir/new*.csv; do
    for old in $dir/old*.csv; do
        echo python3 "scripts/reports/diff.py" "$old" "$new"
        python3 "$PROJECT_ROOT/scripts/reports/diff.py" "$old" "$new"
    done
done
