#!/bin/sh
# SUMMARY: Check that both server and client are experimental or not
# LABELS:
# AUTHOR: David Gageot <david.gageot@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

client=$(docker version -f "{{.Client.Experimental}}")
server=$(docker version -f "{{.Server.Experimental}}")
echo "${client}"
echo "${server}"
[ "${client}" != "${server}" ] && exit 1

exit 0
