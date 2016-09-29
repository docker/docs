#!/bin/bash

set -e
export PS4='< '
set -x

dir="$1"
host="$2"
port="$3"
bs="$4"
count="$5"

sock=$(printf "$dir/*%08x.%08x" "$host" "$port")

socat -b$bs EXEC:"dd if=/dev/zero bs=$bs count=$count" UNIX-LISTEN:"$sock",unlink-early
