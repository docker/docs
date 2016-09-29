#!/bin/sh
# SUMMARY: Start the application and make it's running
# LABELS:
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

if ! d4x_app_installed; then
    echoerr "Docker is supposed to be installed here"
    exit 1
fi

if d4x_app_running; then
    echoerr "Docker is not supposed to be running"
    exit "${RT_TEST_CANCEL}"
fi

d4x_app_start
exit 0
