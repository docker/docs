#!/bin/sh
# SUMMARY: Transfer and verify fixed sized UDP datagrams to container
# LABELS: !win
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>
# AUTHOR: David Scott <dave.scott@docker.com>
# see https://github.com/docker/pinata/issues/5237

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

IMAGE_NAME=udp_ping
PORT=6000
SIZES="1 4 511 1023 2034 2035 4095 8191 9215 9216"
# on OSX the maximum we can send is given by $(sysctl net.inet.udp.maxdgram)
# which seems to be 9216 for me

clean_up() {
    rm -rf testfile outfile || true
    docker stop "${IMAGE_NAME}" || true
    docker rm "${IMAGE_NAME}" || true
    d4x_cleanup_image "${IMAGE_NAME}"
}
trap clean_up EXIT

# in case of debugging leftovers
clean_up

docker build -t "${IMAGE_NAME}" .


for size in ${SIZES}; do
    echo "UDP datagram size ${size}"

    docker run -d -p "${PORT}":"${PORT}"/udp \
           --name "${IMAGE_NAME}" "${IMAGE_NAME}" \
           0.0.0.0 "${PORT}"

    python sender.py "${D4X_HOST_NAME}" "${PORT}" "${size}"

    docker kill "${IMAGE_NAME}"
    docker rm "${IMAGE_NAME}"
done

exit 0
