#!/bin/sh
# SUMMARY: Check various docker inspect invocations
# LABELS:
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>
# Add due to: https://github.com/docker/pinata/issues/4028

set -e # Exit on error
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

CONTAINER_NAME=inspect_test
IMAGE_NAME=hello-world

clean_up() {
    docker rm --force ${CONTAINER_NAME}
    docker rmi ${IMAGE_NAME} || true # container might be in use
}
trap clean_up EXIT

docker run --name ${CONTAINER_NAME} ${IMAGE_NAME}
docker inspect -s ${CONTAINER_NAME}
docker inspect -s ${IMAGE_NAME}
exit 0
