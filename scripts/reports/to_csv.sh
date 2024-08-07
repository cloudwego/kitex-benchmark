#!/bin/bash

# csv
# usage: to_csv.sh xxx.log
output_dir=$(dirname "$1")
# Kind,Concurrency,Data_Size,TPS,AVG,P99,Server_CPU,Client_CPU
grep TPS "$1" | awk '{print $2" "$13" "$11" "$4" "$6" "$8}' | awk '{gsub(/[:c=,(b=ms%]/, "")} 1' > $output_dir/tps.out
grep '@Server' "$1" | grep CPU | awk '{print $14}' | awk '{gsub(/[%:]||AVG/, "")} 1' > $output_dir/server.out
grep '@Client' "$1" | grep CPU | awk '{print $14}' | awk '{gsub(/[%:]||AVG/, "")} 1'  > $output_dir/client.out
# combine each line, replace space by comma
awk '{ lines[FNR] = lines[FNR] $0 " " } END { for (i=1; i<=FNR; i++) print lines[i] }' $output_dir/tps.out $output_dir/server.out $output_dir/client.out | awk '{ print substr($0, 1, length($0)-1) }' | awk '{gsub(" ", ",")} 1'
