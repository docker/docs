#!/usr/bin/env bash

function cleanup {
	docker-compose -f development.yml stop
	# if we're in CircleCI, we cannot remove any containers
	if [[ -z "${CIRCLECI}" ]]; then
		docker-compose -f development.yml rm -f
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

docker-compose -f development.yml config
docker-compose -f development.yml build ${BUILDOPTS} --pull | tee
docker-compose -f development.yml up --abort-on-container-exit

trap cleanupAndExit SIGINT SIGTERM EXIT
