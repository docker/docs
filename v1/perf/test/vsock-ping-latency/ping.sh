#!/bin/bash

set -e

if [ $# -ne 1 ] ; then
    echo >&2 "usage: [SRV-CONTAINER]"
    exit 1
fi

server=$1

port=0x4563686f
# Pre and post PR#3106
#connect=/var/tmp/com.docker.vsock/connect
connect=$HOME/Library/Containers/com.docker.docker/Data/@connect

container=$(docker run --net=host -d vsock-ping ${port})
echo >&2 "Server $server running as $container"

./client ${port} ${connect}

docker stop $container 1>/dev/null
docker rm $container 1>/dev/null
