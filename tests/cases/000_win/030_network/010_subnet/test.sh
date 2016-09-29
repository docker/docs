#!/bin/sh
# SUMMARY: Change the subnet configuration
# LABELS: win
# AUTHOR: David Gageot <david.gageot@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    d4w_backend_cli -SetIP=10.0.75.0/255.255.255.0
}
trap clean_up EXIT

d4w_load_nsenter

d4w_nsenter ifconfig hvint0 | grep -q 'inet addr:10.0.75.2'

d4w_backend_cli -SetIP=10.0.76.0/255.255.255.248

d4w_nsenter ifconfig hvint0 | grep -q 'inet addr:10.0.76.2'

exit 0
