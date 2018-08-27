---
title: Manage Kubernetes api-server flags
description: Learn how to tweak Kubernetes api-server flags managed by Docker Enterprise Edition.
keywords: Docker EE, UCP, Kubernetes, orchestration, cluster, api-server
---

Docker Enterprise Edition (EE) sets pre-defined configuration of Kubernetes components. [Docker EE configuration file]((./ucp-configuration-file.md)) provides setting to adjust Kubernetes `api-server` configuration flags when necessary.
This topic outlines steps necessary to modify default configuration of Kubernetes `api-server` component managed by Docker EE.

## Fetch and save current Docker EE configuration

Follow steps outlined in [inspect and modify existing configuration](./ucp-configuration-file.md#inspect-and-modify-existing-configuration) section to save Docker EE configuration into a `.toml` file.

## Modify `api-server` configuration flags

Set necessary configuration flags for Kubernetes `api-server` in saved `.toml` configuration file.
For example:
```
.......
[cluster_config]
  .......
  custom_kube_api_server_flags = ["--event-ttl=1h0m0s", "--service-node-port-range=30000-32767"]
.......
```

> Note, given example of `--service-node-port-range` configuration flag assumes that `Swarm` workloads will not be executed on Docker EE platform.

## Push updated configuration into Docker EE

Follow steps outlined in second part of [inspect and modify existing configuration](./ucp-configuration-file.md#inspect-and-modify-existing-configuration) section to determine new name of Docker EE config object and push new configuration.
Configuration change will restart `api-server` component and apply configuration flags provided via `custom_kube_api_server_flags` setting.

> Follow [example configuration file](./ucp-configuration-file.md#example-configuration-file) instructions to view extensive list of Docker EE configuration settings.
