---
description: Restore to a running cluster
keywords: documentation, docs, docker, cluster, infrastructure, automation
title: docker cluster restore
notoc: true
---

## Usage
```
docker cluster restore [OPTIONS] cluster
```
Use the following options as needed to restore a cluster from a backup:

- `--dry-run`: Skips resource provisioning.
- `--file string`: Specifies a cluster backup filename. Defaults to `backup.tar.gz`.
- `--log-level string`: Specifies the logging level. Valid values include: 
`trace`,`debug`,`info`,`warn`,`error`, and `fatal`. Defaults to `warn`.
- `--passphrase string`: Specifies a cluster backup passphrase.

The restore command performs a full Docker Cluster restore following the steps found in [Backup and Restore Best Practices](https://success.docker.com/article/backup-restore-best-practices).
