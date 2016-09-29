#!/bin/sh
# SUMMARY: Change the cpu/memory configuration
# LABELS: win
# AUTHOR: David Gageot <david.gageot@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
    d4w_backend_cli -SetCpus=2 -SetMemory=2048
}
trap clean_up EXIT

d4w_load_nsenter

out=$(docker run --rm --privileged --pid=host d4w/nsenter /bin/sh -c "cat /proc/cpuinfo | grep -c 'processor'")
[ "${out}" != "2" ] && exit 1

out=$(docker run --rm --privileged --pid=host d4w/nsenter /bin/sh -c "cat /proc/meminfo | grep 'MemTotal' | grep -c '2013176'")
[ "${out}" != "1" ] && exit 1

d4w_backend_cli -SetCpus=1 -SetMemory=3000

out=$(docker run --rm --privileged --pid=host d4w/nsenter /bin/sh -c "cat /proc/cpuinfo | grep -c 'processor'")
[ "${out}" != "1" ] && exit 1

out=$(docker run --rm --privileged --pid=host d4w/nsenter /bin/sh -c "cat /proc/meminfo | grep 'MemTotal' | grep -c '2968548'")
[ "${out}" != "1" ] && exit 1

exit 0
