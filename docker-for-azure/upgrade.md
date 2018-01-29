---
description: Upgrading your stack
keywords: azure, microsoft, iaas, tutorial
title: Docker for Azure upgrades
---

Docker for Azure supports upgrading from one version to the next within a specific channel. Upgrading across channels (for example, from `edge` to `stable` or `test` to `edge`) is not supported. To upgrade to a specific version, run the upgrade container corresponding to the target version for the upgrade. An upgrade of Docker for Azure involves:

 * Upgrading the VHD backing the manager and worker nodes (Docker ships in the VHD)
 * Upgrading service containers in the manager and worker nodes
 * Changing any other resources in the Azure Resource Group that hosts Docker for Azure

## Prerequisites

 * We recommend only attempting upgrades of swarms with at least 3 managers. A 1-manager swarm can't maintain quorum during an upgrade.
 * You can only upgrade one version at a time. Skipping a version during an upgrade is not supported.
 * Downgrades are not tested.
 * Upgrading across channels (`stable`, `edge`, or `testing`) is not supported.
 * If the swarm contains nodes in the `down` state, remove them from the swarm before attempting the upgrade, using `docker node rm <node-id>`.

## Upgrading

New releases are announced on the [Release Notes](release-notes.md) page.

To initiate an upgrade, connect a manager node using SSH and run the container corresponding to the version you want to upgrade to:

```bash
$ docker run \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /usr/bin/docker:/usr/bin/docker \
  -ti \
  docker4x/upgrade-azure:version-tag
```

For example, this command upgrades from 17.03 CE or 17.06.0 CE stable release to 17.06.1 CE stable:

```bash
$ docker run \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /usr/bin/docker:/usr/bin/docker \
  -ti \
  docker4x/upgrade-azure:17.06.1-ce-azure1
```

If you are already on a version more recent than 17.06.1 CE or 17.07.0 CE, you can use `upgrade.sh` to initiate the upgrade:

```bash
$ upgrade.sh YY.MM.X-ce-azureC
```

This example upgrades Docker for Azure from 17.06.1 CE Edge to 17.07.0 CE Edge:

```bash
$ upgrade.sh 17.07.0-ce-azure1
```

This initiates a rolling upgrade of the Docker swarm. Service state is maintained during and after the upgrade. Appropriately scaled services should not experience downtime during an upgrade. Single containers which are not part of services (for example, containers started with `docker run`) are **not** preserved during an upgrade. This is because they are not Docker services but are known only to the individual Docker engine where they are running.

## Monitoring

The upgrade process may take several minutes to complete. You can follow the progress of the upgrade either from the output of the upgrade command or from the Azure UI by going to the Azure UI blades corresponding to the Virtual Machine Scale Sets hosting your manager and worker nodes. The URL for the Azure UI blades corresponding to your subscription and resource group is printed in the output of upgrade command above. They follow the following format:

```none
https://portal.azure.com/#resource/subscriptions/[subscription-id]/resourceGroups/[resource-group]/providers/Microsoft.Compute/virtualMachineScaleSets/swarm-manager-vmss/instances
```

```none
https://portal.azure.com/#resource/subscriptions/[subscription-id]/resourceGroups/[resource-group]/providers/Microsoft.Compute/virtualMachineScaleSets/swarm-worker-vmss/instances
```
`[subscription-id]` and `[resource-group]` are placeholders which are replaced by
real values.

In the last stage of the upgrade, the manager node where the upgrade is initiated from needs to be shut down, upgraded and reimaged. During this time, you can't access the node and if you were logged in, your SSH connection drop.

## Post upgrade

After the upgrade, the IP address and the port needed to SSH into the manager nodes do not change. However, the host identity of the manager nodes does change as the VMs get reimaged. So when you try to SSH in after a successful upgrade, your SSH client may suspect a Man-In-The-Middle attack. You need to delete the old entries in your SSH client's store [typically `~/.ssh/known_hosts`] for new SSH connections to succeed to the upgraded manager nodes.

## Changing instance sizes and other template parameters

This is not supported for Azure at the moment.
