#!/bin/sh
# SUMMARY: Run `echo hello` in busy box
# LABELS:
# AUTHOR: Dave Scott <dave.scott@docker.com>

set -e # Exit on error
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

out=$(docker run --rm busybox echo hello)
[ -z "$out" ] && exit 1

exit 0
