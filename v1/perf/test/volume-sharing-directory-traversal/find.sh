#!/bin/bash
test=$1
find=/usr/bin/find
/usr/bin/time -f '%e' -o volumes/find-${test}.time \
	sh -ec "$find volumes/${test} >/dev/null"
/usr/bin/time -f '%e' -o volumes/mapfind-${test}.time \
        sh -ec "$find volumes/${test} -name w-\\* | xargs $find >/dev/null"
