#!/bin/bash

# table
# usage: to_table.sh xxx.log
grep TPS "$1" | awk -F '[ :,]+' '{print "<tr><td> "$2" </td><td> 传输 </td><td> "$4" </td><td> "$6" </td><td> "$8" </td></tr>"}'
