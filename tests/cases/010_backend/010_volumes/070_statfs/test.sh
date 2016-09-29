#!/bin/sh
# SUMMARY: Test statfs via df
# LABELS: !win
# AUTHOR: David Sheets <david.sheets@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

docker run --rm -v /tmp:/tmp alpine df
exit 0
