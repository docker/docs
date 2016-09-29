#!/bin/sh
# SUMMARY: Basic docker-compose test
# LABELS:
# AUTHOR: Dave Scott <dave.scott@docker.com>
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up () {
    docker-compose down --rmi all || true
}

trap clean_up EXIT

docker-compose pull
docker-compose up
docker-compose stop
docker-compose down --rmi all

exit 0
