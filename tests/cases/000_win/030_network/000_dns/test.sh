#!/bin/sh
# SUMMARY: Change the DNS server configuration
# LABELS: win
# AUTHOR: David Gageot <david.gageot@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    d4w_backend_cli -SetDNS=automatic
}
trap clean_up EXIT

d4w_backend_cli -SetDNS=8.8.8.8

docker run --rm alpine wget http://www.google.com

exit 0
