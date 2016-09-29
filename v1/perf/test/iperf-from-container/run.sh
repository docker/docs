#!/bin/bash

. ../../lib/functions

if [ "${network}" = "slirp" ]; then
  echo Skipping test because we are using user-space networking
  exit 0
fi

trap "{ docker ps | grep nettest:iperf | cut -f1 -d' ' | xargs docker kill; pkill iperf; }" EXIT

iperf -s > /dev/null &
echo "iperf running with pid $!"

docker run -p 5001:5001 nettest:iperf iperf -c $hostip -t 30 -i 1 -f m
