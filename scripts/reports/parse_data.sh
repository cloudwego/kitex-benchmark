#!/bin/bash

# table
grep TPS output.log | awk -F '[ :,]+' '{print "<tr><td> "$2" </td><td> 传输 </td><td> "$4" </td><td> "$6" </td><td> "$8" </td></tr>"}'

# csv
grep TPS output.log | awk -F '[ :,]+' '{split($6,a,"m");split($8,b,"m");print $2","substr($11,3)","substr($9,4)","$4","a[1]","b[1]}'
