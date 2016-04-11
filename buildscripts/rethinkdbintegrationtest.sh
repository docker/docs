#!/usr/bin/env bash

function cleanup {
	docker-compose -f development.rethink.yml stop
	# if we're in CircleCI, we cannot remove any containers
	if [[ -z "${CIRCLECI}" ]]; then
		docker-compose -f development.rethink.yml rm -f
	fi
}

if [[ -z "${CIRCLECI}" ]]; then
	BUILDOPTS="--force-rm"
fi

set -e
set -x

cleanup

docker-compose -f development.rethink.yml build ${BUILDOPTS} --pull | tee
docker-compose -f development.rethink.yml up --abort-on-container-exit

trap cleanup SIGINT SIGTERM EXIT
