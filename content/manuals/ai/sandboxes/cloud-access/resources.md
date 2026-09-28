---
title: Cloud compute sizes and quotas
linkTitle: Compute and quotas
description: Choose cloud compute resources and understand shared account quotas.
keywords: docker sandboxes, cloud compute, sizes, account quotas
weight: 10
---

Cloud CLI, Console, and API workflows share account quotas. Each interface
selects resources for the sandbox it creates.

## Compute sizes

| Size | CPUs | Memory | Memory in MiB |
| --- | ---: | ---: | ---: |
| `micro` | 1 | 2 GiB | 2048 |
| `small` | 2 | 4 GiB | 4096 |
| `medium` | 4 | 8 GiB | 8192 |
| `large` | 8 | 16 GiB | 16384 |
| `xl` | 16 | 32 GiB | 32768 |

Your account access and available capacity determine whether a request can be
accepted. Resource selection differs by interface: use
[CLI resource flags](/manuals/ai/sandboxes/cli/manage/cloud.md#choose-resources-and-platform), the
[Console compute picker](/manuals/ai/sandboxes/console/sandboxes.md), or
[SDK resource options](/manuals/ai/sandboxes/api/limits.md#compute-sizes).

## Account quotas

The default quotas apply across an account:

| Resource | Default limit |
| --- | ---: |
| Concurrent sandboxes | 10 |
| Stored sandboxes | 50 |
| Volumes | 100 |
| Secrets | 100 |
| Images being prepared at the same time | 3 |

Your account can have different quotas. Confirm your account's limits with
Docker before planning a workload that depends on a particular allowance.

Stopping a sandbox releases its concurrency slot unless the sandbox is
configured as always-on. Restarting a stopped sandbox requires a concurrency
slot. A stopped sandbox still counts toward the stored sandbox quota.
Delete sandboxes you no longer need to reduce stored usage.

Handle quota errors even if you checked usage before creating a resource.
Other applications can consume the remaining allowance between requests.
Reduce concurrency or remove unused resources before retrying.

