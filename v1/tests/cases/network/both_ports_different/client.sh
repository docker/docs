#!/bin/bash

CONTAINER=$1
EXPECTED=$2

ID=$(docker ps | grep both_ports_different | cut -f 1 -d " ")
IP=$(docker port ${ID} | cut -f 3 -d " "| cut -f 1 -d ":")
PORT=$(docker port ${ID} | cut -f 3 -d " "| cut -f 2 -d ":")

RETRIES=60
while [ ${RETRIES} -ne 0 ]; do
  RESULT=$(nc ${IP} $PORT)
  if [ ! -z "${RESULT}" ]; then
    if [ "${EXPECTED}" = "${RESULT}" ]; then
      exit 0
    fi
    echo Unexpected response from server, expected ${EXPECTED}, got ${RESULT}
    exit 1
  fi
  sleep 1
  echo -n .
  RETRIES=$(( ${RETRIES} - 1 ))
done
echo Timed out waiting for response from server
exit 1
