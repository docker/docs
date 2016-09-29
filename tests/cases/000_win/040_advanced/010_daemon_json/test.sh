#!/bin/sh
# SUMMARY: Change the debug daemon settings
# LABELS: win
# AUTHOR: David Gageot <david.gageot@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    d4w_backend_cli -SetDaemonJson="{\"registry-mirrors\":[],\"insecure-registries\":[]}"
}
trap clean_up EXIT

docker info | grep -q 'Debug Mode (server): true'

d4w_backend_cli -SetDaemonJson="{\"debug\":false}"

docker info | grep -q 'Debug Mode (server): false'

exit 0
