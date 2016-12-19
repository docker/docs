---
description: Learn how to uninstall your Docker Trusted Registry installation.
keywords: docker, dtr, install, uninstall
title: Uninstall Docker Trusted Registry
---

Use the `remove` command, to remove a DTR replica from an existing deployment.
To uninstall a DTR cluster you remove all DTR replicas one at a time.

The remove command informs the DTR cluster that the node is about to be removed,
then it removes the replica, stops and removes all DTR containers from that node,
and deletes all DTR volumes.

To uninstall a DTR replica, run:

```none
docker run -it --rm \
  docker/dtr remove \
  --ucp-insecure-tls
```

You will be prompted for:

* Existing replica id: the id of any healthy DTR replica of that cluster
* Replica id: the id of the DTR replica you want to remove. It can be the id of an
unhealthy replica that you want to remove from your deployment
* UCP username and password: the administrator credentials for UCP

To ensure you don't loose data, DTR will not remove the last replica from your
deployment. To confirm you really want to remove that replica, use the
`--force-remove` flag.

To see what options are available in the uninstall command, check the
[uninstall command reference documentation](../../reference/cli/remove.md).

## Where to go next

* [Scale your deployment](scale-your-deployment.md)
* [Install DTR](index.md)
