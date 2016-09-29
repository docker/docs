#!/bin/bash

set -e
export PS4='< '
set -x

dir="$1"
guest="$2"
port="$3"
bs="$4"
count="$5"

expected=$(($bs * $count))

sock="$dir/@connect"

rm -f received.dat

printf "*%08x.%08x\n" "$guest" "$port" | socat -b$bs UNIX-CONNECT:"$sock",retry=10 STDIO | dd bs=$bs of=received.dat

size=$(stat -f %z received.dat)

rm -f received.dat

: expected $expected bytes, got $size bytes

if [ $size -eq $expected ] ; then
    exit 0
else
    exit 1
fi
