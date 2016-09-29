#!/bin/sh
# SUMMARY: Run a never ending docker build, interrupt it, check there's no leak'
# LABELS:
# AUTHOR: David Gageot <david.gageot@docker.com>

set -e # Exit on error
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
	docker rm -f "$(docker ps | awk '/sleep/ {printf $1}')" || true
}
trap clean_up EXIT

clean_up
docker pull alpine:3.3
docker build --no-cache -t sleeper . &
pid=$!

i=0
while [ "$(docker ps | grep -c sleep)" != "1" ]; do
    if [ $i -gt 10 ]; then
        echo "Couldn't find running build"
        exit 1
    fi

    sleep 1s
    i=$((i+1))
done

kill -3 $pid

i=0
while [ "$(docker ps | grep -c sleep)" != "0" ]; do
    if [ $i -gt 10 ]; then
        echo "The build is still running after 10s"
        exit 1
    fi

    sleep 1s
    i=$((i+1))
done

exit 0
