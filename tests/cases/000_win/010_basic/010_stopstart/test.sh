#!/bin/sh
# SUMMARY: Stop backend, check that docker isn't working anymore, restart the backend and check that docker is ok
# LABELS: win
# AUTHOR: Emmanuel Briney <emmanuel.briney@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

d4w_backend_cli -Stop

set +e
docker ps
[ $? -eq 0 ] && exit 1
set -e

d4w_backend_cli -Start

out=$(docker ps)
[ -z "$out" ] && exit 1

exit 0
