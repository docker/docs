#!/bin/sh
# SUMMARY: Test dentry cache invalidation
# LABELS: !win
# AUTHOR: David Sheets <david.sheets@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"
IMAGE_NAME=test_volumes_dentry

clean_up() {
    rm -rf "${D4X_LOCAL_TMPDIR}/dentry_test"
    docker rmi "${IMAGE_NAME}" || true
}
trap clean_up EXIT

clean_up
docker build -t "${IMAGE_NAME}" .
docker run --rm -v "${D4X_LOCAL_TMPDIR}:/tmp" -t "$IMAGE_NAME" /dentry_test.sh

exit 0
