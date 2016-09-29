#!/bin/sh
# SUMMARY: Concurrent tests of simple docker commands
# LABELS:
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

# default values (5 threads, 5 iterations each)
CONC=5
ITER=10

if rt_label_set "master"; then
    CONC=10
    ITER=20
elif rt_label_set "nightly"; then
    CONC=10
    ITER=100
elif rt_label_set "release"; then
    CONC=10
    ITER=100
fi


"${RT_UTILS}/rt-crexec" -c "${CONC}" -i "${ITER}" commands.txt
exit 0
