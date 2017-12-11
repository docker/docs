---
title: Backup Docker EE
description: Learn how to create a backup of your Docker Enterprise Edition, and
  how to restore from a backup
keywords: enterprise, backup, restore
---

To backup Docker Enterprise Edition you need to create individual backups
for each of the following components:

1. Docker Swarm. [Backup Swarm resources like service and network definitions](/engine/swarm/admin_guide.md#back-up-the-swarm).
2. Universal Control Plane (UCP). [Backup UCP configurations](/datacenter/ucp/2.2/guides/admin/backups-and-disaster-recovery.md).
3. Docker Trusted Registry (DTR). [Backup DTR configurations and images](/datacenter/dtr/2.3/guides/admin/backups-and-disaster-recovery.md).

Before proceeding to backup the next component, you should test the backup you've
created to make sure it's not corrupt. One way to test your backups is to do
a fresh installation in a separate infrastructure and restore the new installation
using the backup you've created.

If you create backups for a single component, you won't be able to restore your
deployment to its previous state.

## Restore Docker Enterprise Edition

You should only restore from a backup as a last resort. If you're running Docker
Enterprise Edition in high-availability you can remove unhealthy nodes from the
swarm and join new ones to bring the swarm to an healthy state.

To restore Docker Enterprise Edition, you need to restore the individual
components one by one:

1. Docker Engine. [Learn more](/engine/swarm/admin_guide.md#recover-from-disaster).
2. Universal Control Plane (UCP). [Learn more](/datacenter/ucp/2.2/guides/admin/backups-and-disaster-recovery.md#restore-your-swarm).
3. Docker Trusted Registry (DTR). [Learn more](/datacenter/dtr/2.3/guides/admin/backups-and-disaster-recovery.md#restore-dtr-data).

## Where to go next

* [Upgrade Docker EE](upgrade.md)
