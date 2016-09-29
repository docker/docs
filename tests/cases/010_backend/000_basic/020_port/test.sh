#!/bin/sh -x
# SUMMARY: Test the various command handling ips returns the correct ip
# LABELS:
# AUTHOR: Jean-Laurent de Morlhon <jeanlaurent@docker.com>

CONTAINER_NAME=port_test
IMAGE_NAME="${CONTAINER_NAME}_image"

set -e
. "${RT_PROJECT_ROOT}/_lib/lib.sh"

clean_up() {
  docker kill "${CONTAINER_NAME}" || true
  docker rm "${CONTAINER_NAME}" || true
  docker images -a | grep "${IMAGE_NAME}" | cut -f1 -d" " | xargs docker rmi || true
}

get_hello() {
  ip="${1}"
  [ "${ip}" = "0.0.0.0" ] && ip="localhost"
  echo "ip is ${ip}"
  hello=$(curl -s "${ip}":8081)
  if [[ "$hello" != 'hello' ]]; then
    echo "hello expected, got '${hello}'"
    exit -1
  fi
}

trap clean_up EXIT

docker build -t "$IMAGE_NAME" .

docker run -d --name "$CONTAINER_NAME" -p 8081:8080 -t "$IMAGE_NAME"
sleep 2 # give container some time to start

echo "ip from docker port"
ip_from_port="$(docker port "$CONTAINER_NAME" 8080 | cut -f 1 -d ":")"
get_hello "$ip_from_port"

echo "ip from docker inspect"
ip_from_inspect="$(docker inspect -f '{{((index .NetworkSettings.Ports "8080/tcp" 0).HostIp)}}' $CONTAINER_NAME)"
get_hello "$ip_from_inspect"

echo "ip from docker ps"
ip_from_ps="$(docker ps -l --format '{{.Ports}}' | cut -f 3 -d " "| cut -f 1 -d ":")"
get_hello "$ip_from_ps"

exit 0
