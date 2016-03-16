#!/usr/bin/env bash



set -e
set -x

function finish {
	docker-compose stop
	if [[ -n "${COMPOSE_PID}" ]]; then
		kill "${COMPOSE_PID}" || true
	fi
	if [[ -n "${TEST_PID}" ]]; then
		kill "${TEST_PID}" || true
	fi
}

docker-compose stop

# if we're in CircleCI, we cannot remove any containers
if [[ -z "${CIRCLECI}" ]]; then
	docker-compose rm -f
fi

docker-compose build
docker-compose up --abort-on-container-exit >> /dev/null &
COMPOSE_PID=$!


.build/testclient.sh &
TEST_PID=$!

set +x

wait ${TEST_PID}

(docker-compose logs &)

trap finish SIGINT SIGTERM EXIT
