#!/bin/bash

# csv
# usage: to_csv.sh xxx.log
grep TPS "$1" | awk -F '[ :,]+' '{split($6,a,"m");split($8,b,"m");print $2","substr($11,3)","substr($9,4)","$4","a[1]","b[1]}'
