#!/bin/sh
# SUMMARY: Verify that the host volume is mounted in VM
# LABELS: osx
# Test is currently OS X specific
# AUTHOR: David Sheets <david.sheets@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

CMD="echo -e '\n************\n' > sep && \
    mount > mtab && \
    cat sep mtab sep && \
    apk update && \
    apk add util-linux && \
    echo -e '\n**** VM ****\n' && \
    nsenter -t 1 -m cat /etc/mtab && \
    cat sep && \
    nsenter -t 1 ps aucx && \
    cat sep && \
    grep osxfs mtab"

docker run --privileged --pid=host --rm -v "${D4X_LOCAL_TMPDIR}:/tmp" alpine:3.3 sh -c "$CMD"
