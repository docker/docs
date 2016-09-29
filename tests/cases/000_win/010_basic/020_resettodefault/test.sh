#!/bin/sh
# SUMMARY: run a container, reset to default, check if no images are left
# LABELS:
# AUTHOR: Emmanuel Briney <emmanuel.briney@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    d4w_backend_cli -Mount=C
}
trap clean_up EXIT

docker run alpine ls

d4w_backend_cli -ResetToDefault

out=$(docker images)
[[ "$out" == *"alpine"* ]] && exit 1

exit 0
