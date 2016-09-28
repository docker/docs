#!/bin/bash
set -e
TAG=${TAG:-latest}

IFS='//' read -a HOST <<< "$DOCKER_HOST"
ADDR="${HOST[2]}"
IFS=':' read -a ADDR_PARTS <<< "$ADDR"
IP="${ADDR_PARTS[0]}"

if [ -z "$IP" ]; then
    IP="127.0.0.1"
fi

# rethinkdb
echo "-> Running RethinkDB"
docker run -ti -d --restart=always --name orca-rethinkdb dockerorca/rethinkdb

# docker proxy
echo "-> Running Docker Proxy"
docker run \
    -ti \
    -d \
    --restart=always \
    --name orca-proxy \
    -v /var/run/docker.sock:/var/run/docker.sock \
    ehazlett/docker-proxy:latest

# swarm
echo "-> Running Docker Swarm"
docker run \
    -ti \
    -d \
    --restart=always \
    --name orca-swarm \
    --link orca-proxy:proxy \
    swarm:latest \
    m --host tcp://0.0.0.0:3375 proxy:2375

# generate certs
echo "-> Generating Certs"
docker run \
    --rm \
    -v /etc/certs:/certs \
    ehazlett/certm \
    -d /certs \
    bundle generate \
    --org orca \
    --host localhost \
    --host 127.0.0.1 \
    --host $IP \
    --overwrite

# copy ca cert
docker run \
    --rm \
    -v /etc/certs:/certs \
    alpine \
    sh -c "cat /certs/ca.pem" > ca.pem

# copy client cert
docker run \
    --rm \
    -v /etc/certs:/certs \
    alpine \
    sh -c "cat /certs/cert.pem" > cert.pem

# copy client cert key
docker run \
    --rm \
    -v /etc/certs:/certs \
    alpine \
    sh -c "cat /certs/key.pem" > key.pem

# generate public key for cert
docker run \
    --rm \
    -v /etc/certs:/certs \
    alpine \
    sh -c "apk update > /dev/null && apk add openssl > /dev/null && openssl x509 -in /certs/cert.pem -pubkey -noout" > cert.pub

# controller
echo "-> Running Controller"
docker run \
    -ti \
    -d \
    --restart=always \
    --name orca-controller \
    --link orca-rethinkdb:rethinkdb \
    --link orca-swarm:swarm \
    -v /etc/certs:/certs \
    -p 8080:8080 \
    dockerorca/orca:$TAG \
    -D server \
    -d tcp://swarm:3375 \
    --orca-tls-ca-cert /certs/ca.pem \
    --orca-tls-cert /certs/server.pem \
    --orca-tls-key /certs/server-key.pem

echo "Deploy complete.  You should be able to access Orca at https://$IP:8080"
