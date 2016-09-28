#!/bin/sh

set -e

: "${NAMESPACE?NAMESPACE was not set}"
: "${DOCKER_STACK_BUNDLE?DOCKER_STACK_BUNDLE was not set}"

: "${CONTROLLER_HOST?CONTROLLER_HOST was not set}"
: "${CONTROLLER_CA_CERT?CONTROLLER_CA_CERT was not set}"
: "${SESSION_TOKEN?SESSION_TOKEN was not set}"

echo "${DOCKER_STACK_BUNDLE}" | base64 -d > bundle.dsb
echo "${CONTROLLER_CA_CERT}" | base64 -d > ca.pem

(mkdir -p ${HOME}/.docker/ && \
	echo -e "{\n\t\"HttpHeaders\": {\n\t\t\"Authorization\": \"Bearer ${SESSION_TOKEN}\"\n\t}\n}\n" > ${HOME}/.docker/config.json)

docker -H "${CONTROLLER_HOST}" --tlsverify --tlscacert ca.pem ${COMPOSE_OPTS} deploy -f bundle.dsb ${NAMESPACE}
