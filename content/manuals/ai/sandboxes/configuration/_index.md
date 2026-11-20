---
title: Configure Docker Sandboxes
linkTitle: Configuration
weight: 60
description: Configure Docker Sandboxes settings, credentials, models, project environments, GPU passthrough, registry mirrors, and upstream proxies.
keywords: docker sandboxes, sbx, configuration, settings, credentials, models, environment files, gpu passthrough, registry mirror, upstream proxy
---

Configure credentials and how Docker Sandboxes run for a project, host, or
network environment. These settings control sandbox creation, authentication,
and connectivity. To change the tools and agent configuration inside a
sandbox, see [Customize](../customize/).

- [Settings](settings.md) lists host-level settings, environment variable
  equivalents, and commands to inspect and change values.
- [Credentials](credentials.md) configures API keys, authentication
  credentials, and registry access for sandboxed agents.
- [Models](models.md) selects local models, hosted providers, or custom
  inference endpoints for sandboxed agents.
- [Environment files](environment-files.md) declare reusable project
  configuration in `sbxenv.yaml` for local or cloud sandboxes.
- [GPU passthrough](gpu-passthrough.md) configures a Linux host and sandbox for
  NVIDIA GPU workloads.
- [Registry mirror](registry-mirror.md) routes Docker Hub template, kit, and
  in-sandbox Docker image pulls through an organization's registry mirror.
- [Upstream proxy](upstream-proxy.md) routes sandbox and daemon traffic through
  an operating system or corporate proxy.

## Guest page size

For local sandboxes, use `--page-size` with `sbx create` or `sbx run` to choose
the guest kernel's memory page size. For example, select `4k` for workloads
that require 4 KB pages:

```console
$ sbx create --page-size 4k claude .
```

Without the flag, the guest uses the host's page size when the matching kernel
is installed, and falls back to `4k` otherwise.

Supported values depend on the host platform:

| Host platform | Page sizes |
| :------------ | :--------- |
| Linux arm64 | `4k`, `16k`, `64k` |
| macOS on Apple silicon | `4k`, `16k` |
| Other platforms | `4k` |

The CPU must support the selected page size, and the corresponding guest kernel
must be installed. The experimental `--4k` flag is equivalent to `--page-size 4k`;
combining `--4k` with another page size is rejected.

Page size is set when you create a sandbox and persists across restarts.
When `sbx run` reconnects to an existing sandbox, `--page-size` has no effect.
To choose a different page size, create another sandbox.
