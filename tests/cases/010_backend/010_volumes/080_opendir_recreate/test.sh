#!/bin/sh
# SUMMARY: Tests correct dcache unlinking behavior
# LABELS: !win
# AUTHOR: David Sheets <david.sheets@docker.com>

set -e # Exit on error
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

CONTAINER_NAME=opendir_holder

clean_up() {
    rm -rf dir
    docker kill $CONTAINER_NAME || true
    docker rm $CONTAINER_NAME || true
}
trap clean_up EXIT

clean_up

mkdir -p dir/dir

docker run --rm -v "$PWD/dir:/dir" --name "$CONTAINER_NAME" alpine \
    sh -c "ls -l /dir && cd /dir/dir && touch /dir/start && sleep 10000" &

# wait for container to signal readiness
while ! [ -f dir/start ]; do sleep 0.1; done

rmdir dir/dir

mkdir dir/dir

touch dir/dir/file

FILE_COUNT=$(docker run --rm -v "$PWD/dir:/dir" alpine ls dir/dir | wc -l)

if [ "$FILE_COUNT" -eq "1" ]; then exit 0; fi

exit 1
