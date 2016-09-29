#!/bin/sh
# SUMMARY: Run docker-diagnose after install
# LABELS: installer
# AUTHOR: Magnus Skjegstad <magnus.skjegstad@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

# This test is racy so adding a sleep to prevent false positives
sleep 10
"${OSX_APP_DIR}/Contents/Resources/bin/docker-diagnose" -n -v -o "$RT_RESULTS/diagnose.tar.gz"
