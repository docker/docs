#!/usr/bin/env bash

set -e

make client

set +e

RANDOMSTRING="$(cat /dev/urandom | env LC_CTYPE=C tr -dc 'a-zA-Z0-9' | fold -w 10 | head -n 1)"
HOST="${REMOTE_SERVER_URL:-https://notary-server:4443}"

REPONAME="docker.com/notary/${RANDOMSTRING}"

OPTS="-c cmd/notary/config.json -d /tmp/${RANDOMSTRING}"

export NOTARY_ROOT_PASSPHRASE=ponies
export NOTARY_TARGETS_PASSPHRASE=ponies
export NOTARY_SNAPSHOT_PASSPHRASE=ponies

echo "Notary Host: ${HOST}"
echo "Repo Name: ${REPONAME}"

echo

rm -rf "/tmp/${RANDOMSTRING}"

iter=0
until (curl -s -S -k "${HOST}")
do
	((iter++))
	if (( iter > 30 )); then
		echo "notary service failed to come up within 30 seconds"
		exit 1;
	fi
	echo "waiting for notary service to come up."
	sleep 1
done

set -e
set -x

bin/notary ${OPTS} init ${REPONAME}
bin/notary ${OPTS} delegation add ${REPONAME} targets/releases fixtures/secure.example.com.crt --all-paths
bin/notary ${OPTS} add ${REPONAME} readmetarget README.md
bin/notary ${OPTS} publish ${REPONAME}
bin/notary ${OPTS} delegation list ${REPONAME} | grep targets/releases
cat README.md | bin/notary ${OPTS} verify $REPONAME readmetarget > /test_output/SUCCESS

# Make this file accessible for CI
chmod -R 777 /test_output
