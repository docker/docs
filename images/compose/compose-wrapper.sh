#!/bin/sh

set -e

: "${PROJECT_NAME?PROJECT_NAME was not set}"
: "${CONTROLLER_HOST?CONTROLLER_HOST was not set}"
: "${CONTROLLER_CA_CERT?CONTROLLER_CA_CERT was not set}"
: "${DOCKER_COMPOSE_YML?DOCKER_COMPOSE_YML was not set}"
: "${SESSION_TOKEN?SESSION_TOKEN was not set}"

echo "${DOCKER_COMPOSE_YML}" | base64 -d > docker-compose.yml
echo "${CONTROLLER_CA_CERT}" | base64 -d > ca.pem

(mkdir -p $HOME/.docker/ && \
	echo -e "{\n\t\"HttpHeaders\": {\n\t\t\"Authorization\": \"Bearer ${SESSION_TOKEN}\"\n\t}\n}\n" > $HOME/.docker/config.json)

docker-compose -H "${CONTROLLER_HOST}" --tlsverify --tlscacert ca.pem --project-name="${PROJECT_NAME}" ${COMPOSE_OPTS} up -d
