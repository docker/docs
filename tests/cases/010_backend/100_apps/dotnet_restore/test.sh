#!/bin/sh
# SUMMARY: Restore a .NET project
# LABELS: master, release, nightly, apps
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>
# Created to repro #3886

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

IMAGE_NAME=dotnet_restore

clean_up() {
    d4x_cleanup_image ${IMAGE_NAME}
}
trap clean_up EXIT

docker build -t $IMAGE_NAME .
exit 0
