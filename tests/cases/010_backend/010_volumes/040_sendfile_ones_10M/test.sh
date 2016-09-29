#!/bin/sh
# SUMMARY: Test sendfile on a mounted volume (10M)
# LABELS:
# AUTHOR: David Sheets <david.sheets@docker.com>

set -e # Don't exit on error
IMAGE_NAME=volumes_sendfile_ones_tenmeg
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    docker rmi "$IMAGE_NAME" || true
}
trap clean_up EXIT

docker build -t "$IMAGE_NAME" .
docker run --rm -v "${D4X_LOCAL_TMPDIR}:/tmp" -t "$IMAGE_NAME" /sendfile_ones_10M.py
exit 0
