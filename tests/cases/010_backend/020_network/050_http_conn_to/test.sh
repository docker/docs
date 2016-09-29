#!/bin/sh
# SUMMARY: Create many HTTP connections to a container as fast as possible
# LABELS:
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"
IMAGE_NAME=http_conn

clean_up() {
    d4x_cleanup_image ${IMAGE_NAME}
}

# default
NUM_CONN=1000

# Would love to run with more connections but it takes ages already
if rt_label_set "master"; then
    NUM_CONN=20000
elif rt_label_set "nightly"; then
    NUM_CONN=60000
elif rt_label_set "release"; then
    NUM_CONN=60000
fi

set -e
trap clean_up EXIT

docker run -d --name "${IMAGE_NAME}" -p 9080:80 nginx:alpine
sleep 5
./http_conn.py -c "${NUM_CONN}" "http://${D4X_HOST_NAME}:9080"

exit 0
