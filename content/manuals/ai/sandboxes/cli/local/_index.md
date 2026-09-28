---
title: Local configuration
linkTitle: Local configuration
weight: 130
description: Configure Docker Sandboxes settings, credentials, models, project environments, GPU passthrough, registry mirrors, and upstream proxies.
keywords: docker sandboxes, sbx, configuration, settings, credentials, models, environment files, gpu passthrough, registry mirror, upstream proxy
aliases:
  - /ai/sandboxes/configuration/
---

These settings apply to the machine running local sandboxes. They do not
configure Docker-managed cloud infrastructure.

- [Runtime settings](/manuals/ai/sandboxes/cli/local/settings.md) control the local daemon and CLI behavior.
- [Model providers](/manuals/ai/sandboxes/cli/local/models.md) configure local or hosted inference endpoints.
- [Agent skills](/manuals/ai/sandboxes/cli/local/agent-skills.md) share a host-side skill store with agents.
- [GPU passthrough](/manuals/ai/sandboxes/cli/local/gpu-passthrough.md) configures supported host hardware.
- [Upstream proxy](/manuals/ai/sandboxes/cli/local/upstream-proxy.md) routes traffic through your host network.
- [Registry mirror](/manuals/ai/sandboxes/cli/local/registry-mirror.md) configures local image pulls.

For tasks that also have cloud workflows, see [CLI credentials](/manuals/ai/sandboxes/cli/credentials.md),
[environment files](/manuals/ai/sandboxes/cli/environment-files.md), and [network controls](/manuals/ai/sandboxes/concepts/network.md).
