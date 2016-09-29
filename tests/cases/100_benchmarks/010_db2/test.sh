#!/bin/sh
# SUMMARY: Create a DB2 database
# LABELS: skip, benchmarks
# AUTHOR: David Scott <dave.scott@docker.com>
# see https://github.com/docker/for-mac/issues/668

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

IMAGE_NAME=db2_create

clean_up() {
    docker stop "${IMAGE_NAME}" || true
    docker rm "${IMAGE_NAME}" || true
    d4x_cleanup_image "${IMAGE_NAME}"
}
trap clean_up EXIT

# in case of debugging leftovers
clean_up

# with disk/full-sync-on-flush=1 this takes 300s
# with disk/full-sync-on-flush=0 this takes 30s
time docker build -t "${IMAGE_NAME}" .

