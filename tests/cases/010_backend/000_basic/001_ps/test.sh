#!/bin/sh
# SUMMARY: Runs `docker ps` and checks that it works
# LABELS:
# REPEAT: 2
# AUTHOR: Dave Scott <dave.scott@docker.com>

set -e # exit on error
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

out=$(docker ps)

[ -z "$out" ] && exit 1

exit 0
