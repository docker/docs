# Remote logging configuration

To setup remote logging in UCP you must edit the API directly.
You'll need to run explicit curl commands described below.  This
assumes you've already set up your environment with a downloaded
bundle.

## Display the current settings
```sh
export UCP_URL="https://$(echo $DOCKER_HOST | cut -f3 -d/ )"
curl -s \
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    ${UCP_URL}/api/config/logging | jq "."
```

## Setup remote logging
```sh
curl -s \
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    -XPOST -d '{"host":"mylogger:514","protocol":"tcp","level":"INFO"}' \
    ${UCP_URL}/api/config/logging | jq "."
```

## Stopping remote logging

If you set the host to an empty string, remote logging will be disabled.

```sh
curl -s \
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    -XPOST -d '{"host":"","level":"DEBUG"}' \
    ${UCP_URL}/api/config/logging | jq "."
```

# Setting up an ELK stack

One popular logging stack is composed of Elasticsearch, Logstash and
Kibana.  The following example demonstrates how to set up an example
deployment which can be used for logging.  Once you have these containers
running, configure UCP to send logs to the logstash container.


```sh
docker volume create --name orca-elasticsearch-data

docker run -d \
    --name elasticsearch \
    -v orca-elasticsearch-data:/usr/share/elasticsearch/data \
    elasticsearch elasticsearch -Des.network.host=0.0.0.0

docker run -d \
    -p 514:514 \
    --name logstash \
    --link elasticsearch:es \
    logstash \
    sh -c "logstash -e 'input { syslog { } } output { stdout { } elasticsearch { hosts => [ \"es\" ] } } filter { json { source => \"message\" } }'"

docker run -d \
    --name kibana \
    --link elasticsearch:elasticsearch \
    -p 5601:5601 \
    kibana
```

You can then browse to port 5601 on the system running kibana and browse log/event entries.  You should specify the "time" field
for indexing.

Note: When deployed in production, you should secure kibana (not described in this doc)

## Example Searches

Here are a few examples demonstrating some ways to view the aggregated log data:

* `type:"api" AND (tags:"post" OR tags:"put" OR tags:"delete")` -- Show all the modifications on the system
* `username:"admin"` -- Show all access from a given user
* `type:"auth fail"` -- Show all authentication failures on the system
