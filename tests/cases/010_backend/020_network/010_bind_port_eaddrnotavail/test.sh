#!/bin/sh
# SUMMARY: Test that binding a port to an interface which doesn't exist fails with a decent error
# LABELS:
# AUTHOR: Dave Scott <dave.scott@docker.com>
# AUTHOR: Magnus Skjegstad <magnus.skjegstad@docker.com>

# This tests [docker/pinata#4256]

# Plain docker responds like this:
# $ docker run -p 1.2.3.4:8080:80 nginx
# docker: Error response from daemon: driver failed programming external
# connectivity on endpoint prickly_snyder
# (8c51934f707883edcac3adcfef6ddd6ba3b49c5be327b8f5873913a2d4e04a75):
# Error starting userland proxy: listen tcp 1.2.3.4:8080: bind: cannot assign
# requested address.


set -x
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

IMAGE_NAME=bind_port_eaddrnotavail

clean_up() {
    CONTAINER="${IMAGE_NAME}"
    d4x_cleanup_image "${CONTAINER}"
}
trap clean_up EXIT

build() {
    CONTAINER="${IMAGE_NAME}"
    docker build -t "$CONTAINER" .
}

test_address() {
    CONTAINER="${IMAGE_NAME}"

    echo "starting container"
    OUTPUT=$(docker run --rm -p "$1":8081:8081 -t "$CONTAINER" echo hello 2>&1)
    if [ $? -eq 0 ]; then
        # This shouldn't be possible: the netcat should have the port
        echo "container was able to bind IP successfully to IP $1"
        exit 1
    fi
    case "$OUTPUT" in
        *"errno 526"*)
            echo "error message contained Linux 9P errno 526"
            exit 1
            ;;
        *"Unix_error"*)
            echo "error message contained OCaml exception name"
            exit 1
            ;;
        *"EADDRNOTAVAIL"*)
            echo "error message contained EADDRNOTAVAIL"
            exit 1
            ;;
        *"bind: cannot assign requested address"*)
            echo "error response looks OK"
            ;;
        *)
            echo "unknown error from docker"
            exit 1
            ;;
    esac
}

build

test_address "1.2.3.4"
test_address "192.168.65.2"
