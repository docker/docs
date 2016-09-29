#!/bin/sh
# SUMMARY: Test chown persisting
# LABELS: osx
# AUTHOR: David Sheets <david.sheets@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"
IMAGE_NAME=test_volumes_chown

clean_up() {
    docker rmi "${IMAGE_NAME}" || true
}
trap clean_up EXIT

docker build -t "${IMAGE_NAME}" .
docker run --rm -v "${D4X_LOCAL_TMPDIR}:/tmp" -t "${IMAGE_NAME}" /chown_test.sh
exit 0
