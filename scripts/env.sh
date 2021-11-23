#!/bin/bash
set -e

if ! [ -x "$(command -v taskset)" ]; then
  echo "Error: taskset is not installed." >&2
  exit 1
fi

nprocs=$(getconf _NPROCESSORS_ONLN)
if [ $nprocs -lt 4 ]; then
  echo "Error: your environment should have at least 4 processors."
  exit 1
elif [ $nprocs -gt 20 ]; then
  nprocs=20
fi

# GO
GOEXEC=${GOEXEC:-"go"}
GOROOT=$GOROOT

USER=$(whoami)
n=5000000
body=(1024)
concurrent=(100 200 400 600 800 1000)
sleep=0

scpu=$((nprocs > 16 ? 3 : nprocs / 4 - 1)) # max is 3(4 cpus)
nice_cmd=''
if [ "$USER" == "root" ]; then
    nice_cmd='nice -n -20'
fi
cmd_server="${nice_cmd} taskset -c 0-$scpu"
cmd_client="${nice_cmd} taskset -c $((scpu + 1))-$((nprocs - 1))"
