#!/bin/sh
# SUMMARY: Test copy over a unix domain socket
# LABELS: !win
# Doesn't work over CIFS
# AUTHOR: David Sheets <david.sheets@docker.com>

set -e
IMAGE_NAME=volumes_socket_copy
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    docker rmi "$IMAGE_NAME" || true
    rm -rf "${D4X_LOCAL_TMPDIR}/sock"
}
trap clean_up EXIT

rm -rf "${D4X_LOCAL_TMPDIR}/sock"
docker build -t "$IMAGE_NAME" .
echo "twelve bytes" > "${D4X_LOCAL_TMPDIR}/socket_copy_input"
docker run --rm -v "${D4X_LOCAL_TMPDIR}:/tmp" -t "$IMAGE_NAME" /socket_copy.sh

diff "${D4X_LOCAL_TMPDIR}/socket_copy_input" "${D4X_LOCAL_TMPDIR}/socket_copy_output"
[ $? -ne 0 ] && exit 1

exit 0
