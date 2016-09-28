# Using Prometheus to instrument etcd

Set up the configuration (replace the IP address with your hosts address

(Ideally you'd set this up so that you have all the etcd IP's listed)

```
cat << EOF > /tmp/prometheus.yml
scrape_configs:
  - job_name: 'etcd'
    scrape_interval: 5s
    scheme: https
    static_configs:
      - targets: ['159.203.196.115:12379']
    tls_config:
      ca_file: /etc/docker/ssl/ca.pem
      cert_file: /etc/docker/ssl/cert.pem
      key_file: /etc/docker/ssl/key.pem
EOF
```

Then run prometheus on the local node with volume mounts so the certs pass through:
```
docker run --rm -it \
    -p 9090:9090 \
    -v ucp-kv-certs:/etc/docker/ssl:ro \
    -v /tmp/prometheus.yml:/etc/prometheus/prometheus.yml \
    prom/prometheus:latest
```


## Key change rate

If you want to diagnose possible high rate of churn of keys, a few things to look at:

Assuming you've sourced an admin bundle
```
export KV_URL="https://$(echo $DOCKER_HOST | cut -f3 -d/ | cut -f1 -d:):12379"
```

Now check out the overall churn rate in the etcd cluster:
curl -D - -s \
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    ${KV_URL}/v2/keys
sleep 10s
curl -D - -s \
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    ${KV_URL}/v2/keys
```

Then compare: `X-Etcd-Index` (key writes) and `X-Raft-Index` raft updates (used by snapshot algorithm)

If the key writes are high, you can try:
```
curl -s \     
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    "${KV_URL}/v2/keys/?recursive=true&sorted=true" | jq "." > t1
sleep 10s
curl -s \     
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    "${KV_URL}/v2/keys/?recursive=true&sorted=true" | jq "." > t2
```

Then something like the following (use your favoring diff'ing tool)
```
diffuse t1 t2
```

