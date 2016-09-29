#!/bin/bash

set -e
export PS4='< '
set -x

dir="$1"
guest="$2"
port="$3"
bs="$4"
count="$5"

sock=$dir/@connect

( printf "*%08x.%08x\n" "$guest" "$port" ; dd if=/dev/zero bs=$bs count=$count ) | socat -b$bs STDIO UNIX-CONNECT:"$sock",retry=10
