#!/bin/sh
# SUMMARY: Uninstall and re-Install Docker for Windows while App is running
# LABELS: master, release
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

if ! d4x_app_installed; then
    echoerr "Docker is supposed to be installed here"
    exit "${RT_TEST_CANCEL}"
fi

if ! d4x_app_running; then
    echoerr "Docker is supposed to be running"
    exit "${RT_TEST_CANCEL}"
fi

d4x_app_uninstall
d4x_app_install
d4x_app_start

exit 0
