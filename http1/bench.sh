#!/bin/bash
set -e
CURDIR=$(cd $(dirname $0); pwd)

cd $CURDIR && rm -rf bin/ log && mkdir bin log

go build -o bin/receiver ./

go build -o bin/bencher ./client

GOGC=1000 nohup taskset -c 0-3 ./bin/receiver > log/nohup.log 2>&1 &

sleep 2

GOGC=1000 taskset -c 4-15 ./bin/bencher -b=1024 -c=30 -qps=0 -n=1000000