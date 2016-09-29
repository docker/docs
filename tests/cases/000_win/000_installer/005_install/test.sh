#!/bin/sh
# SUMMARY: Install Docker for Windows
# LABELS:
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

d4x_app_install
# Installer does not return an error code (can't get it...)
exit 0
