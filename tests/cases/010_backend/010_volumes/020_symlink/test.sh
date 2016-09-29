#!/bin/sh
# SUMMARY: Verifies that symlinks work on mounted volumes
# LABELS: !win
# Doesn't work over CIFS
# AUTHOR: Dave Scott <dave.scott@docker.com>

set -e
IMAGE_NAME=volumes_symlink
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    docker rmi "$IMAGE_NAME" || true
    rm -rf "${D4X_LOCAL_TMPDIR}/foo" "${D4X_LOCAL_TMPDIR}/bar"
}
trap clean_up EXIT

rm -rf "${D4X_LOCAL_TMPDIR}/foo" "${D4X_LOCAL_TMPDIR}/bar"
docker build -t "$IMAGE_NAME" .
docker run --rm -v "${D4X_LOCAL_TMPDIR}:/tmp" -t "$IMAGE_NAME" /symlink.sh
ls -l "${D4X_LOCAL_TMPDIR}/foo" "${D4X_LOCAL_TMPDIR}/bar"
out=$(readlink "${D4X_LOCAL_TMPDIR}/bar")
[ "$out" != "foo" ] && exit 1
exit 0
