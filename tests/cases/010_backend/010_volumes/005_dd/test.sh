#!/bin/sh
# SUMMARY: `dd` /dev/zero to a volume mounted file with various block sizes
# LABELS:
# AUTHOR: Justin Cormack <justin.cormack@docker.com>
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    rm -rf "${D4X_LOCAL_TMPDIR}/foo"
}
trap clean_up EXIT


# TODO we could check the size of the files

rm -rf "${D4X_LOCAL_TMPDIR}/foo"

docker run --rm -v "${D4X_LOCAL_TMPDIR}:/tmp" alpine dd if=/dev/zero bs=8k count=1 of=/tmp/foo
ls -l "${D4X_LOCAL_TMPDIR}/foo"

docker run --rm -v "${D4X_LOCAL_TMPDIR}:/tmp" alpine dd if=/dev/zero bs=16k count=1 of=/tmp/foo
ls -l "${D4X_LOCAL_TMPDIR}/foo"

docker run --rm -v "${D4X_LOCAL_TMPDIR}:/tmp" alpine dd if=/dev/zero bs=32k count=1 of=/tmp/foo
ls -l "${D4X_LOCAL_TMPDIR}/foo"

docker run --rm -v "${D4X_LOCAL_TMPDIR}:/tmp" alpine dd if=/dev/zero bs=64k count=1 of=/tmp/foo
ls -l "${D4X_LOCAL_TMPDIR}/foo"

docker run --rm -v "${D4X_LOCAL_TMPDIR}:/tmp" alpine dd if=/dev/zero bs=128k count=1 of=/tmp/foo
ls -l "${D4X_LOCAL_TMPDIR}/foo"

exit 0
