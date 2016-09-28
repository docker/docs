#!/bin/sh
set -e

trap 'echo killing pid:$PID in 1 second; sleep 1; kill -9 $PID' EXIT

# start server
moshpit --debug server --config-file remote-test.yml &
PID=$!

# first client
moshpit --debug client --server 127.0.0.1:1337 --name client1 &
PID2=$!

# second client
moshpit --debug client --server 127.0.0.1:1337 --name client2

wait $PID2
