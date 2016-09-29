#!/bin/sh
# SUMMARY: Transfer and verify fixed sized files to container over TCP
# LABELS:
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>
# see https://github.com/docker/pinata/issues/4491

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

IMAGE_NAME=ncat_fixed
PORT=6000
SIZES="1 4 4095 4096 4097 4098 4099 5000 5001 5002 1048575 1048576 1048577"

clean_up() {
    rm -rf testfile outfile || true
    docker stop "${IMAGE_NAME}" || true
    docker rm "${IMAGE_NAME}" || true
    d4x_cleanup_image "${IMAGE_NAME}"
}
trap clean_up EXIT

docker build -t "${IMAGE_NAME}" .


for size in ${SIZES}; do
    echo "File size ${size}"
    "${RT_UTILS}/rt-filegen" "${size}" testfile
    hostmd5=$("${RT_UTILS}/rt-filemd5" testfile)

    docker run -d -p "${PORT}":"${PORT}" \
           --name "${IMAGE_NAME}" "${IMAGE_NAME}" \
           sh -c "/usr/bin/ncat -l -p ${PORT} > outfile"

    # Give the container a chance to start. Do not remove this! There
    # is a race with slirp accepting a connection when the other end
    # is not yet up. This test is not testing this.
    sleep 5

    python xfer.py "${D4X_HOST_NAME}" "${PORT}" testfile
    # Wait for the container to exit
    tries=0
    while true; do
        status=$(docker inspect --format='{{.State.Status}}' "${IMAGE_NAME}")
        [ "${status}" = "exited" ] && break
        tries=$((tries+1))
        if [ "${tries}" -gt 30 ]; then
            echoerr "ncat container did not exit"
            # Gather some debug.
            docker ps
            docker inspect "${IMAGE_NAME}"
            docker run -it --rm --privileged --pid=host debian nsenter -t 1 -m -u -n -i ps -ef
            # shellcheck disable=SC2016
            docker run -it --rm --privileged --pid=host debian nsenter -t 1 -m -u -n -i /bin/sh -c 'for pid in $(pgrep ncat); do ls -l /proc/"${pid}"/fd; done'
            exit 1
        fi
        sleep 1
    done
    docker cp "${IMAGE_NAME}":/outfile ./outfile
    contmd5=$("${RT_UTILS}/rt-filemd5" outfile)

    docker stop "${IMAGE_NAME}"
    docker rm "${IMAGE_NAME}"

    ls -l ./*file

    [ "${hostmd5}" != "${contmd5}" ] && exit 1
done

exit 0
