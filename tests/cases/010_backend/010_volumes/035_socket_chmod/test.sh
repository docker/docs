#!/bin/sh
# SUMMARY: Test chmod on a unix doamin on a mounted volume
# LABELS: !win
# Doesn't work over CIFS
# AUTHOR: David Sheets <david.sheets@docker.com>

set -e
IMAGE_NAME=volumes_socket_chmod
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    docker rmi "$IMAGE_NAME"
    rm -rf "${D4X_LOCAL_TMPDIR}/sock"
}
trap clean_up EXIT

rm -rf "${D4X_LOCAL_TMPDIR}/sock"
docker build -t "$IMAGE_NAME" .
docker run --rm -v "${D4X_LOCAL_TMPDIR}:/tmp" -t "$IMAGE_NAME" /socket_chmod.sh
exit 0
