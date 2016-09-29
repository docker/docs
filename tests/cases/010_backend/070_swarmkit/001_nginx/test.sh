#!/bin/sh
# SUMMARY: Basic docker 1.12 swarm mode test
# LABELS:
# AUTHOR: Gaetan de Villele <gaetan.devillele@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"
IMAGE_NAME=webserver

clean_up() {
    docker service rm "${IMAGE_NAME}" || true
}
trap clean_up EXIT

# create docker service
docker service create --replicas 1 --name "${IMAGE_NAME}" --publish 80:80/tcp nginx:alpine
sleep 5

# try to load a web page every second for 120 seconds
"${RT_UTILS}/rt-urltest" -r 120 "http://${D4X_HOST_NAME}:80"
exit 0
