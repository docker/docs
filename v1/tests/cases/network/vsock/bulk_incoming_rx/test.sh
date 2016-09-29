#!/bin/sh

set -e
export PS4='> '
set -x

guest="$1"
port="$2"
bs="$3"
count="$4"

expected=$(($bs * $count))

: expected $expected bytes total

nc-vsock -r -l "$port" | dd of=received.dat bs=$bs

size=$(stat -c %s received.dat)

rm -f received.dat

: expected $expected bytes, got $size bytes

if [ $size -eq $expected ] ; then
    return 0
else
    return 1
fi
