---
description: Learn how to uninstall a Docker Universal Control Plane cluster.
keywords: docker, ucp, uninstall
title: Uninstall UCP
---

Docker UCP is designed to scale as your applications grow in size and usage.
You can [add and remove nodes](scale-your-cluster.md) from the cluster, to make
it scale to your needs.

You can also uninstall Docker Universal Control plane from your cluster. In this
case the UCP services are stopped and removed, but your Docker Engines will
continue running in swarm mode. You applications will continue running normally.

After you uninstall UCP from the cluster, you can't enforce role-based access
control to the cluster or centrally monitor and manage the cluster. You can
no longer join new nodes to the cluster using `docker swarm join`.

To uninstall UCP, log in into a manager node using ssh, and run:

```bash
$ docker run --rm -it \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --name ucp \
  docker/ucp uninstall-ucp --interactive
```

This runs the uninstall command in interactive mode, so that you are prompted
for any necessary configuration values.
[Check the reference documentation](../../reference/cli/index.md) to learn the options
available in the `uninstall-ucp` command.

## Where to go next

* [Scale your cluster](scale-your-cluster.md)
