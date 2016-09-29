#!/bin/sh -e

now=$(date +%s)
mkdir -p results
find=results/find.$now
mapfind=results/mapfind.$now
for test in $(ls testcases); do
    printf "%s %s\n" "${test}" $(cat logs/find-$test.time) >> $find
    printf "%s %s\n" "${test}" $(cat logs/mapfind-$test.time) >> $mapfind
done
echo
