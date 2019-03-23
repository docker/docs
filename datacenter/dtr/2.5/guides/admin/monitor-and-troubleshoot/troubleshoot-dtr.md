---
title: Troubleshoot Docker Trusted Registry
description: Learn how to troubleshoot your DTR installation.
keywords: registry, monitor, troubleshoot
redirect_from: /ee/dtr/admin/monitor-and-troubleshoot/troubleshoot-with-logs
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
  --entrypoint sh {{ page.dtr_org }}/{{ page.dtr_repo }}
```

Then use SSH to log into another node and run:

```bash
docker run -it --rm \
  --net dtr-ol --name overlay-test2 \
  --entrypoint ping {{ page.dtr_org }}/{{ page.dtr_repo }} -c 3 overlay-test1
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

### via RethinkCLI

As of v2.5.5, the [RethinkCLI has been removed](/ee/dtr/release-notes/#255) from the RethinkDB image along with other unused components. You can now run RethinkCLI from a separate image in the `dockerhubenterprise` organization. Note that the commands below are using separate tags for non-interactive and interactive modes.

#### Non-interactive

Use SSH to log into a node that is running a DTR replica, and run the following: 

{% raw %}
```bash
# List problems in the cluster detected by the current node.
REPLICA_ID=$(docker container ls --filter=name=dtr-rethink --format '{{.Names}}' | cut -d'/' -f2 | cut -d'-' -f3 | head -n 1) && echo 'r.db("rethinkdb").table("current_issues")' | docker run --rm -i --net dtr-ol -v "dtr-ca-${REPLICA_ID}:/ca" -e DTR_REPLICA_ID=$REPLICA_ID dockerhubenterprise/rethinkcli:v2.2.0-ni non-interactive
```
{% endraw %}

On a healthy cluster the output will be `[]`.

#### Interactive

Starting in DTR 2.5.5, you can run RethinkCLI from a separate image. First, set an environment variable for your DTR replica ID:

{% raw %}
```bash
REPLICA_ID=$(docker inspect -f '{{.Name}}' $(docker ps -q -f name=dtr-rethink) | cut -f 3 -d '-')
```
{% endraw %}

RethinkDB stores data in different databases that contain multiple tables. Run the following command to get into interactive mode
and query the contents of the DB:

```bash
docker run -it --rm --net dtr-ol -v dtr-ca-$REPLICA_ID:/ca dockerhubenterprise/rethinkcli:v2.3.0 $REPLICA_ID
```

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
[ 'blob_links',
  'blobs',
  'client_tokens',
  'content_caches',
  'events',
  'layer_vuln_overrides',
  'manifests',
  'metrics',
  'namespace_team_access',
  'poll_mirroring_policies',
  'promotion_policies',
  'properties',
  'pruning_policies',
  'push_mirroring_policies',
  'repositories',
  'repository_team_access',
  'scanned_images',
  'scanned_layers',
  'tags',
  'user_settings',
  'webhooks' ]

# List the entries in the repositories table
> r.db('dtr2').table('repositories')
[ { enableManifestLists: false,
    id: 'ac9614a8-36f4-4933-91fa-3ffed2bd259b',
    immutableTags: false,
    name: 'test-repo-1',
    namespaceAccountID: 'fc3b4aec-74a3-4ba2-8e62-daed0d1f7481',
    namespaceName: 'admin',
    pk: '3a4a79476d76698255ab505fb77c043655c599d1f5b985f859958ab72a4099d6',
    pulls: 0,
    pushes: 0,
    scanOnPush: false,
    tagLimit: 0,
    visibility: 'public' },
  { enableManifestLists: false,
    id: '9f43f029-9683-459f-97d9-665ab3ac1fda',
    immutableTags: false,
    longDescription: '',
    name: 'testing',
    namespaceAccountID: 'fc3b4aec-74a3-4ba2-8e62-daed0d1f7481',
    namespaceName: 'admin',
    pk: '6dd09ac485749619becaff1c17702ada23568ebe0a40bb74a330d058a757e0be',
    pulls: 0,
    pushes: 0,
    scanOnPush: false,
    shortDescription: '',
    tagLimit: 1,
    visibility: 'public' } ]
```

Individual DBs and tables are a private implementation detail and may change in DTR
from version to version, but you can always use `dbList()` and `tableList()` to explore
the contents and data structure.

[Learn more about RethinkDB queries](https://www.rethinkdb.com/docs/guide/javascript/).

### via API

To check on the overall status of your DTR cluster without interacting with RethinkCLI, run the following API request:

```bash
curl -u admin:$TOKEN -X GET "https://<dtr-url>/api/v0/meta/cluster_status" -H "accept: application/json"
```

#### Example API Response
```none
{
  "rethink_system_tables": {
    "cluster_config": [
      {
        "heartbeat_timeout_secs": 10,
        "id": "heartbeat"
      }
    ],
    "current_issues": [],
    "db_config": [
      {
        "id": "339de11f-b0c2-4112-83ac-520cab68d89c",
        "name": "notaryserver"
      },
      {
        "id": "aa2e893f-a69a-463d-88c1-8102aafebebc",
        "name": "dtr2"
      },
      {
        "id": "bdf14a41-9c31-4526-8436-ab0fed00c2fd",
        "name": "jobrunner"
      },
      {
        "id": "f94f0e35-b7b1-4a2f-82be-1bdacca75039",
        "name": "notarysigner"
      }
    ],
    "server_status": [
      {
        "id": "9c41fbc6-bcf2-4fad-8960-d117f2fdb06a",
        "name": "dtr_rethinkdb_5eb9459a7832",
        "network": {
          "canonical_addresses": [
            {
              "host": "dtr-rethinkdb-5eb9459a7832.dtr-ol",
              "port": 29015
            }
          ],
          "cluster_port": 29015,
          "connected_to": {
            "dtr_rethinkdb_56b65e8c1404": true
          },
          "hostname": "9e83e4fee173",
          "http_admin_port": "<no http admin>",
          "reql_port": 28015,
          "time_connected": "2019-02-15T00:19:22.035Z"
        },
       }
     ...
    ]
  }
}

```

## Recover from an unhealthy replica

When a DTR replica is unhealthy or down, the DTR web UI displays a warning:

```none
Warning: The following replicas are unhealthy: 59e4e9b0a254; Reasons: Replica reported health too long ago: 2017-02-18T01:11:20Z; Replicas 000000000000, 563f02aba617 are still healthy.
```

To fix this, you should remove the unhealthy replica from the DTR cluster,
and join a new one. Start by running:

```bash
docker run -it --rm \
  {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version }} remove \
  --ucp-insecure-tls
```

And then:

```bash
docker run -it --rm \
  {{ page.dtr_org }}/{{ page.dtr_repo }}:{{ page.dtr_version }} join \
  --ucp-node <ucp-node-name> \
  --ucp-insecure-tls
```
