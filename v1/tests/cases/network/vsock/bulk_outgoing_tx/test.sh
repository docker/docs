#!/bin/sh

set -e
export PS4='> '
set -x

host="$1"
port="$2"
bs="$3"
count="$4"

dd if=/dev/zero bs="$bs" count="$count" | nc-vsock -w "$host" "$port"
