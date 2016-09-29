#!/bin/sh
# SUMMARY: Verifies various inotify operations
# LABELS: !win
# AUTHOR: David Sheets <david.sheets@docker.com>

set -ex
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

IMAGE_NAME=test_volumes_inotify

clean_up() {
    d4x_cleanup_image ${IMAGE_NAME}
    rm -rf /tmp/inotify*
}
clean_up
trap clean_up EXIT

mkdir /tmp/inotify
docker build -t $IMAGE_NAME .
./run_inotify.sh $IMAGE_NAME

