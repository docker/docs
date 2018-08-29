---
description: Learn how to troubleshoot your DTR installation.
keywords: docker, registry, monitor, troubleshoot
title: Troubleshoot Docker Trusted Registry
---

This guide contains tips and tricks for troubleshooting DTR problems.

## Troubleshoot overlay networks

High availability in DTR depends on swarm overlay networking.  One way to test
if overlay networks are working correctly is to deploy containers to the same
overlay network on different nodes and see if they can ping one another.

Use SSH to log into a node and run:

```bash
docker run -it --rm \
  --net dtr-ol --name overlay-test1 \
  --entrypoint sh docker/dtr
```

Then use SSH to log into another node and run:

```bash
docker run -it --rm \
  --net dtr-ol --name overlay-test2 \
  --entrypoint ping docker/dtr -c 3 overlay-test1
```

If the second command succeeds, it indicates overlay networking is working
correctly between those nodes.

You can run this test with any attachable overlay network and any Docker image
that has `sh` and `ping`.


## Access RethinkDB directly

DTR uses RethinkDB for persisting data and replicating it across replicas.
It might be helpful to connect directly to the RethinkDB instance running on a
DTR replica to check the DTR internal state.

> **Warning**: Modifying RethinkDB directly is not supported and may cause
> problems.
{: .warning }

Use SSH to log into a node that is running a DTR replica, and run the following
commands:

{% raw %}
```bash
# DTR_REPLICA_ID will be the replica ID for the current node.
DTR_REPLICA_ID=$(docker ps -lf name='^/dtr-rethinkdb-.{12}$' --format '{{.Names}}' | cut -d- -f3)
# List problems in the cluster detected by the current node.
echo 'r.db("rethinkdb").table("current_issues")' | \
  docker run -i --rm \
    --net dtr-ol \
    -e DTR_REPLICA_ID=${DTR_REPLICA_ID} \
    -v dtr-ca-$DTR_REPLICA_ID:/ca \
    dockerhubenterprise/rethinkcli:v2.2.0-ni non-interactive; \
    echo
```
{% endraw %}

On a healthy cluster the output will be `[]`.

RethinkDB stores data in different databases that contain multiple tables. This
container can also be used to connect to the local DTR replica and
interactively query the contents of the DB.

{% raw %}
```bash
# DTR_REPLICA_ID will be the replica ID for the current node.
DTR_REPLICA_ID=$(docker ps -lf name='^/dtr-rethinkdb-.{12}$' --format '{{.Names}}' | cut -d- -f3)
# This command will start a RethinkDB client attached to the database
# on the current node.
docker run -it --rm \
  --net dtr-ol \
  -e DTR_REPLICA_ID=${DTR_REPLICA_ID} \
  -v dtr-ca-$DTR_REPLICA_ID:/ca \
  dockerhubenterprise/rethinkcli:v2.2.0-ni
```
{% endraw %}


```none
# List problems in the cluster detected by the current node.
> r.db("rethinkdb").table("current_issues")
[]

# List all the DBs in RethinkDB
> r.dbList()
[ 'dtr2',
  'jobrunner',
  'notaryserver',
  'notarysigner',
  'rethinkdb' ]

# List the tables in the dtr2 db
> r.db('dtr2').tableList()
[ 'client_tokens',
  'events',
  'manifests',
  'namespace_team_access',
  'properties',
  'repositories',
  'repository_team_access',
  'tags' ]

# List the entries in the repositories table
> r.db('dtr2').table('repositories')
[ { id: '19f1240a-08d8-4979-a898-6b0b5b2338d8',
    name: 'my-test-repo',
    namespaceAccountID: '924bf131-6213-43fa-a5ed-d73c7ccf392e',
    pk: 'cf5e8bf1197e281c747f27e203e42e22721d5c0870b06dfb1060ad0970e99ada',
    visibility: 'public' },
...
```

Individual DBs and tables are a private implementation detail and may change in DTR
from version to version, but you can always use `dbList()` and `tableList()` to explore
the contents and data structure.

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
