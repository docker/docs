#!/bin/sh
# SUMMARY: Test various open flags on a mounted volume
# LABELS:
# AUTHOR: David Sheets <david.sheets@docker.com>

set -e
IMAGE_NAME=volumes_open_flags
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    docker rmi "$IMAGE_NAME" || true
}
trap clean_up EXIT

docker build -t "$IMAGE_NAME" .
docker run --rm -v "${D4X_LOCAL_TMPDIR}:/tmp" -t "$IMAGE_NAME" /open_flags.py
exit 0
