#!/bin/sh
# SUMMARY: Install from specified binary. This test requires passwordless sudo.
# LABELS: installer
# AUTHOR: Magnus Skjegstad <magnus.skjegstad@docker.com>

set -e

. "${RT_PROJECT_ROOT}/_lib/lib.sh"

d4x_app_install
# Wait for app to start
d4x_wait_for_docker || (echo "Docker process not running." ; syslog -k Facility com.docker.docker | tail -n 5 ; exit 1)
