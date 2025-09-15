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

IFS='|' read -ra keys <<< "$2"
no_title=0
for key in "${keys[@]}"; do
  old="$dir/old-$key.csv"
  new="$dir/new-$key.csv"
  if [ "$no_title" -eq 0 ]; then
    python3 "$PROJECT_ROOT/scripts/reports/diff.py" "$old" "$new" "$key"
    no_title=1
  else
    python3 "$PROJECT_ROOT/scripts/reports/diff.py" "$old" "$new" "$key" "True"
  fi
done
