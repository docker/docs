#!/bin/sh
# SUMMARY: Backend regression tests
# NAME: pinata
# LABELS: !benchmarks

set -e # Exit on error
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

group_init() {
    dockerpath=$(which docker)
    echo "Using docker from: $dockerpath"

    if ! d4x_app_running; then
        echoerr "Docker is supposed to be running here"
        return 1
    fi

    return 0
}

group_deinit() {
    return 0
}

CMD=$1
case $CMD in
init)
    group_init
    res=$?
    ;;
deinit)
    group_deinit
    res=$?
    ;;
*)
    res=1
    ;;
esac

echoerr "About to return"
exit $res
