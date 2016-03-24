#!/usr/bin/env bash

function cleanup {
	docker-compose -f development.yml stop
	# if we're in CircleCI, we cannot remove any containers
	if [[ -z "${CIRCLECI}" ]]; then
		docker-compose -f development.yml rm -f
	fi
}

if [[ -z "${CIRCLECI}" ]]; then
	BUILDOPTS="--force-rm"
fi

set -e
set -x

cleanup

docker-compose -f development.yml build ${BUILDOPTS}
docker-compose -f development.yml up --abort-on-container-exit

trap cleanup SIGINT SIGTERM EXIT
