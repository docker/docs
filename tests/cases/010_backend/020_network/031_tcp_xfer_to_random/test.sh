#!/bin/sh
# SUMMARY: Transfer and verify randomly sized files to container over TCP
# LABELS:
# AUTHOR: Rolf Neugebauer <rolf.neugebauer@docker.com>
# see https://github.com/docker/pinata/issues/4491

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

IMAGE_NAME=ncat_random
PORT=6001

# default runs: 5 iterations, max file size 1M
ITERS=5
MAX_SIZE=0x100000

if rt_label_set "master"; then
    # 10 iterations with up to 50MB files
    ITERS=10
    MAX_SIZE=0x3200000
elif rt_label_set "nightly"; then
    # 30 iterations with up to 100MB files
    ITERS=30
    MAX_SIZE=0x6400000
elif rt_label_set "release"; then
    # 30 iterations with up to 100MB files
    ITERS=30
    MAX_SIZE=0x6400000
fi

clean_up() {
    rm -rf testfile outfile || true
    docker stop "${IMAGE_NAME}" || true
    docker rm "${IMAGE_NAME}" || true
    d4x_cleanup_image "${IMAGE_NAME}" || true
}
trap clean_up EXIT

docker build -t "${IMAGE_NAME}" .

iter=0
while true; do
    if [ "$iter" -gt ${ITERS} ]; then
        break
    fi

    "${RT_UTILS}/rt-filerandgen" "${MAX_SIZE}" testfile
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

    docker stop "${IMAGE_NAME}" || true
    docker rm "${IMAGE_NAME}"

    ls -l ./*file

    [ "${hostmd5}" != "${contmd5}" ] && exit 1

    iter=$((iter+1))
done

exit 0
