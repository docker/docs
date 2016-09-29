#!/bin/sh
# SUMMARY: Repeat uninstall/install tests
# LABELS: installer
# AUTHOR: Magnus Skjegstad <magnus.skjegstad@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

d4x_app_uninstall
d4x_app_install
