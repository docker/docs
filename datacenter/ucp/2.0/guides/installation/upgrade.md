---
description: Learn how to upgrade Docker Universal Control Plane with minimal impact
  to your users.
keywords: Docker, UCP, upgrade, update
redirect_from:
- /ucp/upgrade-ucp/
- /ucp/installation/upgrade/
title: Upgrade to UCP 2.0
---

This page guides you in upgrading Docker Universal Control Plane (UCP) to
version 2.0.

Before upgrading to a new version of UCP, check the
[release notes](../release-notes.md) for the version you are upgrading to,
for information about the features, breaking changes, and
other relevant information for upgrading to a particular version.

## Plan the upgrade

As part of the upgrade process, you upgrade the CS Docker Engine
installed in each node of the cluster to version 1.12. If you're currently
running CS Docker Engine 1.11.2-cs3, all containers will be stopped during the
upgrade, causing some downtime to UCP and your applications.

You should plan for the upgrade to take place outside business hours, to ensure
there's minimal impact to your users.

Also, don't make changes to UCP configurations while you're upgrading it. That
can lead to misconfigurations that are difficult to troubleshoot.

## Backup your cluster

Before starting an upgrade, make sure your cluster is healthy. If a problem
occurs that will make it easier to find and troubleshoot any problems.

Then, [create a backup](../high-availability/backups-and-disaster-recovery.md)
of your cluster. This will allow you to recover from an existing backup if
something goes wrong during the upgrade process.

## Upgrade CS Docker Engine

For each node that is part of your cluster, upgrade the CS Docker Engine
installed on that node to CS Docker Engine version 1.12 or higher.

Starting with the controller nodes, and then worker nodes:

1. Log into the node using ssh.
2. Upgrade the CS Docker Engine to version 1.12 or higher.

    If you're upgrading from CS Docker Engine 1.11.3 or previous this will cause
    some downtime on that node, since all containers will be stopped.

    Containers that have a restart policy set to
    'always', are automatically started after the upgrade. This is the case of
    UCP and DTR components. All other containers need to be started manually.

3. Make sure the node is healthy.

    In your browser, navigate to the **UCP web UI**, and validate that the
    node is healthy and is part of the cluster.

> Swarm mode
>
> UCP 2.0 and higher requires swarm mode. Upgrading from a UCP 1.x version
> enables swarm mode in CS Docker Engine.

## Upgrade the first controller node

Start by upgrading a controller node that has valid root CA material. This
can be the first node where you installed UCP or any controller replica
that you've installed using that node's root CA material.

1. Log into the controller node using ssh.
2.  Pull the docker/ucp image for the version you want to upgrade to.

    ```bash
    # Check on Docker Hub which versions are available
    $ docker pull docker/ucp:<version>
    ```

3.  Upgrade UCP by running:

    ```bash
    $ docker run --rm -it \
      --name ucp \
      -v /var/run/docker.sock:/var/run/docker.sock \
      docker/ucp:<version> \
      upgrade --interactive
    ```

    This runs the upgrade command in interactive mode, so that you are prompted
    for any necessary configuration values.

    The upgrade command makes configuration changes to Docker Engine.
    You are prompted to restart the Docker Engine and run the upgrade
    command again, to continue the upgrade.

4. Make sure the node is healthy.

    In your browser, navigate to the **UCP web UI**, and validate that the
    node is healthy and is part of the cluster.

## Upgrade other nodes

Follow the procedure described above to upgrade other nodes in the cluster.
Start by upgrading the remaining controller nodes, and then upgrade any worker
nodes.

## Where to go next

* [UCP release notes](../release-notes.md)
