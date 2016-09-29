#!/bin/sh
# SUMMARY: Test setting atime and mtime
# LABELS: !win
# stat -f %m is not working on windows
# AUTHOR: David Sheets <david.sheets@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"
TEST_FILE="${D4X_LOCAL_TMPDIR}/test_utimes"

clean_up() {
    rm -rf "${TEST_FILE}"
}
clean_up
trap clean_up EXIT

docker run --rm -v "${D4X_LOCAL_TMPDIR}:/tmp" -e TZ=UTC busybox \
       touch -t 197001010000.00 "${TEST_FILE}"

[ "$(stat -f %m "${TEST_FILE}")" = 0 ]
[ "$(stat -f %a "${TEST_FILE}")" = 0 ]

rm "${TEST_FILE}"
mkdir "${TEST_FILE}"

docker run --rm -v "${D4X_LOCAL_TMPDIR}:/tmp" -e TZ=UTC busybox \
       touch -t 197001010000.00 "${TEST_FILE}"

[ "$(stat -f %m "${TEST_FILE}")" = 0 ]
[ "$(stat -f %a "${TEST_FILE}")" = 0 ]

rmdir "${TEST_FILE}"
ln -s "${TEST_FILE}" "${TEST_FILE}"

docker run --rm -v "${D4X_LOCAL_TMPDIR}:/tmp" -e TZ=UTC busybox \
       touch -h -t 197001010000.00 "${TEST_FILE}"

[ "$(stat -f %m "${TEST_FILE}")" = 0 ]
[ "$(stat -f %a "${TEST_FILE}")" = 0 ]

exit 0
