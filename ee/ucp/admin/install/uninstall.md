---
title: Uninstall UCP
description: Learn how to uninstall a Docker Universal Control Plane swarm.
keywords: UCP, uninstall, install, Docker EE
---

Docker UCP is designed to scale as your applications grow in size and usage.
You can [add and remove nodes](../configure/scale-your-cluster.md) from the
cluster to make it scale to your needs.

You can also uninstall Docker Universal Control Plane from your cluster. In this
case the UCP services are stopped and removed, but your Docker Engines will
continue running in swarm mode. You applications will continue running normally.

If you wish to remove a single node from the UCP cluster, you should instead
[Remove that node from the cluster](../configure/scale-your-cluster.md).

After you uninstall UCP from the cluster, you'll no longer be able to enforce
role-based access control to the cluster, or have a centralized way to monitor
and manage the cluster.

After uninstalling UCP from the cluster, you will no longer be able to join new
nodes using `docker swarm join`, unless you reinstall UCP.

To uninstall UCP, log in to a manager node using ssh, and run the following
command:

```bash
docker container run --rm -it \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --name ucp \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} uninstall-ucp --interactive
```

This runs the uninstall command in interactive mode, so that you are prompted
for any necessary configuration values. Running this command on a single manager
node will uninstall UCP from the entire cluster. [Check the reference
documentation](/reference/ucp/3.0/cli/index.md) to learn the options available
in the `uninstall-ucp` command.

## Swarm mode CA

After uninstalling UCP, the nodes in your cluster will still be in swarm mode,
but you can't join new nodes until you reinstall UCP, because swarm mode
relies on UCP to provide the CA certificates that allow nodes in the cluster
to identify one another. Also, since swarm mode is no longer controlling its
own certificates, if the certificates expire after you uninstall UCP, the nodes
in the swarm won't be able to communicate at all. To fix this, either reinstall
UCP before the certificates expire or disable swarm mode by running
`docker swarm leave --force` on every node.

## Restore IP tables

When you install UCP, the Calico network plugin changes the host's IP tables.
When you uninstall UCP, the IP tables aren't reverted to their previous state.
After you uninstall UCP, restart the node to restore its IP tables.

## Where to go next

- [Join nodes to your cluster](../configure/join-nodes.md)
