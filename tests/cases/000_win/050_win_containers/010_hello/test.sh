#!/bin/sh
# SUMMARY: Hello world with windows containers
# LABELS: skip
# AUTHOR: David Gageot <david.gageot@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    d4w_backend_cli -SwitchDaemon
}
trap clean_up EXIT

# Assume points to linux containers
d4w_backend_cli -SwitchDaemon
out=$(docker run --rm microsoft/nanoserver powershell.exe 'Write-Host "Hello"')
[ "${out}" != "Hello" ] && exit 1

exit 0
