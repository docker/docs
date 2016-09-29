#!/bin/sh
# SUMMARY: Test that variable sizes contexts are transferred correctly
# LABELS:
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

IMAGE_NAME=context_fixed

SIZES="1 4 4095 4096 4097 1048575 1048576 1048577"

clean_up() {
    rm -rf rt-filemd5 testfile
    d4x_cleanup_image "${IMAGE_NAME}"
}
trap clean_up EXIT

cp "${RT_UTILS}/rt-filemd5" .

for size in ${SIZES}; do
    echo "Trying file size ${size}"
    "${RT_UTILS}/rt-filegen" "${size}" testfile
    hostmd5=$("${RT_UTILS}/rt-filemd5" testfile)
    docker build -t "${IMAGE_NAME}" . || exit 1
    contmd5=$(docker run --rm "${IMAGE_NAME}" /rt-filemd5 testfile) || exit 1
    # strip carriage return
    contmd5="$(echo "${contmd5}" | tr -d '\r')"
    [ "${hostmd5}" != "${contmd5}" ] && exit 1
done
exit 0
