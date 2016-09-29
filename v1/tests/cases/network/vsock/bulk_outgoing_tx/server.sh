#!/bin/bash

set -e
export PS4='< '
set -x

dir="$1"
host="$2"
port="$3"
bs="$4"
count="$5"

expected=$(($bs * $count))

sock=$(printf "$dir/*%08x.%08x" "$host" "$port")

rm -f received.dat

socat -u -b$bs UNIX-LISTEN:"$sock",unlink-early,shut-close CREATE:received.dat

size=$(stat -f %z received.dat)

rm -f received.dat

: expected $expected bytes, got $size bytes

if [ $size -eq $expected ] ; then
    exit 0
else
    exit 1
fi
