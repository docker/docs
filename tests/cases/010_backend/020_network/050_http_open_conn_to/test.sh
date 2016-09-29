#!/bin/sh
# SUMMARY: Create many connections to a container and test some
# LABELS: skip
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>

# This test does not pass on windows at the moment. needs investigating

. "${RT_PROJECT_ROOT}/_lib/lib.sh"
IMAGE_NAME=http_conn_open

clean_up() {
    d4x_cleanup_image ${IMAGE_NAME}
}

# Setting this to 500 works
# Setting it to 512 or higher fails.
NUM_CONN=600

trap clean_up EXIT

docker run -d --name "${IMAGE_NAME}" -p 9081:80 nginx:alpine
sleep 5
./http_open_conn.py -c "${NUM_CONN}" -t 100 -s "${D4X_HOST_NAME}" -p 9081
