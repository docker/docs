#!/usr/bin/env bash

set -e

RANDOMSTRING="$(cat /dev/urandom | env LC_CTYPE=C tr -dc 'a-zA-Z0-9' | fold -w 10 | head -n 1)"
HOST="${REMOTE_SERVER_URL:-https://localhost:4443}"

OPTS="-c cmd/notary/config.json -d /tmp/${RANDOMSTRING}"
if [[ -n "${DOCKER_HOST}" ]]; then
	if [[ "$(resolveip -s notary-server)" == *"${DOCKER_HOST}"* ]]; then
		echo "This test is going to fail since the client doesn't have a trusted CA root for ${HOST}"
		exit 1
	fi
	HOST="${REMOTE_SERVER_URL:-https://notary-server:4443}"
	OPTS="$OPTS -s ${HOST}"
fi

REPONAME="docker.com/notary/${RANDOMSTRING}"

export NOTARY_ROOT_PASSPHRASE=ponies
export NOTARY_TARGETS_PASSPHRASE=ponies
export NOTARY_SNAPSHOT_PASSPHRASE=ponies

echo "Notary Host: ${HOST}"
echo "Repo Name: ${REPONAME}"

echo

rm -rf "/tmp/${RANDOMSTRING}"

iter=0
until curl -k "${HOST}"
do
	((iter++))
	if (( iter > 30 )); then
		echo "notary service failed to come up within 30 seconds"
		exit 1;
	fi
	echo "waiting for notary service to come up."
	sleep 1
done

set -x

make client

bin/notary ${OPTS} init ${REPONAME}
bin/notary ${OPTS} delegation add ${REPONAME} targets/releases fixtures/secure.example.com.crt --all-paths
bin/notary ${OPTS} add ${REPONAME} readmetarget README.md
bin/notary ${OPTS} publish ${REPONAME}
bin/notary ${OPTS} delegation list ${REPONAME} | grep targets/releases
cat README.md | bin/notary ${OPTS} verify $REPONAME readmetarget > /dev/null
