#!/bin/sh
# SUMMARY: Basic check for windows containers
# LABELS: skip
# AUTHOR: David Gageot <david.gageot@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    d4w_backend_cli -SwitchDaemon
}
trap clean_up EXIT

# Assume we point to linux containers
out=$(docker info | grep "Operating System")
[ "$(echo "${out}" | grep -c 'Linux')" != "1" ] && exit 1
[ "$(echo "${out}" | grep -c 'Windows')" != "0" ] && exit 1

# Switch to windows containers
d4w_backend_cli -SwitchDaemon
out=$(docker info | grep "Operating System")
[ "$(echo "${out}" | grep -c 'Linux')" != "0" ] && exit 1
[ "$(echo "${out}" | grep -c 'Windows')" != "1" ] && exit 1

exit 0
