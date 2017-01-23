---
description: Upgrading your stack
keywords: azure, microsoft, iaas, tutorial
title: Docker for Azure Upgrades
---

Docker for Azure supports upgrading from one version to the next. To upgrade, apply a new version of the Azure ARM template that powers Docker for Azure. An upgrade of Docker for Azure involves:

 * Upgrading the VHD backing the manager and worker nodes (the Docker engine ships in the VHD)
 * Upgrading service containers in the manager and worker nodes
 * Changing any other resources in the Azure Resource Group that hosts Docker for Azure

## Prerequisites

 * We recommend only attempting upgrades of swarms with at least 3 managers. A 1-manager swarm may not be able to maintain quorum during the upgrade
 * You can only upgrade one version at a time. Skipping a version during an upgrade is not supported. Downgrades are not tested.
 * Ensure there are no nodes in the swarm in "down" status. If there are such nodes in the swarm, please remove them from the swarm using `docker node rm node-id`

## Upgrading

New releases are announced on the [Release Notes](release-notes.md) page.

To initiate an upgrade, SSH into a manager node and issue the following command:

    upgrade.sh https://download.docker.com/azure/stable/Docker.tmpl

This will initiate a rolling upgrade of the Docker swarm and service state will be maintained during and after the upgrade. Appropriately scaled services should not experience downtime during an upgrade. Note that single containers started (for example) with `docker run -d` are **not** preserved during an upgrade. This is because they are not Docker Swarm services but are known only to the individual Docker engines.

## Monitoring

The upgrade process may take several minutes to complete. You can follow the progress of the upgrade either from the output of the upgrade command or from the Azure UI by going to the Azure UI blades corresponding to the Virtual Machine Scale Sets hosting your manager and worker nodes. The URL for the Azure UI blades corresponding to your subscription and resource group is printed in the output of upgrade command above. They follow the following format:

https://portal.azure.com/#resource/subscriptions/[subscription-id]/resourceGroups/[resource-group]/providers/Microsoft.Compute/virtualMachineScaleSets/swarm-manager-vmss/instances

https://portal.azure.com/#resource/subscriptions/[subscription-id]/resourceGroups/[resource-group]/providers/Microsoft.Compute/virtualMachineScaleSets/swarm-worker-vmss/instances

Note that in the last stage of the upgrade, the manager node where the upgrade is initiated from needs to be shut down, upgraded and reimaged. During this time, you won't be able to access the node and if you were logged in, your SSH connection will drop.

## Post Upgrade

After the upgrade, the IP address and the port needed to SSH into the manager nodes do not change. However, the host identity of the manager nodes does change as the VMs get reimaged. So when you try to SSH in after a successful upgrade, your SSH client may suspect a Man-In-The-Middle attack. You will need to delete the old entries in your SSH client's store [typically `~/.ssh/known_hosts`] for new SSH connections to succeed to the upgraded manager nodes.

## Changing instance sizes and other template parameters

This is not supported for Azure at the moment.
