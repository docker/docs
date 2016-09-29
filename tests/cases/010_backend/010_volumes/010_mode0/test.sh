#!/bin/sh
# SUMMARY: Test setattr on arbitrary mode files
# LABELS: !win
# AUTHOR: David Sheets <dsheets@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    rm -f /tmp/mode0
}
trap clean_up EXIT

rm -f /tmp/mode0
touch /tmp/mode0 || exit 1
chmod 0 /tmp/mode0 || exit 1
docker pull busybox || exit 1
docker run --rm -v /tmp:/tmp busybox chmod 600 /tmp/mode0 || exit 1
echo "checking permissions after container exit"
[ "$(stat -f '%Lp' /tmp/mode0)" == "600" ] || exit 1
echo success
exit 0
