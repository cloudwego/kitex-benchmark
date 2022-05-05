#!/bin/bash
set -e

pid=$(ps -ef | grep reciever | grep -v grep | awk '{print $2}')
if [ -n "${pid}" ]; then
  kill -9 $pid
fi
