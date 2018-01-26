---
title: Upgrade to UCP 2.2
description: Learn how to upgrade Docker Universal Control Plane with minimal impact to your users.
keywords: UCP, upgrade, update
---

This page guides you in upgrading Docker Universal Control Plane (UCP) to
version 2.2.

Before upgrading to a new version of UCP, check the
[release notes](../../release-notes/index.md) for this version for information
about new features, breaking changes, and other relevant information about
upgrading to a particular version.

## Plan the upgrade

As part of the upgrade process, you upgrade the Docker Engine instance installed
in each node of the swarm to version 17.06 Enterprise Edition. Plan for the
upgrade to take place outside of business hours, to ensure there's minimal
impact to your users.

Also, don't make changes to UCP configurations while you're upgrading it.
This can lead to misconfigurations that are difficult to troubleshoot.

## Back up your swarm

Before starting an upgrade, make sure that your swarm is healthy. If a problem
occurs, this makes it easier to find and troubleshoot it.

[Create a backup](../backups-and-disaster-recovery.md) of your swarm.
This allows you to recover if something goes wrong during the upgrade process.

> Upgrading and backup archives
>
> The backup archive is version-specific, so you can't use it during the
> upgrade process. For example, if you create a backup archive for a UCP 2.1
> swarm, you can't use the archive file after you upgrade to UCP 2.2.

## Upgrade Docker Engine

For each node that is part of your swarm, upgrade the Docker Engine
installed on that node to Docker Engine version 17.06 or higher. Be sure
to install the Docker Enterprise Edition.

Starting with the manager nodes, and then worker nodes:

1. Log into the node using ssh.
2. Upgrade the Docker Engine to version 17.06 or higher.
3. Make sure the node is healthy.

    In your browser, navigate to the **Nodes** page in the UCP web UI,
    and check that the node is healthy and is part of the swarm.

> Swarm mode
>
> UCP 2.0 and higher requires swarm mode. Upgrading from a UCP 1.x version
> enables swarm mode in Docker EE Engine.

## Upgrade UCP

You can upgrade UCP from the web UI or the CLI.

### Use the UI to perform an upgrade

When an upgrade is available for a UCP installation, a banner appears.

![](../../images/upgrade-ucp-1.png){: .with-border}

Clicking this message takes an admin user directly to the upgrade process.
It can be found under the **Cluster Configuration** tab of the **Admin
 Settings** section.

![](../../images/upgrade-ucp-2.png){: .with-border}

Select a version to upgrade to using the **Available UCP Versions** dropdown,
then click to upgrade.

Before the upgrade happens, a confirmation dialog along with important
information regarding swarm and UI availability is displayed.

![](../../images/upgrade-ucp-3.png){: .with-border}

During the upgrade, the UI is unavailable, so wait until the upgrade is complete
before trying to use the UI. When the upgrade completes, a notification alerts
you that a newer version of the UI is available, and you can see the new UI
after you refresh your browser.

### Use the CLI to perform an upgrade

To upgrade from the CLI, log into a UCP manager node using ssh, and run:

```
# Get the latest version of UCP
$ docker image pull {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }}

$ docker container run --rm -it \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} \
  upgrade --interactive
```

This runs the upgrade command in interactive mode, so that you are prompted
for any necessary configuration values.

Once the upgrade finishes, navigate to the UCP web UI and make sure that
all the nodes managed by UCP are healthy.

## Recommended upgrade paths

If you're running a UCP version that's lower than 2.1, first upgrade to the
latest 2.1 version, then upgrade to 2.2. Use these rules for your upgrade
path to UCP 2.2:

- From UCP 1.1: UCP 1.1 -> UCP 2.1 -> UCP 2.2
- From UCP 2.0: UCP 2.0 -> UCP 2.1 -> UCP 2.2
- From UCP 2.1: UCP 2.1 -> UCP 2.2

## Where to go next

* [UCP release notes](../../release-notes/index.md)
