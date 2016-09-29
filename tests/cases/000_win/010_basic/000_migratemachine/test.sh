#!/bin/sh
# SUMMARY: migrate a docker-machine disk and check that the images are there
# LABELS: win
# AUTHOR: Emmanuel Briney <emmanuel.briney@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

bashpath=$(pwd)
winpath=$(cygpath -w "$bashpath")

d4w_backend_cli -MigrateVolume="$winpath/disk.vmdk"

out=$(docker images)
[[ "$out" != *"test"* ]] && exit 1

out=$(docker ps -a)
[[ "$out" != *"test"* ]] && exit 1

exit 0
