#!/bin/sh
# SUMMARY: Make sure Docker for Windows is not installed/running
# LABELS:
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

if d4x_app_running; then
    echoerr "Docker is not supposed to be installed here"
    exit 1
fi

if d4x_app_installed; then
    echoerr "Docker is not supposed to be installed here"
    exit 1
fi

exit 0
