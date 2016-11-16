---
description: Learn how to uninstall your Docker Trusted Registry installation.
keywords: docker, dtr, install, uninstall
redirect_from:
- /docker-trusted-registry/install/uninstall/
title: Uninstall Docker Trusted Registry
---

Use the `docker/dtr remove` command, to remove a DTR replica from a cluster.
To uninstall a DTR cluster you remove all DTR replicas one at a time.
The remove command:

* Removes the replica from the cluster,
* Stops and removes all DTR containers,
* Deletes all DTR volumes.

To see what options are available in the uninstall command, check the
[uninstall command reference](../reference/remove.md), or run:

```bash
$ docker run --rm -it docker/dtr remove --help
```

To remove a replica safely, you must tell the bootstrapper about one healthy
replica using the `--existing-replica-id` flag and the replica to remove with
the `--replica-id` flag. It uses the healthy replica to safely inform your DTR
cluster that the replica is about to be removed before it performs the actual
removal.

## Example

The following example illustrates how use the remove command interactively to
remove a DTR replica from a cluster with multiple replicas:

```bash
# Get the certificates used by UCP
$ curl https://$UCP_HOST/ca > ucp-ca.pem

$ docker run --rm -it docker/dtr remove --ucp-ca "$(cat ucp-ca.pem)"

existing-replica-id (ID of an existing replica in a cluster): 7ae3cb044b70
replica-id (Specify the replica Id. Must be unique per replica, leave blank for random): a701a510126c
ucp-username (Specify the UCP admin username): $UCP_ADMIN
ucp-password: $UCP_PASSWORD
ucp-url (Specify the UCP host using the host[:port] format): $UCP_HOST
```

Where:

* existing-replica-id: is the id of any healthy DTR replica of that cluster,
* replica-id: is the id of the DTR replica you want to remove,
* ucp-username and ucp-password: are the username and password of a UCP administrator.


Now you can confirm on Docker Universal Control Plane that the DTR replica
`a701a510126c` no longer exists.


## Where to go next

* [Install DTR](index.md)
* [Install DTR offline](install-dtr-offline.md)