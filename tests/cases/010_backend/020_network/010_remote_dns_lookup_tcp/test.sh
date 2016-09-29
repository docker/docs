#!/bin/sh
# SUMMARY: Test that DNS using a remote DNS server works
# LABELS:
# AUTHOR: Dave Scott <dave.scott@docker.com>

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"
IMAGE_NAME=remote_dns_lookup_udp

clean_up() {
    docker rmi "${IMAGE_NAME}" || true
}
trap clean_up EXIT

docker build -t "${IMAGE_NAME}" .
docker run --rm -t "${IMAGE_NAME}" /run.sh -t docker.io MX 8.8.8.8 aspmx.l.google.com

exit 0
