---
description: Learn how to troubleshoot your DTR installation.
keywords: docker, registry, monitor, troubleshoot
title: Troubleshoot Docker Trusted Registry
---

## Troubleshooting overlay networks

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

## Accessing rethinkdb directly

To perform operations against rethinkdb directly, one may use the following
custom rethinkdb debugging tool. It connects to one of your rethinkdb servers as
indicated by `$REPLICA_ID` and presents you with an interactive prompt for
running rethinkdb queries. It must be run on the same machine as the replica
it's connecting to.

```none
$ docker run --rm -it --net dtr-ol -v dtr-ca-$REPLICA_ID:/ca dockerhubenterprise/rethinkcli:v2.2.0 $REPLICA_ID
```

You can use javascript syntax to execute rethinkdb queries like so:

```none
> r.db('dtr2').table('repositories')
```
