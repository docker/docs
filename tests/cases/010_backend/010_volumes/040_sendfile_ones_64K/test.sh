#!/bin/sh
# SUMMARY: Test sendfile on a mounted volume (64K)
# LABELS:
# AUTHOR: David Sheets <david.sheets@docker.com>

set -e
IMAGE_NAME=volumes_sendfile_ones_sixtyfourk
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    docker rmi "${IMAGE_NAME}"
}
trap clean_up EXIT

docker build -t "$IMAGE_NAME" .
docker run --rm -v "${D4X_LOCAL_TMPDIR}:/tmp" -t "$IMAGE_NAME" /sendfile_ones_64K.py

exit 0
