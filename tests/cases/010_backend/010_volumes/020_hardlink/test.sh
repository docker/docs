#!/bin/sh
# SUMMARY: Verifies that hard links work on mounted volumes
# LABELS:
# AUTHOR: Dave Scott <dave.scott@docker.com>

set -e
IMAGE_NAME=hardlink
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    docker rmi "$IMAGE_NAME" || true
    rm -rf "${D4X_LOCAL_TMPDIR}/foo" "${D4X_LOCAL_TMPDIR}/bar"
}
trap clean_up EXIT

rm -rf "${D4X_LOCAL_TMPDIR}/foo" "${D4X_LOCAL_TMPDIR}/bar"
docker build -t "$IMAGE_NAME" .
docker run --rm -v "${D4X_LOCAL_TMPDIR}:/tmp" -t "$IMAGE_NAME" /hardlink.sh
ls -l "${D4X_LOCAL_TMPDIR}/foo" "${D4X_LOCAL_TMPDIR}/bar"
out=$(cat "${D4X_LOCAL_TMPDIR}/bar")
[ "$out" != "hello" ] && exit 1
exit 0
