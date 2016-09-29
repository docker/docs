#!/bin/bash

. ../../lib/functions

trap "{ docker ps | grep nettest:iperf | cut -f1 -d' ' | xargs docker kill; }" EXIT

docker run -p 5001:5001 -d nettest:iperf iperf -s

# iperf -s is running asynchronously. Note that iperf will exit 0 if it gets
# a connection refused error which is irritating.
sleep 1

iperf -c ${vmip} -t 30 -i 1 -f m
