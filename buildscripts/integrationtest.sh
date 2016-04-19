#!/usr/bin/env bash

composeFile="$1"

function cleanup {
	docker-compose -f $composeFile stop
	# if we're in CircleCI, we cannot remove any containers
	if [[ -z "${CIRCLECI}" ]]; then
		docker-compose -f $composeFile rm -f
	fi
}

function cleanupAndExit {
    cleanup
    # Check for existence of SUCCESS
    ls test_output/SUCCESS
    exitCode=$?
    # Clean up test_output dir (if not in CircleCI) and exit
    if [[ -z "${CIRCLECI}" ]]; then
        rm -rf test_output
    fi
    exit $exitCode
}

if [[ -z "${CIRCLECI}" ]]; then
	BUILDOPTS="--force-rm"
fi

set -e
set -x

cleanup

docker-compose -f $composeFile config
docker-compose -f $composeFile build ${BUILDOPTS} --pull | tee
docker-compose -f $composeFile up --abort-on-container-exit

trap cleanupAndExit SIGINT SIGTERM EXIT
