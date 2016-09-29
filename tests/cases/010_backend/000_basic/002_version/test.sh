#!/bin/sh
# SUMMARY: Runs `docker version` (TODO: check that they match)
# LABELS:
# AUTHOR: Dave Scott <dave.scott@docker.com>
# AUTHOR: Magnus Skjegstad <magnus.skjegstad@docker.com>

set -e # Exit on error
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

docker version | grep Version | awk '
BEGIN {
    last="unset"; 
    cnt=0
}
{ 
    cnt++
    if (last != "unset" && last != $2) { 
        print "Version mismatch: ", $2, "vs", last; 
        exit 1 
    } 
    last = $2 
}
END {
    if (cnt != 2) {
        print "Expected two version numbers, got",cnt
        exit 1
    }
}'

exit 0
