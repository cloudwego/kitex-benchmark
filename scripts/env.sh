#!/bin/bash

nprocs=$(getconf _NPROCESSORS_ONLN)
if [ $nprocs -lt 4 ]; then
  echo "Your environment should have at least 4 processors"
  exit 1
elif [ $nprocs -gt 20 ]; then
  nprocs=20
fi

n=5000000
body=(1024)
concurrent=(100 200 400 600 800 1000)

scpu=$((nprocs > 16 ? 3 : nprocs / 4 - 1)) # max is 3(4 cpus)
taskset_server="taskset -c 0-$scpu"
taskset_client="taskset -c $((scpu + 1))-$((nprocs - 1))"
