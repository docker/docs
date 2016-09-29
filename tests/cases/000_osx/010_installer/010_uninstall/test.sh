#!/bin/sh
# SUMMARY: Uninstall existing Docker.app if installed.
# LABELS: installer
# AUTHOR: Magnus Skjegstad <magnus.skjegstad@docker.com>

set -e

. "${RT_PROJECT_ROOT}/_lib/lib.sh"

# Call uninstall
echo "### Uninstalling existing version from ${OSX_APP_DIR}"

d4x_app_installed

while pgrep -u root -q loginwindow; do
    echo "User is not yet logged in. Waiting..."
    sleep 2
done

d4x_app_uninstall
