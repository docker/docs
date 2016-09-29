#!/bin/sh
# SUMMARY: Run nmap from container against the host interface (longer/more aggressive for master/release)
# LABELS:
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>
# Created because of #3882

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"
IMAGE_NAME=nmap_host

clean_up() {
    d4x_cleanup_image ${IMAGE_NAME}
}
trap clean_up EXIT

# defaults
CONC=0
PORTS=

if rt_label_set "master"; then
    PORTS="-p-"
elif rt_label_set "nightly"; then
    CONC=9
    PORTS="-p-"
elif rt_label_set "release"; then
    CONC=9
    PORTS="-p-"
fi

case "${RT_OS}" in
    osx)
        IP=$(ifconfig | grep "inet " | grep -v "127.0.0.1" | head -1 | awk '{print $2}')
        ;;
    win)
        # Make sure we don't use the Hyper-V network interface
        IP=$(rt_ps_cmd "((Get-NetIPConfiguration).IPv4Address).IPAddress" | grep -v "10.0.75.1" | head -1)
        ;;
esac

docker build -t "${IMAGE_NAME}" .
docker run --rm -t "${IMAGE_NAME}" /run.sh "${IP}" "${CONC}" "${PORTS}"
