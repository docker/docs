---
description: Learn how to troubleshoot your Docker Universal Control Plane cluster.
keywords: ectd, rethinkdb, key, value, store, database, ucp
title: Troubleshoot cluster configurations
---

This article describes how to troubleshoot a UCP cluster. When experiencing
unexpected behavior, we highly recommend that you follow these steps before
attempting any disruptive operations:

1. Obtain a [support dump](../../get-support.md) of the cluster, so that the
   original failure you are experiencing is isolated before any subsequent
   operations are performed.
2. Make sure that the failure you are experiencing is not an expected
   [UCP node state](./troubleshoot-node-messages.md).
3. Set the cluster-wide log severity level to `DEBUG`. This is an operation that
   will temporarily restart all UCP system components and will introduce a small
   downtime window to the UCP interface. You can configure this by logging in to
   UCP as an Admin user, navigating to the **Admin Settings** section, selecting
   the **Logs** tab, and setting the **Log Severity Level** to `DEBUG`. After
   that, wait until the UCP web UI is responsive again.
4. Obtain a new [support dump](./../../get-support.md) of the cluster which will
   now include debug-level logs.
5. Depending on the problem you are experiencing, it's more likely that you'll
   find related messages in the logs of specific components on manager nodes: 

	* If the problem occurs after a node was added or removed, review the logs
	of the `ucp-reconcile` container.
	* If the problem occurs in the normal state of the system, review the logs
	of the `ucp-controller` container.
	* If you are able to visit the UCP web UI but unable to log in, review the
	logs of the `ucp-auth-api` and `ucp-auth-store` containers.


## Troubleshooting individual UCP components

This section describes how to troubleshoot individual UCP components. This is
typically not required, as UCP uses a node reconciler to bring nodes to a
healthy state. In most cases, if a single UCP component is persistently in a
failed state, you should be able to restore the cluster to a healthy state by
[removing the node from the cluster and joining
again](../configure/scale-your-cluster.md). However, you may wish to inspect the
logs of 

As an admin user, you're able to view the logs of individual UCP containers.
You can [learn more about the UCP architecture](../../architecture.md)

It's normal for the `ucp-reconcile` container to be in a stopped state. This
container is only started when the `ucp-agent` detects that a node needs to transition
to a different state, and it is responsible for creating and removing
containers, issuing certificates and pulling missing images. 

Docker UCP persists configuration data on an [etcd](https://coreos.com/etcd/)
key-value store and [RethinkDB](https://rethinkdb.com/) database that are
replicated on all manager nodes of the UCP cluster. These data stores are for
internal use only, and should not be used by other applications. 

### Troubleshooting the etcd Key-Value Store 

#### Using the RESTful HTTP API
In this example we'll be using `curl` for making requests to the key-value
store REST API, and `jq` to process the responses.

You can install these tools on a Ubuntu distribution by running:

```bash
$ sudo apt-get update && apt-get install curl jq
```

1. Use a client bundle to authenticate your requests.
[Learn more](../../user/access-ucp/cli-based-access.md).

2. Use the REST API to access the cluster configurations.

```bash
# $DOCKER_HOST and $DOCKER_CERT_PATH are set when using the client bundle
$ export KV_URL="https://$(echo $DOCKER_HOST | cut -f3 -d/ | cut -f1 -d:):12379"

$ curl -s \
    --cert ${DOCKER_CERT_PATH}/cert.pem \
    --key ${DOCKER_CERT_PATH}/key.pem \
    --cacert ${DOCKER_CERT_PATH}/ca.pem \
    ${KV_URL}/v2/keys | jq "."
```

To learn more about the key-value store REST API check the
[etcd official documentation](https://coreos.com/etcd/docs/latest/).

#### Using a CLI client

The containers running the key-value store, include `etcdctl`, a command line
client for etcd. You can run it using the `docker exec` command.

The examples below assume you are logged in with ssh into a UCP manager node.

```bash
$ docker exec -it ucp-kv etcdctl \
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

On failure the command exits with an error code, and no output.

To learn more about the `etcdctl` utility, check the
[etcd official documentation](https://coreos.com/etcd/docs/latest/).

## RethinkDB Database

User and organization data for Docker Datacenter is stored in a RetinkDB
database which is replicated across all manager nodes in the UCP cluster.

Replication and failover of this database is typically handled automatically by
UCP's own configuration management processes, but detailed database status and
manual reconfiguration of database replication is available through a command
line tool available as part of UCP.

The examples below assume you are logged in with ssh into a UCP manager node.

### Check the status of the database

```bash
{% raw %}
# NODE_ADDRESS will be the IP address of this Docker Swarm manager node
NODE_ADDRESS=$(docker info --format '{{.Swarm.NodeAddr}}')
# VERSION will be your most recent version of the docker/ucp-auth image
VERSION=$(docker image ls --format '{{.Tag}}' docker/ucp-auth | head -n 1)
# This command will output detailed status of all servers and database tables
# in the RethinkDB cluster.
docker run --rm -v ucp-auth-store-certs:/tls docker/ucp-auth:${VERSION} --db-addr=${NODE_ADDRESS}:12383 db-status
{% endraw %}
```

### Manually reconfigure database replication

```bash
{% raw %}
# NODE_ADDRESS will be the IP address of this Docker Swarm manager node
NODE_ADDRESS=$(docker info --format '{{.Swarm.NodeAddr}}')
# NUM_MANAGERS will be the current number of manager nodes in the cluster
NUM_MANAGERS=$(docker node ls --filter role=manager -q | wc -l)
# VERSION will be your most recent version of the docker/ucp-auth image
VERSION=$(docker image ls --format '{{.Tag}}' docker/ucp-auth | head -n 1)
# This reconfigure-db command will repair the RethinkDB cluster to have a
# number of replicas equal to the number of manager nodes in the cluster.
docker run --rm -v ucp-auth-store-certs:/tls docker/ucp-auth:${VERSION} --db-addr=${NODE_ADDRESS}:12383 --debug reconfigure-db --num-replicas ${NUM_MANAGERS} --emergency-repair
{% endraw %}
```

## Where to go next

* [Get support](../../get-support.md)
