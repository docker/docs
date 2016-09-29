#!/bin/sh
# SUMMARY: Checks that trans-directory rename unlinks mostly work
# LABELS:
# AUTHOR: David Sheets <dsheets@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"
IMAGE_NAME=trans_rename_unlink
DIR="${D4X_LOCAL_TMPDIR}/trans_rename_unlink"

clean_up() {
       docker rmi "$IMAGE_NAME" || true
       rm -rf "${DIR}"
}
trap clean_up EXIT

rm -rf "${DIR}"
mkdir "${DIR}"
docker build -t "$IMAGE_NAME" .
docker run --rm -v "${DIR}":/tmp/trans_rename_unlink -t "$IMAGE_NAME" /trans_rename_unlink
exit 0
