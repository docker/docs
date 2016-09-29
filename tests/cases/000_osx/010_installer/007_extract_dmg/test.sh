#!/bin/sh
# SUMMARY: Extract DMG into ${D4X_LOCAL_DIR}, if ${D4X_LOCAL_INSTALLER} exists and not ${D4X_LOCAL_DIR}. If neither exists, try to download and extract from ${OSX_DOWNLOAD_URL}.
# LABELS: installer
# AUTHOR: Magnus Skjegstad <magnus.skjegstad@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"


if [ ! "$(d4x_app_installed)" ]; then
    echo "App already installed, skipped."
    exit "${RT_TEST_CANCEL}"
fi

d4m_extract_dmg
