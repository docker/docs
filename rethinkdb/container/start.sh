#!/bin/sh

set -e

echo Starting rethinkdb for replica: ${DTR_REPLICA_ID}
DOMAIN=dtr-rethinkdb-${DTR_REPLICA_ID}.dtr-br

ADMIN_IP=
until [ ${ADMIN_IP} ]
do
    echo Trying to resolve own ip address ${DOMAIN}...
    ADMIN_IP=$(getent hosts ${DOMAIN} | awk '{print $1}')
done

echo Admin interface host: ${ADMIN_IP}

# this prints so we can see logs, then we capture its output
arggen > arggen.out

ARGS=$(cat arggen.out)

CMD="/usr/local/bin/rethinkdb --bind-http ${ADMIN_IP} ${ARGS}"
echo Starting with command: ${CMD}
eval ${CMD}
