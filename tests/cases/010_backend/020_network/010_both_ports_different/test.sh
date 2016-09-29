#!/bin/sh
# SUMMARY: Test that port forwarding works for different ports
# LABELS:
# AUTHOR: Dave Scott <dave.scott@docker.com>
# AUTHOR: David Gageot <david.gageot@docker.com>

set -e 
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

IMAGE_NAME=both_ports_different

clean_up() {
    docker rmi "${IMAGE_NAME}" || true
}
trap clean_up EXIT

docker build -t ${IMAGE_NAME} .

CONTAINER=$(docker run -d -p 8082:8081 "${IMAGE_NAME}")

"${RT_UTILS}/rt-urltest" -r 60 -s hello "http://${D4X_HOST_NAME}:8082"

docker kill "${CONTAINER}" || true

exit 0
