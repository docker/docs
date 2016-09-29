#!/bin/sh
# SUMMARY: Test that random sized contexts are transferred correctly
# LABELS:
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"
IMAGE_NAME=context_random

# default runs: 5 iterations, max file size 1M
ITER=5
MAX_SIZE=0x100000



if rt_label_set "master"; then
    # 10 iterations with up to 50MB files
    ITER=10
    MAX_SIZE=0x3200000
elif rt_label_set "nightly"; then
    # 50 iterations with up to 100MB files
    ITER=50
    MAX_SIZE=0x6400000
elif rt_label_set "release"; then
    # 50 iterations with up to 100MB files
    ITER=50
    MAX_SIZE=0x6400000
fi

clean_up() {
    rm -f rt-filemd5 testfile
    d4x_cleanup_image ${IMAGE_NAME}
}
trap clean_up EXIT

cp "${RT_UTILS}/rt-filemd5" .

tries=0
while true; do
    if [ "$tries" -gt ${ITER} ]; then
        break
    fi

    "${RT_UTILS}/rt-filerandgen" "${MAX_SIZE}" testfile
    ls -l testfile
    hostmd5=$("${RT_UTILS}/rt-filemd5" testfile)
    docker build -t "${IMAGE_NAME}" . || exit 1
    contmd5=$(docker run --rm -t "${IMAGE_NAME}" /rt-filemd5 testfile) || exit 1
    # strip carriage return
    contmd5="$(echo "${contmd5}" | tr -d '\r')"
    [ "${hostmd5}" != "${contmd5}" ] && exit 1
    # delete the image each iteration to avoid dangling images
    docker rmi "${IMAGE_NAME}"

    tries=$((tries+1))
done

exit 0
