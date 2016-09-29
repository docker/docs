#!/bin/sh
# SUMMARY: Run many concurrent curl processes in a container
# LABELS: skip
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>
# See: https://github.com/docker/pinata/issues/4755

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"
IMAGE_NAME=http_conn_from

clean_up() {
    d4x_cleanup_image ${IMAGE_NAME}
}

# default
NUM_THDS=50
NUM_CONN=10

# Would love to run with more connections but it takes ages already
if rt_label_set "master"; then
    NUM_THDS=100
    NUM_CONN=50
elif rt_label_set "nightly"; then
    NUM_THDS=100
    NUM_CONN=200
elif rt_label_set "release"; then
    NUM_THDS=100
    NUM_CONN=200
fi

set -e
trap clean_up EXIT

docker build -t "${IMAGE_NAME}" .
docker run --rm -t "${IMAGE_NAME}" /http_conn.py -c "${NUM_CONN}" -t "${NUM_THDS}"
