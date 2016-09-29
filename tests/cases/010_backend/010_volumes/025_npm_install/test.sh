#!/bin/sh
# SUMMARY: Verifies that npm can create bin symlinks on a mount
# LABELS:
# AUTHOR: David Gageot <david.gageot@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
	rm -rf "${D4X_LOCAL_TMPDIR}/node" || true
}
trap clean_up EXIT

clean_up
docker run --rm -w /tmp -v "${D4X_LOCAL_TMPDIR}/node":/tmp node:6.3.1-slim \
       npm install coffee-script@1.10.0 standard-format@2.2.2
exit 0
