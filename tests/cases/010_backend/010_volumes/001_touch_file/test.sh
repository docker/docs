#!/bin/sh
# SUMMARY: Volume mount a directory and touch a file
# LABELS:
# AUTHOR: Dave Scott <dave.scott@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    rm -f "${D4X_LOCAL_TMPDIR}/touch_a_file"
}
trap clean_up EXIT

rm -f "${D4X_LOCAL_TMPDIR}/touch_a_file"

docker pull busybox
docker run --rm -v "${D4X_LOCAL_TMPDIR}":/tmp busybox touch /tmp/touch_a_file
ls -l "${D4X_LOCAL_TMPDIR}/touch_a_file"

exit 0
