#!/bin/sh
# SUMMARY: Test sharing and unsharing of the C drive
# LABELS: win
# AUTHOR: Emmanuel Briney <emmanuel.briney@docker.com>
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

TESTDIR="${D4X_LOCAL_TMPDIR}/modify_share"

clean_up() {
    rm -rf "${TESTDIR}"
    d4w_backend_cli -Mount=C
}
trap clean_up EXIT

rm -rf "${TESTDIR}" || true
mkdir -p "${TESTDIR}"
touch "${TESTDIR}"/foo

out=$(docker run --rm -v "${TESTDIR}":/testdir alpine ls /testdir)
[ -z "$out" ] && exit 1

d4w_backend_cli -Unmount=C

set +e
out=$(docker run --rm -v "${TESTDIR}":/testdir alpine ls /testdir)
EXIT_CODE=$?
set -e

[ "$EXIT_CODE" -eq 0 ] && exit 1
[ "$out" != "${out/drive is not shared/}" ] && exit 1

d4w_backend_cli -Mount=C

out=$(docker run --rm -v "${TESTDIR}":/testdir alpine ls /testdir)
[ -z "$out" ] && exit 1

exit 0
