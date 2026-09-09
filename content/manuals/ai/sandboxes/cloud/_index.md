---
title: Cloud sandboxes
description: Run Docker Sandboxes on Docker-managed cloud infrastructure and understand the cloud-specific command, storage, and billing model.
keywords: docker sandboxes, cloud sandboxes, sbx cloud, ai agents, agentic platform
weight: 35
---

Cloud sandboxes run AI agents on Docker-managed infrastructure instead of your
local machine. Use them when you need an isolated environment that doesn't
depend on the compute resources or virtualization support of your host.

Cloud sandboxes use the same `sbx` CLI as local sandboxes. Add the global
`--cloud` flag to send a supported command to the Cloud Sandboxes API:

```console
$ sbx --cloud ls
```

Cloud and local sandboxes have separate state and different capabilities. A
cloud sandbox can't mount a host workspace or use host hardware, and its
secrets, network policy, ports, and lifecycle are managed in the cloud. See
[Local and cloud differences](local-vs-cloud.md) before adapting a local
workflow.

## Prerequisites

To use cloud sandboxes, you need:

- The [`sbx` CLI](../install.md), version 0.42.0 or later
- A Docker account signed in through `sbx login`
- An active [Docker Agentic Platform plan](/manuals/subscription-billing/plans/docker-agentic-platform.md)

To subscribe, open [Docker Agentic Platform](https://agentic-platform.docker.com/)
and sign in. The plan is available for Docker Personal and Docker Pro accounts.

Cloud sandbox compute is metered through the Docker Agentic Platform
pay-as-you-go plan. Inference charges aren't included. Your model provider
charges for requests made with the API keys or OAuth credentials that you
configure.

## Get started

Credentials configured for local sandboxes aren't available to cloud
sandboxes. Configure a cloud credential for your agent before launching it.
For Claude Code, store an Anthropic API key:

```console
$ sbx --cloud secret set anthropic
```

Cloud sandboxes expire after one hour by default and are deleted when they
expire. Copy out work you want to keep before expiration. For other timeout
options, see [Configure expiration](usage.md#configure-expiration).

Create a sandbox without attaching, allowing access to GitHub for this example:

```console
$ sbx --cloud create --name cloud-project --allow-network github.com:443 claude
```

Cloud sandboxes don't accept a local workspace path. Clone the public
[Welcome to Docker repository](https://github.com/docker/welcome-to-docker)
inside the sandbox:

```console
$ sbx --cloud exec cloud-project git clone \
    https://github.com/docker/welcome-to-docker.git /home/agent/workspace/project
```

Attach to the agent:

```console
$ sbx --cloud attach cloud-project
```

Ask Claude to inspect `/home/agent/workspace/project` and write a description
of the application to `/home/agent/workspace/review.md`. When the file is ready,
press `Ctrl+\` to detach and leave the agent running.

Copy the result to your machine:

```console
$ sbx --cloud cp cloud-project:/home/agent/workspace/review.md ./review.md
```

Read the result, then remove the sandbox when you're finished:

```console
$ sbx --cloud rm cloud-project
```

Removal deletes files stored only in the sandbox. For your own projects, see
[Transfer files](usage.md#transfer-files) and
[Authenticate cloud agents](credentials.md) before cloning private repositories.

## Learn more

- [Local and cloud differences](local-vs-cloud.md) compares the two execution
  environments
- [Use cloud sandboxes](usage.md) covers creation, files, ports, and lifecycle
- [Authenticate cloud agents](credentials.md) covers cloud-specific secrets,
  API keys, and OpenAI OAuth
- [Manage cloud network policy](network-policy.md) covers account-level and
  sandbox-level network access
- [Move a sandbox](move.md) explains filesystem transfers between local and
  cloud environments
- [`sbx` CLI reference](/reference/cli/sbx/) lists commands and options
