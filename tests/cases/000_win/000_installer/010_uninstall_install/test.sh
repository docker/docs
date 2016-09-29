#!/bin/sh
# SUMMARY: Uninstall and re-Install Docker for Windows
# LABELS: master, release
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

if ! d4x_app_installed; then
    echoerr "Docker is supposed to be installed here"
    exit "${RT_TEST_CANCEL}"
fi

d4x_app_uninstall
d4x_app_install
# Installer does not return an error code (can't get it...)
exit 0
