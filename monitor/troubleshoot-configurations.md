<!--[metadata]>
+++
aliases = ["/ucp/kv_store/"]
title ="Troubleshoot cluster configurations"
keywords= ["ectd, key, value, store, ucp"]
description="Learn how to troubleshoot your Docker Universal Control Plane cluster."
[menu.main]
parent="mn_monitor_ucp"
weight=20
+++
<![end-metadata]-->

# Troubleshoot cluster configurations

Docker UCP persists configuration data on an [etcd](https://coreos.com/etcd/)
key-value store. This key-value store is replicated on all controller nodes of
the UCP cluster. The key-value store is for internal use only, and should not
be used by other applications.

This article shows how you can access the key-value store, for
troubleshooting configuration problems in your cluster.

## Using the REST API

In this example we'll be using `curl` for making requests to the key-value
store REST API, and `jq` to process the responses.

You can install these tools on a Ubuntu distribution by running:

```bash
$ sudo apt-get update && apt-get install curl jq
```

To access the cluster configurations, run:

```bash
export KV_URL="https://$(echo $DOCKER_HOST | cut -f3 -d/ | cut -f1 -d:):12379"

curl -s \
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    ${KV_URL}/v2/keys | jq "."
```

To learn more about the key-value store API, check the
[etcd official documentation](https://coreos.com/etcd/docs/latest/api.html).


## Using a CLI client

The containers running the key-value store, include `etcdctl`, a command line
client for etcd. You can run it using the `docker exec` command.

The example below assumes you have the Docker CLI client pointing to the Docker
Engine of a UCP controller. If you are running the example below through UCP,
you should specify the node-specific container name.

These commands assume you are running directly against the Docker Engine in
question.  If you are running these commands through UCP, you should specify the
node specific container name.

Check the health of the etcd cluster. On failure the command exits with an
error code, and no output:

```bash
docker exec -it ucp-kv etcdctl \
        --endpoint https://127.0.0.1:2379 \
        --ca-file /etc/docker/ssl/ca.pem \
        --cert-file /etc/docker/ssl/cert.pem \
        --key-file /etc/docker/ssl/key.pem \
        cluster-health

member 16c9ae1872e8b1f0 is healthy: got healthy result from https://192.168.122.64:12379
member c5a24cfdb4263e72 is healthy: got healthy result from https://192.168.122.196:12379
member ca3c1bb18f1b30bf is healthy: got healthy result from https://192.168.122.223:12379
cluster is healthy
```

List the current members of the cluster:

```bash
docker exec -it ucp-kv etcdctl \
        --endpoint https://127.0.0.1:2379 \
        --ca-file /etc/docker/ssl/ca.pem \
        --cert-file /etc/docker/ssl/cert.pem \
        --key-file /etc/docker/ssl/key.pem \
        member list

16c9ae1872e8b1f0: name=orca-kv-192.168.122.64 peerURLs=https://192.168.122.64:12380 clientURLs=https://192.168.122.64:12379
c5a24cfdb4263e72: name=orca-kv-192.168.122.196 peerURLs=https://192.168.122.196:12380 clientURLs=https://192.168.122.196:12379
ca3c1bb18f1b30bf: name=orca-kv-192.168.122.223 peerURLs=https://192.168.122.223:12380 clientURLs=https://192.168.122.223:12379
```

Remove a failed member:

```bash
docker exec -it ucp-kv etcdctl \
        --endpoint https://127.0.0.1:2379 \
        --ca-file /etc/docker/ssl/ca.pem \
        --cert-file /etc/docker/ssl/cert.pem \
        --key-file /etc/docker/ssl/key.pem \
        member remove c5a24cfdb4263e72

Removed member c5a24cfdb4263e72 from cluster
```

Show the current value of a key:

```bash
docker exec -it ucp-kv etcdctl \
        --endpoint https://127.0.0.1:2379 \
        --ca-file /etc/docker/ssl/ca.pem \
        --cert-file /etc/docker/ssl/cert.pem \
        --key-file /etc/docker/ssl/key.pem \
        ls /docker/nodes

/docker/nodes/192.168.122.196:12376
/docker/nodes/192.168.122.64:12376
/docker/nodes/192.168.122.223:12376
```


## Where to go next

* [Monitor your cluster](monitor-ucp.md)
* [Get support](../support.md)
