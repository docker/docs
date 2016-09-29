#!/bin/sh
# SUMMARY: Start a UCP instance
# LABELS: master, release, nightly, apps
# AUTHOR: David Scott <dave.scott@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    echo y | docker run --rm -i -v /var/run/docker.sock:/var/run/docker.sock --name ucp docker/ucp uninstall --interactive || true
    docker kill ucp || true
    docker rm ucp || true
}
trap clean_up EXIT

docker run --rm -i -v /var/run/docker.sock:/var/run/docker.sock --name ucp docker/ucp install --fresh-install

"${RT_UTILS}/rt-urltest" -r 60 -s "Universal Control Plane" "https://${D4X_HOST_NAME}:443"
exit 0
