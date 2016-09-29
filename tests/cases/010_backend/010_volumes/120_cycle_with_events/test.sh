#!/bin/sh
# SUMMARY: Rapidly restart containers with un-normalized mounts + events
# LABELS: !win
# AUTHOR: David Sheets <david.sheets@docker.com>
# Disabled on windows as it uses symlinks on the host

set -e # Exit on error
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    rm -rf test_dir foo
}
trap clean_up EXIT

clean_up

TEST_PATH=test_dir/with/very/long/path/so/that/we/can/receive/shorter/events

mkdir -p $TEST_PATH
ln -s "$(pwd)" $TEST_PATH/symlink

for i in $(seq 1 10)
do
    echo "$i" > /dev/null # makes shellcheck shutup
    docker run --rm -v "$(pwd)"/$TEST_PATH/symlink/:/host/ alpine ls /host
    touch foo
    rm foo
done

exit 0
