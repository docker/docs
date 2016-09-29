#!/bin/sh
# SUMMARY: Performs FS operations during inotify event delivery
# LABELS: !win
# AUTHOR: David Sheets <david.sheets@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

IMAGE_NAME=interfere

CONTAINER_NAME=inotify-interference-test

clean_up() {
    docker rmi $IMAGE_NAME || true
    docker rm $CONTAINER_NAME || true
    rm -f success start newfile
}
trap clean_up EXIT

clean_up

docker build -t $IMAGE_NAME .
docker run --rm --name $CONTAINER_NAME \
       -v "$(pwd):/host" -w /host $IMAGE_NAME ../interfere &

if [ -f start ]
then
    echo "Container finished before interference could begin."
    exit 1
fi

tries=0
while true ; do
    if [[ "$tries" -gt 10 ]]; then
        echo "Polled for container start file more than 10 times. Failure."
        exit 1
    elif [ -f start ]; then
        echo "Container created container start file"
        break
    else
        tries=$((tries+1))
        echo "Did not find start file on attempt $tries"
        sleep 1
    fi
done

for i in $(seq 1 500); do
    echo "$i" > /dev/null # not using $i upsets shellcheck
    touch newfile || true
    chmod +x newfile || true
    sleep 0.01
done
if ! [ -f success ]
then
    echo "Interference test failed. Killing container."
    docker kill $CONTAINER_NAME
fi

[ -f success ]
