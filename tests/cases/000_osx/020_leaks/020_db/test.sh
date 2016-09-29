#!/bin/sh
# SUMMARY: Check db for resource leaks
# LABELS: osx
# REPEAT:
# AUTHOR: David Sheets <dsheets@docker.com>

set -e # Exit on error
# Source libraries. Uncomment if needed/defined
#. "${RT_ROOT}/lib/lib.sh"
. "${RT_PROJECT_ROOT}/_lib/lib.sh"
# IMAGE_NAME=  # Use a env variable to name images/containers

clean_up() {
    # remove any files, containers, images etc
    true
}
trap clean_up EXIT

PID=$(pgrep com.docker.db)

# Warm up the process
for i in $(seq 1 10); do
    echo "$i" > /dev/null
    docker run --rm -p 80:80 -v "$(pwd):/host" alpine true
done

# Rest the process
sleep 4

# Sample resource counts
THREADS=$(ps -M "$PID" | wc -l)
FDS=$(lsof -p "$PID" | wc -l)
#MEMORY_RSS=$(ps u "$PID" | tail -1 | awk '{ print $6 }')

# Exercise the process
for i in $(seq 1 10); do
    echo "$i" > /dev/null
    docker run --rm -p 80:80 -v "$(pwd):/host" alpine true
done

# Rest the process
sleep 4

# Re-sample resource counts
THREADS2=$(ps -M "$PID" | wc -l)
FDS2=$(lsof -p "$PID" | wc -l)
#MEMORY_RSS2=$(ps u "$PID" | tail -1 | awk '{ print $6 }')

# Compare
FAILED=0
if [ "$THREADS" -lt "$THREADS2" ]; then
    echo "Possible thread leak: $THREADS -> $THREADS2"
    FAILED=1
fi

if [ "$FDS" -lt "$FDS2" ]; then
    echo "Possible fd leak: $FDS -> $FDS2"
    FAILED=1
fi

#if [ "$MEMORY_RSS" -lt "$MEMORY_RSS2" ]; then
#    echo "Possible memory leak: $MEMORY_RSS -> $MEMORY_RSS2"
#    FAILED=1
#fi

exit "$FAILED"
