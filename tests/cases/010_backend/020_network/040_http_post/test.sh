#!/bin/sh
# SUMMARY: Start a webserver in a container and issue HTTP POST requests
# LABELS:
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>
# See: https://github.com/docker/pinata/issues/4399

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"
IMAGE_NAME=http_post
PORT=9090
SIZES="128 4096 32768 65536 131072"

clean_up() {
    rm -f tesfile
    docker kill ${IMAGE_NAME} || true
    docker rm ${IMAGE_NAME} || true
    d4x_cleanup_image ${IMAGE_NAME}
}
trap clean_up EXIT

docker build -t "${IMAGE_NAME}" .
docker run -d --name "${IMAGE_NAME}" -p "${PORT}":"${PORT}" "${IMAGE_NAME}" /http_post.py "${PORT}"
sleep 5

# create files of different sizes and POST them using curl
for size in ${SIZES}; do
    echo "Trying file size ${size}"
    "${RT_UTILS}/rt-filegen" "${size}" testfile
    curl --retry 5 -X POST --data-binary @testfile "${D4X_HOST_NAME}":"${PORT}"
    rm -f testfile
done

exit 0
