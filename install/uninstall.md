<!--[metadata]>
+++
title = "Uninstall"
description = "Learn how to uninstall your Docker Trusted Registry installation."
keywords = ["docker, dtr, install, uninstall"]
[menu.main]
parent="workw_dtr_install"
identifier="dtr_uninstall"
weight=50
+++
<![end-metadata]-->

# Uninstall Docker Trusted Registry

Use the `docker/dtr remove` command, to remove a DTR replica from a cluster.
To uninstall a DTR cluster you remove all DTR replicas one at a time.
The remove command:

* Removes the replica from the cluster,
* Stops and removes all DTR containers,
* Deletes all DTR volumes.

To see what options are available in the uninstall command, check the
[uninstall command reference](../reference/remove.md), or run:

```bash
$ docker run --rm -it docker/dtr uninstall --help
```

## Example

The following example illustrates how use the remove command interactively to
remove a DTR replica from a cluster with multiple replicas:

```bash
# Get the certificates used by UCP
$ curl https://$UCP_HOST/ca > ucp-ca.pem

$ docker run --rm -it docker/dtr remove --ucp-ca "$(cat ucp-ca.pem)"

existing-replica-id (ID of an existing replica in a cluster): 7ae3cb044b70
replica-id (Specify the replica Id. Must be unique per replica, leave blank for random): a701a510126c
username (Specify the UCP admin username): $UCP_ADMIN
password: $UCP_PASSWORD
host (Specify the UCP host using the host[:port] format): $UCP_HOST
```

Where:
* existing-replica-id: is the id of any DTR replica of that cluster,
* replica-id: is the id of the DTR replica you want to remove,
* username and password: are the username and password of a UCP administrator.


Now you can confirm on Docker Universal Control Plane that the DTR replica
`a701a510126c` no longer exists.


## Where to go next

* [Install DTR](install-dtr.md)
* [Install DTR offline](install-dtr-offline.md)
