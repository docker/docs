---
description: Inspect clusters
keywords: documentation, docs, docker, cluster, infrastructure, automation
title: docker cluster inspect
notoc: true
---

## Usage
```
docker cluster inspect [OPTIONS] cluster
```
Use the following options as needed to display detailed information about a cluster:

- `-a, --all`: Displays complete information about the cluster. 
- `--dry-run`: Skips resource provisioning.
- `--log-level string`: Specifies the logging level. Valid values include: `trace`,`debug`,`info`,`warn`,`error`, and `fatal`. Defaults to `warn`.
