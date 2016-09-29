#!/bin/sh

MODE=$1
FQDN=$2
RECORD_TYPE=$3
DNS_REMOTE=$4
EXPECTED=$5

RETRIES=3
while [ ${RETRIES} -gt 0 ]; do
  RESULT=$(drill "${MODE}" "${FQDN}" "@${DNS_REMOTE}" "${RECORD_TYPE}")
  if [ ! -z "${RESULT}" ]; then
    if echo "$RESULT" | grep "${EXPECTED}" ; then
      exit 0
    fi
    echo "Unexpected response from server, expected ${EXPECTED}, got ${RESULT}"
    exit 1
  fi
  sleep 1
  echo -n .
  RETRIES=$(( RETRIES - 1 ))
done
echo Timed out waiting for response from server
exit 1
