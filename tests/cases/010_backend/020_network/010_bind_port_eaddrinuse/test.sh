#!/bin/sh
# SUMMARY: Test that binding a port which is already in use fails with a decent error
# LABELS: !win
# AUTHOR: Dave Scott <dave.scott@docker.com>
# AUTHOR: Magnus Skjegstad <magnus.skjegstad@docker.com>

# Test is OS-X/Unix specific as it uses nc to bind to local port
# This tests [docker/pinata#4256]

set -x
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

IMAGE_NAME=bind_port_eaddrinuse

PID=0

clean_up() {
    CONTAINER="${IMAGE_NAME}"
    d4x_cleanup_image "${CONTAINER}"
    if [ $PID -ne 0 ]; then
      echo "killing netcat with pid $PID"
      kill $PID
      wait $PID
    fi
}
trap clean_up EXIT

build() {
    CONTAINER="${IMAGE_NAME}"
    docker build -t "$CONTAINER" .
}

test_bind() {
    CONTAINER="${IMAGE_NAME}"

    echo "binding port on the Mac"
    nc -l 8081 &
    PID=$!
    sleep 2
    echo "starting container"
    OUTPUT=$(docker run --rm -p 8081:8081 -t "$CONTAINER" echo hello 2>&1)
    if [ $? -eq 0 ]; then
        # This shouldn't be possible: the netcat should have the port
        echo "container was able to bind port successfully even though netcat is running"
        exit 1
    fi
    case "$OUTPUT" in
        *"errno 526"*)
            echo "error message contained Linux 9P errno 526"
            exit 1
            ;;
        *"port is already allocated"*)
            echo "error response looks OK"
            exit 0
            ;;
        *)
            echo "unknown error from docker"
            exit 1
            ;;
    esac
}

build

test_bind
