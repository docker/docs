#!/bin/sh
# SUMMARY: Checks that child creation in open directories mostly works
# LABELS:
# AUTHOR: David Sheets <dsheets@docker.com>

set -e
IMAGE_NAME=create_child
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

DIR="${D4X_LOCAL_TMPDIR}/create_child"

clean_up() {
    docker rmi "$IMAGE_NAME" || true
    rm -rf "${DIR}"
}
trap clean_up EXIT

rm -rf "${DIR}"
mkdir "${DIR}"

docker build -t "$IMAGE_NAME" .
docker run --rm -v "${DIR}":/tmp/create_child -t "$IMAGE_NAME" /create_child

exit 0
