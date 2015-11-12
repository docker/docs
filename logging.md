# Remote logging configuration

To setup remote logging in Orca you must edit the KV store directly.
You'll need to run explicit curl commands described below.  This
assumes you've already set up your environment per the [KV Store
instructions](kv_store.md)

## Display the current settings
```sh
curl -s \
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    ${KV_URL}/v2/keys/orca/v1/config/logging | jq "."
```

## Setup remote logging
```sh
curl -s \
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    -XPUT -d value='{"host":"mylogger:514","protocol":"tcp","level":"INFO"}' \
    ${KV_URL}/v2/keys/orca/v1/config/logging | jq "."
```


## Stopping remote logging

If you're simply changing the target, use the "set" example above with a new host.  If you no longer want to send logging
to any remote syslogger, use the following:

```sh
curl -s \
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    -XDELETE ${KV_URL}/v2/keys/orca/v1/config/logging | jq "."
```

# Setting up an ELK stack

One popular logging stack is composed of Elasticsearch, Logstash and
Kibana.  The following example demonstrates how to set up an example
deployment which can be used for logging.  Once you have these containers
running, configure Orca to send logs to the logstash container.


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
    sh -c "logstash -e 'input { syslog { } } output { stdout { } elasticsearch { hosts => [ \"es\" ] } }'"

docker run -d \
    --name kibana \
    --link elasticsearch:elasticsearch \
    -p 5601:5601 \
    kibana
```

You can then browse to port 5601 on the system running kibana and browse log/event entries.

Note: When deployed in production, you should secure kibana (not described in this doc)
