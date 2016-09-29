#!/bin/sh
# SUMMARY: Test that port forwarding works for the same port on containers started after eachother
# LABELS:
# AUTHOR: Dave Scott <dave.scott@docker.com>
# AUTHOR: Magnus Skjegstad <magnus.skjegstad@docker.com>
# AUTHOR: David Gageot <david.gageot@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

IMAGE_NAME=bind_priv_port_twice
REP=5

clean_up() {
    d4x_cleanup_image "${IMAGE_NAME}"
}
trap clean_up EXIT

docker build -t "${IMAGE_NAME}" .

for x in $(seq 1 "$REP"); do 
    echo "starting container $x"
    CONTAINER=$(docker run -d -p 81:81 "${IMAGE_NAME}")

    echo "connecting ..."
    "${RT_UTILS}/rt-urltest" -r 60 -s hello "http://${D4X_HOST_NAME}:81"

    echo "stop container"
    docker kill "${CONTAINER}" || true
done

exit 0
