---
description: Learn how to troubleshoot your DTR installation.
keywords: docker, registry, monitor, troubleshoot
title: Troubleshoot Docker Trusted Registry
---

This guide contains tips and tricks for troubleshooting DTR problems.

## Troubleshoot overlay networks

High availability in DTR depends on having overlay networking working in UCP.
One way to test if overlay networks are working correctly you can deploy
containers in different nodes, that are attached to the same overlay network
and see if they can ping one another.

Use SSH to log into a UCP node, and run:

```none
docker run -it --rm \
  --net dtr-ol --name overlay-test1 \
  --entrypoint sh docker/dtr
```

Then use SSH to log into another UCP node and run:

```none
docker run -it --rm \
  --net dtr-ol --name overlay-test2 \
  --entrypoint ping docker/dtr -c 3 overlay-test1
```

If the second command succeeds, it means that overlay networking is working
correctly.

You can run this test with any overlay network, and any Docker image that has
`sh` and `ping`.


## Access RethinkDB directly

DTR uses RethinkDB for persisting data and replicating it across replicas.
It might be helpful to connect directly to the RethinkDB instance running on a
DTR replica to check the DTR internal state.

Use SSH to log into a node that is running a DTR replica, and run the following
command, replacing `$REPLICA_ID` by the ID of the DTR replica running on that
node:

```none
docker run -it --rm \
  --net dtr-ol \
  -v dtr-ca-$REPLICA_ID:/ca dockerhubenterprise/rethinkcli:v2.2.0 \
  $REPLICA_ID
```

This starts an interactive prompt where you can run RethinkDB queries like:

```none
> r.db('dtr2').table('repositories')
```

[Learn more about RethinkDB queries](https://www.rethinkdb.com/docs/guide/javascript/).

## Recover from an unhealthy replica

When a DTR replica is unhealthy or down, the DTR web UI displays a warning:

```none
Warning: The following replicas are unhealthy: 59e4e9b0a254; Reasons: Replica reported health too long ago: 2017-02-18T01:11:20Z; Replicas 000000000000, 563f02aba617 are still healthy.
```

To fix this, you should remove the unhealthy replica from the DTR cluster,
and join a new one. Start by running:

```none
docker run -it --rm \
  {{ page.docker_image }} remove \
  --ucp-insecure-tls
```

And then:

```none
docker run -it --rm \
  {{ page.docker_image }} join \
  --ucp-node <ucp-node-name> \
  --ucp-insecure-tls
```
