---
title: Docker Sandboxes
description: Run AI coding agents in isolated environments
keywords: docker sandboxes, sbx, ai agents, sandboxed agents, microVM
weight: 10
params:
  sidebar:
    group: AI and agents
    badge:
      color: blue
      text: Early Access
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

Docker Sandboxes run AI coding agents in isolated microVM sandboxes. Each
sandbox gets its own Docker daemon, filesystem, and network — the agent can
build containers, install packages, and modify files without touching your host
system.

Organization admins can
[centrally manage sandbox network and filesystem policies](security/governance.md)
from the Docker Admin Console, so the same rules apply uniformly across every
developer's machine. Available on a separate paid subscription.

## Get started

Install the `sbx` CLI and sign in:

{{< tabs >}}
{{< tab name="macOS" >}}

```console
$ brew install docker/tap/sbx
$ sbx login
```

{{< /tab >}}
{{< tab name="Windows" >}}

```powershell
> winget install -h Docker.sbx
> sbx login
```

{{< /tab >}}
{{< tab name="Linux (Ubuntu)" >}}

```console
$ curl -fsSL https://get.docker.com | sudo REPO_ONLY=1 sh
$ sudo apt-get install docker-sbx
$ sudo usermod -aG kvm $USER
$ newgrp kvm
$ sbx login
```

{{< /tab >}}
{{< /tabs >}}

Then launch an agent in a sandbox:

```console
$ cd ~/my-project
$ sbx run claude
```

See the [get started guide](get-started.md) for a full walkthrough, or jump to
the [usage guide](usage.md) for common patterns.

## Learn more

- [Agents](agents/) — supported agents and per-agent configuration
- [Customize](customize/) — reusable templates and declarative kits for
  extending or tailoring sandboxes
- [Architecture](architecture.md) — microVM isolation, workspace mounting,
  networking
- [Security](security/) — isolation model, credential handling, network
  policies, workspace trust
- [CLI reference](/reference/cli/sbx/) — full list of `sbx` commands and options
- [Troubleshooting](troubleshooting.md) — common issues and fixes
- [FAQ](faq.md) — login requirements, telemetry, etc

## Feedback

Your feedback shapes what gets built next. If you run into a bug, hit a
missing feature, or have a suggestion, open an issue at
[github.com/docker/sbx-releases/issues](https://github.com/docker/sbx-releases/issues).
