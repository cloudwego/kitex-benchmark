#!/bin/bash

nprocs=$(getconf _NPROCESSORS_ONLN)
if [ $nprocs -lt 4 ]; then
  echo "Your environment should have at least 4 processors"
  exit 1
elif [ $nprocs -gt 20 ]; then
  nprocs=20
fi

# GO
GOEXEC=${GOEXEC:-"go"}
GOROOT=$GOROOT

USER=$(whoami)
REPORT=${REPORT:-"$(date +%F-%H-%M)"}
n=5000000
body=(1024)
concurrent=(100 200 400 600 800 1000)
sleep=0

scpu=$((nprocs > 16 ? 3 : nprocs / 4 - 1)) # max is 3(4 cpus)
nice_cmd=''
tee_cmd="tee -a output/${REPORT}.log"
function finish_cmd() {
  # to csv report
  ./scripts/reports/to_csv.sh output/"$REPORT.log" > output/"$REPORT".csv
  echo "output reports: output/$REPORT.log, output/$REPORT.csv"
  cat ./output/"$REPORT.csv"
}
if [ "$USER" == "root" ]; then
    nice_cmd='nice -n -20'
fi
cmd_server="${nice_cmd} taskset -c 0-$scpu"
cmd_client="${nice_cmd} taskset -c $((scpu + 1))-$((nprocs - 1))"
