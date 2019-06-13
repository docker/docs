---
description: Remove a cluster 
keywords: documentation, docs, docker, cluster, infrastructure, automation
title: docker cluster rm
notoc: true
---

## Usage
```
docker cluster rm [OPTIONS] cluster
```
Use the following options as needed when removing a cluster:

- `--dry-run`: Skips resource provisioning.
- `-f`, `--force`: Forces removal of the cluster files.
- `--log-level string`: Specifies the logging level. Valid values include: `trace`,`debug`,`info`,`warn`,`error`, and `fatal`. Defaults to `warn`.
