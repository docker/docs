#!/bin/sh
# SUMMARY: Remove a sub-directory on a mounted volume
# LABELS:
# AUTHOR: Dave Scott <dave.scott@docker.com>

set -e
IMAGE_NAME=volumes_rmdir_nonempty
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    docker rmi "$IMAGE_NAME" || true
    rm -rf "${D4X_LOCAL_TMPDIR}/foo" "${D4X_LOCAL_TMPDIR}/bar"
}
trap clean_up EXIT

rm -rf "${D4X_LOCAL_TMPDIR}/foo" "${D4X_LOCAL_TMPDIR}/bar"

docker build -t "$IMAGE_NAME" .
docker run --rm -v "${D4X_LOCAL_TMPDIR}:/tmp" -t "$IMAGE_NAME" /rmdir_nonempty.sh

exit 0
