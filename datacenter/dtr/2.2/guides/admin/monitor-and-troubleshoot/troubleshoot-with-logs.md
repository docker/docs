---
description: Learn how to troubleshoot your DTR installation.
keywords: docker, registry, monitor, troubleshoot
title: Troubleshoot Docker Trusted Registry
---

## Troubleshoot overlay networks

High availability in DTR depends on having overlay networking working in UCP.
To manually test that overlay networking is working in UCP run the following
commands on two different UCP machines.

```none
docker run -it --rm --net dtr-ol --name overlay-test1 --entrypoint sh docker/dtr
docker run -it --rm --net dtr-ol --name overlay-test2 --entrypoint ping docker/dtr -c 3 overlay-test1
```

You can create new overlay network for this test with `docker network create -d overaly network-name`.
You can also use any images that contain `sh` and `ping` for this test.

If the second command succeeds, overlay networking is working.

## Access rethinkdb directly

You can connect directly to the rethinkdb data store internal to DTR for
troubleshooting. To do this, you can use this custom rethinkdb debugging
tool. It connects to one of your rethinkdb servers as
indicated by `$REPLICA_ID` and presents you with an interactive prompt for
running rethinkdb queries. It must be run on the same machine as the replica
it's connecting to.

```none
$ docker run --rm -it --net dtr-ol -v dtr-ca-$REPLICA_ID:/ca dockerhubenterprise/rethinkcli:v2.2.0 $REPLICA_ID
```

You can use [javascript
syntax](https://www.rethinkdb.com/docs/guide/javascript/) to execute rethinkdb queries like so:

```none
> r.db('dtr2').table('repositories')
```

## Recover from a lost replica

When one of DTR's replicas is lost, the UI will start showing a warning that
looks something like the following:

```none
Warning: The following replicas are unhealthy: 59e4e9b0a254; Reasons: Replica reported health too long ago: 2017-02-18T01:11:20Z; Replicas 000000000000, 563f02aba617 are still healthy.
```

To remedy this situation, you need to use the `remove` command to tell
the cluster that the lost replica should be treated as permanently removed.
After that you can use the `join` command to grow your cluster back to the
desired number of replicas. In this example you would run the following
commands (and follow the prompts for the UCP connection parameters):

```none
$ docker run --rm -it docker/dtr remove \
  --ucp-insecure-tls --replica-id 59e4e9b0a254 --existing-replica-id 000000000000
$ docker run --rm -it docker/dtr join \
  --ucp-insecure-tls --existing-replica-id 000000000000
```
