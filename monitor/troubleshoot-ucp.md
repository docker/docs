<!--[metadata]>
+++
aliases = ["/ucp/kv_store/"]
title ="Troubleshoot your cluster"
keywords= ["ectd, key, value, store, ucp"]
description="Docker Universal Control Plane"
[menu.main]
parent="mn_monitor_ucp"
+++
<![end-metadata]-->

# Troubleshoot your cluster

In this release, UCP leverages the [etcd](https://github.com/coreos/etcd/) KV
store internally for node discovery and high availability. This use is specific
to UCP. The services you deploy on UCP can use whichever key-store is
appropriate for the service.

Under normal circumstances, you should not have to access the KV store
directly.  To mitigate unforeseen problems or change advanced settings,
you may be required by Docker support or your SE to change configuration
values or data in the store.

The following example demonstrates basic `curl` usage assuming you
have set up your environment with the downloaded ucp bundle.

The example below uses the [jq](https://stedolan.github.io/jq/) tool to
pretty print the resulting json.  This can be omitted for raw json output.


```bash
export KV_URL="https://$(echo $DOCKER_HOST | cut -f3 -d/ | cut -f1 -d:):12379"

curl -s \
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    ${KV_URL}/v2/keys | jq "."
```


You can browse the UCP keys under `/v2/keys/ucp/` and swarm under
`/v2/keys/swarm` as well as modify by `POST`ing updated values to
workaround problems.  Further documentation for the etcd API is available
at https://github.com/coreos/etcd/blob/master/Documentation/api.md


### Troubleshooting with etcdctl

The `ucp-kv` container(s) running on the primary controller (and replicas in an
HA configuration) contain the `etcdctl` binary, which can be accessed using
`docker exec`.  The examples (and their output) using the tool to perform
various tasks on the `etcd` cluster.

These commands assume you are running directly against the Docker Engine in
question.  If you are running these commands through UCP, you should specify the
node specific container name.

Check the health of the etcd cluster (on failure it will exit with an error code, and no output)

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

List the current members of the cluster

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

Remove a failed member (use the list above first to get the ID)

```bash
docker exec -it ucp-kv etcdctl \
        --endpoint https://127.0.0.1:2379 \
        --ca-file /etc/docker/ssl/ca.pem \
        --cert-file /etc/docker/ssl/cert.pem \
        --key-file /etc/docker/ssl/key.pem \
        member remove c5a24cfdb4263e72

Removed member c5a24cfdb4263e72 from cluster
```

Show the current value of a key

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


### Learn about the certificates

The store is configured with mutual TLS to prevent unauthorized access.

All components in the system that require access to the KV store use
client certificates signed by the Swarm Root CA.  As admin account
certificates are also signed by this Swarm Root CA, administrators can
access the KV store using `curl` or other tools, provided the admin's
certificate is used as the client certificate.

## Where to go next

* [Monitor your cluster](monitor-ucp.md)
* [Get support](../support.md)
