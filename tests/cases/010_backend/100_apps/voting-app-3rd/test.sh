#!/bin/sh
# SUMMARY: Docker 3rd birthday voting app
# LABELS: master, release, nightly, apps
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>
# Based on: https://github.com/docker/docker-birthday-3/

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    docker-compose down || true
}
trap clean_up EXIT


docker-compose up -d
sleep 5

"${RT_UTILS}/rt-urltest" -r 10 -s docker "http://${D4X_HOST_NAME}:5000"
"${RT_UTILS}/rt-urltest" -r 5 -s votes "http://${D4X_HOST_NAME}:5001"

# For debug: Sometimes either the voting or the result container don't
# seem to work. "Check" if they are up
docker ps

exit 0
