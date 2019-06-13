---
description: Update a cluster
keywords: documentation, docs, docker, cluster, infrastructure, automation
title: docker cluster update
notoc: true
---

## Usage
```
docker cluster update [Options] cluster
```
Use the following options as needed to update a running cluster's desired state:

Options:

- `--dry-run`: Skips resource provisioning.
- `-f`, `--file string`: Specfies cluster definition.
- `--log-level string`: Specifies the logging level. Valid values include: `trace`,`debug`,`info`,`warn`,`error`, and `fatal`. Defaults to `warn`.
