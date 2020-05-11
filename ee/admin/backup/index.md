---
title: Back up Docker Enterprise
description: Learn how to create a backup of your Docker Enterprise.
keywords: enterprise, backup
redirect_from:
  - /enterprise/backup/
  - /ee/backup/
  - /ee/ucp/admin/backups-and-disaster-recovery/
---

## Introduction
This document provides instructions and best practices for Docker Enterprise backup procedures for all components of the platform.

> **Important**: Make sure you perform regular backups for Docker Enterprise, including prior to an upgrade or uninstallation.

## Prerequisites

- Before performing a backup or restore operation for any component of Docker Enterprise, you must have healthy managers. Otherwise, disaster recovery procedures are involved.   
- Have adequate space available for backup contents.

## Procedure
To back up Docker Enterprise, you must create individual backups
for each of the following components:

1. [Back up Docker Swarm](back-up-swarm.md). Back up Swarm resources like service and network definitions.
2. [Back up Universal Control Plane (UCP)](back-up-ucp.md). Back up UCP configurations.
3. [Back up Docker Trusted Registry (DTR)](back-up-dtr.md). Back up DTR configurations, images, and metadata.

If you do not create backups for all components, you cannot restore your deployment to its previous state. 

Test each backup you create. One way to test your backups is to do
a fresh installation on a separate infrastructure with the backup. Refer to [Restore Docker Enterprise](/ee/admin/restore/)  for additional information.

**Note**: Application data backup is **not** included in this information. Persistent storage data backup is the responsibility of the storage provider for the storage plugin or driver.

### Where to go next

- [Back up Docker Swarm](back-up-swarm.md)
