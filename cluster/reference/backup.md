---
description: Back up a running cluster
keywords: documentation, docs, docker, cluster, infrastructure, automation
title: docker cluster backup
notoc: true
---

## Usage
```
docker cluster backup [OPTIONS] cluster
```

Use the following options as needed to back up a running cluster:

- `--dry-run`: Skips resource provisioning.
- `--file string`: Specifies a cluster backup filename. Defaults to `backup.tar.gz`.
- `--log-level string`: Specifies the logging level. Valid values include: `trace`,`debug`,`info`,`warn`,`error`, and `fatal`. 
Defaults to `warn`.
- `--passphrase string`: Specifies a cluster backup passphrase.

The backup command performs a full Docker Cluster backup following the steps found in [Backup and Restore Best Practices](https://success.docker.com/article/backup-restore-best-practices).
