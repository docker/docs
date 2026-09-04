---
title: Configure Docker Sandboxes
linkTitle: Configuration
weight: 60
description: Configure credentials, project environments, GPU passthrough, and upstream proxy settings for Docker Sandboxes.
keywords: docker sandboxes, sbx, configuration, credentials, environment files, gpu passthrough, upstream proxy
---

Configure credentials and how Docker Sandboxes run for a project, host, or
network environment. These settings control sandbox creation, authentication,
and connectivity. To change the tools and agent configuration inside a
sandbox, see [Customize](../customize/).

- [Credentials](credentials.md) configures API keys, authentication
  credentials, and registry access for sandboxed agents.
- [Environment files](environment-files.md) declare reusable project
  configuration in `sbxenv.yaml`.
- [GPU passthrough](gpu-passthrough.md) configures a Linux host and sandbox for
  NVIDIA GPU workloads.
- [Upstream proxy](upstream-proxy.md) routes sandbox and daemon traffic through
  an operating system or corporate proxy.
