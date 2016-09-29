#!/bin/sh
# SUMMARY: Create a Fifo on a shared volume and verify it works
# LABELS: !win
# Doesn't work over CIFS
# AUTHOR: David Sheets <david.sheets@docker.com>

set -e # Don't exit on error
IMAGE_NAME=volumes_fifo_copy
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    docker rmi $IMAGE_NAME
    rm -f "${D4X_LOCAL_TMPDIR}/fifo_copy_input"
    rm -f "${D4X_LOCAL_TMPDIR}/fifo_copy_output"
}
trap clean_up EXIT


rm -rf "${D4X_LOCAL_TMPDIR}/fifo"
docker build -t "$IMAGE_NAME" .
echo "twelve bytes" > "${D4X_LOCAL_TMPDIR}/fifo_copy_input"
docker run --rm -v "${D4X_LOCAL_TMPDIR}:/tmp" -t "$IMAGE_NAME" /fifo_copy.sh

diff "${D4X_LOCAL_TMPDIR}/fifo_copy_input" "${D4X_LOCAL_TMPDIR}/fifo_copy_output"
[ $? -ne 0 ] && exit 1

exit 0
