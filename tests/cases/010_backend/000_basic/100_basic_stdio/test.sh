#!/bin/sh
# SUMMARY: Checks that stdio to `docker run` works
# LABELS:
# AUTHOR: Ian Campbell <ian.campbell@docker.com>

set -e # exit on error
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

sizes="1024 4096 65536 131072"

for size in $sizes ; do
    echo "stdout bs=$size"
    docker run --rm alpine timeout -t 10 dd if=/dev/zero bs="$size" count=1 >/dev/null
    [ $? -eq 0 ]  || exit 1
    echo
done

for size in $sizes ; do
    echo "stdin bs=$size"
    dd if=/dev/zero bs="$size" count=1 | docker run -i --rm alpine timeout -t 10 sh -c "dd bs=$size of=/null"
    [ $? -eq 0 ]  || exit 1
    echo
done

exit 0
